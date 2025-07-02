package handlers

import (
	r_models "backend/internal/api/routes/models"
	"backend/internal/db/dao"
	"backend/internal/models"
)

func HandleLogin(req r_models.LoginRequest) (*models.User, error) {
	_, err := dao.GetValidatedUser(req.Email, req.Password)

	if err != nil {
		dao.IncrementFailedLogins(req.Email)

		return nil, err
	}

	// Reset failed login attempts and generate two-factor authentication code for the user
	dao.ResetFailedLogins(req.Email)

	// Update user's data
	user, _ := dao.GetValidatedUser(req.Email, req.Password)

	return user, nil
}
