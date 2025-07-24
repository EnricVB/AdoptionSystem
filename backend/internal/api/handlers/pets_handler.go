// Package handlers implements HTTP request handlers for the pet management API.
// This layer is responsible for:
// - HTTP request/response handling and validation
// - Input sanitization and basic validation
// - Calling appropriate service layer functions
// - Converting service errors to HTTP responses
// - Ensuring consistent API response formatting
package handlers

import (
	r_models "backend/internal/api/routes/models"
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
	if pet.Name == "" || pet.SpeciesID == 0 {
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

	if pet.Name == "" || pet.SpeciesID == 0 {
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

// HandleGetFilteredPets processes requests to retrieve pets based on filtering criteria.
// Returns a filtered list of pets matching the specified search parameters.
//
// Validation:
// - Validates status against allowed values: "", "Available", "Adopted", "FosterHome"
// - Validates gender against allowed values: "", "Male", "Female"
// - Ensures species ID is non-negative (0 means no species filter)
// - Validates vaccination status parameter
// - Delegates filtering logic to service layer
//
// Parameters:
//   - name: Pet name filter (empty string for no filter)
//   - status: Pet status filter (empty string for no filter)
//   - speciesID: Species ID filter (0 for no filter)
//   - gender: Pet gender filter (empty string for no filter)
//   - vaccinated: Vaccination status filter (empty string for no filter)
//
// Returns:
//   - []m.SimplifiedPet: List of pets matching the filter criteria
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleGetFilteredPets(name string, status string, speciesID int, gender string, vaccinated string, urgent *bool) ([]m.SimplifiedPet, response.HTTPError) {
	// Define valid status values
	validStatuses := map[string]bool{
		"":           true, // No filter
		"Available":  true, // Pet available for adoption
		"Adopted":    true, // Pet has been adopted
		"FosterHome": true, // Pet is in foster home
	}

	// Define valid gender values
	validGenders := map[string]bool{
		"":       true, // No filter
		"Male":   true, // Male pet
		"Female": true, // Female pet
	}

	// Input validation
	if !validStatuses[status] {
		return nil, response.Error(http.StatusBadRequest, "Estado de mascota no válido")
	}

	if speciesID < 0 {
		return nil, response.Error(http.StatusBadRequest, "ID de especie no válido")
	}

	if !validGenders[gender] {
		return nil, response.Error(http.StatusBadRequest, "Género de mascota no válido")
	}

	// Delegate filtered pet retrieval to service layer
	pets, err := s.GetFilteredPets(name, status, speciesID, gender, vaccinated, urgent)
	if err != nil {
		return nil, response.Error(http.StatusInternalServerError, err.Error())
	}

	return pets, response.EmptyError
}

// ========================================
// PET ADOPTION HANDLERS
// ========================================

// HandleAdoptPet processes pet adoption requests.
// Updates the pet's status to adopted and records adoption details.
//
// Validation:
// - Ensures pet ID is valid (greater than 0)
// - Ensures user ID is valid (greater than 0)
// - Delegates adoption logic and business rules to service layer
//
// Parameters:
//   - petID: Unique identifier of the pet to adopt
//   - userID: Unique identifier of the adopting user
//
// Returns:
//   - *m.Pet: Updated pet data with adoption information
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleAdoptPet(petID uint, userID uint) (*m.Pet, response.HTTPError) {
	// Input validation
	if petID <= 0 {
		return nil, response.Error(http.StatusBadRequest, "ID de mascota no válido")
	}

	if userID <= 0 {
		return nil, response.Error(http.StatusBadRequest, "ID de usuario no válido")
	}

	// Delegate pet adoption to service layer
	pet, err := s.AdoptPet(petID, userID)
	if err != nil {
		return nil, response.Error(http.StatusInternalServerError, err.Error())
	}

	return pet, response.EmptyError
}

// ========================================
// EMAIL NOTIFICATION HANDLERS
// ========================================

// HandleSendPetFosterHomeRequest processes requests to send foster home request emails.
// Sends notification emails for pet foster home requests.
//
// Validation:
// - Ensures recipient email is provided
// - Ensures pet data is valid (pet object must have required fields)
// - Ensures user data is valid (user object must have required fields)
// - Delegates email sending logic to service layer
//
// Parameters:
//   - to: Recipient email address
//   - pet: Pet data for which the foster home request is being sent
//   - user: User data of the person making the request
//
// Returns:
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleSendPetFosterHomeRequest(to string, pet m.Pet, user m.SimplifiedUser) response.HTTPError {
	// Input validation
	if to == "" {
		return response.Error(http.StatusBadRequest, "dirección de correo destinatario es obligatoria")
	}

	if pet.ID <= 0 || pet.Name == "" {
		return response.Error(http.StatusBadRequest, "datos de mascota no válidos")
	}

	if user.ID <= 0 || user.Name == "" || user.Email == "" {
		return response.Error(http.StatusBadRequest, "datos de usuario no válidos")
	}

	// Delegate email sending to service layer
	err := s.SendPetFosterHomeRequest(to, pet, user)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.EmptyError
}

// HandleSendPetAdoptionRequest processes requests to send adoption request emails.
// Sends notification emails for pet adoption requests.
//
// Validation:
// - Ensures recipient email is provided
// - Ensures pet data is valid (pet object must have required fields)
// - Ensures user data is valid (user object must have required fields)
// - Delegates email sending logic to service layer
//
// Parameters:
//   - to: Recipient email address
//   - pet: Pet data for which the adoption request is being sent
//   - user: User data of the person making the request
//
// Returns:
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleSendPetAdoptionRequest(to string, pet m.Pet, user m.SimplifiedUser) response.HTTPError {
	// Input validation
	if to == "" {
		return response.Error(http.StatusBadRequest, "dirección de correo destinatario es obligatoria")
	}

	if pet.ID <= 0 || pet.Name == "" {
		return response.Error(http.StatusBadRequest, "datos de mascota no válidos")
	}

	if user.ID <= 0 || user.Name == "" || user.Email == "" {
		return response.Error(http.StatusBadRequest, "datos de usuario no válidos")
	}

	// Delegate email sending to service layer
	err := s.SendPetAdoptionRequest(to, pet, user)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.EmptyError
}

// HandleSendPetFosterHomeContact processes requests to send foster home contact emails.
// Sends notification emails for contacting foster families about pets in their care.
//
// Validation:
// - Ensures recipient email is provided
// - Ensures pet data is valid (pet object must have required fields)
// - Ensures contact user data is valid (user making the inquiry)
// - Ensures foster family data is valid (foster family information)
// - Ensures reason and message are provided
// - Delegates email sending logic to service layer
//
// Parameters:
//   - to: Recipient email address
//   - pet: Pet data for which the contact is being made
//   - user: User data of the person making the inquiry
//   - contact: User data of the foster family
//   - reason: Reason for contacting the foster family
//   - message: Message content for the foster family
//
// Returns:
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleSendPetFosterHomeContact(to string, pet m.Pet, user m.SimplifiedUser, contact m.SimplifiedUser, reason string, message string) response.HTTPError {
	// Input validation
	if to == "" {
		return response.Error(http.StatusBadRequest, "dirección de correo destinatario es obligatoria")
	}

	if pet.ID <= 0 || pet.Name == "" {
		return response.Error(http.StatusBadRequest, "datos de mascota no válidos")
	}

	if contact.ID <= 0 || contact.Name == "" || contact.Email == "" {
		return response.Error(http.StatusBadRequest, "datos de contacto no válidos")
	}

	if reason == "" {
		return response.Error(http.StatusBadRequest, "motivo de contacto es obligatorio")
	}

	if message == "" {
		return response.Error(http.StatusBadRequest, "mensaje es obligatorio")
	}

	// Delegate email sending to service layer
	err := s.SendPetFosterHomeContact(to, pet, user, contact, reason, message)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.EmptyError
}

func HandleUploadPetImage(base64 r_models.Base64Image) response.HTTPError {
	// Call the service to handle image upload
	err := s.UploadPetImage(base64)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.EmptyError
}

// HandleDeletePetImageByURL processes requests to delete a pet image by its URL.
// This handler manages the deletion of pet images from the file system.
//
// Parameters:
//   - imageUrl: The URL or file path of the image to delete
//
// Returns:
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleDeletePetImageByURL(imageUrl string) response.HTTPError {
	// Validate input
	if imageUrl == "" {
		return response.Error(http.StatusBadRequest, "URL de imagen es requerida")
	}

	// Call the service to handle image deletion
	err := s.DeletePetImageByURL(imageUrl)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.EmptyError
}
