package api

import (
	"fmt"
	"os"
)

// EnableDebugLogging enables detailed request/response logging for debugging
var EnableDebugLogging = os.Getenv("BC_CLI_DEBUG") == "1"

func logRequest(method, url string, body any) {
	if !EnableDebugLogging {
		return
	}
	fmt.Printf("\n=== REQUEST ===\n")
	fmt.Printf("Method: %s\n", method)
	fmt.Printf("URL: %s\n", url)
	if body != nil {
		fmt.Printf("Body: %+v\n", body)
	}
	fmt.Printf("===============\n\n")
}

func logResponse(statusCode int, body []byte) {
	if !EnableDebugLogging {
		return
	}
	fmt.Printf("\n=== RESPONSE ===\n")
	fmt.Printf("Status: %d\n", statusCode)
	fmt.Printf("Body: %s\n", string(body))
	fmt.Printf("================\n\n")
}
