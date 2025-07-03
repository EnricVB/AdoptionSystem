package api

import (
	"backend/internal/api/handlers"
	m "backend/internal/models"
	response "backend/internal/utils/rest"
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
	species, httpErr := handlers.HandleListSpecies()
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}
	return response.MarshalResponse(c, species)
}

func handleGetSpeciesByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "ID de especie inválido")
	}

	species, httpErr := handlers.HandleGetSpeciesByID(uint(id))
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}
	return response.MarshalResponse(c, species)
}

func handleCreateSpecies(c echo.Context) error {
	var species m.Species
	if err := c.Bind(&species); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "datos de especie inválidos")
	}

	created, httpErr := handlers.HandleCreateSpecies(&species)
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}
	return response.MarshalResponse(c, created)
}

func handleDeleteSpecies(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "ID de especie inválido")
	}

	httpErr := handlers.HandleDeleteSpecies(uint(id))
	if httpErr.Code != 0 {
		return response.ConvertToErrorResponse(c, httpErr)
	}
	return response.MarshalResponse(c, map[string]string{"status": "deleted"})
}
