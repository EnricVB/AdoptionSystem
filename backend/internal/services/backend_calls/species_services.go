package services

import (
	"backend/internal/db/dao"
	m "backend/internal/models"
	"fmt"
)

func ListAllSpecies() ([]m.Species, error) {
	species, err := dao.GetAllSpecies()
	if err != nil {
		return nil, fmt.Errorf("error al obtener especies: %v", err)
	}

	return species, nil
}

func GetSpeciesByID(id uint) (*m.Species, error) {
	species, err := dao.GetSpeciesByID(uint(id))
	if err != nil {
		return nil, fmt.Errorf("especie no encontrada: %v", err)
	}

	return species, nil
}

func CreateSpecies(species *m.Species) error {
	err := dao.CreateSpecies(species)
	if err != nil {
		return fmt.Errorf("error al crear especie: %v", err)
	}

	return nil
}

func DeleteSpecies(id uint) error {
	if err := dao.DeleteSpeciesByID(uint(id)); err != nil {
		return fmt.Errorf("error al eliminar especie: %v", err)
	}

	return nil
}
