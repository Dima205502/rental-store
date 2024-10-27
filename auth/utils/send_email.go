package utils

import (
	"auth_service/config"
	"context"
	"net/smtp"
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

func Send(ctx context.Context, email, msg string) error {
	auth := smtp.PlainAuth("", sender, appPassword, smtpHost)
	return smtp.SendMail(smtpHost+":"+strconv.Itoa(smtpPort), auth, sender, []string{email}, []byte(msg))
}
