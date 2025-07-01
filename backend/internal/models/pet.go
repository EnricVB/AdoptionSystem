package models

import "time"

type Pet struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null"`
	Species     string    `json:"species" gorm:"type:varchar(100);not null"`
	Breed       string    `json:"breed" gorm:"type:varchar(100)"`
	IsAdopted   bool      `json:"is_adopted" gorm:"default:false"`
	BirthDate   time.Time `json:"birth_date"`
	AdoptDate   time.Time `json:"adopt_date"`
	Description string    `json:"description" gorm:"type:text"`
	AdoptUserID uint      `json:"adopt_user_id"`
	AdoptUser   User      `json:"adopt_user" gorm:"foreignKey:AdoptUserID"`
	CrtDate     time.Time `json:"crt_date" gorm:"autoCreateTime"`
	UptDate     time.Time `json:"upt_date" gorm:"autoUpdateTime"`
}

type SimplifiedPet struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Species   string `json:"species"`
	Breed     string `json:"breed"`
	IsAdopted bool   `json:"is_adopted"`
	AdoptUser User   `json:"adopt_user"`
}
