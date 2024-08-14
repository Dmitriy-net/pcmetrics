package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Dmitriy-net/pcmetrics/internal/logger"
	"github.com/Dmitriy-net/pcmetrics/internal/repository/memstorage"
	"github.com/Dmitriy-net/pcmetrics/internal/router"
)

func main() {
	logFile := logger.SetupLogging("server.log")
	defer logFile.Close()

	address := flag.String("a", "localhost:8080", "HTTP server address")

	flag.Parse()

	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		*address = envAddress
	}
	log.Printf("Server is starting on http://%s\n", *address)
	repo := memstorage.NewMemStorage()
	r := router.SetupRouter(repo)

	err := http.ListenAndServe(*address, r)
	if err != nil {
		fmt.Printf("Failed to start server: %v", err)
	} else {
		fmt.Printf("Server is listening on http://%s\n", *address)
	}
}
