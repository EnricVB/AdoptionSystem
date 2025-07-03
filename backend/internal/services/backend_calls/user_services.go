package services

import (
	r_models "backend/internal/api/routes/models"
	"backend/internal/db/dao"
	m "backend/internal/models"
	"fmt"
)

func AuthenticateUser(userData r_models.LoginRequest) (*m.User, error) {
	_, err := dao.GetValidatedUser(userData.Email, userData.Password)

	if err != nil {
		dao.IncrementFailedLogins(userData.Email)

		return nil, err
	}

	// Reset failed login attempts and generate two-factor authentication code for the user
	dao.ResetFailedLogins(userData.Email)

	// Update user's data
	user, _ := dao.GetValidatedUser(userData.Email, userData.Password)

	return user, nil
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

func RegisterUser(user *m.User) error {
	err := dao.CreateUser(user)
	if err != nil {
		return fmt.Errorf("error al crear usuario: %v", err)
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
