//utils/validateRequestHash.go
package utils

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

type ValidationResult struct {
	BotId          string
	OrganizationId string
	UserId         string
}

// ValidateRequestHash checks the request hash and returns bot, organization, and user IDs, or an error.
func ValidateRequestHash(botID, sessionID, hash string, postBody interface{}, db *sql.DB) (*ValidationResult, error) {
	referenceId := GlobalVarInstance.GetReferenceId()

	Log(referenceId, "Executing method: ValidateRequestHash")
	Log(referenceId, "Session ID Received:", sessionID)
	Log(referenceId, "Hash Received:", hash)

	// Validate fields
	var missingFields []string
	if sessionID == "" {
		missingFields = append(missingFields, "session_id")
	}
	if hash == "" {
		missingFields = append(missingFields, "hash")
	}
	if botID == "" {
		missingFields = append(missingFields, "bot_id")
	}

	if len(missingFields) > 0 {
		Log(referenceId, fmt.Sprintf("Missing fields: %v", missingFields), postBody)
		return nil, errors.New("missing fields")
	}

	// Query to find session data
	query := `SELECT a.session_secret, a.user_id, b.organization_id 
	          FROM servouser.session a 
	          LEFT JOIN servouser.user b ON b.id = a.user_id 
	          WHERE a.session_id = $1 
	          LIMIT 1`
	Log(referenceId, "Query to find session data and organization ID:", query)

	var sessionSecret, userID, organizationID string
	err := db.QueryRow(query, sessionID).Scan(&sessionSecret, &userID, &organizationID)
	if err == sql.ErrNoRows {
		Log(referenceId, fmt.Sprintf("Session data with id = [%s] not found in database", sessionID))
		return nil, errors.New("session data not found")
	} else if err != nil {
		Log(referenceId, "Error executing query:", err)
		return nil, err
	}

	Log(referenceId, "Session Secret:", sessionSecret)
	Log(referenceId, "User ID:", userID)
	Log(referenceId, "Organization ID:", organizationID)
	Log(referenceId, "Post Body:", postBody)

	// Convert postBody to JSON string
	postBodyBytes, err := json.Marshal(postBody)
	if err != nil {
		Log(referenceId, "Error marshalling postBody:", err)
		return nil, err
	}
	postBodyString := string(postBodyBytes)

	// Generate expected hash
	hashExpected := CreateHMACSHA256HashBase64(postBodyString, sessionSecret)
	Log(referenceId, "Expected Hash:", hashExpected)
	Log(referenceId, "Received Hash:", hash)

	// Validate hash
	if hash != hashExpected {
		Log(referenceId, "Hash validation failed. Expected:", hashExpected, "Received:", hash)
		return nil, errors.New("hash validation failed")
	}

	Log(referenceId, "Bot ID from request body:", botID)
	Log(referenceId, "User ID from DB:", userID)

	// Return validated data
	return &ValidationResult{
		BotId:          botID,
		OrganizationId: organizationID,
		UserId:         userID,
	}, nil
}
