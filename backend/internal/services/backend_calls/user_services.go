// Package services provides business logic services for user management and authentication.
// This layer sits between handlers and DAOs, implementing the core business rules
// and orchestrating database operations and external service calls.
package services

import (
	"context"
	r_models "backend/internal/api/routes/models"
	"backend/internal/db/dao"
	m "backend/internal/models"
	mailer "backend/internal/services/mail"
	"fmt"
	"google.golang.org/api/idtoken"
)

// ========================================
// AUTHENTICATION SERVICES
// ========================================

// AuthenticateUser performs the first step of user authentication (login).
// It validates user credentials and prepares the user for 2FA verification.
//
// Process:
// 1. Validates email and password against database
// 2. Increments failed login attempts on failure
// 3. Resets failed login attempts on success
// 4. Generates a new session ID for the user
// 5. Returns the authenticated user data
//
// Parameters:
//   - userData: LoginRequest containing email and password
//
// Returns:
//   - *m.User: User data if authentication successful
//   - error: Authentication error or nil on success
func AuthenticateUser(userData r_models.LoginRequest) (*m.User, error) {
	// Validate user credentials against database
	_, err := dao.GetValidatedUser(userData.Email, userData.Password)

	if err != nil {
		// Increment failed login attempts for security tracking
		dao.IncrementFailedLogins(userData.Email)
		return nil, err
	}

	// Reset failed login attempts and generate new session ID for successful login
	dao.ResetFailedLogins(userData.Email)
	dao.GenerateSessionID(userData.Email)

	// Retrieve updated user data with new session information
	user, _ := dao.GetValidatedUser(userData.Email, userData.Password)

	return user, nil
}

// AuthenticateUser2FA performs the second step of user authentication (2FA verification).
// It validates the 2FA code against the user's session and completes the authentication process.
//
// Process:
// 1. Retrieves user by session ID
// 2. Validates 2FA code against stored value
// 3. Resets failed login attempts on successful verification
// 4. Returns the verified user data
//
// Parameters:
//   - userData: TwoFactorRequest containing session ID and 2FA code
//
// Returns:
//   - *m.NonValidatedUser: User data if 2FA verification successful
//   - error: Verification error or nil on success
func AuthenticateUser2FA(userData r_models.TwoFactorRequest) (*m.NonValidatedUser, error) {
	// Retrieve user by session ID
	user, err := dao.GetUserBySessionID(userData.SessionID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario: %v", err)
	}

	// Retrieve stored 2FA code for validation
	_2fa, err := dao.Get2FA(userData.SessionID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener 2fa: %v", err)
	}

	// Validate 2FA code
	if _2fa == "" || _2fa != userData.Code {
		return nil, fmt.Errorf("código de autenticación de dos factores inválido")
	}

	// Reset failed login attempts after successful 2FA authentication
	dao.ResetFailedLogins(user.Email)

	// Update user's data
	validatedUser, _ := dao.GetUserBySessionID(userData.SessionID)

	return validatedUser, nil
}

// RefreshUser2FAToken generates and sends a new 2FA token to the user's email.
// This is used when the user needs a new 2FA code (expired, lost, etc.).
//
// Process:
// 1. Generates a new 2FA token in the database
// 2. Sends the token via email to the user
// 3. Returns the generated token for verification
//
// Parameters:
//   - userData: RefreshTokenRequest containing user email
//
// Returns:
//   - string: Generated 2FA token
//   - error: Generation or sending error, nil on success
func RefreshUser2FAToken(userData r_models.RefreshTokenRequest) (string, error) {
	// Generate new 2FA token and update in database
	generated2FAToken, _2faErr := dao.UpdateTwoFactorCode(userData.Email)

	if generated2FAToken == "" || _2faErr != nil {
		return "", fmt.Errorf("error al generar el token 2FA: %v", _2faErr)
	}

	// Send 2FA token via email
	mailerErr := mailer.Send2FAToken(userData.Email, generated2FAToken)
	if mailerErr != nil {
		return "", fmt.Errorf("error al enviar el token 2FA al email %s: %v", userData.Email, mailerErr)
	}

	return generated2FAToken, nil
}

