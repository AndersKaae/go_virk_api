package handlers

import (
	"encoding/json"
	"github.com/AndersKaae/go_virk_api/internal/database"
	"github.com/AndersKaae/go_virk_updater/pkg"
	"net/http"
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
