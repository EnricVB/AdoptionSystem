// Package r_models contains request models for API endpoints.
// These models define the structure of data expected in HTTP request bodies.
package r_models

// LoginRequest represents the request payload for user authentication.
// Used for standard email/password login operations.
//
// Validation Requirements:
//   - Email: Must be a valid email format and exist in the system
//   - Password: Must match the stored password hash for the user
//
// Security Notes:
//   - Password is transmitted in plain text (ensure HTTPS is used)
//   - Failed login attempts are tracked to prevent brute force attacks
type LoginRequest struct {
	Email    string `json:"email"`    // User's email address for authentication
	Password string `json:"password"` // User's plain text password (will be hashed for comparison)
}

// GoogleLoginRequest represents the request payload for Google OAuth authentication.
// Used for authenticating users via Google OAuth 2.0 flow.
//
// Validation Requirements:
//   - Email: Must match the email in the ID token
//   - IDToken: Must be a valid Google ID token that can be verified
//
// Security Notes:
//   - ID token is verified against Google's public keys
//   - Token expiration and audience are validated
//   - If user doesn't exist, a new account may be created automatically
type GoogleLoginRequest struct {
	Email   string `json:"email"`    // User's email address from Google account
	IDToken string `json:"id_token"` // Google OAuth ID token for verification
}

// TwoFactorRequest represents the request payload for two-factor authentication verification.
// Used to verify the second factor in the authentication process.
//
// Validation Requirements:
//   - SessionID: Must be a valid active session ID
//   - Code: Must be a 6-digit numeric code that matches the expected 2FA code
//
// Security Notes:
//   - 2FA codes have a limited time window for validity
//   - Invalid attempts are tracked and may result in session invalidation
//   - Session must be in a pending 2FA state
type TwoFactorRequest struct {
	SessionID string `json:"session_id"` // Session identifier for the pending 2FA verification
	Code      string `json:"code"`       // 6-digit two-factor authentication code
}

// RefreshTokenRequest represents the request payload for refreshing 2FA tokens.
// Used to generate and send a new 2FA code to the user.
//
// Validation Requirements:
//   - Email: Must be a valid email address of an existing user
//   - User must be in a state that allows 2FA token refresh
//
// Business Rules:
//   - Rate limiting may apply to prevent abuse
//   - Previous 2FA codes may be invalidated when new ones are generated
type RefreshTokenRequest struct {
	Email string `json:"email"` // User's email address for 2FA token refresh
}

// CreateUserRequest represents the request payload for user registration.
// Used for creating new user accounts in the system.
//
// Validation Requirements:
//   - Name, Surname: Required, non-empty strings
//   - Email: Must be valid format and unique in the system
//   - Password: Must meet security requirements (if provider is 'local')
//   - Provider: Must be a supported authentication provider
//
// Business Rules:
//   - Email addresses must be unique across all users
//   - For 'local' provider, password is required and will be hashed
//   - For external providers (Google, etc.), ProviderID is required
//   - Default user permissions and settings are applied during creation
type CreateUserRequest struct {
	Name    string `json:"name"`    // User's first name (required)
	Surname string `json:"surname"` // User's last name (required)
	Email   string `json:"email"`   // User's email address (required, must be unique)
	Address string `json:"address"` // User's physical address (optional)

	Password   string `json:"password"`    // User's password (required for 'local' provider)
	Provider   string `json:"provider"`    // Authentication provider ('local', 'google', etc.)
	ProviderID string `json:"provider_id"` // Provider-specific user identifier (for external providers)
}

// ResetPasswordRequest represents the request payload for initiating a password reset process.
// Used when users forget their password and need to reset it via email verification.
//
// Validation Requirements:
//   - Email: Must be a valid email format and exist in the system
//
// Business Rules:
//   - A password reset token is generated and sent to the provided email
//   - Rate limiting may apply to prevent abuse of the reset functionality
//   - Previous reset tokens for the same user may be invalidated
//   - Only works for users with 'local' authentication provider
//
// Security Notes:
//   - Reset tokens have a limited time window for validity
//   - Email verification is required before password can be changed
//   - Process is logged for security auditing purposes
type ResetPasswordRequest struct {
	Email string `json:"email"` // User's email address (required, must be unique)
}

// ChangePasswordRequest represents the request payload for changing a user's password.
// Used when authenticated users want to update their current password.
//
// Validation Requirements:
//   - Email: Must be a valid email format and exist in the system
//   - Password: Must meet security requirements and be different from current password
//
// Business Rules:
//   - User must be authenticated to perform this operation
//   - Only works for users with 'local' authentication provider
//   - Password will be hashed before storage
//   - Previous sessions may be invalidated for security
//
// Security Notes:
//   - New password is transmitted in plain text (ensure HTTPS is used)
//   - Password change events are logged for security auditing
//   - May trigger email notifications to the user about the password change
type ChangePasswordRequest struct {
	Email    string `json:"email"`    // User's email address (required for identification)
	Password string `json:"password"` // New password (required, will be hashed for storage)
}
