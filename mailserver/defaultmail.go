package mailserver

import (
	"github.com/bitparx/common/config"
	"strings"
	"net/smtp"
	"log"
)

func sendDefaultMail(to, subject, body string) {
	from := config.EMAIL
	pass := config.PASS
	stringArr := []string{
		"From: " 	+ from,
		"To: " 		+ to,
		"Subject: " + subject,
		"", 		body,
	}
	msg := strings.Join(stringArr, `\n`)

	err :=
		smtp.SendMail("smtp.gmail.com:587",
			smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
			from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("sent")
}
