package models

import "time"

func (User) TableName() string {
	return "Users"
}

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
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string    `json:"name" gorm:"type:varchar(100);not null"`
	Surname       string    `json:"surname" gorm:"type:varchar(100);not null"`
	Email         string    `json:"email" gorm:"type:varchar(150);uniqueIndex;not null"`
	Address       string    `json:"address" gorm:"type:varchar(255)"`
	TwoFactorAuth uint      `json:"two_factor_auth" gorm:"default:0"`
	FailedLogins  uint      `json:"failed_logins" gorm:"default:0"`
	IsBlocked     bool      `json:"is_blocked" gorm:"default:false"`
	CrtDate       time.Time `json:"crt_date" gorm:"autoCreateTime"`
	UptDate       time.Time `json:"upt_date" gorm:"autoUpdateTime"`
}

type SimplifiedUser struct {
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name    string `json:"name" gorm:"type:varchar(100);not null"`
	Surname string `json:"surname" gorm:"type:varchar(100);not null"`
	Email   string `json:"email" gorm:"type:varchar(150);uniqueIndex;not null"`
	Address string `json:"address" gorm:"type:varchar(255)"`
}
