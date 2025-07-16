// Package services provides business logic services for company information management.
// This layer sits between handlers and DAOs, implementing the core business rules
// and orchestrating database operations for company-related functionality.
package services

import (
	"backend/internal/db/dao"
	"backend/internal/models"
	"fmt"
)

// ========================================
// COMPANY INFORMATION SERVICES
// ========================================

// GetCompanyInfo retrieves the company information from the database.
// Returns the primary company data for system configuration and display.
//
// Business Logic:
// - Retrieves the first (and should be only) company record
// - Used for system branding and contact information
// - Provides company details for email templates and UI display
//
// Returns:
//   - *models.Company: Company information with all details
//   - error: Database error or nil on success
func GetCompanyInfo() (*models.Company, error) {
	// Retrieve company information from database
	company, err := dao.GetCompanyInfo()
	if err != nil {
		return nil, fmt.Errorf("error al obtener información de la empresa: %v", err)
	}

	return company, nil
}
