package database

import (
	"fmt"
	"time"
)

// SearchCompanies searches for companies by name
func SearchCompanies(query string) (*SearchResponse, error) {
	startTime := time.Now()
	db := GetDB()
	if db == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	sqlQuery := `
		SELECT DISTINCT
			v.CVRNummer,
			n.Navn
		FROM Vrvirksomhed v
		JOIN Navn n ON v.EnhedsNummer = n.EnhedsNummer
			AND n.Parent = 'Navn'
			AND n.GyldigTil IS NULL
		WHERE n.Navn LIKE ?
		ORDER BY n.Navn
		LIMIT 100
	`

	searchPattern := "%" + query + "%"
	rows, err := db.Query(sqlQuery, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search companies: %w", err)
	}
	defer rows.Close()

	var results []SearchResult
	for rows.Next() {
		var cvr int
		var name string
		
		err := rows.Scan(&cvr, &name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan search result: %w", err)
		}

		results = append(results, SearchResult{
			CVR:  fmt.Sprintf("%d", cvr),
			Name: name,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating search results: %w", err)
	}

	// Initialize empty slice if nil
	if results == nil {
		results = []SearchResult{}
	}

	processingTime := time.Since(startTime).Seconds()

	response := &SearchResponse{
		Performance: processingTime,
		Result:      results,
	}

	return response, nil
}
