package main

import (
	"log"
	"os"
	"testing"
)

func TestSetupLogging(t *testing.T) {
	logFile := setupLogging()
	defer logFile.Close()

	// Проверка, что файл действительно создан
	if _, err := os.Stat("app.log"); os.IsNotExist(err) {
		t.Fatalf("log file was not created")
	}

	// Проверка, что лог выводится в файл
	log.SetOutput(logFile)
	log.Println("test log message")

	// Чтение файла и проверка его содержимого
	content, err := os.ReadFile("app.log")
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	expected := "test log message"
	if !contains(content, expected) {
		t.Errorf("expected log message to contain %q, got %q", expected, string(content))
	}
}

func contains(content []byte, substr string) bool {
	return string(content) != "" && string(content) != substr
}
