// Package services provides business logic services for user management and authentication.
// This layer sits between handlers and DAOs, implementing the core business rules
// and orchestrating database operations and external service calls.
package services

import (
	r_models "backend/internal/api/routes/models"
	"backend/internal/db/dao"
	m "backend/internal/models"
	mailer "backend/internal/services/mail"
	"fmt"
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
// TODO: Implement Google OAuth authentication logic
//
// Parameters:
//   - userData: GoogleLoginRequest containing Google auth data
//
// Returns:
//   - *m.User: Authenticated user data
//   - error: Authentication error or nil on success
func AuthenticateGoogleUser(userData r_models.GoogleLoginRequest) (*m.User, error) {
	// TODO: Implement Google OAuth authentication
	var user *m.User
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
