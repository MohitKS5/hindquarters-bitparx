package mailserver

import (
	"net/smtp"
	"log"
	"net/mail"
	"fmt"
	"net"
	"crypto/tls"
)

func sendSingleMail(body string) {
	from := "mohitkumarsingh907@gmail.com"
	pass := "TripurariSingh"
	to := "admin@antaragni.in"
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Hello there\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("sent")
}

func sendMultiple() {

	from := mail.Address{"", "mohitkumarsingh907@gmail.com"}
	to := mail.Address{"", "admin@antaragni.in"}
	pass:="TripurariSingh"
	subj := "This is the email subject"
	body := "This is an example body.\n With two lines."

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	servername := "smtp.gmail.com:587"

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", from.Address, pass, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	client, err := smtp.Dial(servername)
	if err != nil {
		log.Panic(err)
	}

	client.StartTLS(tlsconfig)

	// Auth
	if err = client.Auth(auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = client.Mail(from.Address); err != nil {
		log.Panic(err)
	}

	if err = client.Rcpt(to.Address); err != nil {
		log.Panic(err)
	}
	if err = client.Rcpt("mohitks@iitk.ac.in"); err != nil {
		log.Panic(err)
	}

	// Data
	w, err := client.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	client.Quit()

}