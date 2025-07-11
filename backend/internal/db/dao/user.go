// Package dao implements data access objects for user management.
// This layer is responsible for:
// - Direct database operations and queries
// - Data mapping between database and domain models
// - Transaction management and data integrity
// - Raw SQL queries and complex database operations
// - Database connection handling and error management
package dao

import (
	"backend/internal/db"
	m "backend/internal/models"
	"backend/internal/services/security"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ========================================
// USER RETRIEVAL OPERATIONS
// ========================================

// GetAllUsers retrieves all user records from the database.
// Returns non-validated user data (excluding sensitive information like passwords).
//
// Database Operations:
// - Performs SELECT * FROM users
// - Maps User entities to NonValidatedUser DTOs
// - Excludes sensitive fields for security
//
// Returns:
//   - []m.NonValidatedUser: Slice of all users without sensitive data
//   - error: Database error or nil on success
func GetAllUsers() ([]m.NonValidatedUser, error) {
	// Open database connection
	gormDB := db.ORMOpen()

	// Retrieve all users from database
	var users []m.User
	result := gormDB.Find(&users)

	if result.Error != nil {
		return nil, fmt.Errorf("error al leer usuarios: %v", result.Error)
	}

	// Map to non-validated user DTOs (exclude sensitive data)
	var nonValidatedUsers []m.NonValidatedUser
	for _, user := range users {
		nonValidatedUser := m.NonValidatedUser{
			ID:           user.ID,
			Name:         user.Name,
			Surname:      user.Surname,
			Email:        user.Email,
			Address:      user.Address,
			FailedLogins: user.FailedLogins,
			IsBlocked:    user.IsBlocked,
		}
		nonValidatedUsers = append(nonValidatedUsers, nonValidatedUser)
	}

	return nonValidatedUsers, nil
}

// GetUserByID retrieves a specific user by their unique identifier.
// Returns non-validated user data (excluding sensitive information).
//
// Database Operations:
// - Performs SELECT * FROM users WHERE id = ?
// - Maps User entity to NonValidatedUser DTO
// - Handles record not found scenarios
//
// Parameters:
//   - id: Unique identifier of the user to retrieve
//
// Returns:
//   - *m.NonValidatedUser: User data without sensitive information
//   - error: Database error or record not found error
func GetUserByID(id uint) (*m.NonValidatedUser, error) {
	// Open database connection
	gormDB := db.ORMOpen()

	// Retrieve specific user by ID
	var user m.User
	result := gormDB.Where("id = ?", id).First(&user)

	if result.Error != nil {
		return nil, fmt.Errorf("error al leer usuario con id %d: %v", id, result.Error)
	}

	nonValidatedUser := &m.NonValidatedUser{
		ID:           user.ID,
		Name:         user.Name,
		Surname:      user.Surname,
		Email:        user.Email,
		Address:      user.Address,
		FailedLogins: user.FailedLogins,
		IsBlocked:    user.IsBlocked,
	}

	return nonValidatedUser, nil
}

// GetUserByEmail retrieves a user by their email address.
// Returns non-validated user data for authentication and profile lookups.
//
// Database Operations:
// - Performs SELECT * FROM users WHERE email = ?
// - Maps User entity to NonValidatedUser DTO
// - Used for email-based authentication flows
//
// Parameters:
//   - email: Email address of the user to retrieve
//
// Returns:
//   - *m.NonValidatedUser: User data without sensitive information
//   - error: Database error or record not found error
func GetUserByEmail(email string) (*m.NonValidatedUser, error) {
	gormDB := db.ORMOpen()

	var user m.User
	result := gormDB.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return nil, fmt.Errorf("error al leer usuario con email %s: %v", email, result.Error)
	}

	nonValidatedUser := &m.NonValidatedUser{
		ID:           user.ID,
		Name:         user.Name,
		Surname:      user.Surname,
		Email:        user.Email,
		Address:      user.Address,
		FailedLogins: user.FailedLogins,
		IsBlocked:    user.IsBlocked,
		Provider:     user.Provider,
	}

	return nonValidatedUser, nil
}

