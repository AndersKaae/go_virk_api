package database

import (
	"database/sql"
	"fmt"
	"time"
)

// GetTopInvestors retrieves investors ranked by number of companies they own
func GetTopInvestors() (*TopInvestorsResponse, error) {
	startTime := time.Now()
	db := GetDB()
	if db == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	query := `
		SELECT 
			n.Navn as InvestorName,
			vr.CVRNummer as InvestorCVR,
			COUNT(DISTINCT o.VirksomhesNummer) as CompanyCount
		FROM Organisationer o
		JOIN Navn n ON o.DeltagerNummer = n.EnhedsNummer 
			AND n.Parent = 'Navn' 
			AND n.GyldigTil IS NULL
		LEFT JOIN Vrvirksomhed vr ON o.DeltagerNummer = vr.EnhedsNummer
		WHERE o.Hovedtype = 'REGISTER'
		AND o.Navn = 'EJERREGISTER'
		AND o.GyldigTil IS NULL
		GROUP BY o.DeltagerNummer, n.Navn, vr.CVRNummer
		ORDER BY CompanyCount DESC
		LIMIT 100
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query top investors: %w", err)
	}
	defer rows.Close()

	var investors []TopInvestor
	for rows.Next() {
		var investorName string
		var investorCVR sql.NullInt64
		var companyCount int
		
		err := rows.Scan(&investorName, &investorCVR, &companyCount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan investor row: %w", err)
		}

		investor := TopInvestor{
			Name:      investorName,
			Companies: companyCount,
		}

		if investorCVR.Valid {
			cvrInt := int(investorCVR.Int64)
			investor.CVR = &cvrInt
		}

		investors = append(investors, investor)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating investor rows: %w", err)
	}

	// Initialize empty slice if nil
	if investors == nil {
		investors = []TopInvestor{}
	}

	processingTime := time.Since(startTime).Seconds()

	response := &TopInvestorsResponse{
		Performance: processingTime,
		Result:      investors,
	}

	return response, nil
}
