package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/AndersKaae/go_virk_api/internal/database"
	"github.com/AndersKaae/go_virk_updater/pkg"
	"net/http"
	"time"
)

// GetCompanyInfo returns basic company information by CVR number
func GetCompanyInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get CVR parameter
	cvr := r.URL.Query().Get("cvr")
	if cvr == "" {
		http.Error(w, "CVR parameter is required", http.StatusBadRequest)
		return
	}

	companyInfo, err := database.GetCompanyInfo(cvr)
	if err != nil {
		if err.Error() == "company not found" {
			http.Error(w, "Company not found", http.StatusNotFound)
			return
		}
		pkg.LogError("Failed to get company info: " + err.Error())
		http.Error(w, "Failed to retrieve company info", http.StatusInternalServerError)
		return
	}

	pkg.LogInfo("Company info retrieved successfully for CVR " + cvr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(companyInfo)
}

// GetManagement returns board and management information by CVR number
func GetManagement(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get CVR parameter
	cvr := r.URL.Query().Get("cvr")
	if cvr == "" {
		http.Error(w, "CVR parameter is required", http.StatusBadRequest)
		return
	}

	management, err := database.GetManagement(cvr)
	if err != nil {
		pkg.LogError("Failed to get management: " + err.Error())
		http.Error(w, "Failed to retrieve management", http.StatusInternalServerError)
		return
	}

	pkg.LogInfo("Management retrieved successfully for CVR " + cvr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(management)
}

// GetOwners returns ownership information by CVR number
func GetOwners(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get CVR parameter
	cvr := r.URL.Query().Get("cvr")
	if cvr == "" {
		http.Error(w, "CVR parameter is required", http.StatusBadRequest)
		return
	}

	owners, err := database.GetOwners(cvr)
	if err != nil {
		pkg.LogError("Failed to get owners: " + err.Error())
		http.Error(w, "Failed to retrieve owners", http.StatusInternalServerError)
		return
	}

	pkg.LogInfo("Owners retrieved successfully for CVR " + cvr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(owners)
}

// GetTopInvestors returns top investors by number of companies
func GetTopInvestors(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	investors, err := database.GetTopInvestors()
	if err != nil {
		pkg.LogError("Failed to get top investors: " + err.Error())
		http.Error(w, "Failed to retrieve top investors", http.StatusInternalServerError)
		return
	}

	pkg.LogInfo("Top investors retrieved successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(investors)
}

// SearchCompanies searches for companies by name
func SearchCompanies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get query parameter
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "query parameter is required", http.StatusBadRequest)
		return
	}

	results, err := database.SearchCompanies(query)
	if err != nil {
		pkg.LogError("Failed to search companies: " + err.Error())
		http.Error(w, "Failed to search companies", http.StatusInternalServerError)
		return
	}

	pkg.LogInfo("Search completed for query: " + query)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

// GetSitemap returns sitemap for all companies
func GetSitemap(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	entries, err := database.GetSitemap()
	if err != nil {
		pkg.LogError("Failed to get sitemap: " + err.Error())
		http.Error(w, "Failed to retrieve sitemap", http.StatusInternalServerError)
		return
	}

	pkg.LogInfo("Sitemap retrieved successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entries)
}

// GetCapitalIncreases returns capital increases for a company
func GetCapitalIncreases(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cvr := r.URL.Query().Get("cvr")
	if cvr == "" {
		http.Error(w, "cvr parameter is required", http.StatusBadRequest)
		return
	}

	increases, err := database.GetCapitalIncreases(cvr)
	if err != nil {
		pkg.LogError("Failed to get capital increases for CVR " + cvr + ": " + err.Error())
		http.Error(w, "Failed to retrieve capital increases", http.StatusInternalServerError)
		return
	}

	duration := time.Since(startTime).Seconds()
	pkg.LogInfo(fmt.Sprintf("Capital increases retrieved for CVR %s", cvr))

	response := database.IncreaseResponse{
		Stats: database.Stats{
			URL:         r.URL.String(),
			Performance: duration,
		},
		Increases: increases,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
