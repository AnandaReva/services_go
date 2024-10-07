package controllers

import (
	"Report_service_go/utils"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// @Model requestBody3
type requestBody3 struct {
	Data         Data3  `json:"data" binding:"required" extensions:"x-order=1"`        // x-order=1 untuk Data
	FromDate     int    `json:"from_date" binding:"required" extensions:"x-order=3"`   // x-order=3 untuk FromDate
	ToDate       int    `json:"to_date" binding:"required" extensions:"x-order=4"`     // x-order=4 untuk ToDate
	SearchFilter string `json:"search_filter" binding:"max=50" extensions:"x-order=5"` // Max 50 karakter, x-order=5 untuk SearchFilter
	DateMode     int    `json:"date_mode" binding:"required" extensions:"x-order=6"`   // x-order=6 untuk DateMode
}

// @Model Data
type Data3 struct {
	BotID string `json:"bot_id" binding:"required" extensions:"x-order=1"` // x-order=5 untuk BotID
}

// GetBotConversationTopicChart retrieves the bot conversation topic chart based on the provided parameters.
//
// @Summary      3. Retrieve bot conversation topic chart
// @Description  This endpoint fetches the conversation topic chart of a bot from a backend service
//
//	based on the provided bot ID, date range, and other filters.
//
// @Tags         Bot Conversation
// @Accept       json
// @Produce      json
// @Param        requestBody3   body  requestBody3 true  "Request Body"
// @Param        ecwx-session-id header  string true  "Session ID for authentication"
// @Param        ecwx-hash       header  string true  "Hash for request validation"
// @Success      200  {object}    map[string]interface{}  "Successful response"
// @Failure      400  {object}    map[string]interface{}  "Invalid request"
// @Failure      401  {object}    map[string]interface{}  "Unauthenticated"
// @Failure      403  {object}    map[string]interface{}  "Forbidden"
// @Failure      500  {object}    map[string]interface{}  "Internal server error"
// @Router       /get_bot_conversation_topic_chart [post]
// @Extensions x-order=3
func GetBotConversationTopicChart(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	referenceId := utils.GlobalVarInstance.GetReferenceId()
	if referenceId == "" {
		referenceId = "undefined"
	}

	utils.Log(referenceId, "\nExecuting method: GetBotConversationTopicChart")

	// Ambil endpoint backend dari environment variables
	realBackendURL := os.Getenv("endpoint3")
	if realBackendURL == "" {
		errorResponse := map[string]interface{}{
			"error_code":    5000001,
			"error_message": "internal server error",
		}
		errorResponseJSON, _ := json.Marshal(errorResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorResponseJSON)
		utils.Log(referenceId, "Real Backend URL is not defined")
		utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))
		return
	}

	// Ambil bot_id dari body request
	var reqBody requestBody3
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		errorResponse := map[string]interface{}{
			"error_code":    4000000,
			"error_message": "invalid request",
		}
		errorResponseJSON, _ := json.Marshal(errorResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponseJSON)
		utils.Log(referenceId, "Failed to decode request body: "+err.Error())
		utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))
		return
	}

	// Extract botID from the request body
	botID := reqBody.Data.BotID
	if botID == "" {
		errorResponse := map[string]interface{}{
			"error_code":    4000001,
			"error_message": "invalid request. invalid field value",
		}
		errorResponseJSON, _ := json.Marshal(errorResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponseJSON)
		utils.Log(referenceId, "Bot ID not found in request body")
		utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))
		return
	}

	// Ambil session_id dan hash dari header
	sessionID := r.Header.Get("ecwx-session-id")
	hash := r.Header.Get("ecwx-hash")

	// Validate request
	validationResult, err := utils.ValidateRequestHash(botID, sessionID, hash, reqBody, db)
	if err != nil || validationResult == nil {
		errorResponse := map[string]interface{}{
			"error_code":    4010000,
			"error_message": "unauthenticated",
		}
		errorResponseJSON, _ := json.Marshal(errorResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(errorResponseJSON)
		utils.Log(referenceId, "Unauthenticated request")
		utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))
		return
	}

	// Get userId and organizationId from validation
	userID := validationResult.UserId
	organizationID := validationResult.OrganizationId

	utils.Log(referenceId, fmt.Sprintf("Bot ID received: %s", botID))
	utils.Log(referenceId, fmt.Sprintf("User ID from session data: %s", userID))
	utils.Log(referenceId, fmt.Sprintf("Organization ID from session data: %s", organizationID))

	isOrganization := utils.CheckBotOrganization(db, botID, userID, organizationID)
	if !isOrganization {
		errorResponse := map[string]interface{}{
			"error_code":    4030000,
			"error_message": "forbidden",
		}
		errorResponseJSON, _ := json.Marshal(errorResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		w.Write(errorResponseJSON)
		utils.Log(referenceId, "Bot ID does not match organization ID")
		utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))
		return
	}

	utils.Log(referenceId, "Hash is valid")
	utils.Log(referenceId, fmt.Sprintf("Continuing request to real backend URL: %s", realBackendURL))

	// Send request to real backend
	client := &http.Client{}
	requestBody3, err := json.Marshal(reqBody) // Use reqBody for backend
	if err != nil {
		errorResponse := map[string]interface{}{
			"error_code":    5000002,
			"error_message": "internal server error",
		}
		errorResponseJSON, _ := json.Marshal(errorResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorResponseJSON)
		utils.Log(referenceId, "Error encoding request body")
		utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))
		return
	}

	req, err := http.NewRequest("POST", realBackendURL, bytes.NewBuffer(requestBody3))
	if err != nil {
		errorResponse := map[string]interface{}{
			"error_code":    5000003,
			"error_message": "internal server error",
		}
		errorResponseJSON, _ := json.Marshal(errorResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorResponseJSON)
		utils.Log(referenceId, "Error creating request to backend")
		utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("ecwx-session-id", sessionID)
	req.Header.Set("ecwx-hash", hash)

	// Kirim request ke backend
	resp, err := client.Do(req)
	if err != nil {
		errorResponse := map[string]interface{}{
			"error_code":    5000004,
			"error_message": "internal server error",
		}
		errorResponseJSON, _ := json.Marshal(errorResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorResponseJSON)
		utils.Log(referenceId, fmt.Sprintf("Error forwarding request to backend: %v", err))
		utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))
		return
	}
	defer resp.Body.Close()

	// Baca response dari backend
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errorResponse := map[string]interface{}{
			"error_code":    5000005,
			"error_message": "internal server error",
		}
		errorResponseJSON, _ := json.Marshal(errorResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorResponseJSON)
		utils.Log(referenceId, "Error reading response from backend")
		utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))
		return
	}

	// Kirim response dari backend ke FE (frontend)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)

	utils.Log(referenceId, fmt.Sprintf("Response from real backend: res.status(%d).json(%s);", resp.StatusCode, string(body)))
}
