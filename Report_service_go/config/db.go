//config/db.go
package config

import (
	"Report_service_go/utils"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func DBConnect(referenceId string) (*sql.DB, error) {
	// Load database connection variables from environment
	dbUser := os.Getenv("DB_USER")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")

	// Database connection string
	connUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	utils.Log(referenceId, "DB URL", connUrl)

	// Open the database connection
	db, err := sql.Open("postgres", connUrl)
	if err != nil {
		utils.Log(referenceId, "Failed to connect to the database:", err)
		return nil, err
	}

	// Test the database connection
	if err := db.Ping(); err != nil {
		utils.Log(referenceId, "Database connection failed:", err)
		return nil, err
	}

	utils.Log(referenceId, "Database connection successful!")
	return db, nil
}
