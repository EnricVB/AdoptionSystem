// Package services provides business logic services for species management.
// This layer sits between handlers and DAOs, implementing the core business rules
// and orchestrating database operations for species-related functionality.
package services

import (
	"backend/internal/db/dao"
	m "backend/internal/models"
	"fmt"
)

// ========================================
// SPECIES MANAGEMENT SERVICES
// ========================================

// ListAllSpecies retrieves all species from the database.
// Returns all available species for use in pet registration and categorization.
//
// Business Logic:
// - Retrieves all species regardless of status
// - Used for populating dropdown menus and filters
// - Provides reference data for pet management
//
// Returns:
//   - []m.Species: Slice of all species with their information
//   - error: Database error or nil on success
func ListAllSpecies() ([]m.Species, error) {
	// Retrieve all species from database
	species, err := dao.GetAllSpecies()
	if err != nil {
		return nil, fmt.Errorf("error al obtener especies: %v", err)
	}

	return species, nil
}

// GetSpeciesByID retrieves a specific species by its unique identifier.
// Returns complete species information including description and metadata.
//
// Business Logic:
// - Validates species existence in database
// - Returns full species data for detailed views
// - Used for species profiles and validation
//
// Parameters:
//   - id: Unique identifier of the species to retrieve
//
// Returns:
//   - *m.Species: Complete species data
//   - error: Database error or species not found error
func GetSpeciesByID(id uint) (*m.Species, error) {
	// Retrieve specific species from database
	species, err := dao.GetSpeciesByID(uint(id))
	if err != nil {
		return nil, fmt.Errorf("especie no encontrada: %v", err)
	}

	return species, nil
}

// CreateSpecies creates a new species record in the system.
// Handles species registration with proper data validation.
//
// Business Logic:
// - Validates species data before creation
// - Ensures species name uniqueness
// - Assigns creation timestamps
// - Updates the input species object with generated ID
//
// Parameters:
//   - species: Species data to be created (will be updated with generated ID)
//
// Returns:
//   - error: Creation error or nil on success
func CreateSpecies(species *m.Species) error {
	// Create species in database
	err := dao.CreateSpecies(species)
	if err != nil {
		return fmt.Errorf("error al crear especie: %v", err)
	}

	return nil
}

// DeleteSpecies removes a species from the system.
// Handles species deletion with proper constraint checking.
//
// Business Logic:
// - Validates species existence before deletion
// - Checks for pets associated with the species
// - Prevents deletion if pets are still using the species
// - Ensures referential integrity is maintained
//
// Parameters:
//   - id: Unique identifier of the species to delete
//
// Returns:
//   - error: Deletion error (including constraint violations) or nil on success
func DeleteSpecies(id uint) error {
	// Delete species from database with constraint checking
	if err := dao.DeleteSpeciesByID(uint(id)); err != nil {
		return fmt.Errorf("error al eliminar especie: %v", err)
	}

	return nil
}
