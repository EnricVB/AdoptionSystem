package utils

import (
	"log"

	"github.com/go-mail/mail"
)

func SendMail(to string, subject string, body string) error {
	m := mail.NewMessage()
	m.SetHeader("From", "Adoption System <zanckor002@gmail.com>")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := mail.NewDialer("smtp.gmail.com", 465, "zanckor002@gmail.com", "caib nqve pbrw gqjq")
	d.StartTLSPolicy = mail.MandatoryStartTLS

	if err := d.DialAndSend(m); err != nil {
		log.Printf("could not send email: %v", err)
		return err
	}

	log.Println("Email sent!")
	return nil
}
