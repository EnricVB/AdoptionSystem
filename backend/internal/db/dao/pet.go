package dao

import (
	"backend/internal/db"
	m "backend/internal/models"
	"fmt"
	"time"
)

func GetAllPets() ([]m.SimplifiedPet, error) {
	gormDB := db.ORMOpen()

	var pets []m.SimplifiedPet
	result := gormDB.Preload("AdoptUser").Find(&pets)
	if result.Error != nil {
		return nil, fmt.Errorf("error al leer mascotas: %v", result.Error)
	}

	return pets, nil
}

func GetPetByID(id uint) (*m.Pet, error) {
	gormDB := db.ORMOpen()

	var pet m.Pet
	result := gormDB.Preload("AdoptUser").Where("id = ?", id).First(&pet)
	if result.Error != nil {
		return nil, fmt.Errorf("error al leer mascota con id %d: %v", id, result.Error)
	}

	return &pet, nil
}

func CreatePet(pet *m.Pet) (*m.Pet, error) {
	gormDB := db.ORMOpen()

	now := time.Now()
	pet.CrtDate = now
	pet.UptDate = now

	result := gormDB.Create(pet)
	if result.Error != nil {
		return nil, fmt.Errorf("error al crear mascota: %v", result.Error)
	}

	return pet, nil
}

func UpdatePet(pet *m.Pet) error {
	gormDB := db.ORMOpen()

	pet.UptDate = time.Now()

	result := gormDB.Model(&m.Pet{}).
		Where("id = ?", pet.ID).
		Select("*").
		Updates(pet)

	if result.Error != nil {
		return fmt.Errorf("error al actualizar mascota con id %d: %v", pet.ID, result.Error)
	}

	return nil
}

func DeletePetByID(id uint) error {
	gormDB := db.ORMOpen()

	result := gormDB.Delete(&m.Pet{}, id)
	if result.Error != nil {
		return fmt.Errorf("error al eliminar mascota con id %d: %v", id, result.Error)
	}

	return nil
}
