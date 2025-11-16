package main

import (
	"fmt"
	"github.com/AndersKaae/go_virk_api/internal/database"
	"github.com/AndersKaae/go_virk_api/internal/handlers"
	"github.com/AndersKaae/go_virk_updater/pkg"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize logger from go_virk_updater package
	pkg.InitLoggers(pkg.INFO)
	pkg.LogInfo("Starting Virk API Server")

	// Connect to database
	_, err := database.Connect()
	if err != nil {
		pkg.LogError(fmt.Sprintf("Failed to connect to database: %v", err))
		log.Fatal(err)
	}
	defer database.Close()
	pkg.LogInfo("Database connection established")

	// Get port from environment or use default
	port := getEnv("PORT", "8080")

	// Setup routes
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/api/stats/database", handlers.GetDatabaseStats)
	http.HandleFunc("/api/stats/table", handlers.GetTableStats)
	http.HandleFunc("/api/v1/feed", handlers.GetFeed)
	http.HandleFunc("/api/v1/company_info", handlers.GetCompanyInfo)

	// Start server
	addr := fmt.Sprintf(":%s", port)
	pkg.LogInfo(fmt.Sprintf("Server listening on %s", addr))

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Virk API","version":"1.0.0"}`))
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
