package logger

import (
	"fmt"
	"log"
	"os"
)

func SetupLogging(filepath string) *os.File {
	logFile, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		os.Exit(1)
	}

	log.SetOutput(logFile)
	return logFile
}
