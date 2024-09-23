package utils

func CalculateChallengeResponse(full_nonce string, salted_password string) string{
	referenceId := GlobalVarInstance.GetReferenceId()

	Log(referenceId, "Executing method: calculateChallengeResponse")
	Log(referenceId, "Full Nonce: ", full_nonce)
	Log(referenceId, "Salted Password: ", salted_password)

	challengeResponse := CreateHMACSHA256HashBase64(full_nonce, salted_password)

	Log(referenceId, "Challenge Response generated:", challengeResponse)
	return challengeResponse
}
