// Package api implements HTTP route handlers and endpoint registration for species management.
// This layer is responsible for:
// - HTTP endpoint registration and routing for species operations
// - Request binding and basic input validation
// - Calling appropriate handler functions for species management
// - HTTP response formatting and status code management
// - RESTful API design compliance for species resources
package api

import (
	"backend/internal/api/handlers"
	m "backend/internal/models"
	response "backend/internal/utils/rest"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// ========================================
// ROUTE REGISTRATION
// ========================================

// RegisterSpeciesRoutes registers all species-related HTTP endpoints with the Echo router.
// Implements standard RESTful API design patterns for species resource management.
//
// Endpoint Organization:
// - GET /api/species: List all species
// - GET /api/species/:id: Get specific species by ID
// - POST /api/species: Create new species
// - DELETE /api/species/:id: Delete species by ID
//
// Note: PUT endpoint not implemented as species updates are typically restricted
// to maintain data integrity with existing pet records.
//
// Parameters:
//   - e: Echo router instance for endpoint registration
func RegisterSpeciesRoutes(e *echo.Echo) {
	e.GET("/api/species", handleListSpecies)
	e.GET("/api/species/:id", handleGetSpeciesByID)
	e.POST("/api/species", handleCreateSpecies)
	e.DELETE("/api/species/:id", handleDeleteSpecies)
}

// ========================================
// SPECIES MANAGEMENT ROUTE HANDLERS
// ========================================

// handleListSpecies processes requests to retrieve all species in the system.
// Returns all available species for use in pet registration and filtering.
//
// HTTP Method: GET
// Endpoint: /api/species
//
// Response:
//   - Success: Array of all species with complete information
//   - Error: HTTP error with appropriate status code
func handleListSpecies(c echo.Context) error {
	// Delegate species listing to handler layer
	species, httpErr := handlers.HandleListSpecies()
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}

	return response.MarshalResponse(c, species)
}

// handleGetSpeciesByID processes requests to retrieve a specific species by its ID.
// Returns complete species information including description and metadata.
//
// HTTP Method: GET
// Endpoint: /api/species/:id
// Path Parameters:
//   - id: Species ID to retrieve
//
// Response:
//   - Success: Complete species data
//   - Error: HTTP error with appropriate status code
func handleGetSpeciesByID(c echo.Context) error {
	// Extract and validate species ID from path parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "ID de especie inválido")
	}

	// Delegate species retrieval to handler layer
	species, httpErr := handlers.HandleGetSpeciesByID(uint(id))
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}
	return response.MarshalResponse(c, species)
}

// handleCreateSpecies processes species creation requests.
// Creates new species records with proper validation.
//
// HTTP Method: POST
// Endpoint: /api/species
// Content-Type: application/json
//
// Request Body:
//   - Species data for registration (name, description)
//
// Response:
//   - Success: Created species data with assigned ID
//   - Error: HTTP error with appropriate status code
func handleCreateSpecies(c echo.Context) error {
	var species m.Species

	// Bind and validate request body
	if err := c.Bind(&species); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "datos de especie inválidos")
	}

	// Delegate species creation to handler layer
	created, httpErr := handlers.HandleCreateSpecies(&species)
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}

	return response.MarshalResponse(c, created)
}

// handleDeleteSpecies processes species deletion requests.
// Removes species from the system with proper constraint checking.
//
// HTTP Method: DELETE
// Endpoint: /api/species/:id
// Path Parameters:
//   - id: Species ID to delete
//
// Business Rules:
// - Species cannot be deleted if pets are associated with it
// - Deletion will fail with appropriate error if constraints are violated
// - Ensures referential integrity across the pet adoption system
//
// Response:
//   - Success: Deletion confirmation message
//   - Error: HTTP error with appropriate status code (including constraint violations)
func handleDeleteSpecies(c echo.Context) error {
	// Extract and validate species ID from path parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "ID de especie inválido")
	}

	// Delegate species deletion to handler layer
	httpErr = handlers.HandleDeleteSpecies(uint(id))
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}

	// Return deletion confirmation
	return response.MarshalResponse(c, map[string]string{"status": "deleted"})
}
