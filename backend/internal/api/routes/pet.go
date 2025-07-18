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
	r_models "backend/internal/api/routes/models"
	m "backend/internal/models"
	response "backend/internal/utils/rest"
	"backend/internal/utils/time"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
	e.GET("/api/filtered-pets", handleGetFilteredPets)

	// Adoption and foster home requests
	e.POST("/api/pets/adopt", handlePetAdoptionRequest)
	e.POST("/api/pets/foster-home", handlePetFosterHomeRequest)
	e.POST("/api/pets/foster-home/contact", handlePetFosterHomeContact)
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

// handleGetPetByID retrieves detailed information for a specific pet by ID.
// Returns complete pet data including species, images, vaccination history, and fostering information.
//
// HTTP Method: GET
// Endpoint: /api/pets/:id
// Path Parameters:
//   - id: Unique identifier for the pet (numeric)
//
// Response:
//   - Success: Complete pet information with all related data
//   - Error: HTTP error with appropriate status code and error message
//
// Business Flow:
// 1. Extracts and validates pet ID from URL path
// 2. Requests detailed pet information from handler layer
// 3. Returns complete pet data with all relationships
func handleGetPetByID(c echo.Context) error {
	// Extract and validate pet ID from path parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest,
			"ID de mascota inválido: debe ser un número")
	}

	// Retrieve detailed pet information from handler layer
	pet, httpErr := handlers.HandleGetPetByID(uint(id))
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}

	return response.MarshalResponse(c, pet)
}

// handleCreatePet processes requests to create new pet records in the system.
// Creates new pet entries with proper validation and data integrity checks.
//
// HTTP Method: POST
// Endpoint: /api/pets
// Content-Type: application/json
//
// Request Body:
//   - Pet data for registration (name, species, breed, age, gender, etc.)
//   - Required fields: Name, SpeciesID
//   - Optional fields: BirthDate, AdoptDate (defaults applied if missing)
//
// Response:
//   - Success: Created pet data with assigned ID and timestamps
//   - Error: HTTP error with appropriate status code and error message
//
// Business Flow:
// 1. Binds and validates incoming pet data
// 2. Validates required fields (name, species)
// 3. Applies default dates for birth and adoption if missing
// 4. Delegates creation to handler layer
// 5. Returns newly created pet information
func handleCreatePet(c echo.Context) error {
	var pet m.Pet

	// Bind and validate request body
	if err := c.Bind(&pet); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest,
			fmt.Sprintf("Datos de mascota inválidos: %v", err))
	}

	// Validate required fieldsº
	if pet.Name == "" || pet.SpeciesID == 0 {
		return response.ErrorResponse(c, http.StatusBadRequest,
			"El nombre y la especie son campos requeridos")
	}

	// Apply default dates if missing
	if pet.BirthDate.IsZero() {
		pet.BirthDate = time.DefaultDate()
	}

	if pet.AdoptDate != nil && pet.AdoptDate.IsZero() {
		adoptDate := time.DefaultDate()
		pet.AdoptDate = &adoptDate
	}

	// Delegate pet creation to handler layer
	createdPet, httpErr := handlers.HandleCreatePet(&pet)
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}

	return response.MarshalResponse(c, createdPet)
}

// handleUpdatePet processes requests to update existing pet information.
// Updates pet records with proper validation and data integrity checks.
//
// HTTP Method: PUT
// Endpoint: /api/pets/:id
// Path Parameters:
//   - id: Unique identifier for the pet to update (numeric)
//
// Content-Type: application/json
//
// Request Body:
//   - Updated pet data (all fields that need modification)
//   - ID in body will be overridden by path parameter
//
// Response:
//   - Success: Updated pet data with all modifications applied
//   - Error: HTTP error with appropriate status code and error message
//
// Business Flow:
// 1. Extracts and validates pet ID from URL path
// 2. Binds and validates updated pet data from request body
// 3. Ensures URL ID matches the pet being updated
// 4. Applies default dates if missing
// 5. Delegates update to handler layer
// 6. Returns updated pet information
func handleUpdatePet(c echo.Context) error {
	// Extract and validate pet ID from path parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest,
			"ID de mascota inválido: debe ser un número")
	}

	var pet m.Pet

	// Bind and validate request body
	if err := c.Bind(&pet); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest,
			fmt.Sprintf("Datos de mascota inválidos: %v", err))
	}

	// Ensure URL ID matches request body for consistency
	pet.ID = uint(id)

	// Apply default dates if missing
	if pet.BirthDate.IsZero() {
		pet.BirthDate = time.DefaultDate()
	}
	if pet.AdoptDate.IsZero() {
		adoptDate := time.DefaultDate()
		pet.AdoptDate = &adoptDate
	}

	// Delegate pet update to handler layer
	updatedPet, httpErr := handlers.HandleUpdatePet(&pet)
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}

	return response.MarshalResponse(c, updatedPet)
}

