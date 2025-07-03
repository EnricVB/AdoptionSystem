// Package handlers implements HTTP request handlers for the species management API.
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
// SPECIES MANAGEMENT HANDLERS
// ========================================

// HandleListSpecies processes requests to retrieve all species in the system.
// Returns all available species for use in pet registration and filtering.
//
// Returns:
//   - []m.Species: List of all species with their information
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleListSpecies() ([]m.Species, response.HTTPError) {
	// Delegate species listing to service layer
	species, err := s.ListAllSpecies()
	if err != nil {
		return nil, response.Error(http.StatusInternalServerError, err.Error())
	}

	return species, response.EmptyError
}

// HandleGetSpeciesByID processes requests to retrieve a specific species by its ID.
// Returns complete species information including description and metadata.
//
// Validation:
// - Ensures species ID is valid (greater than 0)
// - Delegates species retrieval to service layer
//
// Parameters:
//   - id: Species ID to retrieve
//
// Returns:
//   - *m.Species: Complete species data
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleGetSpeciesByID(id uint) (*m.Species, response.HTTPError) {
	// Input validation
	if id <= 0 {
		return nil, response.Error(http.StatusBadRequest, "ID de especie no válido")
	}

	// Delegate species retrieval to service layer
	species, err := s.GetSpeciesByID(id)
	if err != nil {
		return nil, response.Error(http.StatusNotFound, err.Error())
	}

	return species, response.EmptyError
}

// HandleCreateSpecies processes species creation requests.
// Creates new species records with proper validation.
//
// Validation:
// - Ensures required fields are provided (name is mandatory)
// - Delegates creation logic and business rules to service layer
//
// Parameters:
//   - species: Species data for the new species to be created
//
// Returns:
//   - *m.Species: Created species data with assigned ID
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleCreateSpecies(species *m.Species) (*m.Species, response.HTTPError) {
	// Input validation
	if species.Name == "" {
		return nil, response.Error(http.StatusBadRequest, "nombre de especie es obligatorio")
	}

	// Delegate species creation to service layer
	err := s.CreateSpecies(species)
	if err != nil {
		return nil, response.Error(http.StatusInternalServerError, err.Error())
	}

	return species, response.EmptyError
}

// HandleDeleteSpecies processes species deletion requests.
// Performs deletion with proper validation and constraint checking.
//
// Note: Species deletion may be restricted if there are pets associated with the species.
// The service layer handles these business rules and constraints.
//
// Validation:
// - Ensures species ID is valid (greater than 0)
// - Delegates deletion logic and constraints to service layer
//
// Parameters:
//   - id: Species ID to delete
//
// Returns:
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleDeleteSpecies(id uint) response.HTTPError {
	// Input validation
	if id <= 0 {
		return response.Error(http.StatusBadRequest, "ID de especie no válido")
	}

	// Delegate species deletion to service layer
	err := s.DeleteSpecies(id)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.EmptyError
}
