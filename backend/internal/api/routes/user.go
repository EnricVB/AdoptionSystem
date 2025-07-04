// Package api implements HTTP route handlers and endpoint registration for user management.
// This layer is responsible for:
// - HTTP endpoint registration and routing
// - Request binding and basic input validation
// - Calling appropriate handler functions
// - HTTP response formatting and status code management
// - Middleware integration for authentication and logging
package api

import (
	"backend/internal/api/handlers"
	r_models "backend/internal/api/routes/models"
	m "backend/internal/models"
	response "backend/internal/utils/rest"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// ========================================
// ROUTE REGISTRATION
// ========================================

// RegisterUserRoutes registers all user-related HTTP endpoints with the Echo router.
// Organizes endpoints by functionality for better maintainability.
//
// Endpoint Organization:
// - User CRUD operations: Standard REST endpoints for user management
// - Authentication endpoints: Login and 2FA verification endpoints
//
// Parameters:
//   - e: Echo router instance for endpoint registration
func RegisterUserRoutes(e *echo.Echo) {
	// User CRUD operations
	e.GET("/api/users", handleListUsers)
	e.GET("/api/users/:id", handleGetUserByID)
	e.POST("/api/users", handleCreateUser)
	e.PUT("/api/users/:id", handleUpdateUser)
	e.DELETE("/api/users/:id", handleDeleteUser)

	// Authentication endpoints
	e.POST("/api/auth/login", handleLoginUser)
	e.POST("/api/auth/login/google", handleLoginWithGoogle) // Use dedicated Google login handler
	e.POST("/api/auth/verify-2fa", handle2FAAuth)
	e.POST("/api/auth/refresh-token", handleRefresh2FAToken)
}

// ========================================
// AUTHENTICATION ROUTE HANDLERS
// ========================================

// handleLoginUser processes user login requests.
// Implements the first step of two-factor authentication.
//
// HTTP Method: POST
// Endpoint: /api/auth/login
// Content-Type: application/json
//
// Request Body:
//   - email: User's email address
//   - password: User's password
//
// Response:
//   - Success: User data with session information
//   - Error: HTTP error with appropriate status code
func handleLoginUser(c echo.Context) error {
	var req r_models.LoginRequest

	// Bind and validate request body
	if err := c.Bind(&req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	// Delegate authentication to handler layer
	user, err := handlers.HandleManualLogin(req)
	if err != response.EmptyError {
		return response.ConvertToErrorResponse(c, err)
	}

	return response.MarshalResponse(c, user)
}

// handle2FAAuth processes two-factor authentication verification requests.
// Implements the second step of two-factor authentication.
//
// HTTP Method: POST
// Endpoint: /api/auth/verify-2fa
// Content-Type: application/json
//
// Request Body:
//   - sessionID: Session ID from login step
//   - code: 2FA verification code
//
// Response:
//   - Success: "OK" message indicating successful verification
//   - Error: HTTP error with appropriate status code
func handle2FAAuth(c echo.Context) error {
	var req r_models.TwoFactorRequest

	// Bind and validate request body
	if err := c.Bind(&req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	// Delegate 2FA verification to handler layer
	_, err := handlers.Handle2FAAuth(req)
	if err != response.EmptyError {
		return response.ConvertToErrorResponse(c, err)
	}

	// Return simple success message for 2FA verification
	return response.MarshalResponse(c, "OK")
}

// handleRefresh2FAToken processes requests to generate and resend 2FA tokens.
// Used when users need a new 2FA code (expired, lost, etc.).
//
// HTTP Method: POST
// Endpoint: /api/auth/refresh-token
// Content-Type: application/json
//
// Request Body:
//   - email: User's email address
//
// Response:
//   - Success: "OK" message indicating token was sent
//   - Error: HTTP error with appropriate status code
func handleRefresh2FAToken(c echo.Context) error {
	var req r_models.RefreshTokenRequest

	if err := c.Bind(&req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	_, err := handlers.HandleRefresh2FAToken(req)

	if err != response.EmptyError {
		return response.ConvertToErrorResponse(c, err)
	}

	return response.MarshalResponse(c, "OK")
}

// handleLoginWithGoogle processes Google OAuth login requests.
// Validates input, calls Google login handler, and returns user data or error.
//
// Validations:
//   - Ensures email and ID token are provided
//   - Calls handler to validate Google login
//   - Returns user data on success or error response on failure
//
// Parameters:
//   - c: Echo context containing the HTTP request and response
//
// Returns:
//   - HTTP response with user data on successful Google login
//   - Error response on failure
func handleLoginWithGoogle(c echo.Context) error {
	var req r_models.GoogleLoginRequest

	if err := c.Bind(&req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	user, err := handlers.HandleGoogleLogin(req)

	if err != response.EmptyError {
		return response.ConvertToErrorResponse(c, err)
	}

	return response.MarshalResponse(c, user)
}

// handleListUsers retrieves all users in the system.
// Calls the handler to get the list of users and returns it in the response.
//
// Validations:
//   - No specific validations required for listing users
func handleListUsers(c echo.Context) error {
	users, httpErr := handlers.HandleListUsers()
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}
	return response.MarshalResponse(c, users)
}

// handleGetUserByID retrieves a specific user by their unique identifier.
// This endpoint provides access to individual user data.
//
// Business Rules:
//   - Requires valid user ID as URL parameter
//   - Returns complete user information if found
//   - Should validate user access permissions (implement authorization middleware)
//
// Parameters:
//   - c: Echo context containing the HTTP request and response
//   - id: User ID as URL parameter (e.g., /api/users/{id})
//
// Returns:
//   - HTTP 200 with user object on success
//   - HTTP 400 if user ID is invalid or non-numeric
//   - HTTP 404 if user not found
//   - HTTP 500 on internal server error
//   - Error response with appropriate status code on failure
func handleGetUserByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "ID de usuario inválido")
	}

	user, httpErr := handlers.HandleGetUserByID(uint(id))
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}
	return response.MarshalResponse(c, user)
}

