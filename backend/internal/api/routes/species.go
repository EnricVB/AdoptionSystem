package api

import (
	"backend/internal/db/dao"
	m "backend/internal/models"
	response "backend/internal/utils/rest"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func RegisterSpeciesRoutes(e *echo.Echo) {
	e.GET("/api/species", handleListSpecies)
	e.GET("/api/species/:id", handleGetSpeciesByID)
	e.POST("/api/species", handleCreateSpecies)
	e.DELETE("/api/species/:id", handleDeleteSpecies)
}

func handleListSpecies(c echo.Context) error {
	species, err := dao.GetAllSpecies()
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error al obtener especies: %v", err))
	}
	return response.MarshalResponse(c, species)
}

func handleGetSpeciesByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	s, err := dao.GetSpeciesByID(uint(id))
	if err != nil {
		return response.ErrorResponse(c, http.StatusNotFound, fmt.Sprintf("especie no encontrada: %v", err))
	}
	return response.MarshalResponse(c, s)
}

func handleCreateSpecies(c echo.Context) error {
	var s m.Species
	if err := c.Bind(&s); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "datos de especie inválidos")
	}
	created, err := dao.CreateSpecies(&s)
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error al crear especie: %v", err))
	}
	return response.MarshalResponse(c, created)
}

func handleDeleteSpecies(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := dao.DeleteSpeciesByID(uint(id)); err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error al eliminar especie: %v", err))
	}
	return response.MarshalResponse(c, map[string]string{"status": "deleted"})
}
