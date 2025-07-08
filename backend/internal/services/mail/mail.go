package mailer

import (
	"bytes"
	_ "embed"
	"log"
	"text/template"

	"github.com/go-mail/mail"
)

//go:embed templates/password.html
var passwordTemplate string

type PasswordData struct {
	Password string
}

//go:embed templates/2fa.html
var twoFATemplate string

type TwoFAData struct {
	Code string
}

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

	return nil
}

func Send2FAToken(to string, _2fa string) error {
	m := mail.NewMessage()
	m.SetHeader("From", "Adoption System <zanckor002@gmail.com>")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Código de Autenticación 2FA")

	data := TwoFAData{Code: _2fa}

	tmpl, err := template.New("2fa").Parse(twoFATemplate)
	if err != nil {
		log.Printf("error parsing 2fa template: %v", err)
		return err
	}

	var htmlBody bytes.Buffer
	if err := tmpl.Execute(&htmlBody, data); err != nil {
		log.Printf("error executing 2fa template: %v", err)
		return err
	}

	plainBody := "Tu código de autenticación 2FA es: " + _2fa

	m.SetBody("text/plain", plainBody)
	m.AddAlternative("text/html", htmlBody.String())

	d := mail.NewDialer("smtp.gmail.com", 465, "zanckor002@gmail.com", "caib nqve pbrw gqjq")
	d.StartTLSPolicy = mail.MandatoryStartTLS

	if err := d.DialAndSend(m); err != nil {
		log.Printf("could not send email: %v", err)
		return err
	}

	return nil
}

func SendPassword(to string, password string) error {
	m := mail.NewMessage()
	m.SetHeader("From", "Adoption System <zanckor002@gmail.com>")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Tu nueva contraseña")

	data := PasswordData{Password: password}

	tmpl, err := template.New("password").Parse(passwordTemplate)
	if err != nil {
		log.Printf("error parsing password template: %v", err)
		return err
	}

	var htmlBody bytes.Buffer
	if err := tmpl.Execute(&htmlBody, data); err != nil {
		log.Printf("error executing password template: %v", err)
		return err
	}

	plainBody := "Tu nueva contraseña es: " + password

	m.SetBody("text/plain", plainBody)
	m.AddAlternative("text/html", htmlBody.String())

	d := mail.NewDialer("smtp.gmail.com", 465, "zanckor002@gmail.com", "caib nqve pbrw gqjq")
	d.StartTLSPolicy = mail.MandatoryStartTLS

	if err := d.DialAndSend(m); err != nil {
		log.Printf("could not send email: %v", err)
		return err
	}

	return nil
}
