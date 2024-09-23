// utils/generateRandomString.go
package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomString(length int) string {
	referenceId := GlobalVarInstance.GetReferenceId()

	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}

	Log(referenceId, "Random String generated", result)
	return string(result)
}
