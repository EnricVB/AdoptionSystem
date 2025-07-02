package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HTTPMarshal struct {
	Code    int `json:"code"`
	Content any `json:"content"`
}

/*
MarshalResponse sends a JSON response with the provided objects and a status code of 200.
It is used to standardize successful responses in the API.
*/
func MarshalResponse(c echo.Context, objects any) error {
	return c.JSON(http.StatusOK, HTTPMarshal{
		Code:    http.StatusOK,
		Content: objects,
	})
}
