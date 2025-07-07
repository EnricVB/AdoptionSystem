package handlers

import (
	r_models "backend/internal/api/routes/models"
	"backend/internal/models"
	s "backend/internal/services/backend_calls"
	response "backend/internal/utils/rest"
	"net/http"
)

func HandleLogin(req r_models.LoginRequest) (*models.User, response.HTTPError) {
	if req.Email == "" || req.Password == "" {
		return nil, response.Error(http.StatusBadRequest, "email y contraseña son obligatorios")
	}

	user, err := s.AuthenticateUser(req)
	if err != nil {
		return nil, response.Error(http.StatusUnauthorized, err.Error())
	}

	return user, response.EmptyError
}

func Handle2FAAuth(req r_models.TwoFactorRequest) (*models.NonValidatedUser, response.HTTPError) {
	if req.SessionID == "" || req.Code == "" {
		return nil, response.Error(http.StatusBadRequest, "sessionID y código de 2FA son obligatorios")
	}

	user, err := s.AuthenticateUser2FA(req)
	if err != nil {
		return nil, response.Error(http.StatusUnauthorized, err.Error())
	}

	return user, response.EmptyError
}

func HandleRefresh2FAToken(req r_models.RefreshTokenRequest) (string, response.HTTPError) {
	if req.Email == "" {
		return "", response.Error(http.StatusBadRequest, "email es obligatorio")
	}

	token, err := s.RefreshUser2FAToken(req)
	if err != nil {
		return "", response.Error(http.StatusUnauthorized, err.Error())
	}

	return token, response.EmptyError
}

func HandleListUsers() (*[]models.NonValidatedUser, response.HTTPError) {
	users, err := s.ListAllUsers()
	if err != nil {
		return nil, response.Error(http.StatusInternalServerError, err.Error())
	}

	return users, response.EmptyError
}

func HandleGetUserByID(id uint) (*models.NonValidatedUser, response.HTTPError) {
	if id <= 0 {
		return nil, response.Error(http.StatusBadRequest, "ID de usuario no válido")
	}

	user, err := s.GetUserProfile(id)
	if err != nil {
		return nil, response.Error(http.StatusNotFound, err.Error())
	}

	return user, response.EmptyError
}

func HandleCreateUser(user *models.User) response.HTTPError {
	if user.Email == "" || user.Password == "" {
		return response.Error(http.StatusBadRequest, "email y contraseña son obligatorios")
	}

	if user.Name == "" {
		return response.Error(http.StatusBadRequest, "nombre es obligatorio")
	}

	userByEmail, _ := s.GetUserByEmail(user.Email)

	if userByEmail != nil {
		return response.Error(http.StatusConflict, "el email ya está registrado")
	}

	err := s.RegisterUser(user)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.EmptyError
}

func HandleUpdateUser(user *models.User) response.HTTPError {
	err := s.UpdateUserProfile(user)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.EmptyError
}

func HandleDeleteUser(id uint) (*models.SimplifiedUser, response.HTTPError) {
	if id <= 0 {
		return nil, response.Error(http.StatusBadRequest, "ID de usuario no válido")
	}

	deleted, err := s.DeactivateUser(id)
	if err != nil {
		return nil, response.Error(http.StatusInternalServerError, err.Error())
	}

	return deleted, response.EmptyError
}
