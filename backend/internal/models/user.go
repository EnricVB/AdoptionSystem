// Package models contains data models for the adoption system.
// These models define the structure of database entities and their relationships.
package models

import "time"

// TableName returns the database table name for the User model.
// This method implements the GORM Tabler interface to specify custom table names.
func (User) TableName() string {
	return "Users"
}

// FullUser represents the complete user entity with all fields including sensitive data.
// This model is used internally for authentication and administrative operations.
// It includes password, provider information, and two-factor authentication data.
//
// Security Note: This model contains sensitive information and should be used
// carefully. Never return this model directly to API clients.
//
// Database Table: Users
type FullUser struct {
	ID            uint   `json:"id" gorm:"primaryKey;autoIncrement"`                               // Unique identifier for the user
	Name          string `json:"name" gorm:"type:varchar(100);not null"`                           // User's first name
	Surname       string `json:"surname" gorm:"type:varchar(100);not null"`                        // User's last name
	Email         string `json:"email" gorm:"type:varchar(150);uniqueIndex;not null"`              // User's email address (unique)
	SessionID     string `json:"session_id" gorm:"type:varchar(50);uniqueIndex;column:Session_ID"` // Current session identifier
	Address       string `json:"address" gorm:"type:varchar(255)"`                                 // User's physical address
	FailedLogins  uint   `json:"failed_logins" gorm:"default:0;column:Failed_Logins"`              // Count of failed login attempts
	IsBlocked     bool   `json:"is_blocked" gorm:"default:false;column:Is_Blocked"`                // Whether the user account is blocked
	TwoFactorAuth string `json:"two_factor_auth" gorm:"type:varchar(6);column:Two_Factor_Auth"`    // Two-factor authentication code

	Password   string `json:"password,omitempty" gorm:"type:varchar(255);column:Password"`       // Hashed password (omitted from JSON)
	Provider   string `json:"provider" gorm:"default:'local';type:varchar(255);column:Provider"` // Authentication provider (local, google, etc.)
	ProviderID string `json:"provider_id" gorm:"type:varchar(255);column:Provider_ID"`           // Provider-specific user ID

	ChangePassword bool `json:"change_password" gorm:"default:false;column:Change_Password"` // Flag indicating if user must change password on next login

	CrtDate time.Time `json:"crt_date" gorm:"autoCreateTime"` // Record creation timestamp
	UptDate time.Time `json:"upt_date" gorm:"autoUpdateTime"` // Record last update timestamp
}

// User represents the standard user entity for API responses and general operations.
// This model excludes sensitive information like passwords and 2FA codes.
// It's safe to return to API clients and use in public-facing operations.
//
// Database Table: Users
type User struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`                                // Unique identifier for the user
	Name         string    `json:"name" gorm:"type:varchar(100);not null"`                            // User's first name
	Surname      string    `json:"surname" gorm:"type:varchar(100);not null"`                         // User's last name
	Email        string    `json:"email" gorm:"type:varchar(150);uniqueIndex;not null"`               // User's email address (unique)
	SessionID    string    `json:"session_id" gorm:"type:varchar(50);uniqueIndex;column:Session_ID"`  // Current session identifier
	Address      string    `json:"address" gorm:"type:varchar(255)"`                                  // User's physical address
	FailedLogins uint      `json:"failed_logins" gorm:"default:0;column:Failed_Logins"`               // Count of failed login attempts
	IsBlocked    bool      `json:"is_blocked" gorm:"default:false;column:Is_Blocked"`                 // Whether the user account is blocked
	Password     string    `json:"password,omitempty" gorm:"type:varchar(255);column:Password"`       // Hashed password (omitted from JSON)
	Provider     string    `json:"provider" gorm:"default:'local';type:varchar(255);column:Provider"` // Authentication provider (local, google, etc.)
	ProviderID   string    `json:"provider_id" gorm:"type:varchar(255);column:Provider_ID"`           // Provider-specific user ID
	CrtDate      time.Time `json:"crt_date" gorm:"autoCreateTime"`                                    // Record creation timestamp
	UptDate      time.Time `json:"upt_date" gorm:"autoUpdateTime"`                                    // Record last update timestamp
}

// NonValidatedUser represents a user entity without session validation.
// This model is used for operations that don't require active session validation,
// such as user profile display or administrative operations.
//
// Database Table: Users
type NonValidatedUser struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`                                // Unique identifier for the user
	Name         string    `json:"name" gorm:"type:varchar(100);not null"`                            // User's first name
	Surname      string    `json:"surname" gorm:"type:varchar(100);not null"`                         // User's last name
	Email        string    `json:"email" gorm:"type:varchar(150);uniqueIndex;not null"`               // User's email address (unique)
	Address      string    `json:"address" gorm:"type:varchar(255)"`                                  // User's physical address
	FailedLogins uint      `json:"failed_logins" gorm:"default:0;column:Failed_Logins"`               // Count of failed login attempts
	IsBlocked    bool      `json:"is_blocked" gorm:"default:false;column:Is_Blocked"`                 // Whether the user account is blocked
	Provider     string    `json:"provider" gorm:"default:'local';type:varchar(255);column:Provider"` // Authentication provider (local, google, etc.)
	CrtDate      time.Time `json:"crt_date" gorm:"autoCreateTime"`                                    // Record creation timestamp
	UptDate      time.Time `json:"upt_date" gorm:"autoUpdateTime"`                                    // Record last update timestamp
}

// SimplifiedUser represents a minimal user entity with only essential information.
// This model is used for operations that require only basic user data,
// such as user lists, search results, or reference lookups.
//
// Database Table: Users
type SimplifiedUser struct {
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`                  // Unique identifier for the user
	Name    string `json:"name" gorm:"type:varchar(100);not null"`              // User's first name
	Surname string `json:"surname" gorm:"type:varchar(100);not null"`           // User's last name
	Email   string `json:"email" gorm:"type:varchar(150);uniqueIndex;not null"` // User's email address (unique)
	Address string `json:"address" gorm:"type:varchar(255)"`                    // User's physical address
}
