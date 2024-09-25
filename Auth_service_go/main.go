package main

import (
	"net/http"

	"Auth_service_go/config"
	"Auth_service_go/docs"
	"Auth_service_go/middlewares"
	"Auth_service_go/routes"
	"Auth_service_go/utils"

	"github.com/lpernett/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Swagger docs

	docs.SwaggerInfo.Title = "API Documentation"
	docs.SwaggerInfo.Description = "Documentation for Auth service endpoints"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:3000"
	docs.SwaggerInfo.BasePath = "/"

	referenceId := "MAIN"
	if err := godotenv.Load(); err != nil {
		utils.Log(referenceId, "Error loading .env file", err)
		return
	}

	db, err := config.DBConnect(referenceId)
	if err != nil {
		return
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.Log(referenceId, "Received request for /")
		w.Write([]byte("Request received"))
	}))

	// Swagger UI
	mux.Handle("/docs/*", httpSwagger.WrapHandler)

	authRouter := routes.AuthRoutes(db)
	mux.Handle("/login", authRouter)
	mux.Handle("/verify-challenge", authRouter)

	loggedMux := middlewares.LogRequestMiddleware(mux)
	utils.Log(referenceId, "Server running on http://localhost:3000/")
	if err := http.ListenAndServe(":3000", loggedMux); err != nil {
		utils.Log(referenceId, "Failed to start server:", err)
	}
}