// handleCreateUser creates a new user account in the system.
// This endpoint handles user registration with validation and business rule enforcement.
//
// Business Rules:
//   - Validates all required user fields
//   - Enforces unique email addresses
//   - Performs password strength validation
//   - Creates user with default settings and permissions
//   - Sends welcome email notification (if mail service is configured)
//
// Parameters:
//   - c: Echo context containing the HTTP request and response
//   - Request body: CreateUserRequest JSON containing user details
//
// Request Body Fields:
//   - name: User's full name (required)
//   - email: User's email address (required, must be unique)
//   - password: User's password (required, must meet security requirements)
//   - Additional fields as defined in CreateUserRequest model
//
// Returns:
//   - HTTP 200 with "OK" message on successful creation
//   - HTTP 400 if request data is invalid or missing required fields
//   - HTTP 409 if email already exists
//   - HTTP 500 on internal server error
//   - Error response with appropriate status code on failure
func handleCreateUser(c echo.Context) error {
	var req r_models.CreateUserRequest

	if err := c.Bind(&req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "datos inválidos")
	}

	httpErr := handlers.HandleCreateUser(&req)
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}

	return response.MarshalResponse(c, "OK")
}

// handleUpdateUser modifies an existing user's information.
// This endpoint handles user profile updates with validation and authorization.
//
// Business Rules:
//   - Requires valid user ID as URL parameter
//   - Validates updated fields according to business rules
//   - Maintains data integrity and referential constraints
//   - Updates only provided fields (partial updates supported)
//   - Logs user modification for audit trail
//
// Parameters:
//   - c: Echo context containing the HTTP request and response
//   - id: User ID as URL parameter (e.g., /api/users/{id})
//   - Request body: User data to update (JSON format)
//
// Returns:
//   - HTTP 200 with updated user object on success
//   - HTTP 400 if user ID is invalid or request data is malformed
//   - HTTP 404 if user not found
//   - HTTP 409 if update would violate unique constraints
//   - HTTP 500 on internal server error
//   - Error response with appropriate status code on failure
func handleUpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "ID de usuario inválido")
	}

	var user m.User
	if err := c.Bind(&user); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "datos inválidos")
	}
	user.ID = uint(id)

	httpErr := handlers.HandleUpdateUser(&user)
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}
	return response.MarshalResponse(c, user)
}

// handleDeleteUser removes a user from the system.
// This endpoint handles user account deletion with proper validation and cleanup.
//
// Business Rules:
//   - Requires valid user ID as URL parameter
//   - Performs soft delete (marks as deleted) rather than hard delete
//   - Validates that user can be safely deleted (no critical dependencies)
//   - Maintains referential integrity with related entities
//   - Logs deletion for audit trail and compliance
//   - May require administrative privileges for execution
//
// Parameters:
//   - c: Echo context containing the HTTP request and response
//   - id: User ID as URL parameter (e.g., /api/users/{id})
//
// Returns:
//   - HTTP 200 with deletion confirmation on success
//   - HTTP 400 if user ID is invalid or non-numeric
//   - HTTP 404 if user not found or already deleted
//   - HTTP 409 if user cannot be deleted due to dependencies
//   - HTTP 500 on internal server error
//   - Error response with appropriate status code on failure
func handleDeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "ID de usuario inválido")
	}

	deleted, httpErr := handlers.HandleDeleteUser(uint(id))
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}
	return response.MarshalResponse(c, deleted)
}
