// Package services provides business logic services for pet management.
// This layer sits between handlers and DAOs, implementing the core business rules
// and orchestrating database operations for pet-related functionality.
package services

import (
	r_models "backend/internal/api/routes/models"
	"backend/internal/db/dao"
	"backend/internal/models"
	m "backend/internal/models"
	"backend/internal/services/mail"
	"fmt"
	"time"
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

// GetFilteredPets retrieves pets that match specified filtering criteria.
// Returns simplified pet data filtered by multiple optional parameters.
//
// Business Logic:
// - Applies multiple filters simultaneously (name, status, species, gender, vaccination)
// - Empty or zero values for parameters are ignored in filtering
// - Returns simplified data to optimize performance for listing views
// - Used for pet search and filtering functionality
//
// Parameters:
//   - name: Pet name filter (partial match, empty string ignores filter)
//   - status: Pet status filter (e.g., "available", "adopted", empty string ignores filter)
//   - speciesID: Species identifier filter (0 ignores filter)
//   - gender: Pet gender filter (empty string ignores filter)
//   - vaccinated: Vaccination status filter (empty string ignores filter)
//
// Returns:
//   - []m.SimplifiedPet: Slice of pets matching the filter criteria
//   - error: Database error or nil on success
func GetFilteredPets(name string, status string, speciesID int, gender string, vaccinated string, urgent *bool) ([]m.SimplifiedPet, error) {
	pets, err := dao.GetFilteredPets(name, status, speciesID, gender, vaccinated, urgent)
	if err != nil {
		return nil, fmt.Errorf("error al obtener mascotas filtradas: %v", err)
	}

	return pets, nil
}

// ========================================
// PET ADOPTION SERVICES
// ========================================

// AdoptPet processes a pet adoption request.
// Updates the pet's status to adopted and records adoption details.
//
// Business Logic:
// - Validates pet is available for adoption
// - Updates pet status to "adopted"
// - Records adoption information
// - Ensures data consistency during the adoption process
//
// Parameters:
//   - petID: Unique identifier of the pet to adopt
//   - userID: Unique identifier of the adopting user
//
// Returns:
//   - error: Adoption error or nil on success
func AdoptPet(petID uint, userID uint) (*models.Pet, error) {
	// Verify pet exists and is available for adoption
	pet, err := dao.GetPetByID(petID)
	if err != nil {
		return nil, fmt.Errorf("mascota no encontrada: %v", err)
	}

	if pet.Status != "available" {
		return nil, fmt.Errorf("mascota no disponible para adopción")
	}

	// Update pet status to adopted
	pet.Status = m.PetStatusAdopted
	pet.AdoptUserID = &userID

	adoptDate := time.Now()
	pet.AdoptDate = &adoptDate

	err = dao.UpdatePet(pet)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar estado de adopción: %v", err)
	}

	// Get updated pet data with relationships
	pet, err = dao.GetPetByID(petID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener mascota actualizada: %v", err)
	}

	return pet, nil
}

func SendPetFosterHomeRequest(to string, pet m.Pet, user m.SimplifiedUser) error {
	err := mail.SendPetFosterHomeRequest(to, pet, user)

	if err != nil {
		return fmt.Errorf("error al enviar solicitud de casa de acogida: %v", err)
	}

	return nil
}

func SendPetAdoptionRequest(to string, pet m.Pet, user m.SimplifiedUser) error {
	err := mail.SendPetAdoptionRequest(to, pet, user)

	if err != nil {
		return fmt.Errorf("error al enviar solicitud de adopción: %v", err)
	}

	return nil
}

func SendPetFosterHomeContact(to string, pet m.Pet, user m.SimplifiedUser, contact m.SimplifiedUser, reason string, message string) error {
	err := mail.SendPetFosterHomeContact(to, pet, user, contact, reason, message)

	if err != nil {
		return fmt.Errorf("error al enviar contacto de casa de acogida: %v", err)
	}

	return nil
}

func UploadPetImage(base64 r_models.Base64Image) error {
	// Call the service to handle image upload
	err := dao.UploadPetImage(base64)
	if err != nil {
		return fmt.Errorf("error al cargar imagen de mascota: %v", err)
	}

	// Return success response
	return nil
}

// DeletePetImageByURL deletes a pet image by its URL or file path.
// This service handles the deletion of pet images from the file system.
//
// Business Logic:
// - Validates the image URL/path format
// - Delegates to DAO layer for actual file deletion
// - Used when editing pets and removing specific images
//
// Parameters:
//   - imageUrl: The URL or file path of the image to delete
//
// Returns:
//   - error: File system error or nil on success
func DeletePetImageByURL(imageUrl string) error {
	// Call the DAO to handle image deletion
	err := dao.DeletePetImageByURL(imageUrl)
	if err != nil {
		return fmt.Errorf("error al eliminar imagen de mascota: %v", err)
	}

	// Return success response
	return nil
}