// handleDeletePet processes requests to remove pets from the system.
// Removes pet records with proper constraint checking and data integrity validation.
//
// HTTP Method: DELETE
// Endpoint: /api/pets/:id
// Path Parameters:
//   - id: Unique identifier for the pet to delete (numeric)
//
// Response:
//   - Success: Deletion confirmation message with status
//   - Error: HTTP error with appropriate status code and error message
//
// Business Flow:
// 1. Extracts and validates pet ID from URL path
// 2. Delegates deletion logic to handler layer (checks constraints)
// 3. Returns confirmation of successful deletion
//
// Note: Handler layer will validate if pet can be safely deleted
// (e.g., no active adoption requests, foster relationships, etc.)
func handleDeletePet(c echo.Context) error {
	// Extract and validate pet ID from path parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest,
			"ID de mascota inválido: debe ser un número")
	}

	// Delegate pet deletion to handler layer
	httpErr := handlers.HandleDeletePet(uint(id))
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}

	// Return deletion confirmation
	return response.MarshalResponse(c, map[string]string{
		"status":  "deleted",
		"message": "Mascota eliminada exitosamente",
	})
}

// handleGetFilteredPets processes requests to retrieve pets with filtering criteria.
// Returns filtered pet data based on query parameters for search and browsing functionality.
//
// HTTP Method: GET
// Endpoint: /api/filtered-pets
// Query Parameters:
//   - name: Pet name filter (partial match, optional)
//   - status: Pet status filter (Available, Adopted, FosterHome, optional)
//   - species_id: Species ID filter (numeric, optional)
//   - gender: Pet gender filter (Male, Female, optional)
//   - vaccinated: Vaccination status filter (optional)
//
// Response:
//   - Success: Array of simplified pet data matching the filter criteria
//   - Error: HTTP error with appropriate status code and error message
//
// Business Flow:
// 1. Extracts and validates query parameters
// 2. Validates species_id parameter format if provided
// 3. Delegates filtering logic to handler layer
// 4. Returns filtered pet results
func handleGetFilteredPets(c echo.Context) error {
	// Extract query parameters
	name := strings.TrimSpace(c.QueryParam("name"))
	status := strings.TrimSpace(c.QueryParam("status"))
	speciesIDStr := strings.TrimSpace(c.QueryParam("species_id"))
	gender := strings.TrimSpace(c.QueryParam("gender"))
	vaccinated := strings.TrimSpace(c.QueryParam("vaccinated"))

	// Parse species_id parameter if provided
	var speciesID int
	if speciesIDStr != "" {
		var err error
		speciesID, err = strconv.Atoi(speciesIDStr)
		if err != nil {
			return response.ErrorResponse(c, http.StatusBadRequest,
				"ID de especie inválido: debe ser un número")
		}
	}

	// Delegate filtering logic to handler layer
	pets, httpErr := handlers.HandleGetFilteredPets(name, status, speciesID, gender, vaccinated)
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}

	return response.MarshalResponse(c, pets)
}

// handlePetAdoptionRequest processes pet adoption request submissions.
// Sends adoption request emails to the organization for review and processing.
//
// HTTP Method: POST
// Endpoint: /api/pets/adopt
// Content-Type: application/json
//
// Request Body:
//   - PetRequest: Adoption request data including pet ID and user ID
//
// Response:
//   - Success: Confirmation message indicating the adoption request was sent
//   - Error: HTTP error with appropriate status code and error message
//
// Business Flow:
// 1. Validates request data binding
// 2. Retrieves company information for email configuration
// 3. Validates pet existence and availability
// 4. Validates adopting user existence
// 5. Converts user data to simplified format for security
// 6. Sends adoption request email to organization
func handlePetAdoptionRequest(c echo.Context) error {
	var request r_models.PetRequest

	// Bind and validate request body
	if err := c.Bind(&request); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest,
			fmt.Sprintf("datos de solicitud de adopción inválidos: %v", err))
	}

	// Retrieve company information for email configuration
	company, httpErr := handlers.HandleGetCompanyInfo()
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}

	// Validate pet existence and availability
	pet, httpErr := handlers.HandleGetPetByID(request.PetID)
	if httpErr.Code != 0 {
		return response.ErrorResponse(c, http.StatusNotFound,
			fmt.Sprintf("mascota no encontrada: %v", httpErr))
	}

	// Validate adopting user existence
	adoptingUser, httpErr := handlers.HandleGetUserByID(request.UserID)
	if httpErr.Code != 0 {
		return response.ErrorResponse(c, http.StatusNotFound,
			fmt.Sprintf("usuario adoptante no encontrado: %v", httpErr))
	}

	// Convert user data to simplified format for security
	simplifiedUser := adoptingUser.ToSimplifiedUser()

	// Send adoption request email to organization
	httpErr = handlers.HandleSendPetAdoptionRequest(company.ContactEmail, *pet, simplifiedUser)
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}

	return response.MarshalResponse(c, map[string]string{
		"status": "adoption request sent",
	})
}

