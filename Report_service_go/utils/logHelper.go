//utils/logHelper.go

package utils

import (
	"fmt"
	"time"
	//	"Auth_service_go/config"
)

func Log(referenceID string, text string, param ...interface{}) {
	timestamp := time.Now().Format(time.RFC3339)

	if len(param) == 0 {
		fmt.Printf("%s - %s - %s\n", timestamp, referenceID, text)
	} else {
		fmt.Printf("%s - %s - %s %v\n", timestamp, referenceID, text, param)
	}
}
