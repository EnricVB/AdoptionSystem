package handlers

import (
	m "backend/internal/models"
	s "backend/internal/services/backend_calls"
	response "backend/internal/utils/rest"
	"net/http"
)

func HandleListSpecies() ([]m.Species, response.HTTPError) {
	species, err := s.ListAllSpecies()
	if err != nil {
		return nil, response.Error(http.StatusInternalServerError, err.Error())
	}

	return species, response.HTTPError{}
}

func HandleGetSpeciesByID(id uint) (*m.Species, response.HTTPError) {
	if id <= 0 {
		return nil, response.Error(http.StatusBadRequest, "ID de especie no válido")
	}

	species, err := s.GetSpeciesByID(id)
	if err != nil {
		return nil, response.Error(http.StatusNotFound, err.Error())
	}

	return species, response.HTTPError{}
}

func HandleCreateSpecies(species *m.Species) (*m.Species, response.HTTPError) {
	if species.Name == "" {
		return nil, response.Error(http.StatusBadRequest, "nombre de especie es obligatorio")
	}

	err := s.CreateSpecies(species)
	if err != nil {
		return nil, response.Error(http.StatusInternalServerError, err.Error())
	}

	return species, response.HTTPError{}
}

func HandleDeleteSpecies(id uint) response.HTTPError {
	if id <= 0 {
		return response.Error(http.StatusBadRequest, "ID de especie no válido")
	}

	err := s.DeleteSpecies(id)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.HTTPError{}
}
