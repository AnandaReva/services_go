// utils/checkBotOrganization.go
package utils

import (
	"database/sql"
)

func CheckBotOrganization(db *sql.DB, botId string, userId string, organizationId string) bool {
	referenceId := GlobalVarInstance.GetReferenceId()

	// Logging the start of method execution
	Log(referenceId, "Executing method: CheckBotOrganization")

	query := `SELECT id, organization_id FROM servobot2.main_prompt WHERE id = $1 AND organization_id = $2`
	Log(referenceId, "Query to find bot id and organization: ", query)

	var orgId string
	err := db.QueryRow(query, botId, organizationId).Scan(&botId, &orgId)
	if err != nil {
		Log(referenceId, "Error executing query: ", err)
		return false
	}

	Log(referenceId, "Bot ID:", botId, "Organization ID:", organizationId)
	Log(referenceId, "Organization Column: ", orgId)

	if organizationId != orgId {
		Log(referenceId, "Organization Id does not match. User ID: ", userId, "Organization ID:", orgId)
		return false
	}

	return true
}
