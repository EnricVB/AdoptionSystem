package api

import (
	"backend/internal/api/handlers"
	m "backend/internal/models"
	response "backend/internal/utils/rest"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func RegisterPetRoutes(e *echo.Echo) {
	e.GET("/api/pets", handleListPets)
	e.GET("/api/pets/:id", handleGetPetByID)
	e.POST("/api/pets", handleCreatePet)
	e.PUT("/api/pets/:id", handleUpdatePet)
	e.DELETE("/api/pets/:id", handleDeletePet)
}

func handleListPets(c echo.Context) error {
	pets, httpErr := handlers.HandleListPets()
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}
	return response.MarshalResponse(c, pets)
}

func handleGetPetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "ID de mascota inválido")
	}

	pet, httpErr := handlers.HandleGetPetByID(uint(id))
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}
	return response.MarshalResponse(c, pet)
}

func handleCreatePet(c echo.Context) error {
	var pet m.Pet
	if err := c.Bind(&pet); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "datos de mascota inválidos")
	}

	created, httpErr := handlers.HandleCreatePet(&pet)
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}
	return response.MarshalResponse(c, created)
}

func handleUpdatePet(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "ID de mascota inválido")
	}

	var pet m.Pet
	if err := c.Bind(&pet); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "datos de mascota inválidos")
	}
	pet.ID = uint(id)

	updated, httpErr := handlers.HandleUpdatePet(&pet)
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}
	return response.MarshalResponse(c, updated)
}

func handleDeletePet(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "ID de mascota inválido")
	}

	httpErr := handlers.HandleDeletePet(uint(id))
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}

	return response.MarshalResponse(c, map[string]string{"status": "deleted"})
}
