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

	IsAdmin bool `json:"is_admin" gorm:"default:false;column:Is_Admin"` // Whether the user has administrative privileges

	CrtDate time.Time `json:"crt_date" gorm:"autoCreateTime"` // Record creation timestamp
	UptDate time.Time `json:"upt_date" gorm:"autoUpdateTime"` // Record last update timestamp
}

// User represents the standard user entity for API responses and general operations.
// This model excludes sensitive information like passwords and 2FA codes.
// It's safe to return to API clients and use in public-facing operations.
//
// Database Table: Users
type User struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name         string    `json:"name" gorm:"type:varchar(100);not null"`
	Surname      string    `json:"surname" gorm:"type:varchar(100);not null"`
	Email        string    `json:"email" gorm:"type:varchar(150);uniqueIndex;not null"`
	SessionID    string    `json:"session_id" gorm:"type:varchar(50);uniqueIndex;column:Session_ID"`
	Address      string    `json:"address" gorm:"type:varchar(255)"`
	Provider     string    `json:"provider" gorm:"default:'local';type:varchar(255);column:Provider"`
	ProviderID   string    `json:"provider_id" gorm:"type:varchar(255);column:Provider_ID"`
	Password     string    `json:"password" gorm:"type:varchar(255);not null"`
	ChangePass   bool      `json:"change_pass" gorm:"default:false;column:Change_Password"`
	FailedLogins uint      `json:"failed_logins" gorm:"default:0;column:Failed_Logins"`
	IsBlocked    bool      `json:"is_blocked" gorm:"default:false;column:Is_Blocked"`
	IsAdmin      bool      `json:"is_admin" gorm:"default:false;column:Is_Admin"`
	CrtDate      time.Time `json:"crt_date" gorm:"autoCreateTime"`
	UptDate      time.Time `json:"upt_date" gorm:"autoUpdateTime"`
}

// NonValidatedUser represents a user entity without session validation.
// This model is used for operations that don't require active session validation,
// such as user profile display or administrative operations.
//
// Database Table: Users
type NonValidatedUser struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name         string    `json:"name" gorm:"type:varchar(100);not null"`
	Surname      string    `json:"surname" gorm:"type:varchar(100);not null"`
	Email        string    `json:"email" gorm:"type:varchar(150);uniqueIndex;not null"`
	Address      string    `json:"address" gorm:"type:varchar(255)"`
	FailedLogins uint      `json:"failed_logins" gorm:"default:0;column:Failed_Logins"`
	Provider     string    `json:"provider" gorm:"default:'local';type:varchar(255);column:Provider"` // Authentication provider (local, google, etc.)
	IsBlocked    bool      `json:"is_blocked" gorm:"default:false;column:Is_Blocked"`
	IsAdmin      bool      `json:"is_admin" gorm:"default:false;column:Is_Admin"`
	CrtDate      time.Time `json:"crt_date" gorm:"autoCreateTime"`
	UptDate      time.Time `json:"upt_date" gorm:"autoUpdateTime"`
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
	IsAdmin bool   `json:"is_admin" gorm:"default:false;column:Is_Admin"`       // Whether the user has administrative privileges
}

// ========================================
// USER CONVERSION METHODS
// ========================================

// ToUser converts a FullUser to a standard User model.
// This method strips sensitive information like two-factor auth codes
// while preserving most user data for API responses.
func (fu FullUser) ToUser() User {
	return User{
		ID:           fu.ID,
		Name:         fu.Name,
		Surname:      fu.Surname,
		Email:        fu.Email,
		SessionID:    fu.SessionID,
		Address:      fu.Address,
		Provider:     fu.Provider,
		ProviderID:   fu.ProviderID,
		Password:     fu.Password,
		ChangePass:   fu.ChangePassword,
		FailedLogins: fu.FailedLogins,
		IsBlocked:    fu.IsBlocked,
		CrtDate:      fu.CrtDate,
		UptDate:      fu.UptDate,
		IsAdmin:      fu.IsAdmin,
	}
}

