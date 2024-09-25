// routes/authRoutes.go
package routes

import (
	"Auth_service_go/controllers"
	"database/sql"
	"net/http"
)



func AuthRoutes(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	// Route for handling login
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			// Pass db, w, and r to HandleLoginRequest
			controllers.HandleLoginRequest(db, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/verify-challenge", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controllers.HandleChallengeResponseVerification(db, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	return mux
}
