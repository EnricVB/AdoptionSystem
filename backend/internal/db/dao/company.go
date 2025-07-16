package dao

import (
	"backend/internal/db"
	"backend/internal/models"
	"fmt"
)

// GetCompanyInfo retrieves the first company record from the database.
// Since there should only be one company record in the system,
// this function returns the primary company information.
//
// Returns:
//   - *models.Company: The company information record
//   - error: Database error or nil on success
func GetCompanyInfo() (*models.Company, error) {
	var company models.Company

	// Open database connection
	gormDB := db.ORMOpen()

	// Query the first company record
	result := gormDB.First(&company)
	if result.Error != nil {
		return nil, fmt.Errorf("error retrieving company information: %v", result.Error)
	}

	return &company, nil
}