// AuthenticateGoogleUser handles Google OAuth authentication.
// It verifies the Google ID token and creates/updates the user account.
//
// Process:
// 1. Verifies Google ID token using Google's public keys
// 2. Extracts user information from the verified token
// 3. Creates new user account if doesn't exist, or updates existing one
// 4. Generates session ID for the user
// 5. Returns user data without requiring 2FA (per requirements)
//
// Parameters:
//   - userData: GoogleLoginRequest containing Google auth data
//
// Returns:
//   - *m.User: Authenticated user data
//   - error: Authentication error or nil on success
func AuthenticateGoogleUser(userData r_models.GoogleLoginRequest) (*m.User, error) {
	// Verify Google ID token
	payload, err := verifyGoogleToken(userData.IDToken)
	if err != nil {
		return nil, fmt.Errorf("token de Google inválido: %v", err)
	}

	// Extract user information from verified token
	email, ok := payload["email"].(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("no se pudo obtener el email del token de Google")
	}

	name, _ := payload["given_name"].(string)
	surname, _ := payload["family_name"].(string)
	googleID, _ := payload["sub"].(string)

	// Check if user exists
	existingUser, err := dao.GetUserByEmail(email)
	if err != nil {
		// User doesn't exist, create new account
		fullUser := &m.FullUser{
			Name:       name,
			Surname:    surname,
			Email:      email,
			Provider:   "google",
			ProviderID: googleID,
			Password:   "", // No password for Google users
		}

		err = dao.CreateUser(fullUser)
		if err != nil {
			return nil, fmt.Errorf("error al crear usuario con Google: %v", err)
		}
	} else {
		// User exists, update provider information if needed
		if existingUser.Provider != "google" {
			// Update existing user to Google provider
			err = dao.UpdateUser(&m.User{
				ID:         existingUser.ID,
				Name:       existingUser.Name,
				Surname:    existingUser.Surname,
				Email:      existingUser.Email,
				Address:    existingUser.Address,
				Provider:   "google",
				ProviderID: googleID,
			})
			if err != nil {
				return nil, fmt.Errorf("error al actualizar usuario con Google: %v", err)
			}
		}
	}

	// Generate session ID for the user
	sessionID, err := dao.GenerateSessionID(email)
	if err != nil {
		return nil, fmt.Errorf("error al generar sessionID: %v", err)
	}

	// Get complete user data with session
	user, err := dao.GetValidatedUser(email, "")
	if err != nil {
		// For Google users, we need to get user data differently since there's no password
		nonValidatedUser, getUserErr := dao.GetUserByEmail(email)
		if getUserErr != nil {
			return nil, fmt.Errorf("error al obtener usuario: %v", getUserErr)
		}
		
		// Convert NonValidatedUser to User for response
		user = &m.User{
			ID:           nonValidatedUser.ID,
			Name:         nonValidatedUser.Name,
			Surname:      nonValidatedUser.Surname,
			Email:        nonValidatedUser.Email,
			Address:      nonValidatedUser.Address,
			Provider:     "google",
			ProviderID:   googleID,
			SessionID:    sessionID,
			FailedLogins: nonValidatedUser.FailedLogins,
			IsBlocked:    nonValidatedUser.IsBlocked,
		}
	}

	return user, nil
}

// ResetPassword resets the password for a user with the given email address.
// It generates a new password and updates it in the database.
//
// Parameters:
//   - email: The email address of the user whose password should be reset
//
// Returns:
//   - *string: A pointer to the new generated password if successful
//   - error: An error if the password reset operation fails
func ResetPassword(email string) (*string, error) {
	user, _ := dao.GetUserByEmail(email)

	if user.Provider != "local" {
		return nil, fmt.Errorf("proveedor debe ser 'local' para restablecer contraseña")
	}

	password, err := dao.ResetPassword(email)
	if err != nil {
		return nil, fmt.Errorf("error al reiniciar la contraseña: %v", err)
	}

	return password, nil
}

