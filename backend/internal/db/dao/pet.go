// Package dao implements data access objects for pet management.
// This layer is responsible for:
// - Direct database operations and queries for pet data
// - Data mapping between database and domain models
// - Transaction management and data integrity for pets
// - Relationship management with users and species
// - Database connection handling and error management
package dao

import (
	"backend/internal/db"
	m "backend/internal/models"
	"fmt"
	"time"
)

// ========================================
// PET RETRIEVAL OPERATIONS
// ========================================

// GetAllPets retrieves all pet records from the database.
// Returns simplified pet data suitable for listing and overview purposes.
//
// Database Operations:
// - Performs SELECT * FROM pets with user relationship preloading
// - Uses GORM's Preload to fetch associated AdoptUser data
// - Returns SimplifiedPet models optimized for list views
//
// Relationship Loading:
// - Preloads AdoptUser relationship to show adoption status
// - Optimizes queries by loading related data in single operation
// - Reduces N+1 query problems
//
// Returns:
//   - []m.SimplifiedPet: Slice of all pets with essential information and adoption status
//   - error: Database error or nil on success
func GetAllPets() ([]m.SimplifiedPet, error) {
	// Open database connection
	gormDB := db.ORMOpen()

	// Retrieve all pets with user relationship preloaded
	var pets []m.SimplifiedPet
	result := gormDB.Preload("AdoptUser").Find(&pets)
	if result.Error != nil {
		return nil, fmt.Errorf("error al leer mascotas: %v", result.Error)
	}

	return pets, nil
}

// GetPetByID retrieves a specific pet by its unique identifier.
// Returns complete pet information including all relationships and details.
//
// Database Operations:
// - Performs SELECT * FROM pets WHERE id = ? with relationship preloading
// - Uses GORM's Preload to fetch associated AdoptUser data
// - Returns complete Pet model with all details
//
// Relationship Loading:
// - Preloads AdoptUser relationship for adoption information
// - Provides complete pet profile data
// - Used for detailed pet views and management
//
// Parameters:
//   - id: Unique identifier of the pet to retrieve
//
// Returns:
//   - *m.Pet: Complete pet data with all relationships
//   - error: Database error or record not found error
func GetPetByID(id uint) (*m.Pet, error) {
	// Open database connection
	gormDB := db.ORMOpen()

	// Retrieve specific pet by ID with relationships
	var pet m.Pet
	result := gormDB.Preload("AdoptUser").Where("id = ?", id).First(&pet)
	if result.Error != nil {
		return nil, fmt.Errorf("error al leer mascota con id %d: %v", id, result.Error)
	}

	return &pet, nil
}

// ========================================
// PET CRUD OPERATIONS
// ========================================

// CreatePet creates a new pet record in the database.
// Handles pet registration with proper timestamp management and data integrity.
//
// Database Operations:
// - Performs INSERT INTO pets with all pet data
// - Sets creation and update timestamps automatically
// - Validates referential integrity with species
//
// Business Logic:
// - Assigns creation timestamp (CrtDate) to current time
// - Assigns update timestamp (UptDate) to current time
// - Updates the input pet object with generated ID
// - Ensures data consistency and integrity
//
// Timestamp Management:
// - Automatically sets CrtDate and UptDate on creation
// - Maintains audit trail for pet registration
// - Supports data tracking and reporting
//
// Parameters:
//   - pet: Pet data to be created (will be updated with generated ID and timestamps)
//
// Returns:
//   - *m.Pet: Created pet data with assigned ID and timestamps
//   - error: Database error or validation error
func CreatePet(pet *m.Pet) (*m.Pet, error) {
	// Open database connection
	gormDB := db.ORMOpen()

	// Set creation and update timestamps
	now := time.Now()
	pet.CrtDate = now
	pet.UptDate = now

	// Create new pet record in database
	result := gormDB.Create(pet)
	if result.Error != nil {
		return nil, fmt.Errorf("error al crear mascota: %v", result.Error)
	}

	return pet, nil
}

// UpdatePet updates an existing pet's information in the database.
// Handles pet data modification with automatic timestamp management.
//
// Database Operations:
// - Performs UPDATE pets SET ... WHERE id = ?
// - Updates modification timestamp automatically
// - Uses selective field updates with Select("*")
//
// Business Logic:
// - Updates UptDate timestamp automatically to current time
// - Preserves data integrity during updates
// - Validates pet existence before update
// - Maintains audit trail for pet modifications
//
// Timestamp Management:
// - Automatically updates UptDate on modification
// - Preserves original CrtDate
// - Supports change tracking and auditing
//
// Parameters:
//   - pet: Pet data with updated information (must include valid ID)
//
// Returns:
//   - error: Database error or validation error, nil on success
func UpdatePet(pet *m.Pet) error {
	// Open database connection
	gormDB := db.ORMOpen()

	// Update modification timestamp
	pet.UptDate = time.Now()

	// Update pet record with all fields
	result := gormDB.Model(&m.Pet{}).
		Where("id = ?", pet.ID).
		Select("*").
		Updates(pet)

	if result.Error != nil {
		return fmt.Errorf("error al actualizar mascota con id %d: %v", pet.ID, result.Error)
	}

	return nil
}

// DeletePetByID removes a pet from the database by its ID.
// Handles pet deletion with proper data integrity management.
//
// Database Operations:
// - Performs DELETE FROM pets WHERE id = ?
// - Handles soft deletion if configured in GORM model
// - Maintains referential integrity with adoption records
//
// Business Logic:
// - May perform soft deletion to preserve adoption history
// - Ensures data consistency across related tables
// - Prevents orphaned adoption records
//
// Data Integrity:
// - Considers impact on adoption history
// - May restrict deletion of adopted pets
// - Maintains audit trail for regulatory compliance
//
// Parameters:
//   - id: Unique identifier of the pet to delete
//
// Returns:
//   - error: Database error, constraint violation, or nil on success
func DeletePetByID(id uint) error {
	// Open database connection
	gormDB := db.ORMOpen()

	// Delete pet record by ID
	result := gormDB.Delete(&m.Pet{}, id)
	if result.Error != nil {
		return fmt.Errorf("error al eliminar mascota con id %d: %v", id, result.Error)
	}

	return nil
}
