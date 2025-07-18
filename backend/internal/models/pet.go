// Package models contains data models for the pet adoption system.
// These models define the structure of pet-related database entities and their relationships.
package models

import "time"

// TableName returns the database table name for the Pet model.
// This method implements the GORM Tabler interface to specify custom table names.
func (Pet) TableName() string {
	return "Pets"
}

// PetStatus represents the status of a pet in the adoption system
type PetStatus string

const (
	PetStatusAvailable  PetStatus = "Available"
	PetStatusFosterHome PetStatus = "FosterHome"
	PetStatusAdopted    PetStatus = "Adopted"
)

// PetGender represents the gender of a pet
type PetGender string

const (
	PetGenderMale   PetGender = "Male"
	PetGenderFemale PetGender = "Female"
)

// Pet represents the complete pet entity in the adoption system.
// This model contains all pet information including adoption status,
// associated user data, and temporal information.
//
// Database Table: Pets
// Relationships:
//   - AdoptUser: Many-to-One relationship with User (foreign key: AdoptUserID)
type Pet struct {
	ID                 uint                 `json:"id" gorm:"primaryKey;autoIncrement"`                                        // Unique identifier for the pet
	Name               string               `json:"name" gorm:"type:varchar(100);not null"`                                    // Pet's name
	SpeciesID          uint                 `json:"species_id" gorm:"column:Species_ID"`                                       // Foreign key to Species
	Species            Species              `json:"species" gorm:"foreignKey:SpeciesID"`                                       // Relationship to Species
	Breed              string               `json:"breed" gorm:"type:varchar(100)"`                                            // Pet's breed (optional)
	Gender             PetGender            `json:"gender" gorm:"type:enum('Male','Female');column:Gender"`                    // Pet's gender
	Weight             float64              `json:"weight" gorm:"column:Weight;type:decimal(5,2);default:0.00"`                // Pet's weight in kg
	Status             PetStatus            `json:"status" gorm:"type:enum('Available','FosterHome','Adopted');column:Status"` // Pet's adoption status
	BirthDate          time.Time            `json:"birthdate" gorm:"column:Birthdate"`                                         // Pet's date of birth
	AdoptDate          *time.Time           `json:"adopt_date" gorm:"column:Adoptdate"`                                        // Date when the pet was adopted
	Description        string               `json:"description" gorm:"column:Description"`                                     // Detailed description of the pet
	AdoptUserID        *uint                `json:"adopt_user_id" gorm:"column:Adopt_User"`                                    // ID of the user who adopted the pet (nullable)
	AdoptUser          *User                `json:"adopt_user" gorm:"foreignKey:AdoptUserID"`                                  // Relationship to User
	ImageURL           string               `json:"image_url" gorm:"column:ImageURL"`                                          // Image URL
	IsVaccinated       bool                 `json:"is_vaccinated" gorm:"column:Vaccinated;default:false"`                      // Vaccination status
	VaccinationHistory []VaccinationHistory `json:"vaccination_history" gorm:"foreignKey:PetID"`                               // Relationship to VaccinationHistory
	CrtDate            time.Time            `json:"crt_date" gorm:"autoCreateTime"`                                            // Record creation timestamp
	UptDate            time.Time            `json:"upt_date" gorm:"autoUpdateTime"`                                            // Record last update timestamp
}

// VaccinationHistory represents a vaccination record for a pet
type VaccinationHistory struct {
	PetID           uint      `json:"pet_id" gorm:"column:Pet_ID;primaryKey"`
	VaccinationDate time.Time `json:"vaccination_date" gorm:"column:Vaccination_Date;primaryKey"`
	VaccineName     string    `json:"vaccine_name" gorm:"column:Vaccine_Name;primaryKey;type:varchar(100)"`
}

// TableName returns the database table name for the VaccinationHistory model
func (VaccinationHistory) TableName() string {
	return "VaccinationHistory"
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
	ID                 uint                 `json:"id"`                                                                        // Unique identifier for the pet
	Name               string               `json:"name"`                                                                      // Pet's name
	Description        string               `json:"description" gorm:"column:Description"`                                     // Detailed description of the pet
	SpeciesID          uint                 `json:"species_id" gorm:"column:Species_ID"`                                       // Foreign key to Species
	Species            Species              `json:"species"`                                                                   // Relationship to Species
	Breed              string               `json:"breed"`                                                                     // Pet's breed (optional)
	Gender             PetGender            `json:"gender" gorm:"type:enum('Male','Female');column:Gender"`                    // Pet's gender
	Weight             float64              `json:"weight" gorm:"column:Weight;type:decimal(5,2);default:0.00"`                // Pet's weight in kg
	Status             PetStatus            `json:"status" gorm:"type:enum('Available','FosterHome','Adopted');column:Status"` // Pet's adoption status
	AdoptUserID        *uint                `json:"adopt_user_id" gorm:"column:Adopt_User"`                                    // ID of the user who adopted the pet
	AdoptUser          *User                `json:"adopt_user"`                                                                // Relationship to User
	BirthDate          time.Time            `json:"birthdate" gorm:"column:Birthdate"`                                         // Pet's date of birth
	ImageURL           string               `json:"image_url" gorm:"column:ImageURL"`                                          // Image URL
	IsVaccinated       bool                 `json:"is_vaccinated" gorm:"column:Vaccinated;default:false"`                      // Vaccination status
	VaccinationHistory []VaccinationHistory `json:"vaccination_history" gorm:"foreignKey:PetID"`                               // Relationship to VaccinationHistory
}