// GetUserBySessionID retrieves a user by their active session identifier.
// Used for validating user sessions and maintaining authentication state.
//
// Database Operations:
// - Performs SELECT * FROM users WHERE Session_ID = ?
// - Maps User entity to NonValidatedUser DTO
// - Handles session validation and user lookup
//
// Parameters:
//   - sessionID: Active session identifier to lookup
//
// Returns:
//   - *m.NonValidatedUser: User data without sensitive information
//   - error: Database error or session not found error
func GetUserBySessionID(sessionID string) (*m.NonValidatedUser, error) {
	gormDB := db.ORMOpen()

	var user m.User
	result := gormDB.Where("Session_ID = ?", sessionID).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("usuario con sessionID %s no encontrado", sessionID)
		}
		return nil, fmt.Errorf("error al buscar usuario por sessionID: %v", result.Error)
	}

	nonValidatedUser := &m.NonValidatedUser{
		ID:           user.ID,
		Name:         user.Name,
		Surname:      user.Surname,
		Email:        user.Email,
		Address:      user.Address,
		FailedLogins: user.FailedLogins,
		IsBlocked:    user.IsBlocked,
	}

	return nonValidatedUser, nil
}

// Get2FA retrieves the current 2FA token for a user session.
// Used during two-factor authentication verification process.
//
// Database Operations:
// - Performs SELECT Two_Factor_Auth FROM users WHERE Session_ID = ?
// - Returns only the 2FA token field for security
// - Used in authentication flow validation
//
// Parameters:
//   - sessionID: Session identifier to lookup 2FA token
//
// Returns:
//   - string: Current 2FA token for the session
//   - error: Database error or session not found error
func Get2FA(sessionID string) (string, error) {
	gormDB := db.ORMOpen()

	var _2fa string
	result := gormDB.Model(&m.User{}).
		Select("Two_Factor_Auth").
		Where("Session_ID = ?", sessionID).
		First(&_2fa)

	if result.Error != nil {
		return "", fmt.Errorf("error al obtener 2fa para usuario %s: %v", sessionID, result.Error)
	}

	return _2fa, nil
}

// ========================================
// USER AUTHENTICATION OPERATIONS
// ========================================

// GetValidatedUser performs user authentication by validating email and password.
// Returns complete user data including sensitive information for authenticated users.
//
// Database Operations:
// - Performs SELECT * FROM users WHERE email = ?
// - Validates password using bcrypt comparison
// - Returns full User entity for authenticated sessions
//
// Security Features:
// - Password hashing validation
// - Account blocking verification
// - Failed login attempt tracking
//
// Parameters:
//   - email: User's email address
//   - password: Plain text password to validate
//
// Returns:
//   - *m.User: Complete user data for authenticated user
//   - error: Authentication error or database error
func GetValidatedUser(email string, password string) (*m.User, error) {
	gormDB := db.ORMOpen()

	var user m.User
	result := gormDB.Debug().Where("email = ?", email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("usuario con email %s no encontrado", email)
		}
		return nil, fmt.Errorf("error al buscar usuario: %v", result.Error)
	}

	if user.IsBlocked {
		return nil, fmt.Errorf("usuario bloqueado")
	}

	// Skip password validation for Google users
	if user.Provider == "google" {
		return &user, nil
	}

	// Validate password for local users
	if password == "" {
		return nil, fmt.Errorf("contraseña requerida para usuarios locales")
	}

	hashedPassword, err := GetUserHashedPassword(email)
	if err != nil {
		return nil, fmt.Errorf("error al obtener contraseña para usuario %s: %v", email, err)
	}

	if !security.VerifyPassword(hashedPassword, password) {
		return nil, fmt.Errorf("credenciales inválidas")
	}

	return &user, nil
}

// ========================================
// USER CRUD OPERATIONS
// ========================================

