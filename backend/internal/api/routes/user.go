package api

import (
	"backend/internal/api/handlers"
	r_models "backend/internal/api/routes/models"
	"backend/internal/db/dao"
	m "backend/internal/models"
	response "backend/internal/utils/rest"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func RegisterUserRoutes(e *echo.Echo) {
	e.GET("/api/users", handleListUsers)
	e.GET("/api/users/:id", handleGetUserByID)
	e.POST("/api/login", handleLoginUser)
	e.POST("/api/users", handleCreateUser)
	e.PUT("/api/users/:id", handleUpdateUser)
	e.DELETE("/api/users/:id", handleDeleteUser)
}

func handleLoginUser(c echo.Context) error {
	var req r_models.LoginRequest

	if err := c.Bind(&req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	user, err := handlers.HandleLogin(req)

	if err != nil {
		return response.ErrorResponse(c, http.StatusUnauthorized, fmt.Sprintf("error al inciar sesion: %v", err))
	}

	return response.MarshalResponse(c, user)
}

func handleListUsers(c echo.Context) error {
	users, err := dao.GetAllUsers()
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error al leer usuarios: %v", err))
	}
	return response.MarshalResponse(c, users)
}

func handleGetUserByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := dao.GetUserByID(uint(id))
	if err != nil {
		return response.ErrorResponse(c, http.StatusNotFound, err.Error())
	}
	return response.MarshalResponse(c, user)
}

func handleCreateUser(c echo.Context) error {
	var user m.User
	if err := c.Bind(&user); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "datos inválidos")
	}
	created, err := dao.CreateUser(&user)
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return response.MarshalResponse(c, created)
}

func handleUpdateUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var user m.User
	if err := c.Bind(&user); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "datos inválidos")
	}
	user.ID = uint(id)
	updated, err := dao.UpdateUser(&user)
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return response.MarshalResponse(c, updated)
}

func handleDeleteUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	deleted, err := dao.DeleteUserByID(uint(id))
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return response.MarshalResponse(c, deleted)
}
