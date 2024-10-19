package utils

import (
	"context"
	"net/smtp"
)

var sender, appPassword, smtpServer, smtpAdr = "kozyrevdmitrij202@gmail.com", "NoNoNoNo", "smtp.gmail.com", "smtp.gmail.com:587"

func Send(ctx context.Context, email, msg string) error {
	auth := smtp.PlainAuth("", sender, appPassword, smtpServer)
	return smtp.SendMail(smtpAdr, auth, sender, []string{email}, []byte(msg))
}
