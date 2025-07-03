package services

import (
	"backend/internal/db/dao"
	m "backend/internal/models"
	"fmt"
)

func ListAllPets() (*[]m.SimplifiedPet, error) {
	pets, err := dao.GetAllPets()
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("error al obtener mascotas: %v", err))
	}
	return pets, nil
}

func GetPetByID(id uint) (*m.Pet, error) {
	pet, err := dao.GetPetByID(id)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("mascota no encontrada: %v", err))
	}

	return pet, nil
}

func CreatePet(pet *m.Pet) error {
	created, err := dao.CreatePet(pet)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("error al crear mascota: %v", err))
	}

	if created == nil {
		return fmt.Errorf("mascota no creada")
	} else {
		pet = created
	}

	return nil
}

func UpdatePet(pet *m.Pet) error {
	err := dao.UpdatePet(pet)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("error al actualizar mascota: %v", err))
	}

	return nil
}

func DeletePet(id uint) error {
	if err := dao.DeletePetByID(id); err != nil {
		return fmt.Errorf(fmt.Sprintf("error al eliminar mascota: %v", err))
	}

	return nil
}
