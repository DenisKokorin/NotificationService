package service

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	From     string
	Password string
	SmtpHost string
	SmtpPort int
}

func NewEmailService(from, password, smtpHost string, smtpPort int) *EmailService {
	return &EmailService{
		From:     from,
		Password: password,
		SmtpHost: smtpHost,
		SmtpPort: smtpPort,
	}
}

func (es *EmailService) SendEmail(to []string, subject, textBody, htmlBody string, attachments []string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", es.From)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)

	if textBody != "" && htmlBody != "" {
		m.SetBody("text/plain", textBody)
		m.AddAlternative("text/html", htmlBody)
	} else if textBody != "" {
		m.SetBody("text/plain", textBody)
	} else if htmlBody != "" {
		m.SetBody("text/html", htmlBody)
	}

	for _, attachment := range attachments {
		m.Attach(attachment)
	}

	d := gomail.NewDialer(es.SmtpHost, es.SmtpPort, es.From, es.Password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("не удалось отправить письмо: %v", err)
	}

	return nil
}
