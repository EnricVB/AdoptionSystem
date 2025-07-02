package dao

import (
	"backend/internal/db"
	m "backend/internal/models"
	"fmt"
	"time"
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

func GetUserByID(id uint) (*m.SimplifiedUser, error) {
	gormDB := db.ORMOpen()

	var user m.SimplifiedUser
	result := gormDB.Model(&m.User{}).
		Select("id, name, surname, email, address").
		Where("id = ?", id).
		First(&user)

	if result.Error != nil {
		return nil, fmt.Errorf("error al leer usuario con id %d: %v", id, result.Error)
	}

	return &user, nil
}

func GetUserByEmail(email string) (*m.SimplifiedUser, error) {
	gormDB := db.ORMOpen()

	var user m.SimplifiedUser
	result := gormDB.Model(&m.User{}).
		Select("id, name, surname, email, address").
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
	result := gormDB.
		Where("email = ?", email).
		Where("password = ?", password).
		First(&user)

	if result.Error != nil {
		return nil, fmt.Errorf("error al leer usuario con email %s usando contraseña: %v", email, result.Error)
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

func UpdateUser(user *m.User) (*m.SimplifiedUser, error) {
	gormDB := db.ORMOpen()

	user.UptDate = time.Now()
	result := gormDB.Model(&m.User{}).
		Where("id = ?", user.ID).
		Updates(user)

	if result.Error != nil {
		return nil, fmt.Errorf("error al actualizar usuario con id %d: %v", user.ID, result.Error)
	}

	simplified := &m.SimplifiedUser{
		ID:      user.ID,
		Name:    user.Name,
		Surname: user.Surname,
		Email:   user.Email,
		Address: user.Address,
	}

	return simplified, nil
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
