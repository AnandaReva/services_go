// controllers/loginController.go
package controllers

import (
	"Auth_service_go/utils"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

// upsertChallengeResponse upserts challenge response into the database
func upsertChallengeResponse(w http.ResponseWriter, db *sql.DB, fullNonce string, userID int64, challengeResponse string, currentTime int64) error {
	referenceId := utils.GlobalVarInstance.GetReferenceId()
	upsertQuery := `
        INSERT INTO servouser.challenge_response (full_nonce, user_id, challenge_response, tstamp)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (user_id) DO UPDATE
        SET full_nonce = EXCLUDED.full_nonce,
            challenge_response = EXCLUDED.challenge_response,
            tstamp = EXCLUDED.tstamp
    `

	utils.Log(referenceId, "DB_EXEC", "Executing upsert query for challenge response:", upsertQuery)
	_, err := db.Exec(upsertQuery, fullNonce, userID, challengeResponse, currentTime)
	if err != nil {
		utils.Log("DB_ERROR", "Error during upsert challenge response", err)
		errorResponse := map[string]interface{}{
			"error_code":    5000000,
			"error_message": "internal server error",
		}

		errorResponseJSON, _ := json.Marshal(errorResponse)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponseJSON) // Write the response body
		utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))
		return err
	}

	utils.Log(referenceId, "DB_SUCCESS: ", "Challenge response upserted successfully.")
	return nil
}

// HandleLoginRequest handles login request by processing username and half_nonce
func HandleLoginRequest(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	referenceId := utils.GlobalVarInstance.GetReferenceId()
	
	var requestBody struct {
		Username  string `json:"username"`
		HalfNonce string `json:"half_nonce"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		utils.Log(referenceId, "REQUEST_ERROR", "Failed to decode request body:", err)

		errorResponse := map[string]interface{}{
			"error_code":    4000000,
			"error_message": "invalid request",
		}

		// Convert error response to JSON for logging
		errorResponseJSON, _ := json.Marshal(errorResponse)

		// Send the error response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponseJSON) // Write the response body

		// Log the response that was sent
		utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))

		return
	}

	username := requestBody.Username
	halfNonce := requestBody.HalfNonce

	missingFields := []string{}
	if username == "" {
		missingFields = append(missingFields, "username")
	}
	if halfNonce == "" {
		missingFields = append(missingFields, "half_nonce")
	}

	if len(missingFields) > 0 {

		utils.Log(referenceId, " Missing fields:", missingFields)
		errorResponse := map[string]interface{}{
			"error_code":    4000001,
			"error_message": "invalid request",
		}

		errorResponseJSON, _ := json.Marshal(errorResponse)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponseJSON) // Write the response body
		utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))
		return
	}

	if len(halfNonce) != 8 {
		utils.Log(referenceId, "VALIDATION_ERROR", "half_nonce must be 8 characters long")
		errorResponse := map[string]interface{}{
			"error_code":    4000002,
			"error_message": "invalid request",
		}

		errorResponseJSON, _ := json.Marshal(errorResponse)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponseJSON) // Write the response body
		utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))
		return
	}

	var user struct {
		ID             int64
		Salt           string
		SaltedPassword string
		Iterations     int
	}

	// Get user data from database
	userQuery := "SELECT id, salt, saltedpassword, iterations FROM servouser.user WHERE username = $1"
	utils.Log(referenceId, "DB_EXEC", "Executing user query:", userQuery)
	err := db.QueryRow(userQuery, username).Scan(&user.ID, &user.Salt, &user.SaltedPassword, &user.Iterations)

	if err == sql.ErrNoRows {
		fakeFullNonce := halfNonce + utils.GenerateRandomString(8)
		fakeSalt := utils.GenerateRandomString(26)
		fakeIterations := 0

		response := map[string]interface{}{
			"full_nonce": fakeFullNonce,
			"salt":       fakeSalt,
			"iterations": fakeIterations,
		}
		json.NewEncoder(w).Encode(response)
		utils.Log(referenceId, "USER_NOT_FOUND", "User not found, returning fake data for username:", username)
		return
	} else if err != nil {
		utils.Log(referenceId, "DB_ERROR", "Error retrieving user data:", err)

		errorResponse := map[string]interface{}{
			"error_code":    5000001,
			"error_message": "internal server error",
		}

		errorResponseJSON, _ := json.Marshal(errorResponse)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponseJSON)
		utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))
		return
	}

	// Check for existing challenge response timestamp
	challengeQuery := "SELECT tstamp FROM servouser.challenge_response WHERE user_id = $1"
	var existingTstamp int64
	err = db.QueryRow(challengeQuery, user.ID).Scan(&existingTstamp)

	currentTime := time.Now().Unix()

	if err == nil && (currentTime-existingTstamp) < 10 {
		utils.Log(referenceId, "REQUEST_LIMIT", "Too many requests, try again later after 10s")
		errorResponse := map[string]interface{}{
			"error_code":    4290000,
			"error_message": "too many requests",
		}

		errorResponseJSON, _ := json.Marshal(errorResponse)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponseJSON)
		utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))

		return
	}

	// Generate new full_nonce and challenge response
	nonce1 := utils.GenerateRandomString(8)
	fullNonce := halfNonce + nonce1
	challengeResponse := utils.CalculateChallengeResponse(fullNonce, user.SaltedPassword)

	utils.Log(referenceId, "CHALLENGE_RESPONSE", "Generated full_nonce:", fullNonce, "for user ID:", user.ID)

	// Upsert challenge response
	if err := upsertChallengeResponse(w, db, fullNonce, user.ID, challengeResponse, currentTime); err != nil {
		utils.Log(referenceId, "DB_ERROR", "Error during upsert challenge response:", err)
		errorResponse := map[string]interface{}{
			"error_code":    5000002,
			"error_message": "internal server error",
		}

		errorResponseJSON, _ := json.Marshal(errorResponse)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponseJSON) // Write the response body
		utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Prepare response
	response := map[string]interface{}{
		"full_nonce": fullNonce,
		"salt":       user.Salt,
		"iterations": user.Iterations,
	}

	// Encode response to JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.Log(referenceId, "Error encoding response to JSON:", err)
		errorResponse := map[string]interface{}{
			"error_code":    5000003,
			"error_message": "internal server error",
		}

		errorResponseJSON, _ := json.Marshal(errorResponse)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponseJSON) // Write the response body
		utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))
		return
	}

	// Log the successful response
	utils.Log(referenceId, "RESPONSE_SENT", "Response sent successfully:", response)
}
