package middlewares

import (
	"Report_service_go/utils"
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func LogRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a unique reference ID
		referenceId := utils.GenerateRandomString(6)
		utils.GlobalVarInstance.SetReferenceId(referenceId) // Set ReferenceId di sini

		// Retrieve request details
		ipAddress := r.Header.Get("X-Forwarded-For")
		if ipAddress == "" {
			ipAddress = r.RemoteAddr
		}
		method := r.Method
		url := r.URL.Path
		headers := r.Header
		contentType := r.Header.Get("Content-Type")

		// Log incoming request details
		fmt.Println("\n-----------------------------------------")
		utils.Log(referenceId, "INCOMING REQUEST FROM IP ADDRESS:", ipAddress)
		utils.Log(referenceId, fmt.Sprintf("Received %s request to url: %s", method, url))
		utils.Log(referenceId, "Headers:", headers)

		// Log request body
		if contentType == "application/json" || contentType == "application/x-www-form-urlencoded" {
			body := make([]byte, r.ContentLength)
			r.Body.Read(body)
			utils.Log(referenceId, "Body:", string(body))
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		} else {
			utils.Log(referenceId, "Body (Other):", "No body or unsupported content-type")
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}