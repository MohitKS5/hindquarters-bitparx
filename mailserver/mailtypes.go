package mailserver

import (
	"net/mail"
	"github.com/bitparx/common/config"
)

type Mail struct {
	ServerName  string
	To          mail.Address
	From        mail.Address
	Receipients []string
	Subject     string
	Body        string
}

func getDefauts() Mail{
	return Mail{
		"smtp.gmail.com:587",
		mail.Address{"", config.EMAIL},
		mail.Address{"", "admin@antaragni.in"},
		[]string{},
		"test",
		"test",
	}
}
