package handlers

import (
	r_models "backend/internal/api/routes/models"
	"backend/internal/models"
	s "backend/internal/services/backend_calls"
	response "backend/internal/utils/rest"
	"net/http"
	"time"
)

func HandleManualLogin(req r_models.LoginRequest) (*models.User, response.HTTPError) {
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

func HandleGoogleLogin(req r_models.GoogleLoginRequest) (*models.User, response.HTTPError) {
	if req.Email == "" || req.IDToken == "" {
		return nil, response.Error(http.StatusBadRequest, "email y ID Token son obligatorios")
	}

	user, err := s.AuthenticateGoogleUser(req)
	if err != nil {
		return nil, response.Error(http.StatusUnauthorized, err.Error())
	}

	return user, response.EmptyError
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

func HandleCreateUser(user *r_models.CreateUserRequest) response.HTTPError {
	if user.Name == "" {
		return response.Error(http.StatusBadRequest, "nombre es obligatorio")
	}

	fullUser := &models.FullUser{
		Name:       user.Name,
		Surname:    user.Surname,
		Email:      user.Email,
		Password:   user.Password,
		Address:    user.Address,
		Provider:   user.Provider,
		ProviderID: user.ProviderID,

		CrtDate: time.Now(),
		UptDate: time.Now(),
	}

	err := s.RegisterUser(fullUser)
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
