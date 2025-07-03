// Package handlers implements HTTP request handlers for the user management API.
// This layer is responsible for:
// - HTTP request/response handling and validation
// - Input sanitization and basic validation
// - Calling appropriate service layer functions
// - Converting service errors to HTTP responses
// - Ensuring consistent API response formatting
//
// Error Handling:
// All handlers return consistent HTTP error responses using the response.HTTPError
// type, ensuring uniform error formatting across the API. Common error scenarios:
//   - 400 Bad Request: Invalid input data or missing required fields
//   - 401 Unauthorized: Authentication failures or invalid credentials
//   - 404 Not Found: Requested user not found
//   - 500 Internal Server Error: Service layer or database errors
package handlers

import (
	r_models "backend/internal/api/routes/models"
	"backend/internal/models"
	s "backend/internal/services/backend_calls"
	"backend/internal/services/security"
	response "backend/internal/utils/rest"
	"net/http"
	"time"
)

// ========================================
// AUTHENTICATION HANDLERS
// ========================================

// HandleManualLogin processes manual login requests (email/password authentication).
// This is the first step in the two-factor authentication process.
//
// Validation:
// - Ensures email and password are provided
// - Delegates credential validation to service layer
//
// Parameters:
//   - req: LoginRequest containing user credentials
//
// Returns:
//   - *models.User: Authenticated user data with session information
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleManualLogin(req r_models.LoginRequest) (*models.User, response.HTTPError) {
	// Delegate authentication to service layer
	user, err := s.AuthenticateUser(req)
	if err != nil {
		return nil, response.Error(http.StatusUnauthorized, err.Error())
	}

	return user, response.EmptyError
}

// Handle2FAAuth processes two-factor authentication verification requests.
// This is the second step in the authentication process, validating the 2FA code.
//
// Validation:
// - Ensures session ID and 2FA code are provided
// - Delegates code verification to service layer
//
// Parameters:
//   - req: TwoFactorRequest containing session ID and 2FA code
//
// Returns:
//   - *models.NonValidatedUser: Verified user data (without sensitive info)
//   - response.HTTPError: HTTP error or EmptyError on success
func Handle2FAAuth(req r_models.TwoFactorRequest) (*models.NonValidatedUser, response.HTTPError) {
	// Input validation
	if req.SessionID == "" || req.Code == "" {
		return nil, response.Error(http.StatusBadRequest, "sessionID y código de 2FA son obligatorios")
	}

	// Delegate 2FA verification to service layer
	user, err := s.AuthenticateUser2FA(req)
	if err != nil {
		return nil, response.Error(http.StatusUnauthorized, err.Error())
	}

	return user, response.EmptyError
}

// HandleRefresh2FAToken processes requests to generate and resend 2FA tokens.
// Used when users need a new 2FA code (expired, lost, etc.).
//
// Validation:
// - Ensures email is provided
// - Delegates token generation and email sending to service layer
//
// Parameters:
//   - req: RefreshTokenRequest containing user email
//
// Returns:
//   - string: Generated 2FA token
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleRefresh2FAToken(req r_models.RefreshTokenRequest) (*string, response.HTTPError) {
	// Input validation
	if req.Email == "" {
		return nil, response.Error(http.StatusBadRequest, "email es obligatorio")
	}

	// Delegate token refresh to service layer
	token, err := s.RefreshUser2FAToken(req)
	if err != nil {
		return nil, response.Error(http.StatusUnauthorized, err.Error())
	}

	return &token, response.EmptyError
}

// HandleGoogleLogin processes Google OAuth authentication requests.
// Validates Google ID tokens and creates/updates user accounts.
//
// Validation:
// - Ensures email and ID token are provided
// - Delegates Google OAuth validation to service layer
//
// Parameters:
//   - req: GoogleLoginRequest containing Google authentication data
//
// Returns:
//   - *models.User: Authenticated user data
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleGoogleLogin(req r_models.GoogleLoginRequest) (*models.User, response.HTTPError) {
	// Input validation
	if req.Email == "" || req.IDToken == "" {
		return nil, response.Error(http.StatusBadRequest, "email y ID Token son obligatorios")
	}

	// Delegate Google authentication to service layer
	user, err := s.AuthenticateGoogleUser(req)
	if err != nil {
		return nil, response.Error(http.StatusUnauthorized, err.Error())
	}

	return user, response.EmptyError
}

// HandleResetPassword processes password reset requests for local authentication users.
// Generates a new temporary password and sends it to the user's email address.
//
// Validation:
// - Ensures email is provided
// - Validates that provider is "local" (password resets only for local accounts)
// - Delegates password generation and email sending to service layer
//
// Parameters:
//   - req: ResetPasswordRequest containing user email and provider type
//
// Returns:
//   - *string: Generated temporary password
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleResetPassword(req r_models.ResetPasswordRequest) (*string, response.HTTPError) {
	// Input validation
	if req.Email == "" {
		return nil, response.Error(http.StatusBadRequest, "email es obligatorio")
	}

	// Delegate reset password to service layer
	password, responseError := s.ResetPassword(req.Email)
	if responseError != nil {
		return nil, response.Error(http.StatusInternalServerError, responseError.Error())
	}

	return password, response.EmptyError
}