// ToNonValidatedUser converts a FullUser to a NonValidatedUser model.
// This method removes session information and password data.
func (fu FullUser) ToNonValidatedUser() NonValidatedUser {
	return NonValidatedUser{
		ID:           fu.ID,
		Name:         fu.Name,
		Surname:      fu.Surname,
		Email:        fu.Email,
		Address:      fu.Address,
		FailedLogins: fu.FailedLogins,
		Provider:     fu.Provider,
		IsBlocked:    fu.IsBlocked,
		IsAdmin:      fu.IsAdmin,
		CrtDate:      fu.CrtDate,
		UptDate:      fu.UptDate,
	}
}

// ToSimplifiedUser converts a FullUser to a SimplifiedUser model.
// This method keeps only essential user information for listings and references.
func (fu FullUser) ToSimplifiedUser() SimplifiedUser {
	return SimplifiedUser{
		ID:      fu.ID,
		Name:    fu.Name,
		Surname: fu.Surname,
		Email:   fu.Email,
		Address: fu.Address,
		IsAdmin: fu.IsAdmin,
	}
}

// ToNonValidatedUser converts a User to a NonValidatedUser model.
// This method removes session information while keeping other data.
func (u User) ToNonValidatedUser() NonValidatedUser {
	return NonValidatedUser{
		ID:           u.ID,
		Name:         u.Name,
		Surname:      u.Surname,
		Email:        u.Email,
		Address:      u.Address,
		FailedLogins: u.FailedLogins,
		Provider:     u.Provider,
		IsBlocked:    u.IsBlocked,
		IsAdmin:      u.IsAdmin,
		CrtDate:      u.CrtDate,
		UptDate:      u.UptDate,
	}
}

// ToSimplifiedUser converts a User to a SimplifiedUser model.
// This method keeps only essential user information for listings and references.
func (u User) ToSimplifiedUser() SimplifiedUser {
	return SimplifiedUser{
		ID:      u.ID,
		Name:    u.Name,
		Surname: u.Surname,
		Email:   u.Email,
		Address: u.Address,
		IsAdmin: u.IsAdmin,
	}
}

// ToSimplifiedUser converts a NonValidatedUser to a SimplifiedUser model.
// This method keeps only essential user information for listings and references.
func (nvu NonValidatedUser) ToSimplifiedUser() SimplifiedUser {
	return SimplifiedUser{
		ID:      nvu.ID,
		Name:    nvu.Name,
		Surname: nvu.Surname,
		Email:   nvu.Email,
		Address: nvu.Address,
		IsAdmin: nvu.IsAdmin,
	}
}

// ========================================
// BATCH CONVERSION METHODS
// ========================================

// ToUserSlice converts a slice of FullUser to a slice of User.
func ToUserSlice(fullUsers []FullUser) []User {
	users := make([]User, len(fullUsers))
	for i, fu := range fullUsers {
		users[i] = fu.ToUser()
	}
	return users
}

// ToNonValidatedUserSlice converts a slice of FullUser to a slice of NonValidatedUser.
func ToNonValidatedUserSlice(fullUsers []FullUser) []NonValidatedUser {
	users := make([]NonValidatedUser, len(fullUsers))
	for i, fu := range fullUsers {
		users[i] = fu.ToNonValidatedUser()
	}
	return users
}

// ToSimplifiedUserSlice converts a slice of FullUser to a slice of SimplifiedUser.
func ToSimplifiedUserSlice(fullUsers []FullUser) []SimplifiedUser {
	users := make([]SimplifiedUser, len(fullUsers))
	for i, fu := range fullUsers {
		users[i] = fu.ToSimplifiedUser()
	}
	return users
}

// ToSimplifiedUserSliceFromUser converts a slice of User to a slice of SimplifiedUser.
func ToSimplifiedUserSliceFromUser(users []User) []SimplifiedUser {
	simplified := make([]SimplifiedUser, len(users))
	for i, u := range users {
		simplified[i] = u.ToSimplifiedUser()
	}
	return simplified
}

// ToSimplifiedUserSliceFromNonValidated converts a slice of NonValidatedUser to a slice of SimplifiedUser.
func ToSimplifiedUserSliceFromNonValidated(users []NonValidatedUser) []SimplifiedUser {
	simplified := make([]SimplifiedUser, len(users))
	for i, u := range users {
		simplified[i] = u.ToSimplifiedUser()
	}
	return simplified
}
