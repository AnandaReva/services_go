package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// CreateHMACSHA256HashBase64 generates an HMAC-SHA256 hash and encodes it in Base64.
func CreateHMACSHA256HashBase64(data string, key string) string {
	// Mengambil referenceId dari GlobalInstance
	referenceId := GlobalVarInstance.GetReferenceId()

	// Logging proses hash
	Log(referenceId, "Executing method: CreateHMACSHA256HashBase64")
	Log(referenceId, "Data:", data)
	Log(referenceId, "Key:", key)

	// Membuat hash HMAC-SHA256
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))

	// Menghasilkan hash dalam bentuk base64
	hmacSHA256Hash := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Logging hasil hash
	Log(referenceId, "Generated HMAC-SHA256 Hash (Base64):", hmacSHA256Hash)

	// Mengembalikan hash yang sudah di-encode base64
	return hmacSHA256Hash
}

// CreateHMACSHA256HashHex generates an HMAC-SHA256 hash and encodes it in Hex.
func CreateHMACSHA256HashHex(data string, key string) string {
	referenceId := GlobalVarInstance.GetReferenceId()

	Log(referenceId, "Executing method: CreateHMACSHA256HashHex")
	Log(referenceId, "Data:", data)
	Log(referenceId, "Key:", key)


	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))

	// Menghasilkan hash dalam bentuk hex
	hmacSHA256Hash := hex.EncodeToString(h.Sum(nil))
	Log(referenceId, "Generated HMAC-SHA256 Hash (Hex):", hmacSHA256Hash)
	return hmacSHA256Hash
}
