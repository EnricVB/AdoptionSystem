package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB      *sql.DB
	GORM_DB *gorm.DB
	once    sync.Once
)

/*
This file handles the connection to the database using GORM and the standard sql package.
It uses environment variables to configure the connection parameters.
*/

// buildDSN builds a MySQL DSN string using environment variables.
func buildDSN() string {
	user := "user"           //os.Getenv("DB_USER")
	password := "1234"       //os.Getenv("DB_PASSWORD")
	host := "127.0.0.1"      //os.Getenv("DB_HOST")
	port := "3306"           //os.Getenv("DB_PORT")
	dbname := "ADOPTION_SYS" //os.Getenv("DB_NAME")

	// charset=utf8mb4 and parseTime=True are standard and recommended for MySQL
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)
}

/*
gormConnect initializes a GORM connection to the database using the provided DSN.
It returns a pointer to the gorm.DB instance or an error if the connection fails.
*/
func gormConnect() (*gorm.DB, error) {
	var err error

	GORM_DB, err = gorm.Open(mysql.Open(buildDSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm.Open failed: %w", err)
	}

	return GORM_DB, nil
}

/*
RawConnect initializes a raw SQL connection to the database using the provided DSN.
It returns a pointer to the sql.DB instance or an error if the connection fails.
*/
func RawConnect() (*sql.DB, error) {
	var err error

	DB, err = sql.Open("mysql", buildDSN())
	if err != nil {
		return nil, fmt.Errorf("sql.Open failed: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return nil, fmt.Errorf("sql.Ping failed: %w", err)
	}

	return DB, nil
}

/*
ORMOpen initializes a GORM connection to the database, ensuring that it is only done once.
It returns a pointer to the gorm.DB instance.
If the connection fails, it logs a fatal error and exits the program.
*/
func ORMOpen() *gorm.DB {
	once.Do(func() {
		var err error
		GORM_DB, err = gormConnect()
		if err != nil {
			log.Fatalf("failed to connect to DataBase: %v", err)
		}
	})

	return GORM_DB
}
