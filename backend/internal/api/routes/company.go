// Package api implements HTTP route handlers and endpoint registration for company information management.
// This layer is responsible for:
// - HTTP endpoint registration and routing for company operations
// - Request binding and basic input validation
// - Calling appropriate handler functions for company information
// - HTTP response formatting and status code management
// - RESTful API design compliance for company resources
package api

import (
	"backend/internal/api/handlers"
	response "backend/internal/utils/rest"

	"github.com/labstack/echo/v4"
)

// ========================================
// ROUTE REGISTRATION
// ========================================

// RegisterCompanyRoutes registers all company-related HTTP endpoints with the Echo router.
// Implements standard RESTful API design patterns for company resource management.
//
// Endpoint Organization:
// - GET /api/company: Get company information
//
// Parameters:
//   - e: Echo router instance for endpoint registration
func RegisterCompanyRoutes(e *echo.Echo) {
	e.GET("/api/company", handleGetCompanyInfo)
}

// ========================================
// COMPANY INFORMATION ROUTE HANDLERS
// ========================================

// handleGetCompanyInfo processes requests to retrieve company information.
// Returns the primary company data for system configuration and display.
//
// HTTP Method: GET
// Endpoint: /api/company
//
// Response:
//   - Success: Company information with all details
//   - Error: HTTP error with appropriate status code
func handleGetCompanyInfo(c echo.Context) error {
	// Delegate company information retrieval to handler layer
	company, httpErr := handlers.HandleGetCompanyInfo()
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}

	return response.MarshalResponse(c, company)
}