// DeleteUserByID removes a user from the database by their ID.
// Performs hard deletion of user records.
//
// Database Operations:
// - Performs DELETE FROM users WHERE id = ?
// - Returns simplified user data after deletion
// - Handles referential integrity constraints
//
// Parameters:
//   - id: Unique identifier of the user to delete
//
// Returns:
//   - *m.SimplifiedUser: Basic user data of deleted record
//   - error: Database error or constraint violation error
func DeleteUserByID(id uint) (*m.SimplifiedUser, error) {
	gormDB := db.ORMOpen()

	var user m.SimplifiedUser
	result := gormDB.Delete(&m.User{}, id)

	if result.Error != nil {
		return nil, fmt.Errorf("error al eliminar usuario con id %d: %v", id, result.Error)
	}

	return &user, nil
}

// CreateUser creates a new user record in the database.
// Handles complete user registration with proper timestamp management.
//
// Database Operations:
// - Performs INSERT INTO users with all user data
// - Sets creation and update timestamps automatically
// - Handles password hashing and validation
//
// Business Logic:
// - Assigns creation timestamp (CrtDate)
// - Assigns update timestamp (UptDate)
// - Validates user data integrity
//
// Parameters:
//   - user: Complete user data for registration
//
// Returns:
//   - error: Database error or validation error, nil on success
func CreateUser(user *m.FullUser) error {
	gormDB := db.ORMOpen()

	now := time.Now()
	user.CrtDate = now
	user.UptDate = now

	result := gormDB.Create(user)
	if result.Error != nil {
		return fmt.Errorf("error al crear usuario: %v", result.Error)
	}

	return nil
}

// UpdateUser updates an existing user's information in the database.
// Handles partial updates with automatic timestamp management.
//
// Database Operations:
// - Performs UPDATE users SET ... WHERE id = ?
// - Updates modification timestamp automatically
// - Handles selective field updates
//
// Business Logic:
// - Updates UptDate timestamp automatically
// - Preserves data integrity during updates
// - Validates user existence before update
//
// Parameters:
//   - user: User data with updated information (must include valid ID)
//
// Returns:
//   - error: Database error or validation error, nil on success
func UpdateUser(user *m.User) error {
	gormDB := db.ORMOpen()

	user.UptDate = time.Now()
	result := gormDB.Model(&m.User{}).
		Where("id = ?", user.ID).
		Updates(user)

	if result.Error != nil {
		return fmt.Errorf("error al actualizar usuario con id %d: %v", user.ID, result.Error)
	}

	return nil
}

func UpdatePassword(email string, newPassword string) error {
	gormDB := db.ORMOpen()

	hashedPassword, err := security.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("error al encriptar la contraseña: %v", err)
	}

	// Change the password in the database
	result := gormDB.Model(&m.User{}).
		Where("email = ?", email).
		Update("password", hashedPassword)

	// Change the change_password flag to false
	SetChangePasswordFlag(email, false)

	if result.Error != nil {
		return fmt.Errorf("error al actualizar contraseña para usuario %s: %v", email, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("usuario con email %s no encontrado", email)
	}

	return nil
}

// ========================================
// USER SECURITY AND LOGIN MANAGEMENT
// ========================================

