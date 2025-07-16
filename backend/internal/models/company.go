package models

import (
	"time"
)

// Company represents the company information schema.
// This model contains the basic information about the organization
// managing the pet adoption system.
type Company struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	CompanyName  string    `json:"company_name" gorm:"column:Company_Name;not null"`
	ContactEmail string    `json:"contact_email" gorm:"column:Contact_Email;not null"`
	Phone        string    `json:"phone" gorm:"column:Phone"`
	Address      string    `json:"address" gorm:"column:Address"`
	City         string    `json:"city" gorm:"column:City"`
	State        string    `json:"state" gorm:"column:State"`
	PostalCode   string    `json:"postal_code" gorm:"column:Postal_Code"`
	Country      string    `json:"country" gorm:"column:Country"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:Created_At;autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:Updated_At;autoUpdateTime"`
}

// TableName specifies the table name for the Company model.
func (Company) TableName() string {
	return "schema_company"
}
