package response

import "github.com/labstack/echo/v4"

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

/*
ErrorResponse sends a JSON response with an error message and status code.
It is used to standardize error responses in the API.
*/
func ErrorResponse(c echo.Context, status int, message string) error {
	return c.JSON(status, HTTPError{
		Code:    status,
		Message: message,
	})
}
