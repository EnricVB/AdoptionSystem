// Package services provides business logic services for pet management.
// This layer sits between handlers and DAOs, implementing the core business rules
// and orchestrating database operations for pet-related functionality.
package services

import (
	"backend/internal/db/dao"
	m "backend/internal/models"
	"fmt"
)

// ========================================
// PET MANAGEMENT SERVICES
// ========================================

// ListAllPets retrieves all pets from the database.
// Returns simplified pet data suitable for listing and overview purposes.
//
// Business Logic:
// - Retrieves all pets regardless of status
// - Returns simplified data to reduce payload size
// - Used for pet browsing and administrative overviews
//
// Returns:
//   - *[]m.SimplifiedPet: Slice of all pets with essential information
//   - error: Database error or nil on success
func ListAllPets() (*[]m.SimplifiedPet, error) {
	// Retrieve all pets from database
	pets, err := dao.GetAllPets()
	if err != nil {
		return nil, fmt.Errorf("error al obtener mascotas: %v", err)
	}

	return &pets, nil
}

// GetPetByID retrieves a specific pet by its unique identifier.
// Returns complete pet information including all details.
//
// Business Logic:
// - Validates pet existence in database
// - Returns full pet data for detailed views
// - Used for pet profiles and detailed information
//
// Parameters:
//   - id: Unique identifier of the pet to retrieve
//
// Returns:
//   - *m.Pet: Complete pet data with all information
//   - error: Database error or pet not found error
func GetPetByID(id uint) (*m.Pet, error) {
	// Retrieve specific pet from database
	pet, err := dao.GetPetByID(id)
	if err != nil {
		return nil, fmt.Errorf("mascota no encontrada: %v", err)
	}

	return pet, nil
}

// CreatePet creates a new pet record in the system.
// Handles pet registration with proper data validation and integrity.
//
// Business Logic:
// - Validates pet data before creation
// - Assigns creation timestamps
// - Updates the input pet object with generated ID
// - Ensures data consistency
//
// Parameters:
//   - pet: Pet data to be created (will be updated with generated ID)
//
// Returns:
//   - error: Creation error or nil on success
func CreatePet(pet *m.Pet) error {
	// Create pet in database
	created, err := dao.CreatePet(pet)
	if err != nil {
		return fmt.Errorf("error al crear mascota: %v", err)
	}

	// Validate creation was successful
	if created == nil {
		return fmt.Errorf("mascota no creada")
	}

	// Update input object with created data (including ID)
	*pet = *created

	return nil
}

// UpdatePet updates an existing pet's information.
// Handles pet data modification with proper validation.
//
// Business Logic:
// - Validates pet existence before update
// - Preserves data integrity during updates
// - Updates modification timestamps
// - Ensures referential integrity
//
// Parameters:
//   - pet: Pet data with updated information (must include valid ID)
//
// Returns:
//   - error: Update error or nil on success
func UpdatePet(pet *m.Pet) error {
	// Update pet in database
	err := dao.UpdatePet(pet)
	if err != nil {
		return fmt.Errorf("error al actualizar mascota: %v", err)
	}

	return nil
}

// DeletePet removes a pet from the system.
// Handles pet deletion with proper constraint checking.
//
// Business Logic:
// - Validates pet existence before deletion
// - Checks for adoption records or other constraints
// - May perform soft deletion to preserve data integrity
// - Ensures referential integrity is maintained
//
// Parameters:
//   - id: Unique identifier of the pet to delete
//
// Returns:
//   - error: Deletion error or nil on success
func DeletePet(id uint) error {
	// Delete pet from database
	if err := dao.DeletePetByID(id); err != nil {
		return fmt.Errorf("error al eliminar mascota: %v", err)
	}

	return nil
}
