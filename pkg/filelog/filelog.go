package filelog

import (
	"fmt"
	"log"
	"os"
	"time"
)

func LogToFile(errorMessage string, request interface{}) {
	logPath := os.Getenv("LOG_PATH")
	currentTime := time.Now().Format(time.RFC3339)
	logMessage := fmt.Sprintf("Time: %s\nError: %s\nRequest: %+v\n---------------------------------------------------------------------------\n\n", currentTime, errorMessage, request)
	file, err := os.OpenFile(logPath+"logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(logMessage); err != nil {
		log.Fatalf("Failed to write to log file: %v", err)
	}
}
