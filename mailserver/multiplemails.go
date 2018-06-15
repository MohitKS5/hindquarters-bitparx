package mailserver

import (
	"net/smtp"
	"log"
	"net/mail"
	"fmt"
	"net"
	"crypto/tls"
)

func sendMultiple(to, from mail.Address, receipients []string,
	pass, subject, body, servername string) {

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", from.Address, pass, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	client, err := smtp.Dial(servername)
	if err != nil {
		log.Println(err)
		return
	}

	client.StartTLS(tlsconfig)

	// Auth
	if err = client.Auth(auth); err != nil {
		log.Println(err)
		return
	}

	// To && From
	if err = client.Mail(from.Address); err != nil {
		log.Println(err)
		return
	}

	if err = client.Rcpt(to.Address); err != nil {
		log.Println(err)
		return
	}

	for r := range receipients {
		if err = client.Rcpt(receipients[r]); err != nil {
			log.Println(err)
			return
		}
	}

	// Data
	w, err := client.Data()
	if err != nil {
		log.Println(err)
		return
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Println(err)
		return
	}

	err = w.Close()
	if err != nil {
		log.Println(err)
		return
	}

	client.Quit()

}
