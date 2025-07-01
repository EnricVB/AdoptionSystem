package api

import (
	"backend/internal/db/dao"
	m "backend/internal/models"
	response "backend/internal/utils"
	"fmt"
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
	pets, err := dao.GetAllPets()
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error al obtener mascotas: %v", err))
	}
	return response.MarshalResponse(c, pets)
}

func handleGetPetByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	pet, err := dao.GetPetByID(uint(id))
	if err != nil {
		return response.ErrorResponse(c, http.StatusNotFound, fmt.Sprintf("mascota no encontrada: %v", err))
	}
	return response.MarshalResponse(c, pet)
}

func handleCreatePet(c echo.Context) error {
	var pet m.Pet
	if err := c.Bind(&pet); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "datos de mascota inválidos")
	}
	created, err := dao.CreatePet(&pet)
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error al crear mascota: %v", err))
	}
	return response.MarshalResponse(c, created)
}

func handleUpdatePet(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var pet m.Pet
	if err := c.Bind(&pet); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "datos de mascota inválidos")
	}
	pet.ID = uint(id)
	updated, err := dao.UpdatePet(&pet)
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error al actualizar mascota: %v", err))
	}
	return response.MarshalResponse(c, updated)
}

func handleDeletePet(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := dao.DeletePetByID(uint(id)); err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error al eliminar mascota: %v", err))
	}
	return response.MarshalResponse(c, map[string]string{"status": "deleted"})
}
