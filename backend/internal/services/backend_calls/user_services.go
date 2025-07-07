package services

import (
	r_models "backend/internal/api/routes/models"
	"backend/internal/db/dao"
	m "backend/internal/models"
	mailer "backend/internal/services/mail"
	"backend/internal/services/security"
	"fmt"
)

func AuthenticateUser(userData r_models.LoginRequest) (*m.User, error) {
	_, err := dao.GetValidatedUser(userData.Email, userData.Password)

	if err != nil {
		dao.IncrementFailedLogins(userData.Email)

		return nil, err
	}

	// Reset failed login attempts and generate two-factor authentication code for the user, and generate new SessionID
	dao.ResetFailedLogins(userData.Email)
	dao.GenerateSessionID(userData.Email)

	// Update user's data
	user, _ := dao.GetValidatedUser(userData.Email, userData.Password)

	return user, nil
}

func AuthenticateUser2FA(userData r_models.TwoFactorRequest) (*m.NonValidatedUser, error) {
	user, err := dao.GetUserBySessionID(userData.SessionID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario: %v", err)
	}

	_2fa, err := dao.Get2FA(userData.SessionID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener 2fa: %v", err)
	}

	if _2fa == "" || _2fa != userData.Code {
		return nil, fmt.Errorf("código de autenticación de dos factores inválido")
	}

	// Reset failed login attempts after successful 2FA authentication
	dao.ResetFailedLogins(user.Email)

	// Update user's data
	validatedUser, _ := dao.GetUserBySessionID(userData.SessionID)

	return validatedUser, nil
}

func RefreshUser2FAToken(userData r_models.RefreshTokenRequest) (string, error) {
	generated2FAToken, _2faErr := dao.UpdateTwoFactorCode(userData.Email)

	if generated2FAToken == "" || _2faErr != nil {
		return "", fmt.Errorf("error al generar el token 2FA: %v", _2faErr)
	}

	mailerErr := mailer.Send2FAToken(userData.Email, generated2FAToken)
	if mailerErr != nil {
		return "", fmt.Errorf("error al enviar el token 2FA al email %s: %v", userData.Email, mailerErr)
	}

	return generated2FAToken, nil
}

func ListAllUsers() (*[]m.NonValidatedUser, error) {
	users, err := dao.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("error al leer usuarios: %v", err)
	}

	return &users, nil
}

func GetUserProfile(id uint) (*m.NonValidatedUser, error) {
	user, err := dao.GetUserByID(uint(id))
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario con id: %d %v", id, err)
	}

	return user, nil
}

func GetUserByEmail(email string) (*m.NonValidatedUser, error) {
	user, err := dao.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario con email %s: %v", email, err)
	}

	return user, nil
}

func RegisterUser(user *m.User) error {
	user.Password, _ = security.HashPassword(user.Password)

	err := dao.CreateUser(user)
	if err != nil {
		return fmt.Errorf("error al crear usuario")
	}

	return nil
}

func UpdateUserProfile(user *m.User) error {
	err := dao.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("error al actualizar usuario: %v", err)
	}

	return nil
}

func DeactivateUser(id uint) (*m.SimplifiedUser, error) {
	deleted, err := dao.DeleteUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("error al eliminar usuario con id: %d %v", id, err)
	}

	return deleted, nil
}
