package handlers

import (
	"encoding/json"
	"github.com/AndersKaae/go_virk_api/internal/database"
	"github.com/AndersKaae/go_virk_updater/pkg"
	"net/http"
	"strconv"
)

// GetFeed returns the feed of companies with capital increases
func GetFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get page parameter (default: 1)
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			http.Error(w, "Invalid page parameter", http.StatusBadRequest)
			return
		}
	}

	// Page size (matching old API behavior: ~46-50 items)
	pageSize := 50

	feed, err := database.GetFeed(page, pageSize)
	if err != nil {
		pkg.LogError("Failed to get feed: " + err.Error())
		http.Error(w, "Failed to retrieve feed", http.StatusInternalServerError)
		return
	}

	pkg.LogInfo("Feed retrieved successfully for page " + pageStr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feed)
}
