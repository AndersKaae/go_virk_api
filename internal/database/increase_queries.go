package database

import (
	"fmt"
	"time"
)

// GetCapitalIncreases retrieves all capital increases for a company
func GetCapitalIncreases(cvr string) ([]CapitalIncrease, error) {
	db := GetDB()
	if db == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	query := `
		SELECT 
			va.Vaerdi,
			va.GyldigFra
		FROM Vrvirksomhed v
		JOIN Attributter a ON v.EnhedsNummer = a.EnhedsNummer
		JOIN Vaerdier va ON a.Id = va.AttributId
		WHERE v.CVRNummer = ?
		AND a.Type = 'KAPITAL'
		AND a.Sekvensnr = 0
		AND a.Vaerditype = 'decimal'
		ORDER BY va.GyldigFra ASC
	`

	rows, err := db.Query(query, cvr)
	if err != nil {
		return nil, fmt.Errorf("failed to query capital increases: %w", err)
	}
	defer rows.Close()

	var increases []CapitalIncrease
	for rows.Next() {
		var capital float64
		var validFrom time.Time

		err := rows.Scan(&capital, &validFrom)
		if err != nil {
			return nil, fmt.Errorf("failed to scan capital increase row: %w", err)
		}

		increases = append(increases, CapitalIncrease{
			Capital:   capital,
			ValidFrom: validFrom.Format("2006-01-02"),
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating capital increase rows: %w", err)
	}

	// Initialize empty slice if nil
	if increases == nil {
		increases = []CapitalIncrease{}
	}

	return increases, nil
}
