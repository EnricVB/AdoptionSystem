// Package response provides utilities for standardized REST API error handling and responses.
// This package implements a consistent error response format across the entire application.
package response

import "github.com/labstack/echo/v4"

// HTTPError represents a standardized error response structure for REST APIs.
// It contains both an HTTP status code and a human-readable error message.
//
// Fields:
//   - Code: HTTP status code (e.g., 400, 404, 500)
//   - Message: Human-readable error message for client consumption
//
// This structure ensures consistent error responses across all API endpoints.
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse sends a JSON error response with the specified status code and message.
// This is the primary function for returning error responses to clients.
//
// Parameters:
//   - c: Echo context for sending the HTTP response
//   - code: HTTP status code (e.g., 400, 404, 500)
//   - message: Human-readable error message
//
// Returns:
//   - error: Echo framework error, typically nil unless JSON marshaling fails
//
// Usage:
//
//	return response.ErrorResponse(c, http.StatusBadRequest, "Invalid input data")
func ErrorResponse(c echo.Context, code int, message string) error {
	return c.JSON(code, HTTPError{
		Code:    code,
		Message: message,
	})
}

// Error creates an HTTPError instance with the specified code and message.
// This is a utility function for creating error structures that can be passed
// between layers of the application before being converted to HTTP responses.
//
// Parameters:
//   - code: HTTP status code
//   - message: Human-readable error message
//
// Returns:
//   - HTTPError: Structured error object
//
// Usage:
//
//	return response.Error(http.StatusNotFound, "User not found")
func Error(code int, message string) HTTPError {
	return HTTPError{
		Code:    code,
		Message: message,
	}
}

// ConvertToErrorResponse converts an HTTPError instance to a JSON HTTP response.
// This function is used throughout the application to convert internal error
// structures into standardized HTTP responses.
//
// Parameters:
//   - c: Echo context for sending the HTTP response
//   - err: HTTPError instance containing code and message
//
// Returns:
//   - error: Echo framework error, typically nil unless JSON marshaling fails
//
// Usage:
//
//	if httpErr.Code != 0 {
//	    return response.ConvertToErrorResponse(c, httpErr)
//	}
func ConvertToErrorResponse(c echo.Context, err HTTPError) error {
	return c.JSON(err.Code, HTTPError{
		Code:    err.Code,
		Message: err.Message,
	})
}

// EmptyError represents an empty/uninitialized error state.
// Used as a sentinel value to indicate no error has occurred.
// An HTTPError with Code 0 is considered "empty" or "no error".
var EmptyError = HTTPError{}