// handleForgotPassword processes requests to send the new generated password.
//
// Validation:
// - Ensures email is provided
// - Ensures user's provider is local
//
// Parameters:
//   - req: ResetPasswordRequest containing user email and provider
//
// Returns:
//   - string: Generated password
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleForgotPassword(req r_models.ResetPasswordRequest) (*string, response.HTTPError) {
	// Input validation
	if req.Email == "" {
		return nil, response.Error(http.StatusBadRequest, "email es obligatorio")
	}

	// Delegate reset password to service layer
	password, responseError := s.ResetPassword(req.Email)
	if responseError != nil {
		return nil, response.Error(http.StatusInternalServerError, responseError.Error())
	}

	responseError = s.SendNewPassword(req.Email, *password)
	if responseError != nil {
		return nil, response.Error(http.StatusInternalServerError, responseError.Error())
	}

	return password, response.EmptyError
}

// ========================================
// USER MANAGEMENT HANDLERS
// ========================================

// HandleListUsers processes requests to retrieve all users.
// Returns non-sensitive user data for administrative purposes.
//
// Returns:
//   - *[]models.NonValidatedUser: List of all users without sensitive data
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleListUsers() (*[]models.NonValidatedUser, response.HTTPError) {
	// Delegate user listing to service layer
	users, err := s.ListAllUsers()
	if err != nil {
		return nil, response.Error(http.StatusInternalServerError, err.Error())
	}

	return users, response.HTTPError{}
}

// HandleGetUserByID processes requests to retrieve a specific user by ID.
// Returns non-sensitive user data for the requested user.
//
// Validation:
// - Ensures user ID is valid (greater than 0)
// - Delegates user retrieval to service layer
//
// Parameters:
//   - id: User ID to retrieve
//
// Returns:
//   - *models.NonValidatedUser: User data without sensitive information
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleGetUserByID(id uint) (*models.NonValidatedUser, response.HTTPError) {
	// Input validation
	if id <= 0 {
		return nil, response.Error(http.StatusBadRequest, "ID de usuario no válido")
	}

	// Delegate user retrieval to service layer
	user, err := s.GetUserProfile(id)
	if err != nil {
		return nil, response.Error(http.StatusNotFound, err.Error())
	}

	return user, response.HTTPError{}
}

// HandleCreateUser processes user registration requests.
// Creates new user accounts with proper data transformation and validation.
//
// Validation:
// - Ensures required fields are provided (name is mandatory)
// - Transforms request data to internal user model
// - Sets creation and update timestamps
//
// Parameters:
//   - user: CreateUserRequest containing user registration data
//
// Returns:
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleCreateUser(user *r_models.CreateUserRequest) response.HTTPError {
	// Input validation
	if user.Name == "" {
		return response.Error(http.StatusBadRequest, "nombre es obligatorio")
	}

	// Transform request data to internal model
	fullUser := &models.FullUser{
		Name:       user.Name,
		Surname:    user.Surname,
		Email:      user.Email,
		Password:   user.Password,
		Address:    user.Address,
		Provider:   user.Provider,
		ProviderID: user.ProviderID,
		CrtDate:    time.Now(),
		UptDate:    time.Now(),
	}

	// Delegate user creation to service layer
	err := s.RegisterUser(fullUser)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.HTTPError{}
}

// HandleUpdateUser processes user profile update requests.
// Updates existing user information with proper validation.
//
// Validation:
// - Delegates validation and update logic to service layer
// - Ensures data integrity during updates
//
// Parameters:
//   - user: User data with updated information
//
// Returns:
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleUpdateUser(user *models.User) response.HTTPError {
	// Hash password if provided
	if user.Password != "" {
		hashedPassword, err := security.HashPassword(user.Password)
		if err != nil {
			return response.Error(http.StatusInternalServerError, "error al hashear la contraseña")
		}
		user.Password = hashedPassword
	}

	// Automatically set update timestamp
	user.UptDate = time.Now()

	// Delegate user update to service layer
	err := s.UpdateUserProfile(user)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.HTTPError{}
}

func HandleUpdateUserPassword(email string, password string) response.HTTPError {
	// Input validation
	if email == "" || password == "" {
		return response.Error(http.StatusBadRequest, "email y contraseña son obligatorios")
	}

	// Delegate password update to service layer
	err := s.UpdateUserPassword(email, password)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.EmptyError
}

// HandleDeleteUser processes user deletion requests.
// Performs soft deletion to preserve data integrity.
//
// Validation:
// - Ensures user ID is valid (greater than 0)
// - Delegates deletion logic to service layer
//
// Parameters:
//   - id: User ID to delete/deactivate
//
// Returns:
//   - *models.SimplifiedUser: Data of the deleted user
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleDeleteUser(id uint) (*models.SimplifiedUser, response.HTTPError) {
	// Input validation
	if id <= 0 {
		return nil, response.Error(http.StatusBadRequest, "ID de usuario no válido")
	}

	// Delegate user deletion to service layer
	deleted, err := s.DeactivateUser(id)
	if err != nil {
		return nil, response.Error(http.StatusInternalServerError, err.Error())
	}

	return deleted, response.HTTPError{}
}
