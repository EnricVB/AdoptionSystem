package models

import "time"

func (User) TableName() string {
	return "Users"
}

type FullUser struct {
	ID            uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string `json:"name" gorm:"type:varchar(100);not null"`
	Surname       string `json:"surname" gorm:"type:varchar(100);not null"`
	Email         string `json:"email" gorm:"type:varchar(150);uniqueIndex;not null"`
	SessionID     string `json:"session_id" gorm:"type:varchar(50);uniqueIndex;column:Session_ID"`
	Address       string `json:"address" gorm:"type:varchar(255)"`
	FailedLogins  uint   `json:"failed_logins" gorm:"default:0;column:Failed_Logins"`
	IsBlocked     bool   `json:"is_blocked" gorm:"default:false;column:Is_Blocked"`
	TwoFactorAuth string `json:"two_factor_auth" gorm:"type:varchar(6);column:Two_Factor_Auth"`

	Password   string `json:"password,omitempty" gorm:"type:varchar(255);column:Password"`
	Provider   string `json:"provider" gorm:"default:'local';type:varchar(255);column:Provider"`
	ProviderID string `json:"provider_id" gorm:"type:varchar(255);column:Provider_ID"`

	CrtDate time.Time `json:"crt_date" gorm:"autoCreateTime"`
	UptDate time.Time `json:"upt_date" gorm:"autoUpdateTime"`
}

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name         string    `json:"name" gorm:"type:varchar(100);not null"`
	Surname      string    `json:"surname" gorm:"type:varchar(100);not null"`
	Email        string    `json:"email" gorm:"type:varchar(150);uniqueIndex;not null"`
	SessionID    string    `json:"session_id" gorm:"type:varchar(50);uniqueIndex;column:Session_ID"`
	Address      string    `json:"address" gorm:"type:varchar(255)"`
	FailedLogins uint      `json:"failed_logins" gorm:"default:0;column:Failed_Logins"`
	IsBlocked    bool      `json:"is_blocked" gorm:"default:false;column:Is_Blocked"`
	CrtDate      time.Time `json:"crt_date" gorm:"autoCreateTime"`
	UptDate      time.Time `json:"upt_date" gorm:"autoUpdateTime"`
}

type NonValidatedUser struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name         string    `json:"name" gorm:"type:varchar(100);not null"`
	Surname      string    `json:"surname" gorm:"type:varchar(100);not null"`
	Email        string    `json:"email" gorm:"type:varchar(150);uniqueIndex;not null"`
	Address      string    `json:"address" gorm:"type:varchar(255)"`
	FailedLogins uint      `json:"failed_logins" gorm:"default:0;column:Failed_Logins"`
	IsBlocked    bool      `json:"is_blocked" gorm:"default:false;column:Is_Blocked"`
	CrtDate      time.Time `json:"crt_date" gorm:"autoCreateTime"`
	UptDate      time.Time `json:"upt_date" gorm:"autoUpdateTime"`
}

type SimplifiedUser struct {
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name    string `json:"name" gorm:"type:varchar(100);not null"`
	Surname string `json:"surname" gorm:"type:varchar(100);not null"`
	Email   string `json:"email" gorm:"type:varchar(150);uniqueIndex;not null"`
	Address string `json:"address" gorm:"type:varchar(255)"`
}
