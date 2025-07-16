package mail

import (
	"backend/internal/models"
	"bytes"
	_ "embed"
	"fmt"
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

//go:embed templates/pet-adoption-request.html
var petAdoptionRequestTemplate string

//go:embed templates/pet-foster-home-request.html
var petFosterHomeRequestTemplate string

//go:embed templates/pet-foster-home-contact.html
var petFosterHomeContactTemplate string

type PetAdoptionRequestData struct {
	Pet  models.Pet
	User models.SimplifiedUser
}

type PetFosterHomeRequestData struct {
	Pet  models.Pet
	User models.SimplifiedUser
}

type PetFosterHomeContactData struct {
	Pet     models.Pet
	User    models.SimplifiedUser
	Contact models.SimplifiedUser
	Reason  string
	Message string
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

func SendPetAdoptionRequest(to string, pet models.Pet, user models.SimplifiedUser) error {
	m := mail.NewMessage()
	m.SetHeader("From", user.Email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Nueva Solicitud de Adopción - "+pet.Name)

	data := PetAdoptionRequestData{
		Pet:  pet,
		User: user,
	}

	tmpl, err := template.New("pet-adoption-request").Parse(petAdoptionRequestTemplate)
	if err != nil {
		log.Printf("error parsing pet adoption request template: %v", err)
		return err
	}

	var htmlBody bytes.Buffer
	if err := tmpl.Execute(&htmlBody, data); err != nil {
		log.Printf("error executing pet adoption request template: %v", err)
		return err
	}

	plainBody := fmt.Sprintf("Nueva solicitud de adopción para %s de %s (%s)", pet.Name, user.Name, user.Email)

	m.SetBody("text/plain", plainBody)
	m.AddAlternative("text/html", htmlBody.String())

	d := mail.NewDialer("smtp.gmail.com", 465, "zanckor002@gmail.com", "caib nqve pbrw gqjq")
	d.StartTLSPolicy = mail.MandatoryStartTLS

	if err := d.DialAndSend(m); err != nil {
		log.Printf("could not send pet adoption request email: %v", err)
		return err
	}

	return nil
}

func SendPetFosterHomeRequest(to string, pet models.Pet, user models.SimplifiedUser) error {
	m := mail.NewMessage()
	m.SetHeader("From", user.Email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Nueva Solicitud de Casa de Acogida - "+pet.Name)

	data := PetFosterHomeRequestData{
		Pet:  pet,
		User: user,
	}

	tmpl, err := template.New("pet-foster-home-request").Parse(petFosterHomeRequestTemplate)
	if err != nil {
		log.Printf("error parsing pet foster home request template: %v", err)
		return err
	}

	var htmlBody bytes.Buffer
	if err := tmpl.Execute(&htmlBody, data); err != nil {
		log.Printf("error executing pet foster home request template: %v", err)
		return err
	}

	plainBody := fmt.Sprintf("Nueva solicitud de casa de acogida para %s de %s (%s)", pet.Name, user.Name, user.Email)

	m.SetBody("text/plain", plainBody)
	m.AddAlternative("text/html", htmlBody.String())

	d := mail.NewDialer("smtp.gmail.com", 465, "zanckor002@gmail.com", "caib nqve pbrw gqjq")
	d.StartTLSPolicy = mail.MandatoryStartTLS

	if err := d.DialAndSend(m); err != nil {
		log.Printf("could not send pet foster home request email: %v", err)
		return err
	}

	return nil
}

func SendPetFosterHomeContact(to string, pet models.Pet, user models.SimplifiedUser, contact models.SimplifiedUser, reason string, message string) error {
	m := mail.NewMessage()
	m.SetHeader("From", user.Email)
	m.SetHeader("To", to)
	m.SetHeader("Cc", contact.Email)
	m.SetHeader("Subject", "Consulta sobre "+pet.Name+" - "+reason)

	data := PetFosterHomeContactData{
		Pet:     pet,
		User:    user,
		Contact: contact,
		Reason:  reason,
		Message: message,
	}

	tmpl, err := template.New("pet-foster-home-contact").Parse(petFosterHomeContactTemplate)
	if err != nil {
		log.Printf("error parsing pet foster home contact template: %v", err)
		return err
	}

	var htmlBody bytes.Buffer
	if err := tmpl.Execute(&htmlBody, data); err != nil {
		log.Printf("error executing pet foster home contact template: %v", err)
		return err
	}

	plainBody := fmt.Sprintf("Consulta sobre %s de %s (%s): %s", pet.Name, contact.Name, contact.Email, message)

	m.SetBody("text/plain", plainBody)
	m.AddAlternative("text/html", htmlBody.String())

	d := mail.NewDialer("smtp.gmail.com", 465, "zanckor002@gmail.com", "caib nqve pbrw gqjq")
	d.StartTLSPolicy = mail.MandatoryStartTLS

	if err := d.DialAndSend(m); err != nil {
		log.Printf("could not send pet foster home contact email: %v", err)
		return err
	}

	return nil
}
