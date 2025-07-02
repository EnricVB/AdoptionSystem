package main

import (
	api "backend/internal/api/routes"
	"backend/internal/db"
	"database/sql"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"
)

/*
Main entry point for the application.
This function initializes the database connection and sets up the CORS middleware for the Echo web framework.
It also registers user routes defined in the API package and starts the Echo server on port 8080.
*/
func main() {
	defer setupDatabase().Close()
	setupCORS()
}

/*
setupCORS initializes the Echo web framework, registers user routes,
and configures CORS middleware to allow requests from a specific origin.
*/
func setupCORS() {
	e := echo.New()
	api.RegisterUserRoutes(e)
	api.RegisterPetRoutes(e)
	api.RegisterSpeciesRoutes(e)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:4200"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	e.Start(":4200")
}

/*
setupDatabase initializes the database connection using GORM and returns the underlying sql.DB instance.
It logs a fatal error if the connection cannot be established.
*/
func setupDatabase() *sql.DB {
	gormDB := db.ORMOpen()
	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("could not get unlerying sql.DB: %v", err)
	}

	return sqlDB
}
