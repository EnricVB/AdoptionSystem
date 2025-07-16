// Package handlers implements HTTP request handlers for company information management.
// This layer is responsible for:
// - HTTP request/response handling and validation for company data
// - Input sanitization and basic validation
// - Calling appropriate service layer functions
// - Converting service errors to HTTP responses
// - Ensuring consistent API response formatting
package handlers

import (
	m "backend/internal/models"
	s "backend/internal/services/backend_calls"
	response "backend/internal/utils/rest"
	"net/http"
)

// ========================================
// COMPANY INFORMATION HANDLERS
// ========================================

// HandleGetCompanyInfo processes requests to retrieve company information.
// Returns the primary company data for system configuration and display.
//
// Business Logic:
// - Retrieves the first (and should be only) company record
// - Used for system branding and contact information display
// - Provides company details for frontend configuration
//
// Returns:
//   - *m.Company: Company information with all details
//   - response.HTTPError: HTTP error or EmptyError on success
func HandleGetCompanyInfo() (*m.Company, response.HTTPError) {
	// Delegate company information retrieval to service layer
	company, err := s.GetCompanyInfo()
	if err != nil {
		return nil, response.Error(http.StatusInternalServerError, err.Error())
	}

	return company, response.EmptyError
}
