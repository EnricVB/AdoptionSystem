package models

import "time"

type User struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string    `json:"name" gorm:"type:varchar(100);not null"`
	Surname       string    `json:"surname" gorm:"type:varchar(100);not null"`
	Email         string    `json:"email" gorm:"type:varchar(150);uniqueIndex;not null"`
	Address       string    `json:"address" gorm:"type:varchar(255)"`
	Password      string    `json:"password" gorm:"type:varchar(255);not null"`
	TwoFactorAuth uint      `json:"two_factor_auth" gorm:"default:0"`
	FailedLogins  uint      `json:"failed_logins" gorm:"default:0"`
	IsBlocked     bool      `json:"is_blocked" gorm:"default:false"`
	CrtDate       time.Time `json:"crt_date" gorm:"autoCreateTime"`
	UptDate       time.Time `json:"upt_date" gorm:"autoUpdateTime"`
}

type NonValidatedUser struct {
	ID      uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name    string    `json:"name" gorm:"type:varchar(100);not null"`
	Surname string    `json:"surname" gorm:"type:varchar(100);not null"`
	Email   string    `json:"email" gorm:"type:varchar(150);uniqueIndex;not null"`
	Address string    `json:"address" gorm:"type:varchar(255)"`
	CrtDate time.Time `json:"crt_date" gorm:"autoCreateTime"`
	UptDate time.Time `json:"upt_date" gorm:"autoUpdateTime"`
}

type SimplifiedUser struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
	Address string `json:"address"`
}
