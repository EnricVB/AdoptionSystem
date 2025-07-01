package utils

import (
	m "backend/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(bytes), err
}

func CheckPassword(user *m.User, checkPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(checkPassword))

	return err != nil
}
