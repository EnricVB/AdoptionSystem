// Package dao implements data access objects for species management.
// This layer is responsible for:
// - Direct database operations and queries for species data
// - Data mapping between database and domain models
// - Transaction management and data integrity for species
// - Referential integrity with pet records
// - Database connection handling and error management
package dao

import (
	"backend/internal/db"
	m "backend/internal/models"
	"fmt"
)

// ========================================
// SPECIES RETRIEVAL OPERATIONS
// ========================================

// GetAllSpecies retrieves all species records from the database.
// Returns complete species information for use in pet registration and filtering.
//
// Database Operations:
// - Performs SELECT * FROM species
// - Returns all species without pagination
// - Used for dropdown menus and reference data
//
// Returns:
//   - []m.Species: Slice of all species with complete information
//   - error: Database error or nil on success
func GetAllSpecies() ([]m.Species, error) {
	// Open database connection
	gormDB := db.ORMOpen()

	// Retrieve all species from database
	var species []m.Species
	result := gormDB.Find(&species)
	if result.Error != nil {
		return nil, fmt.Errorf("error al leer especies: %v", result.Error)
	}

	return species, nil
}

// GetSpeciesByID retrieves a specific species by its unique identifier.
// Returns complete species information including description and metadata.
//
// Database Operations:
// - Performs SELECT * FROM species WHERE id = ?
// - Uses GORM's First method for single record retrieval
// - Handles record not found scenarios
//
// Parameters:
//   - id: Unique identifier of the species to retrieve
//
// Returns:
//   - *m.Species: Complete species data
//   - error: Database error or record not found error
func GetSpeciesByID(id uint) (*m.Species, error) {
	// Open database connection
	gormDB := db.ORMOpen()

	// Retrieve specific species by ID
	var s m.Species
	result := gormDB.First(&s, id)
	if result.Error != nil {
		return nil, fmt.Errorf("error al leer especie con id %d: %v", id, result.Error)
	}

	return &s, nil
}

// ========================================
// SPECIES CRUD OPERATIONS
// ========================================

// CreateSpecies creates a new species record in the database.
// Handles species registration with proper data validation and integrity.
//
// Database Operations:
// - Performs INSERT INTO species with provided data
// - Assigns auto-generated ID to the species
// - Validates data integrity and constraints
//
// Business Logic:
// - Ensures species name uniqueness (handled by database constraints)
// - Updates the input species object with generated ID
// - Maintains referential integrity for future pet associations
//
// Parameters:
//   - s: Species data to be created (will be updated with generated ID)
//
// Returns:
//   - error: Database error or validation error, nil on success
func CreateSpecies(s *m.Species) error {
	// Open database connection
	gormDB := db.ORMOpen()

	// Create new species record
	result := gormDB.Create(s)
	if result.Error != nil {
		return fmt.Errorf("error al crear especie: %v", result.Error)
	}

	return nil
}

// DeleteSpeciesByID removes a species from the database by its ID.
// Handles species deletion with proper constraint checking.
//
// Database Operations:
// - Performs DELETE FROM species WHERE id = ?
// - Handles foreign key constraints with pet records
// - May prevent deletion if pets are associated with the species
//
// Business Logic:
// - Enforces referential integrity with pet records
// - Prevents deletion of species that are in use
// - Maintains data consistency across the system
//
// Constraint Handling:
// - If pets exist with this species, deletion will fail
// - Returns appropriate error messages for constraint violations
// - Preserves data integrity in the adoption system
//
// Parameters:
//   - id: Unique identifier of the species to delete
//
// Returns:
//   - error: Database error, constraint violation, or nil on success
func DeleteSpeciesByID(id uint) error {
	// Open database connection
	gormDB := db.ORMOpen()

	// Delete species record by ID
	result := gormDB.Delete(&m.Species{}, id)
	if result.Error != nil {
		return fmt.Errorf("error al eliminar especie con id %d: %v", id, result.Error)
	}

	return nil
}
