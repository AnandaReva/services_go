package config

import (
	"database/sql"
	"fmt"
	"os"
	"Auth_service_go/utils"

	_ "github.com/lib/pq"
)

// DBConnect initializes and returns a database connection
func DBConnect(referenceID string) (*sql.DB, error) {
	// Load database connection variables from environment
	dbUser := os.Getenv("DB_USER")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")

	// Database connection string
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	utils.Log(referenceID, "DB URL", connStr)

	// Open the database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		utils.Log(referenceID, "Failed to connect to the database:", err)
		return nil, err
	}

	// Test the database connection
	if err := db.Ping(); err != nil {
		utils.Log(referenceID, "Database connection failed:", err)
		return nil, err
	}

	utils.Log(referenceID, "Database connection successful!")
	return db, nil
}
