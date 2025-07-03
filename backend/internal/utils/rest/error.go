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
func ErrorResponse(c echo.Context, code int, message string) error {
	return c.JSON(code, HTTPError{
		Code:    code,
		Message: message,
	})
}

/*
Error creates an HTTPError instance with the given code and message.
This is a utility function to create error responses.
*/
func Error(code int, message string) HTTPError {
	return HTTPError{
		Code:    code,
		Message: message,
	}
}

/*
ConvertToErrorResponse converts an HTTPError to a JSON response.
This function is used to send error responses in a consistent format.
It takes an echo.Context and an HTTPError, and returns an error that can be used in
the Echo framework.
*/
func ConvertToErrorResponse(c echo.Context, err HTTPError) error {
	return c.JSON(err.Code, HTTPError{
		Code:    err.Code,
		Message: err.Message,
	})
}

// EmptyError is a variable for an empty error message.
var EmptyError = HTTPError{}
