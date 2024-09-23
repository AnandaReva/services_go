package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"Auth_service_go/utils"
	"database/sql"
	"strconv"
)

// deleteChallengeResponse deletes a challenge response from the database
func deleteChallengeResponse(db *sql.DB, fullNonce string) {
	referenceId := utils.GlobalVarInstance.GetReferenceId()

	utils.Log(referenceId, "Executing method: deleteChallengeResponse")
	utils.Log(referenceId, "Full nonce:", fullNonce)

	deleteQuery := "DELETE FROM servouser.challenge_response WHERE full_nonce = $1 RETURNING *"
	utils.Log(referenceId, "Delete Challenge Response Query:", deleteQuery)

	_, err := db.Exec(deleteQuery, fullNonce)
	if err != nil {
		utils.Log(referenceId, "Failed to delete challenge response:", err)
		return
	}

	utils.Log(referenceId, "Challenge response successfully deleted for full_nonce:", fullNonce)
}

// getUserPrivileges retrieves user privileges based on user ID
func getUserPrivileges(db *sql.DB, userId int64) string {
	referenceId := utils.GlobalVarInstance.GetReferenceId()
	utils.Log(referenceId, "Executing method: getUserPrivileges")

	getRoleQuery := "SELECT role FROM servouser.user WHERE id = $1 LIMIT 1"
	utils.Log(referenceId, "Get role query:", getRoleQuery)

	var role string
	err := db.QueryRow(getRoleQuery, userId).Scan(&role)
	if err == sql.ErrNoRows {
		utils.Log(referenceId, "No role found for user ID:", userId)
		return "0"
	} else if err != nil {
		utils.Log(referenceId, "Error retrieving role:", err)
		return "0"
	}

	utils.Log(referenceId, "userId:", userId)
	utils.Log(referenceId, "Role from db servouser.user:", role)

	getPrivilegesQuery := "SELECT privileges FROM servouser.role WHERE name = $1 LIMIT 1"
	utils.Log(referenceId, "Get privileges query:", getPrivilegesQuery)

	var privileges string
	err = db.QueryRow(getPrivilegesQuery, role).Scan(&privileges)
	if err == sql.ErrNoRows {
		utils.Log(referenceId, "No privileges found for the role:", role)
		return "0"
	} else if err != nil {
		utils.Log(referenceId, "Error retrieving privileges:", err)
		return "0"
	}

	utils.Log(referenceId, "Privileges from db servouser.role:", privileges)
	return privileges
}

// getOrganizationTier retrieves the organization tier based on user ID
func getOrganizationTier(db *sql.DB, userId int64) string {
	referenceId := utils.GlobalVarInstance.GetReferenceId()
	utils.Log(referenceId, "Executing method: getOrganizationTier")

	getOrgIdQuery := "SELECT organization_id FROM servouser.user WHERE id = $1 LIMIT 1"
	utils.Log(referenceId, "Get organization ID query:", getOrgIdQuery)

	var organizationId int64
	err := db.QueryRow(getOrgIdQuery, userId).Scan(&organizationId)
	if err == sql.ErrNoRows {
		utils.Log(referenceId, "No organization_id found for user ID:", userId)
		return "0"
	} else if err != nil {
		utils.Log(referenceId, "Error retrieving organization ID:", err)
		return "0"
	}

	utils.Log(referenceId, "User ID:", userId)
	utils.Log(referenceId, "Organization ID from servouser.user:", organizationId)

	getOrganizationTierQuery := "SELECT tier FROM servouser.organization WHERE id = $1 LIMIT 1"
	utils.Log(referenceId, "Get organization tier query:", getOrganizationTierQuery)

	var tier string
	err = db.QueryRow(getOrganizationTierQuery, organizationId).Scan(&tier)
	if err == sql.ErrNoRows {
		utils.Log(referenceId, "No tier found for organization_id:", organizationId)
		return "0"
	} else if err != nil {
		utils.Log(referenceId, "Error retrieving organization tier:", err)
		return "0"
	}

	utils.Log(referenceId, "Organization tier from servouser.organization:", tier)
	return tier
}

