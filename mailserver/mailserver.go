package mailserver

import (
	"github.com/bitparx/common/config"
)

func (newMail Mail) SendMails() {
	newMail = getDefauts()
	sendMultiple(
		newMail.To,
		newMail.From,
		newMail.Receipients,
		config.PASS,
		newMail.Subject,
		newMail.Body,
		newMail.ServerName)
}
