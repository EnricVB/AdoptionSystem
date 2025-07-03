package handlers

import (
	m "backend/internal/models"
	s "backend/internal/services/backend_calls"
	response "backend/internal/utils/rest"
	"net/http"
)

func HandleListPets() (*[]m.SimplifiedPet, response.HTTPError) {
	pets, err := s.ListAllPets()
	if err != nil {
		return nil, response.Error(http.StatusInternalServerError, err.Error())
	}

	return pets, response.EmptyError
}

func HandleGetPetByID(id uint) (*m.Pet, response.HTTPError) {
	if id <= 0 {
		return nil, response.Error(http.StatusBadRequest, "ID de mascota no válido")
	}

	pet, err := s.GetPetByID(id)
	if err != nil {
		return nil, response.Error(http.StatusNotFound, err.Error())
	}

	return pet, response.EmptyError
}

func HandleCreatePet(pet *m.Pet) (*m.Pet, response.HTTPError) {
	if pet.Name == "" || pet.Species == "" {
		return nil, response.Error(http.StatusBadRequest, "nombre y especie de mascota son obligatorios")
	}

	err := s.CreatePet(pet)
	if err != nil {
		return nil, response.Error(http.StatusInternalServerError, err.Error())
	}

	return pet, response.EmptyError
}

func HandleUpdatePet(pet *m.Pet) (*m.Pet, response.HTTPError) {
	if pet.ID <= 0 {
		return nil, response.Error(http.StatusBadRequest, "ID de mascota no válido")
	}

	if pet.Name == "" || pet.Species == "" {
		return nil, response.Error(http.StatusBadRequest, "nombre y especie de mascota son obligatorios")
	}

	err := s.UpdatePet(pet)
	if err != nil {
		return nil, response.Error(http.StatusInternalServerError, err.Error())
	}

	return pet, response.EmptyError
}

func HandleDeletePet(id uint) response.HTTPError {
	if id <= 0 {
		return response.Error(http.StatusBadRequest, "ID de mascota no válido")
	}

	err := s.DeletePet(id)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.EmptyError
}
