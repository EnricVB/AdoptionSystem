package dao

import (
	"backend/internal/db"
	m "backend/internal/models"
	"fmt"
)

func GetAllSpecies() ([]m.Species, error) {
	gormDB := db.ORMOpen()

	var species []m.Species
	result := gormDB.Find(&species)
	if result.Error != nil {
		return nil, fmt.Errorf("error al leer especies: %v", result.Error)
	}

	return species, nil
}

func GetSpeciesByID(id uint) (*m.Species, error) {
	gormDB := db.ORMOpen()

	var s m.Species
	result := gormDB.First(&s, id)
	if result.Error != nil {
		return nil, fmt.Errorf("error al leer especie con id %d: %v", id, result.Error)
	}

	return &s, nil
}

func CreateSpecies(s *m.Species) (*m.Species, error) {
	gormDB := db.ORMOpen()

	result := gormDB.Create(s)
	if result.Error != nil {
		return nil, fmt.Errorf("error al crear especie: %v", result.Error)
	}

	return s, nil
}

func DeleteSpeciesByID(id uint) error {
	gormDB := db.ORMOpen()

	result := gormDB.Delete(&m.Species{}, id)
	if result.Error != nil {
		return fmt.Errorf("error al eliminar especie con id %d: %v", id, result.Error)
	}

	return nil
}
