package mailserver

import (
	"net/smtp"
	"log"
)

func send(body string) {
	from := "mohitkumarsingh907@gmail.com"
	pass := ""
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
