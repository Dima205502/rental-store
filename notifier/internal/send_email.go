package internal

import (
	"net/smtp"
	"notifier/config"
	"strconv"
)

var sender, appPassword, smtpHost string
var smtpPort int

func Init(cfg *config.Config) {
	sender = cfg.Sender
	appPassword = cfg.AppPassword
	smtpHost = cfg.SmtpHost
	smtpPort = cfg.SmtpPort
}

func Send(email, msg string) error {
	auth := smtp.PlainAuth("", sender, appPassword, smtpHost)
	return smtp.SendMail(smtpHost+":"+strconv.Itoa(smtpPort), auth, sender, []string{email}, []byte(msg))
}
