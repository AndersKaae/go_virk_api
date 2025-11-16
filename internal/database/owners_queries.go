package database

import (
	"database/sql"
	"fmt"
	"time"
)

// GetOwners retrieves ownership information for a company by CVR
func GetOwners(cvr string) (*OwnersResponse, error) {
	startTime := time.Now()
	db := GetDB()
	if db == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	query := `
		SELECT 
			n.Navn as OwnerName,
			vr.CVRNummer as OwnerCVR,
			vd.Vaerdi as OwnershipPercent
		FROM Vrvirksomhed v
		JOIN Attributter a ON v.EnhedsNummer = a.EnhedsNummer
		JOIN Vaerdier vd ON a.Id = vd.AttributId
		JOIN Navn n ON a.DeltagerNummer = n.EnhedsNummer AND n.Parent = 'Navn' AND n.GyldigTil IS NULL
		LEFT JOIN Vrvirksomhed vr ON a.DeltagerNummer = vr.EnhedsNummer
		WHERE v.CVRNummer = ?
		AND a.Type = 'EJERANDEL_PROCENT'
		AND vd.GyldigTil IS NULL
		ORDER BY CAST(vd.Vaerdi AS DECIMAL(10,5)) DESC
	`

	rows, err := db.Query(query, cvr)
	if err != nil {
		return nil, fmt.Errorf("failed to query owners: %w", err)
	}
	defer rows.Close()

	var owners []Owner
	for rows.Next() {
		var ownerName, ownershipPercent string
		var ownerCVR sql.NullInt64
		
		err := rows.Scan(&ownerName, &ownerCVR, &ownershipPercent)
		if err != nil {
			return nil, fmt.Errorf("failed to scan owner row: %w", err)
		}

		owner := Owner{
			Name:  ownerName,
			Value: ownershipPercent,
		}

		if ownerCVR.Valid {
			cvrInt := int(ownerCVR.Int64)
			owner.CVR = &cvrInt
		}

		owners = append(owners, owner)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating owner rows: %w", err)
	}

	// Initialize empty slice if nil
	if owners == nil {
		owners = []Owner{}
	}

	processingTime := time.Since(startTime).Seconds()

	response := &OwnersResponse{
		Performance: processingTime,
		Result:      owners,
	}

	return response, nil
}
