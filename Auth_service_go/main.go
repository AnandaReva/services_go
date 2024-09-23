package main

import (
	"net/http"

	"Auth_service_go/config"
	"Auth_service_go/middlewares"
	"Auth_service_go/routes"
	"Auth_service_go/utils"

	"github.com/lpernett/godotenv"
)

func main() {
	referenceID := "MAIN"

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		utils.Log(referenceID, "Error loading .env file", err)
		return
	}

	// Open the database connection using the separated utility function
	db, err := config.DBConnect(referenceID)
	if err != nil {
		return
	}
	defer db.Close()

	// Setup the HTTP server with a handler
	mux := http.NewServeMux()

	// Route for home page
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.Log(referenceID, "Received request for /")
		w.Write([]byte("Request received"))
	}))

	// Add routes from authRoutes
	authRouter := routes.NewRouter(db)
	mux.Handle("/login", authRouter)
	mux.Handle("/verify-challenge", authRouter)

	// Add middleware to the router
	loggedMux := middlewares.LogRequestMiddleware(mux)

	// Start the server
	utils.Log(referenceID, "Server running on :3000")
	if err := http.ListenAndServe(":3000", loggedMux); err != nil {
		utils.Log(referenceID, "Failed to start server:", err)
	}
}