// handlePetFosterHomeRequest processes foster home request submissions.
// Sends foster home request emails to the organization for review and processing.
//
// HTTP Method: POST
// Endpoint: /api/pets/foster-home
// Content-Type: application/json
//
// Request Body:
//   - PetRequest: Foster home request data including pet ID and user ID
//
// Response:
//   - Success: Confirmation message indicating the foster home request was sent
//   - Error: HTTP error with appropriate status code and error message
//
// Business Flow:
// 1. Validates request data binding
// 2. Retrieves company information for email configuration
// 3. Validates pet existence and availability
// 4. Validates foster family user existence
// 5. Converts user data to simplified format for security
// 6. Sends foster home request email to organization
func handlePetFosterHomeRequest(c echo.Context) error {
	var request r_models.PetRequest

	// Bind and validate request body
	if err := c.Bind(&request); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest,
			fmt.Sprintf("datos de solicitud de casa de acogida inválidos: %v", err))
	}

	// Retrieve company information for email configuration
	company, httpErr := handlers.HandleGetCompanyInfo()
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}

	// Validate pet existence and availability
	pet, httpErr := handlers.HandleGetPetByID(request.PetID)
	if httpErr.Code != 0 {
		return response.ErrorResponse(c, http.StatusNotFound,
			fmt.Sprintf("mascota no encontrada: %v", httpErr))
	}

	// Validate user existence
	user, httpErr := handlers.HandleGetUserByID(request.UserID)

	if httpErr.Code != 0 {
		return response.ErrorResponse(c, http.StatusNotFound,
			fmt.Sprintf("usuario no encontrado: %v", httpErr))
	}

	// Convert user data to simplified format for security
	simplifiedUser := user.ToSimplifiedUser()

	// Send foster home request email to organization
	httpErr = handlers.HandleSendPetFosterHomeRequest(company.ContactEmail, *pet, simplifiedUser)
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}

	return response.MarshalResponse(c, map[string]string{
		"status": "foster home request sent",
	})
}

// handlePetFosterHomeContact processes requests to send contact messages to foster families.
// Facilitates communication between interested parties and foster families caring for pets.
//
// HTTP Method: POST
// Endpoint: /api/pets/foster-home/contact
// Content-Type: application/json
//
// Request Body:
//   - PetFosterHomeContact: Contact request data including pet ID, user IDs, reason, and message
//
// Response:
//   - Success: Confirmation message indicating the contact request was sent
//   - Error: HTTP error with appropriate status code and error message
//
// Business Flow:
// 1. Validates request data binding
// 2. Retrieves company information for email configuration
// 3. Validates pet existence
// 4. Validates contact user (foster family) existence
// 5. Validates requesting user existence
// 6. Converts user data to simplified format for security
// 7. Sends contact email to foster family
func handlePetFosterHomeContact(c echo.Context) error {
	var request r_models.PetFosterHomeContact

	// Bind and validate request body
	if err := c.Bind(&request); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest,
			fmt.Sprintf("datos de solicitud de contacto inválidos: %v", err))
	}

	// Retrieve company information for email configuration
	company, httpErr := handlers.HandleGetCompanyInfo()
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}

	// Validate pet existence
	pet, httpErr := handlers.HandleGetPetByID(request.PetID)
	if httpErr.Code != 0 {
		return response.ErrorResponse(c, http.StatusNotFound,
			fmt.Sprintf("mascota no encontrada: %v", httpErr))
	}

	// Validate foster family contact user existence
	fosterFamilyUser, httpErr := handlers.HandleGetUserByID(request.ContactUserID)
	if httpErr.Code != 0 {
		return response.ErrorResponse(c, http.StatusNotFound,
			fmt.Sprintf("usuario de familia de acogida no encontrado: %v", httpErr))
	}

	// Validate requesting user existence
	requestingUser, httpErr := handlers.HandleGetUserByID(request.UserID)
	if httpErr.Code != 0 {
		return response.ErrorResponse(c, http.StatusNotFound,
			fmt.Sprintf("usuario solicitante no encontrado: %v", httpErr))
	}

	// Convert user data to simplified format for security
	simplifiedFosterFamily := fosterFamilyUser.ToSimplifiedUser()
	simplifiedRequestingUser := requestingUser.ToSimplifiedUser()

	// Send contact email to foster family
	httpErr = handlers.HandleSendPetFosterHomeContact(
		company.ContactEmail,
		*pet,
		simplifiedRequestingUser,
		simplifiedFosterFamily,
		request.Reason,
		request.Message,
	)
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}

	return response.MarshalResponse(c, map[string]string{
		"status": "foster home contact request sent",
	})
}
