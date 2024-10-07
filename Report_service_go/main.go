package main

import (
	"Report_service_go/config"
	"Report_service_go/docs"
	"Report_service_go/middlewares"
	"Report_service_go/routes"
	"Report_service_go/utils"
	"net/http"

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

	reportRouter := routes.ReportRoutes(db)
	mux.Handle("/get_bot_conversation_history_table", reportRouter)
	mux.Handle("/get_bot_executive_summary", reportRouter)
	mux.Handle("/get_bot_conversation_topic_chart", reportRouter)
	mux.Handle("/get_bot_conversation", reportRouter)
	mux.Handle("/get_bot_internal_knowledge", reportRouter)
	mux.Handle("/update_bot_internal_knowledge", reportRouter)
	mux.Handle("/get_bot_internal_greeting", reportRouter)
	mux.Handle("/update_bot_internal_greeting", reportRouter)
	mux.Handle("/get_initial_data", reportRouter)

	loggedMux := middlewares.LogRequestMiddleware(mux)
	utils.Log(referenceId, "Server running on http://localhost:3000/")
	if err := http.ListenAndServe(":3000", loggedMux); err != nil {
		utils.Log(referenceId, "Failed to start server:", err)
	}

}
