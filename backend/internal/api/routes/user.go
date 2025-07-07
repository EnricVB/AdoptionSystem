package api

import (
	"backend/internal/api/handlers"
	r_models "backend/internal/api/routes/models"
	m "backend/internal/models"
	response "backend/internal/utils/rest"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func RegisterUserRoutes(e *echo.Echo) {
	// User CRUD operations
	e.GET("/api/users", handleListUsers)
	e.GET("/api/users/:id", handleGetUserByID)
	e.POST("/api/register", handleCreateUser)
	e.PUT("/api/users/:id", handleUpdateUser)
	e.DELETE("/api/users/:id", handleDeleteUser)

	// Authentication endpoints
	e.POST("/api/auth/login", handleLoginUser)
	e.POST("/api/auth/verify-2fa", handle2FAAuth)
	e.POST("/api/auth/refresh-token", handleRefresh2FAToken)
}

func handleLoginUser(c echo.Context) error {
	var req r_models.LoginRequest

	if err := c.Bind(&req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	user, err := handlers.HandleLogin(req)

	if err != response.EmptyError {
		return response.ConvertToErrorResponse(c, err)
	}

	return response.MarshalResponse(c, user)
}

func handle2FAAuth(c echo.Context) error {
	var req r_models.TwoFactorRequest

	if err := c.Bind(&req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	_, err := handlers.Handle2FAAuth(req)

	if err != response.EmptyError {
		return response.ConvertToErrorResponse(c, err)
	}

	return response.MarshalResponse(c, "OK")
}

func handleRefresh2FAToken(c echo.Context) error {
	var req r_models.RefreshTokenRequest

	if err := c.Bind(&req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	_, err := handlers.HandleRefresh2FAToken(req)

	if err != response.EmptyError {
		return response.ConvertToErrorResponse(c, err)
	}

	return response.MarshalResponse(c, "OK")
}

func handleListUsers(c echo.Context) error {
	users, httpErr := handlers.HandleListUsers()
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}
	return response.MarshalResponse(c, users)
}

func handleGetUserByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "ID de usuario inválido")
	}

	user, httpErr := handlers.HandleGetUserByID(uint(id))
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}
	return response.MarshalResponse(c, user)
}

func handleCreateUser(c echo.Context) error {
	var user m.User
	if err := c.Bind(&user); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "datos inválidos")
	}

	httpErr := handlers.HandleCreateUser(&user)
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}
	return response.MarshalResponse(c, user)
}

func handleUpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "ID de usuario inválido")
	}

	var user m.User
	if err := c.Bind(&user); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "datos inválidos")
	}
	user.ID = uint(id)

	httpErr := handlers.HandleUpdateUser(&user)
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}
	return response.MarshalResponse(c, user)
}

func handleDeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "ID de usuario inválido")
	}

	deleted, httpErr := handlers.HandleDeleteUser(uint(id))
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}
	return response.MarshalResponse(c, deleted)
}
