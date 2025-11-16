package main

import (
	"encoding/json"
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
	http.HandleFunc("/api/v1/management", handlers.GetManagement)
	http.HandleFunc("/api/v1/owners", handlers.GetOwners)
	http.HandleFunc("/api/v1/top_investors", handlers.GetTopInvestors)
	http.HandleFunc("/api/v1/search", handlers.SearchCompanies)
	http.HandleFunc("/api/v1/sitemap", handlers.GetSitemap)
	http.HandleFunc("/api/v1/increase", handlers.GetCapitalIncreases)

	// Start server with CORS middleware
	addr := fmt.Sprintf(":%s", port)
	pkg.LogInfo(fmt.Sprintf("Server listening on %s", addr))

	handler := corsMiddleware(http.DefaultServeMux)
	if err := http.ListenAndServe(addr, handler); err != nil {
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
	response := map[string]interface{}{
		"message": "Virk API",
		"version": "1.0.0",
		"endpoints": map[string]interface{}{
			"health": map[string]string{
				"path":        "/health",
				"description": "Health check endpoint",
			},
			"stats": map[string]interface{}{
				"database": map[string]string{
					"path":        "/api/stats/database",
					"description": "Database statistics",
				},
				"table": map[string]string{
					"path":        "/api/stats/table",
					"description": "Table statistics",
				},
			},
			"v1": map[string]interface{}{
				"feed": map[string]string{
					"path":        "/api/v1/feed?page=1",
					"description": "Companies with capital increases (paginated)",
				},
				"increase": map[string]string{
					"path":        "/api/v1/increase?cvr=37890235",
					"description": "Capital increases for a specific company by CVR",
				},
				"company_info": map[string]string{
					"path":        "/api/v1/company_info?cvr=37890235",
					"description": "Basic company information",
				},
				"management": map[string]string{
					"path":        "/api/v1/management?cvr=37890235",
					"description": "Company management and board members",
				},
				"owners": map[string]string{
					"path":        "/api/v1/owners?cvr=37890235",
					"description": "Company ownership information",
				},
				"top_investors": map[string]string{
					"path":        "/api/v1/top_investors",
					"description": "Top investors by number of companies",
				},
				"search": map[string]string{
					"path":        "/api/v1/search?query=oaktoad",
					"description": "Search companies by name",
				},
				"sitemap": map[string]string{
					"path":        "/api/v1/sitemap",
					"description": "Sitemap of all companies",
				},
			},
		},
	}
	
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(response)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// corsMiddleware adds CORS headers to all responses
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
