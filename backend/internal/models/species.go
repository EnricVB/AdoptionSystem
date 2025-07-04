// Package models contains data models for the pet adoption system.
// These models define the structure of species-related database entities.
package models

// TableName returns the database table name for the Species model.
// This method implements the GORM Tabler interface to specify custom table names.
func (Species) TableName() string {
	return "Species"
}

// Species represents a pet species entity in the adoption system.
// This model defines the types of animals that can be registered for adoption.
//
// Business Rules:
//   - Species names must be unique across the system
//   - Species are used to categorize pets for better organization and searching
//   - Common species include: Dog, Cat, Bird, Rabbit, etc.
//   - Species cannot be deleted if pets are associated with them
//
// Database Table: Species
// Constraints:
//   - Name field has a unique constraint to prevent duplicates
type Species struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"`            // Unique identifier for the species
	Name string `json:"name" gorm:"type:varchar(100);not null;unique"` // Species name (e.g., "Dog", "Cat", "Bird")
}
