package dao

import (
	"backend/internal/db"
	m "backend/internal/models"
	"backend/internal/utils/security"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func GetAllUsers() ([]m.NonValidatedUser, error) {
	gormDB := db.ORMOpen()

	var users []m.NonValidatedUser
	result := gormDB.Model(&m.User{}).
		Select("id, name, surname, email, address, two_factor_auth, failed_logins, is_blocked, crt_date, upt_date").
		Find(&users)

	if result.Error != nil {
		return nil, fmt.Errorf("error al leer usuarios: %v", result.Error)
	}

	return users, nil
}

func GetUserByID(id uint) (*m.NonValidatedUser, error) {
	gormDB := db.ORMOpen()

	var user m.NonValidatedUser
	result := gormDB.Model(&m.User{}).
		Select("id, name, surname, email, address, two_factor_auth, failed_logins, is_blocked").
		Where("id = ?", id).
		First(&user)

	if result.Error != nil {
		return nil, fmt.Errorf("error al leer usuario con id %d: %v", id, result.Error)
	}

	return &user, nil
}

func GetUserByEmail(email string) (*m.NonValidatedUser, error) {
	gormDB := db.ORMOpen()

	var user m.NonValidatedUser
	result := gormDB.Model(&m.User{}).
		Select("id, name, surname, email, address, two_factor_auth, failed_logins, is_blocked").
		Where("email = ?", email).
		First(&user)

	if result.Error != nil {
		return nil, fmt.Errorf("error al leer usuario con email %s: %v", email, result.Error)
	}

	return &user, nil
}

func GetValidatedUser(email string, password string) (*m.User, error) {
	gormDB := db.ORMOpen()

	var user m.User
	result := gormDB.Debug().Where("email = ?", email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("usuario con email %s no encontrado", email)
		}
		return nil, fmt.Errorf("error al buscar usuario: %v", result.Error)
	}

	if user.IsBlocked {
		return nil, fmt.Errorf("usuario bloqueado")
	}

	if !security.VerifyPassword(user.Password, password) {
		return nil, fmt.Errorf("credenciales inválidas")
	}

	return &user, nil
}

func DeleteUserByID(id uint) (*m.SimplifiedUser, error) {
	gormDB := db.ORMOpen()

	var user m.SimplifiedUser
	result := gormDB.Delete(&m.User{}, id)

	if result.Error != nil {
		return nil, fmt.Errorf("error al eliminar usuario con id %d: %v", id, result.Error)
	}

	return &user, nil
}

func CreateUser(user *m.User) (*m.User, error) {
	gormDB := db.ORMOpen()

	now := time.Now()
	user.CrtDate = now
	user.UptDate = now

	result := gormDB.Create(user)
	if result.Error != nil {
		return nil, fmt.Errorf("error al crear usuario: %v", result.Error)
	}

	return user, nil
}

func UpdateUser(user *m.User) (*m.NonValidatedUser, error) {
	gormDB := db.ORMOpen()

	user.UptDate = time.Now()
	result := gormDB.Model(&m.User{}).
		Where("id = ?", user.ID).
		Updates(user)

	if result.Error != nil {
		return nil, fmt.Errorf("error al actualizar usuario con id %d: %v", user.ID, result.Error)
	}

	simplified := &m.NonValidatedUser{
		ID:            user.ID,
		Name:          user.Name,
		Surname:       user.Surname,
		Email:         user.Email,
		Address:       user.Address,
		TwoFactorAuth: user.TwoFactorAuth,
		FailedLogins:  user.FailedLogins,
		IsBlocked:     user.IsBlocked,
	}

	return simplified, nil
}

func UpdateLoginData(email string, failedLogins int, isBlocked bool) error {
	gormDB := db.ORMOpen()

	result := gormDB.Model(&m.User{}).
		Where("email = ?", email).
		Updates(map[string]interface{}{
			"failed_logins": failedLogins,
			"is_blocked":    isBlocked,
			"upt_date":      time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("error al actualizar datos de login para usuario %s: %v", email, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("usuario con email %s no encontrado", email)
	}

	return nil
}

func IncrementFailedLogins(email string) error {
	gormDB := db.ORMOpen()

	var currentFailedLogins int
	result := gormDB.Model(&m.User{}).
		Select("failed_logins").
		Where("email = ?", email).
		First(&currentFailedLogins)

	if result.Error != nil {
		return fmt.Errorf("error al obtener failed_logins para usuario %s: %v", email, result.Error)
	}

	newFailedLogins := currentFailedLogins + 1

	isBlocked := newFailedLogins >= 5

	return UpdateLoginData(email, newFailedLogins, isBlocked)
}

func ResetFailedLogins(email string) error {
	return UpdateLoginData(email, 0, false)
}

func BlockUser(email string) error {
	gormDB := db.ORMOpen()

	result := gormDB.Model(&m.User{}).
		Where("email = ?", email).
		Updates(map[string]interface{}{
			"is_blocked": true,
			"upt_date":   time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("error al bloquear usuario %s: %v", email, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("usuario con email %s no encontrado", email)
	}

	return nil
}

func UnblockUser(email string) error {
	return UpdateLoginData(email, 0, false)
}

func UpdateTwoFactorCode(email string) (string, error) {
	gormDB := db.ORMOpen()
	_2fa := security.Generate2FA(6)

	result := gormDB.Model(&m.User{}).
		Where("email =?", email).
		Updates(map[string]any{
			"two_factor_auth": _2fa,
		})

	user, err := GetUserByEmail(email)
	if result.Error != nil || err != nil || user.TwoFactorAuth == "" {
		return "", fmt.Errorf("error al actualizarTwoFactorAuth para usuario %s: %v", email, result.Error)
	}

	if result.RowsAffected == 0 {
		return "", fmt.Errorf("usuario con email %s no encontrado", email)
	}

	return _2fa, nil
}
