// Package api implements HTTP route handlers and endpoint registration for pet management.
// This layer is responsible for:
// - HTTP endpoint registration and routing for pet operations
// - Request binding and basic input validation
// - Calling appropriate handler functions for pet management
// - HTTP response formatting and status code management
// - RESTful API design compliance for pet resources
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

// RegisterPetRoutes registers all pet-related HTTP endpoints with the Echo router.
// Implements standard RESTful API design patterns for pet resource management.
//
// Endpoint Organization:
// - GET /api/pets: List all pets
// - GET /api/pets/:id: Get specific pet by ID
// - POST /api/pets: Create new pet
// - PUT /api/pets/:id: Update existing pet
// - DELETE /api/pets/:id: Delete pet by ID
//
// Parameters:
//   - e: Echo router instance for endpoint registration
func RegisterPetRoutes(e *echo.Echo) {
	e.GET("/api/pets", handleListPets)
	e.GET("/api/pets/:id", handleGetPetByID)
	e.POST("/api/pets", handleCreatePet)
	e.PUT("/api/pets/:id", handleUpdatePet)
	e.DELETE("/api/pets/:id", handleDeletePet)
}

// ========================================
// PET MANAGEMENT ROUTE HANDLERS
// ========================================

// handleListPets processes requests to retrieve all pets in the system.
// Returns a simplified view of pets suitable for listing and browsing.
//
// HTTP Method: GET
// Endpoint: /api/pets
//
// Response:
//   - Success: Array of simplified pet data with adoption status
//   - Error: HTTP error with appropriate status code
func handleListPets(c echo.Context) error {
	// Delegate pet listing to handler layer
	pets, httpErr := handlers.HandleListPets()
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}

	return response.MarshalResponse(c, pets)
}

// handleGetPetByID processes requests to retrieve a specific pet by its ID.
// Returns complete pet information including all details and relationships.
//
// HTTP Method: GET
// Endpoint: /api/pets/:id
// Path Parameters:
//   - id: Pet ID to retrieve
//
// Response:
//   - Success: Complete pet data with all information
//   - Error: HTTP error with appropriate status code
func handleGetPetByID(c echo.Context) error {
	// Extract and validate pet ID from path parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "ID de mascota inválido")
	}

	// Delegate pet retrieval to handler layer
	pet, httpErr = handlers.HandleGetPetByID(uint(id))
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}

	return response.MarshalResponse(c, pet)
}

// handleCreatePet processes pet creation requests.
// Creates new pet records with proper validation and data integrity.
//
// HTTP Method: POST
// Endpoint: /api/pets
// Content-Type: application/json
//
// Request Body:
//   - Pet data for registration (name, species, breed, age, etc.)
//
// Response:
//   - Success: Created pet data with assigned ID and timestamps
//   - Error: HTTP error with appropriate status code
func handleCreatePet(c echo.Context) error {
	var pet m.Pet

	// Bind and validate request body
	if err := c.Bind(&pet); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "datos de mascota inválidos")
	}

	// Delegate pet creation to handler layer
	created, httpErr := handlers.HandleCreatePet(&pet)
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}

	return response.MarshalResponse(c, created)
}

// handleUpdatePet processes pet update requests.
// Updates existing pet information with proper validation.
//
// HTTP Method: PUT
// Endpoint: /api/pets/:id
// Path Parameters:
//   - id: Pet ID to update
//
// Content-Type: application/json
//
// Request Body:
//   - Updated pet data
//
// Response:
//   - Success: Updated pet data
//   - Error: HTTP error with appropriate status code
func handleUpdatePet(c echo.Context) error {
	// Extract and validate pet ID from path parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "ID de mascota inválido")
	}

	var pet m.Pet

	// Bind and validate request body
	if err := c.Bind(&pet); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "datos de mascota inválidos")
	}

	// Ensure URL ID matches request body
	pet.ID = uint(id)

	// Delegate pet update to handler layer
	updated, httpErr := handlers.HandleUpdatePet(&pet)
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}

	return response.MarshalResponse(c, updated)
}

// handleDeletePet processes pet deletion requests.
// Removes pets from the system with proper constraint checking.
//
// HTTP Method: DELETE
// Endpoint: /api/pets/:id
// Path Parameters:
//   - id: Pet ID to delete
//
// Response:
//   - Success: Deletion confirmation message
//   - Error: HTTP error with appropriate status code
func handleDeletePet(c echo.Context) error {
	// Extract and validate pet ID from path parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "ID de mascota inválido")
	}

	// Delegate pet deletion to handler layer
	httpErr := handlers.HandleDeletePet(uint(id))
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}

	// Return deletion confirmation
	return response.MarshalResponse(c, map[string]string{"status": "deleted"})
}
