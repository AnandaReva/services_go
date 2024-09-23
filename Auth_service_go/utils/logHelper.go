//utils/logHelper.go

package utils

import (
	"fmt"
	"time"
	//	"Auth_service_go/config"
)

/*
	Unlike JavaScript, Go does not have template literals (like backticks in JS). Instead, Go uses format specifiers (like %s, %d, %f, etc.) with fmt.Printf, fmt.Sprintf, and similar functions to achieve similar functionality.

How % Works:
%s inserts a string.
%d inserts an integer.
%f inserts a floating-point number.
%v is used as a general placeholder for variables of any type.
*/

// func Log(referenceID string, text string, param ...interface{}) {
// 	timestamp := time.Now().Format(time.RFC3339)
// 	version := Confi

// 	if len(param) == 0 {
// 		fmt.Printf("%s - %s - %s - %s\n", timestamp, version, referenceID, text)
// 	} else {
// 		fmt.Printf("%s - %s - %s - %s %v\n", timestamp, version, referenceID, text, param)
// 	}
// 
func Log(referenceID string, text string, param ...interface{}) {
	timestamp := time.Now().Format(time.RFC3339)

	if len(param) == 0 {
		fmt.Printf("%s - %s - %s\n", timestamp, referenceID, text)
	} else {
		fmt.Printf("%s - %s - %s %v\n", timestamp, referenceID, text, param)
	}
}
