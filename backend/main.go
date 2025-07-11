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

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			w := c.Response().Writer
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			return next(c)
		}
	})

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:4200"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	e.Start(":8080")
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
