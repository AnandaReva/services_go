// routes/reportService.go
package routes

import (
	"Report_service_go/controllers"
	"database/sql"
	"net/http"
)

func ReportRoutes(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	// 1. Route for handling getBotConversationHistoryTable
	mux.HandleFunc("/get_bot_conversation_history_table", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			// Pass db, w, and r to HandleLoginRequest

			controllers.GetBotConversationHistoryTable(db, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// 2. Route for handling getBotExecutiveSummary
	mux.HandleFunc("/get_bot_executive_summary", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controllers.GetBotExecutiveSummary(db, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// 3. Route for handling getBotConversationTopicChart
	mux.HandleFunc("/get_bot_conversation_topic_chart", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controllers.GetBotConversationTopicChart(db, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	
	/* 	// Route for handling getBotConversation
		mux.HandleFunc("/get_bot_conversation", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				controllers.(db, w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		}) */

		// 5. Route for handling getBotInternalKnowledge
		mux.HandleFunc("/get_bot_internal_knowledge", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				controllers.GetBotInternalKnowledge(db, w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		})

		// 6. Route for handling updateBotInternalKnowledge
		mux.HandleFunc("/update_bot_internal_knowledge", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				controllers.UpdateBotInternalKnowledge(db, w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		})

		// 7. Route for handling getBotInternalGreeting
		mux.HandleFunc("/get_bot_internal_greeting", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				controllers.GetBotInternalGreeting(db, w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		})

		// 8. Route for handling updateBotInternalGreeting
		mux.HandleFunc("/update_bot_internal_greeting", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				controllers.UpdateBotInternalGreeting(db, w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		})

		// 9. Route for handling getInitialData
		mux.HandleFunc("/get_initial_data", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				controllers.GetInitialData(db, w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		})
	return mux
}
