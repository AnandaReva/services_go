package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// Parameter structures for API calls
type API1Parameters struct {
	Data struct {
		RowLength  int    `json:"row_length"`
		Page       int    `json:"page"`
		SortColumn int    `json:"sort_column"`
		Direction  string `json:"direction"`
		BotID      string `json:"bot_id"`
	} `json:"data"`
	FromDate     int    `json:"from_date"`
	ToDate       int    `json:"to_date"`
	SearchFilter string `json:"search_filter"`
	DateMode     int    `json:"date_mode"`
}

type API2Parameters struct {
	Data struct {
		BotID string `json:"bot_id"`
	} `json:"data"`
	FromDate int `json:"from_date"`
	ToDate   int `json:"to_date"`
	DateMode int `json:"date_mode"`
}

type API3Parameters struct {
	Data struct {
		BotID string `json:"bot_id"`
	} `json:"data"`
	FromDate     int    `json:"from_date"`
	ToDate       int    `json:"to_date"`
	SearchFilter string `json:"search_filter"`
	DateMode     int    `json:"date_mode"`
}

type API4Parameters struct {
	JobID string `json:"job_id"`
}

type API5Parameters struct {
	BotID string `json:"bot_id"`
}

type API6Parameters struct {
	BotID              string `json:"bot_id"`
	ChildPromptID      string `json:"child_prompt_id"`
	KnowledgeText      string `json:"knowledge_text"`
	ClassificationName string `json:"classification_name"`
}

type API7Parameters struct {
	BotID string `json:"bot_id"`
}

type API8Parameters struct {
	BotID    string `json:"bot_id"`
	Greeting string `json:"greeting"`
	Topics   string `json:"topics"`
}

type API9Parameters struct {
	BotID string `json:"bot_id"`
}

func main() {
	sessionSecret := "98a38f3482febe01c7de873eed504667a075e2c75c579778b00c14bb8cd2913a"

	// API parameter initialization and hashing
	hashAPIParameters(sessionSecret)
}

func hashAPIParameters(sessionSecret string) {
	// API 1 parameters
	api1Params := API1Parameters{
		FromDate:     0,
		ToDate:       0,
		SearchFilter: "string",
		DateMode:     0,
	}
	api1Params.Data = struct {
		RowLength  int    `json:"row_length"`
		Page       int    `json:"page"`
		SortColumn int    `json:"sort_column"`
		Direction  string `json:"direction"`
		BotID      string `json:"bot_id"`
	}{
		RowLength:  0,
		Page:       0,
		SortColumn: 0,
		Direction:  "desc",
		BotID:      "1",
	}
	hashAndPrint(api1Params, sessionSecret, "API 1")

	// API 2 parameters
	api2Params := API2Parameters{
		FromDate: 0,
		ToDate:   0,
		DateMode: 0,
	}
	api2Params.Data.BotID = "1"
	hashAndPrint(api2Params, sessionSecret, "API 2")

	// API 3 parameters
	api3Params := API3Parameters{
		FromDate:     0,
		ToDate:       0,
		SearchFilter: "string",
		DateMode:     0,
	}
	api3Params.Data.BotID = "1"
	hashAndPrint(api3Params, sessionSecret, "API 3")

	// API 4 parameters
	api4Params := API4Parameters{
		JobID: "1725331930-c25c0600518e422c8a6ce2ee51998089",
	}
	hashAndPrint(api4Params, sessionSecret, "API 4")

	// API 5 parameters
	api5Params := API5Parameters{
		BotID: "1",
	}
	hashAndPrint(api5Params, sessionSecret, "API 5")

	// API 6 parameters
	api6Params := API6Parameters{
		BotID:              "1",
		ChildPromptID:      "",
		KnowledgeText:      "",
		ClassificationName: "",
	}
	hashAndPrint(api6Params, sessionSecret, "API 6")

	// API 7 parameters
	api7Params := API7Parameters{
		BotID: "1",
	}
	hashAndPrint(api7Params, sessionSecret, "API 7")

	// API 8 parameters
	api8Params := API8Parameters{
		BotID:    "1",
		Greeting: "",
		Topics:   "",
	}
	hashAndPrint(api8Params, sessionSecret, "API 8")

	// API 9 parameters
	api9Params := API9Parameters{
		BotID: "1",
	}
	hashAndPrint(api9Params, sessionSecret, "API 9")
}

func hashAndPrint(params interface{}, key string, apiName string) {
	parametersJSON, err := json.Marshal(params)
	if err != nil {
		fmt.Printf("Error serializing parameters for %s: %v\n", apiName, err)
		return
	}
	hashedBody := generateHmac(string(parametersJSON), key, apiName)
	fmt.Printf("Hashed Body for %s: %s\n", apiName, hashedBody)
}

func generateHmac(message string, key string, notif string) string {
	fmt.Println("-----------------\nStart hashing", notif, "!!")
	fmt.Println("Key:", key)
	fmt.Println("Message:", message)

	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	expectedMAC := mac.Sum(nil)
	expectedMACBase64 := base64.StdEncoding.EncodeToString(expectedMAC)

	fmt.Println(notif, ":", expectedMACBase64)
	return expectedMACBase64
}
