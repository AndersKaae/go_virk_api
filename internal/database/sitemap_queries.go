package database

import (
	"fmt"
	"time"
)

// GetSitemap retrieves all companies for sitemap generation
func GetSitemap() ([]SitemapEntry, error) {
	db := GetDB()
	if db == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	query := `
		SELECT 
			v.CVRNummer,
			v.SidstOpdateret
		FROM Vrvirksomhed v
		WHERE v.CVRNummer IS NOT NULL
		ORDER BY v.SidstOpdateret ASC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query sitemap: %w", err)
	}
	defer rows.Close()

	var entries []SitemapEntry
	for rows.Next() {
		var cvr int
		var lastUpdated time.Time
		
		err := rows.Scan(&cvr, &lastUpdated)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sitemap row: %w", err)
		}

		// Format URL and timestamp
		loc := fmt.Sprintf("https://www.enhjorning.bot/company/%d", cvr)
		lastMod := lastUpdated.Format("2006-01-02T15:04:05.000-07:00")

		entries = append(entries, SitemapEntry{
			Loc:     loc,
			LastMod: lastMod,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sitemap rows: %w", err)
	}

	// Initialize empty slice if nil
	if entries == nil {
		entries = []SitemapEntry{}
	}

	return entries, nil
}
