package handlers

import (
	"encoding/json"
	"github.com/AndersKaae/go_virk_api/internal/database"
	"github.com/AndersKaae/go_virk_updater/pkg"
	"net/http"
)

// GetDatabaseStats returns statistics for all tables in the database
func GetDatabaseStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := database.GetDatabaseStats()
	if err != nil {
		pkg.LogError("Failed to get database stats: " + err.Error())
		http.Error(w, "Failed to retrieve database statistics", http.StatusInternalServerError)
		return
	}

	pkg.LogInfo("Database stats retrieved successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

// GetTableStats returns statistics for a specific table
func GetTableStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get table name from query parameter
	tableName := r.URL.Query().Get("table")
	if tableName == "" {
		http.Error(w, "Missing 'table' query parameter", http.StatusBadRequest)
		return
	}

	stats, err := database.GetTableStats(tableName)
	if err != nil {
		pkg.LogError("Failed to get table stats for " + tableName + ": " + err.Error())
		http.Error(w, "Failed to retrieve table statistics", http.StatusInternalServerError)
		return
	}

	pkg.LogInfo("Table stats retrieved for: " + tableName)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}
