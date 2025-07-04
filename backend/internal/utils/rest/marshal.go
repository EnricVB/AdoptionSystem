// Package response provides utilities for standardized REST API success responses.
// This package implements consistent response formatting for successful API operations.
package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HTTPMarshal represents a standardized success response structure for REST APIs.
// It wraps successful response data in a consistent format with status information.
//
// Fields:
//   - Code: HTTP status code (typically 200 for successful responses)
//   - Content: The actual response data (can be any JSON-serializable type)
//
// This structure ensures consistent success responses across all API endpoints,
// making it easier for clients to parse and handle responses uniformly.
type HTTPMarshal struct {
	Code    int `json:"code"`
	Content any `json:"content"`
}

// MarshalResponse sends a JSON success response with HTTP 200 status.
// This is the primary function for returning successful responses to clients.
// It wraps the provided data in a standardized HTTPMarshal structure.
//
// Parameters:
//   - c: Echo context for sending the HTTP response
//   - objects: The data to include in the response (any JSON-serializable type)
//
// Returns:
//   - error: Echo framework error, typically nil unless JSON marshaling fails
//
// Response Format:
//
//	{
//	  "code": 200,
//	  "content": <provided objects>
//	}
//
// Usage Examples:
//
//	return response.MarshalResponse(c, user)           // Single object
//	return response.MarshalResponse(c, users)          // Array of objects
//	return response.MarshalResponse(c, "OK")           // Simple string message
//	return response.MarshalResponse(c, map[string]any{"id": 1}) // Map/object
func MarshalResponse(c echo.Context, objects any) error {
	return c.JSON(http.StatusOK, HTTPMarshal{
		Code:    http.StatusOK,
		Content: objects,
	})
}