func SendNewPassword(email string, password string) error {
	// Set ChangePassword flag to true
	err := dao.SetChangePasswordFlag(email, true)
	if err != nil {
		return fmt.Errorf("error al establecer la bandera de cambio de contraseña: %v", err)
	}

	// Send new password via email
	mailerErr := mailer.SendPassword(email, password)
	if mailerErr != nil {
		return fmt.Errorf("error al enviar la nueva contraseña al email %s: %v", email, mailerErr)
	}

	return nil
}

// ========================================
// USER MANAGEMENT SERVICES
// ========================================

// ListAllUsers retrieves all users from the database.
// Returns non-validated user data (without sensitive information like passwords).
//
// Returns:
//   - *[]m.NonValidatedUser: Slice of all users without sensitive data
//   - error: Database error or nil on success
func ListAllUsers() (*[]m.NonValidatedUser, error) {
	users, err := dao.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("error al leer usuarios: %v", err)
	}

	return &users, nil
}

// GetUserProfile retrieves a specific user by their ID.
// Returns non-validated user data (without sensitive information like passwords).
//
// Parameters:
//   - id: User ID to retrieve
//
// Returns:
//   - *m.NonValidatedUser: User data without sensitive information
//   - error: Database error or nil on success
func GetUserProfile(id uint) (*m.NonValidatedUser, error) {
	user, err := dao.GetUserByID(uint(id))
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario con id: %d %v", id, err)
	}

	return user, nil
}

// RegisterUser creates a new user account in the system.
// Handles the complete user registration process including validation and storage.
//
// Parameters:
//   - user: FullUser data containing all registration information
//
// Returns:
//   - error: Registration error or nil on success
func RegisterUser(user *m.FullUser) error {
	err := dao.CreateUser(user)
	if err != nil {
		return fmt.Errorf("error al crear usuario: %v", err)
	}

	return nil
}

// UpdateUserProfile updates an existing user's profile information.
// Handles partial updates and validates data integrity.
//
// Parameters:
//   - user: User data with updated information
//
// Returns:
//   - error: Update error or nil on success
func UpdateUserProfile(user *m.User) error {
	err := dao.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("error al actualizar usuario: %v", err)
	}

	return nil
}

func UpdateUserPassword(email string, password string) error {
	// Check if the user exists
	_, err := dao.GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("usuario no encontrado %s: %v", email, err)
	}

	// Update the user's password
	err = dao.UpdatePassword(email, password)
	if err != nil {
		return fmt.Errorf("error al actualizar la contraseña del usuario %s: %v", email, err)
	}

	return nil
}

// DeactivateUser soft-deletes a user by marking them as inactive.
// This preserves data integrity while removing user access.
//
// Parameters:
//   - id: User ID to deactivate
//
// Returns:
//   - *m.SimplifiedUser: Simplified user data of deactivated user
//   - error: Deactivation error or nil on success
func DeactivateUser(id uint) (*m.SimplifiedUser, error) {
	deleted, err := dao.DeleteUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("error al eliminar usuario con id: %d %v", id, err)
	}

	return deleted, nil
}

// ========================================
// GOOGLE AUTHENTICATION HELPERS
// ========================================

// verifyGoogleToken verifies a Google ID token and returns the payload
// Uses Google's public keys to verify the token signature and validity
//
// Parameters:
//   - idToken: The Google ID token to verify
//
// Returns:
//   - map[string]interface{}: The verified token payload containing user info
//   - error: Verification error or nil on success
func verifyGoogleToken(idToken string) (map[string]interface{}, error) {
	// Google Client ID - replace with your actual client ID
	clientID := "1017473621019-9hbmho8kqgq7pjhvjl4nqsjq6kc6q5qv.apps.googleusercontent.com"
	
	// Verify the token using Google's idtoken package
	payload, err := idtoken.Validate(context.Background(), idToken, clientID)
	if err != nil {
		return nil, fmt.Errorf("error al verificar el token de Google: %v", err)
	}

	// Convert the payload to a map for easier access
	claims := make(map[string]interface{})
	
	// Extract common claims
	claims["sub"] = payload.Subject
	claims["email"] = payload.Claims["email"]
	claims["given_name"] = payload.Claims["given_name"]
	claims["family_name"] = payload.Claims["family_name"]
	claims["name"] = payload.Claims["name"]
	claims["picture"] = payload.Claims["picture"]
	
	return claims, nil
}