// upsertSession upserts a session into the database
func upsertSession(db *sql.DB, sessionID string, userID string, sessionSecret string) error {
	referenceId := utils.GlobalVarInstance.GetReferenceId()
	utils.Log(referenceId, "Executing method: upsertSession")

	queryUpsertSession := `
        INSERT INTO servouser.session (session_id, user_id, session_secret, tstamp, st)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (user_id) DO UPDATE
        SET session_id = EXCLUDED.session_id,
            session_secret = EXCLUDED.session_secret,
            tstamp = EXCLUDED.tstamp,
            st = EXCLUDED.st
    `
	_, err := db.Exec(queryUpsertSession, sessionID, userID, sessionSecret, time.Now().Unix(), 1)
	if err != nil {
		utils.Log(referenceId, "Error during upsert session:", err)
		return err
	}

	utils.Log(referenceId, "Session upserted successfully.")
	return nil
}

func HandleChallengeResponseVerification(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	referenceId := utils.GlobalVarInstance.GetReferenceId()
	utils.Log(referenceId, "Execute method: handleChallengeResponseVerification")

	var requestBody struct {
		FullNonce         string `json:"full_nonce"`
		ChallengeResponse string `json:"challenge_response"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		utils.Log(referenceId, "REQUEST_ERROR", "Failed to decode request body:", err)

		errorResponse := map[string]interface{}{
			"error_code":    4000003,
			"error_message": "invalid request",
		}

		errorResponseJSON, _ := json.Marshal(errorResponse)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponseJSON)

		utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))

		return
	}

	fullNonce := requestBody.FullNonce
	challengeResponse := requestBody.ChallengeResponse

	missingFields := []string{}
	if fullNonce == "" {
		missingFields = append(missingFields, "full_nonce")
	}
	if challengeResponse == "" {
		missingFields = append(missingFields, "challenge_response")
	}

	if len(missingFields) > 0 {

		utils.Log(referenceId, " Missing fields:", missingFields)
		errorResponse := map[string]interface{}{
			"error_code":    4000004,
			"error_message": "invalid request",
		}

		errorResponseJSON, _ := json.Marshal(errorResponse)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponseJSON) // Write the response body
		utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))
		return
	}

	if len(fullNonce) != 16 {
		utils.Log(referenceId, "VALIDATION_ERROR", " full_nonce must be 16 characters length")
		errorResponse := map[string]interface{}{
			"error_code":    4000005,
			"error_message": "invalid request",
		}

		errorResponseJSON, _ := json.Marshal(errorResponse)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponseJSON) // Write the response body
		utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))
		return
	}

	// Find challenge response in DB
	/* 	challengeDataQuery := `
	SELECT cr.*, u.saltedpassword, u.full_name, u.id as user_id
	FROM servouser.challenge_response cr
	JOIN servouser.user u ON cr.user_id = u.id
	WHERE cr.full_nonce = $1` */
	challengeDataQuery := `SELECT cr.full_nonce, cr.user_id, cr.challenge_response, cr.tstamp, 
	u.saltedpassword, u.full_name 
	FROM servouser.challenge_response cr
	JOIN servouser.user u ON cr.user_id = u.id
	WHERE cr.full_nonce = $1
`

	challengeDataResult, err := db.Query(challengeDataQuery, fullNonce)
	if err != nil {
		utils.Log(referenceId, "Error querying challenge data:", err)
		errorResponse := map[string]interface{}{
			"error_code":    5000004,
			"error_message": "internal server error",
		}

		errorResponseJSON, _ := json.Marshal(errorResponse)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponseJSON)
		utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))
		return
	}
	defer challengeDataResult.Close()

	if !challengeDataResult.Next() {
		utils.Log(referenceId, "Challenge not valid: The challenge provided is not valid.")
		errorResponse := map[string]interface{}{
			"error_code":    4010000,
			"error_message": "unauthenticated",
		}

		errorResponseJSON, _ := json.Marshal(errorResponse)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponseJSON)
		utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))
		return
	}

	var challengeData struct {
		FullNonce         string
		UserID            int64
		ChallengeResponse string
		Timestamp         int64
		SaltedPassword    string
		FullName          string
	}

	if err := challengeDataResult.Scan(
		&challengeData.FullNonce,         // Full nonce
		&challengeData.UserID,            // User ID
		&challengeData.ChallengeResponse, // Challenge response
		&challengeData.Timestamp,         // Timestamp
		&challengeData.SaltedPassword,    // Salted password
		&challengeData.FullName,          // Full name
	); err != nil {
		utils.Log(referenceId, "Error scanning challenge data:", err)

		errorResponse := map[string]interface{}{
			"error_code":    5000005,
			"error_message": "internal server error",
		}

		errorResponseJSON, _ := json.Marshal(errorResponse)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponseJSON)
		utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))
		return
	}

	currentTime := time.Now().Unix()
	if currentTime-challengeData.Timestamp > 60 {
		utils.Log(referenceId, "Challenge response exceeds 60 seconds, deleting challenge response")

		deleteChallengeResponse(db, fullNonce)

		errorResponse := map[string]interface{}{
			"error_code":    4010001,
			"error_message": "unauthenticated",
		}

		errorResponseJSON, _ := json.Marshal(errorResponse)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponseJSON) // Write the response body
		utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))

		utils.Log(referenceId, "Challenge has expired")
		return
	}

	// Verify challenge response
	expectedChallengeResponse := utils.CalculateChallengeResponse(fullNonce, challengeData.SaltedPassword)
	isValid := expectedChallengeResponse == challengeResponse
	utils.Log(referenceId, "Expected Challenge Response:", expectedChallengeResponse)
	utils.Log(referenceId, "Compared challenge response:", challengeResponse)

	if isValid {
		sessionID := utils.GenerateRandomString(16)
		nonce2 := utils.GenerateRandomString(8)
		sessionSecret := utils.CreateHMACSHA256HashHex(fullNonce+nonce2, challengeData.SaltedPassword)

		utils.Log(referenceId, "Challenge response valid")
		utils.Log(referenceId, "Generate session ID and nonce2:")
		utils.Log(referenceId, "session_id:", sessionID)
		utils.Log(referenceId, "nonce2:", nonce2)
		utils.Log(referenceId, "session_secret:", sessionSecret)

		if err := upsertSession(db, sessionID, strconv.FormatInt(challengeData.UserID, 10), sessionSecret); err != nil {
			utils.Log(referenceId, "Error upsert session", err)

			errorResponse := map[string]interface{}{
				"error_code":    5000006,
				"error_message": "internal server error",
			}

			errorResponseJSON, _ := json.Marshal(errorResponse)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errorResponseJSON)
			utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))
			return
		}

		privileges := getUserPrivileges(db, challengeData.UserID)
		tier := getOrganizationTier(db, challengeData.UserID)

		var privilegesData map[string]interface{}
		if err := json.Unmarshal([]byte(privileges), &privilegesData); err != nil {
			// Handle the error if decoding fails
			utils.Log(referenceId, "Error decoding privileges:", err)
			privilegesData = map[string]interface{}{"error": "invalid privileges format"}
		}

		response := map[string]interface{}{
			"session_id": sessionID,
			"nonce2":     nonce2,
			"user_data": map[string]interface{}{
				"full_name":  challengeData.FullName,
				"privileges": privilegesData, // Use decoded data
				"tier":       tier,
			},
		}
		w.Header().Set("Content-Type", "application/json") // Set content type
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

		// response := map[string]interface{}{
		// 	"session_id": sessionID,
		// 	"nonce2":     nonce2,
		// 	"user_data": map[string]interface{}{
		// 		"full_name":  challengeData.FullName,
		// 		"privileges": privileges,
		// 		"tier":       tier,
		// 	},
		// }
		// w.WriteHeader(http.StatusOK)
		// json.NewEncoder(w).Encode(response)

		utils.Log(referenceId, "Continuing Response Real Backend to FE:", response)
	} else {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		utils.Log(referenceId, "Invalid challenge response.")

		errorResponse := map[string]interface{}{
			"error_code":    4010002,
			"error_message": "unauthenticated",
		}

		errorResponseJSON, _ := json.Marshal(errorResponse)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponseJSON)
		utils.Log(referenceId, "RESPONSE_SENT", string(errorResponseJSON))
		return
	}

	defer deleteChallengeResponse(db, fullNonce)
}
