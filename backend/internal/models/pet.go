// Package models contains data models for the pet adoption system.
// These models define the structure of pet-related database entities and their relationships.
package models

import "time"

// TableName returns the database table name for the Pet model.
// This method implements the GORM Tabler interface to specify custom table names.
func (Pet) TableName() string {
	return "Pets"
}

// Pet represents the complete pet entity in the adoption system.
// This model contains all pet information including adoption status,
// associated user data, and temporal information.
//
// Database Table: Pets
// Relationships:
//   - AdoptUser: Many-to-One relationship with User (foreign key: AdoptUserID)
type Pet struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`        // Unique identifier for the pet
	Name        string    `json:"name" gorm:"type:varchar(100);not null"`    // Pet's name
	Species     string    `json:"species" gorm:"type:varchar(100);not null"` // Pet's species (dog, cat, etc.)
	Breed       string    `json:"breed" gorm:"type:varchar(100)"`            // Pet's breed (optional)
	IsAdopted   bool      `json:"is_adopted" gorm:"default:false"`           // Whether the pet has been adopted
	BirthDate   time.Time `json:"birth_date"`                                // Pet's date of birth
	AdoptDate   time.Time `json:"adopt_date"`                                // Date when the pet was adopted
	Description string    `json:"description" gorm:"type:text"`              // Detailed description of the pet
	AdoptUserID uint      `json:"adopt_user_id"`                             // ID of the user who adopted the pet
	AdoptUser   User      `json:"adopt_user" gorm:"foreignKey:AdoptUserID"`  // User who adopted the pet (relationship)
	CrtDate     time.Time `json:"crt_date" gorm:"autoCreateTime"`            // Record creation timestamp
	UptDate     time.Time `json:"upt_date" gorm:"autoUpdateTime"`            // Record last update timestamp
}

// SimplifiedPet represents a minimal pet entity with essential information.
// This model is used for operations that require only basic pet data,
// such as pet lists, search results, or summary displays.
//
// Business Rules:
//   - Used primarily for listing and summary operations
//   - Includes adoption status and user information for quick reference
//   - Excludes detailed fields like description and dates for performance
type SimplifiedPet struct {
	ID        uint   `json:"id"`         // Unique identifier for the pet
	Name      string `json:"name"`       // Pet's name
	Species   string `json:"species"`    // Pet's species (dog, cat, etc.)
	Breed     string `json:"breed"`      // Pet's breed (optional)
	IsAdopted bool   `json:"is_adopted"` // Whether the pet has been adopted
	AdoptUser User   `json:"adopt_user"` // User who adopted the pet (if adopted)
}
