// Package handlers implements HTTP request handlers for the pet management API.
// This layer is responsible for:
// - HTTP request/response handling and validation
// - Input sanitization and basic validation
// - Calling appropriate service layer functions
// - Converting service errors to HTTP responses
// - Ensuring consistent API response formatting
package handlers

import (
	m "backend/internal/models"
	s "backend/internal/services/backend_calls"
	response "backend/internal/utils/rest"
	"net/http"
)

// ========================================
// PET MANAGEMENT HANDLERS
// ========================================

// HandleListPets processes requests to retrieve all pets in the system.
// Returns a simplified view of pets suitable for listing purposes.
//
// Returns:
//   - *[]m.SimplifiedPet: List of all pets with essential information
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleListPets() (*[]m.SimplifiedPet, response.HTTPError) {
	// Delegate pet listing to service layer
	pets, err := s.ListAllPets()
	if err != nil {
		return nil, response.Error(http.StatusInternalServerError, err.Error())
	}

	return pets, response.EmptyError
}

// HandleGetPetByID processes requests to retrieve a specific pet by its ID.
// Returns complete pet information including all details.
//
// Validation:
// - Ensures pet ID is valid (greater than 0)
// - Delegates pet retrieval to service layer
//
// Parameters:
//   - id: Pet ID to retrieve
//
// Returns:
//   - *m.Pet: Complete pet data with all information
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleGetPetByID(id uint) (*m.Pet, response.HTTPError) {
	// Input validation
	if id <= 0 {
		return nil, response.Error(http.StatusBadRequest, "ID de mascota no válido")
	}

	// Delegate pet retrieval to service layer
	pet, err := s.GetPetByID(id)
	if err != nil {
		return nil, response.Error(http.StatusNotFound, err.Error())
	}

	return pet, response.EmptyError
}

// HandleCreatePet processes pet creation requests.
// Creates new pet records with proper validation and data integrity.
//
// Validation:
// - Ensures required fields are provided (name and species are mandatory)
// - Delegates creation logic and business rules to service layer
//
// Parameters:
//   - pet: Pet data for the new pet to be created
//
// Returns:
//   - *m.Pet: Created pet data with assigned ID and timestamps
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleCreatePet(pet *m.Pet) (*m.Pet, response.HTTPError) {
	// Input validation
	if pet.Name == "" || pet.Species == "" {
		return nil, response.Error(http.StatusBadRequest, "nombre y especie de mascota son obligatorios")
	}

	// Delegate pet creation to service layer
	err := s.CreatePet(pet)
	if err != nil {
		return nil, response.Error(http.StatusInternalServerError, err.Error())
	}

	return pet, response.EmptyError
}

// HandleUpdatePet processes pet update requests.
// Updates existing pet information with proper validation.
//
// Validation:
// - Ensures pet ID is valid (greater than 0)
// - Ensures required fields are provided (name and species are mandatory)
// - Delegates update logic and business rules to service layer
//
// Parameters:
//   - pet: Pet data with updated information (must include valid ID)
//
// Returns:
//   - *m.Pet: Updated pet data
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleUpdatePet(pet *m.Pet) (*m.Pet, response.HTTPError) {
	// Input validation
	if pet.ID <= 0 {
		return nil, response.Error(http.StatusBadRequest, "ID de mascota no válido")
	}

	if pet.Name == "" || pet.Species == "" {
		return nil, response.Error(http.StatusBadRequest, "nombre y especie de mascota son obligatorios")
	}

	// Delegate pet update to service layer
	err := s.UpdatePet(pet)
	if err != nil {
		return nil, response.Error(http.StatusInternalServerError, err.Error())
	}

	return pet, response.EmptyError
}

// HandleDeletePet processes pet deletion requests.
// Performs deletion with proper validation and business rule enforcement.
//
// Validation:
// - Ensures pet ID is valid (greater than 0)
// - Delegates deletion logic and constraints to service layer
//
// Parameters:
//   - id: Pet ID to delete
//
// Returns:
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleDeletePet(id uint) response.HTTPError {
	// Input validation
	if id <= 0 {
		return response.Error(http.StatusBadRequest, "ID de mascota no válido")
	}

	// Delegate pet deletion to service layer
	err := s.DeletePet(id)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.EmptyError
}
