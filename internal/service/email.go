package service

import "net/smtp"

type EmailService struct {
	From     string
	Password string
	SmtpHost string
	SmtpPort int
	SmtpAddr string
}

func NewEmailService(from, password, smtpHost string, smtpPort int, smtpAddr string) *EmailService {
	return &EmailService{
		From:     from,
		Password: password,
		SmtpHost: smtpHost,
		SmtpPort: smtpPort,
		SmtpAddr: smtpAddr,
	}
}

func (es *EmailService) SendEmail(to []string, subject, textBody string) error {
	auth := smtp.PlainAuth("", es.From, es.Password, es.SmtpHost)

	msg := []byte("Subject: " + subject + "\n" + textBody)

	return smtp.SendMail(es.SmtpAddr, auth, es.From, to, msg)
}