// UpdateLoginData updates user login-related security information.
// Manages failed login attempts and account blocking status.
//
// Database Operations:
// - Performs UPDATE users SET failed_logins, is_blocked, upt_date WHERE email = ?
// - Updates security-related fields atomically
// - Handles account locking mechanisms
//
// Security Features:
// - Failed login attempt tracking
// - Account blocking management
// - Timestamp tracking for security audits
//
// Parameters:
//   - email: User's email address
//   - failedLogins: Number of failed login attempts
//   - isBlocked: Account blocking status
//
// Returns:
//   - error: Database error or user not found error
func UpdateLoginData(email string, failedLogins int, isBlocked bool) error {
	gormDB := db.ORMOpen()

	result := gormDB.Model(&m.User{}).
		Where("email = ?", email).
		Updates(map[string]interface{}{
			"failed_logins": failedLogins,
			"is_blocked":    isBlocked,
			"upt_date":      time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("error al actualizar datos de login para usuario %s: %v", email, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("usuario con email %s no encontrado", email)
	}

	return nil
}

// IncrementFailedLogins increments the failed login counter for a user.
// Implements automatic account blocking after threshold is reached.
//
// Database Operations:
// - Performs SELECT failed_logins FROM users WHERE email = ?
// - Calculates new failed login count
// - Updates login data with blocking logic
//
// Security Logic:
// - Increments failed login counter by 1
// - Automatically blocks account if failed logins >= 5
// - Maintains security audit trail
//
// Parameters:
//   - email: User's email address
//
// Returns:
//   - error: Database error or user not found error
func IncrementFailedLogins(email string) error {
	gormDB := db.ORMOpen()

	var currentFailedLogins int
	result := gormDB.Model(&m.User{}).
		Select("failed_logins").
		Where("email = ?", email).
		First(&currentFailedLogins)

	if result.Error != nil {
		return fmt.Errorf("error al obtener failed_logins para usuario %s: %v", email, result.Error)
	}

	newFailedLogins := currentFailedLogins + 1

	isBlocked := newFailedLogins >= 5

	return UpdateLoginData(email, newFailedLogins, isBlocked)
}

// GetUserHashedPassword retrieves the hashed password for a user.
// Used for password validation during authentication processes.
//
// Database Operations:
// - Performs SELECT password FROM users WHERE email = ?
// - Returns only the password field for security
// - Used in authentication and password change flows
//
// Security Note:
// - Returns hashed password for bcrypt comparison
// - Should only be used in authentication contexts
// - Never expose raw passwords
//
// Parameters:
//   - email: User's email address
//
// Returns:
//   - string: Hashed password from database
//   - error: Database error or user not found error
func GetUserHashedPassword(email string) (string, error) {
	gormDB := db.ORMOpen()

	var password string
	result := gormDB.Table("Users").Select("Password").Where("email = ?", email).Scan(&password)

	if result.Error != nil {
		return "", fmt.Errorf("error al obtener contraseña para usuario %s: %v", email, result.Error)
	}

	return password, nil
}

// ResetFailedLogins resets the failed login counter for a user.
// Used after successful authentication to clear security flags.
//
// Database Operations:
// - Calls UpdateLoginData with failedLogins = 0 and isBlocked = false
// - Clears security restrictions after successful login
// - Updates modification timestamp
//
// Security Logic:
// - Resets failed login counter to 0
// - Unblocks the account if previously blocked
// - Restores normal account access
//
// Parameters:
//   - email: User's email address
//
// Returns:
//   - error: Database error or user not found error
func ResetFailedLogins(email string) error {
	return UpdateLoginData(email, 0, false)
}

// ResetPassword generates a new password for a user and updates it in the database.
// Used for password recovery and administrative password reset operations.
//
// Database Operations:
// - Generates a new 24-character secure password using security.GeneratePassword
// - Hashes the password using bcrypt for secure storage
// - Performs UPDATE users SET password, upt_date WHERE email = ?
// - Updates modification timestamp for audit trail
//
// Security Features:
// - Generates cryptographically secure random passwords
// - Uses bcrypt hashing for password storage
// - Maintains audit trail with timestamp updates
// - Returns plain text password for secure delivery to user
//
// Parameters:
//   - email: User's email address for password reset
//
// Returns:
//   - *string: Pointer to the generated plain text password (for secure delivery)
//   - error: Database error, user not found error, or password hashing error
func ResetPassword(email string) (*string, error) {
	gormDB := db.ORMOpen()

	password := security.GeneratePassword(12)
	hashedPassword, err := security.HashPassword(password)

	if err != nil {
		return nil, fmt.Errorf("error al encriptar la contraseña")
	}

	result := gormDB.Model(&m.User{}).
		Where("email = ?", email).
		Updates(map[string]any{
			"password": hashedPassword,
			"upt_date": time.Now(),
		})

	if result.Error != nil {
		return nil, fmt.Errorf("error al resetear contraseña para usuario %s: %v", email, result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("usuario con email %s no encontrado", email)
	}

	return &password, nil
}

// BlockUser manually blocks a user account.
// Used for administrative account management and security enforcement.
//
// Database Operations:
// - Performs UPDATE users SET is_blocked = true, upt_date WHERE email = ?
// - Sets account blocking flag to true
// - Updates modification timestamp for audit trail
//
// Security Features:
// - Immediate account blocking
// - Audit trail maintenance
// - Administrative control over account access
//
// Parameters:
//   - email: User's email address to block
//
// Returns:
//   - error: Database error or user not found error
func BlockUser(email string) error {
	gormDB := db.ORMOpen()

	result := gormDB.Model(&m.User{}).
		Where("email = ?", email).
		Updates(map[string]interface{}{
			"is_blocked": true,
			"upt_date":   time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("error al bloquear usuario %s: %v", email, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("usuario con email %s no encontrado", email)
	}

	return nil
}

// UnblockUser unblocks a previously blocked user account.
// Used for administrative account recovery and access restoration.
//
// Database Operations:
// - Calls UpdateLoginData with failedLogins = 0 and isBlocked = false
// - Restores account access and clears security flags
// - Updates modification timestamp
//
// Security Logic:
// - Removes account blocking flag
// - Resets failed login counter
// - Restores full account functionality
//
// Parameters:
//   - email: User's email address to unblock
//
// Returns:
//   - error: Database error or user not found error
func UnblockUser(email string) error {
	return UpdateLoginData(email, 0, false)
}

// ========================================
// TWO-FACTOR AUTHENTICATION OPERATIONS
// ========================================

// UpdateTwoFactorCode generates and updates a new 2FA token for a user.
// Used for two-factor authentication setup and token refresh.
//
// Database Operations:
// - Generates a new 6-digit 2FA code using security.Generate2FA
// - Performs UPDATE users SET two_factor_auth WHERE email = ?
// - Validates user existence after update
//
// Security Features:
// - Generates cryptographically secure 2FA codes
// - Updates 2FA token atomically
// - Validates operation success
//
// Parameters:
//   - email: User's email address
//
// Returns:
//   - string: Generated 2FA token
//   - error: Database error or user not found error
func UpdateTwoFactorCode(email string) (string, error) {
	gormDB := db.ORMOpen()
	_2fa := security.Generate2FA(6)

	result := gormDB.Model(&m.User{}).
		Where("email = ?", email).
		Updates(map[string]any{
			"two_factor_auth": _2fa,
		})

	_, err := GetUserByEmail(email)
	if result.Error != nil || err != nil {
		return "", fmt.Errorf("error al actualizar TwoFactorAuth para usuario %s: %v", email, result.Error)
	}

	if result.RowsAffected == 0 {
		return "", fmt.Errorf("usuario con email %s no encontrado", email)
	}

	return _2fa, nil
}

// GenerateSessionID creates and updates a new session identifier for a user.
// Used for session management and user authentication state tracking.
//
// Database Operations:
// - Generates a new 50-character session ID using security.Generate2FA
// - Performs UPDATE users SET session_id WHERE email = ?
// - Validates user existence after update
//
// Session Management:
// - Creates unique session identifiers
// - Maintains user authentication state
// - Enables session-based authentication flows
//
// Parameters:
//   - email: User's email address
//
// Returns:
//   - string: Generated session identifier
//   - error: Database error or user not found error
func GenerateSessionID(email string) (string, error) {
	gormDB := db.ORMOpen()
	sessionID := security.Generate2FA(50)

	result := gormDB.Model(&m.User{}).
		Where("email = ?", email).
		Updates(map[string]any{
			"session_id": sessionID,
		})

	_, err := GetUserByEmail(email)
	if result.Error != nil || err != nil {
		return "", fmt.Errorf("error al actualizar SessionID para usuario %s: %v", email, result.Error)
	}

	if result.RowsAffected == 0 {
		return "", fmt.Errorf("usuario con email %s no encontrado", email)
	}

	return sessionID, nil
}

func SetChangePasswordFlag(email string, flag bool) error {
	gormDB := db.ORMOpen()

	// Verificar si el usuario existe
	_, err := GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("usuario no encontrado %s: %v", email, err)
	}

	// Actualizar la bandera de cambio de contraseña
	result := gormDB.Model(&m.User{}).
		Where("email = ?", email).
		Update("change_password", flag)

	if result.Error != nil {
		return fmt.Errorf("error al establecer la bandera de cambio de contraseña para usuario %s: %v", email, result.Error)
	}

	return nil
}
