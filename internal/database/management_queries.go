package database

import (
	"fmt"
	"time"
)

// GetManagement retrieves board and management information for a company by CVR
func GetManagement(cvr string) (*ManagementResponse, error) {
	startTime := time.Now()
	db := GetDB()
	if db == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	query := `
		SELECT n.Navn, o.Rolle, o.Navn as OrgType
		FROM Vrvirksomhed v
		JOIN (
			SELECT o1.*
			FROM Organisationer o1
			JOIN (
				SELECT VirksomhesNummer, DeltagerNummer, Hovedtype, Navn, MAX(GyldigFra) as MaxGyldigFra
				FROM Organisationer
				WHERE Hovedtype = 'LEDELSESORGAN'
				AND GyldigTil IS NULL
				GROUP BY VirksomhesNummer, DeltagerNummer, Hovedtype, Navn
			) o2 ON o1.VirksomhesNummer = o2.VirksomhesNummer
				AND o1.DeltagerNummer = o2.DeltagerNummer
				AND o1.Hovedtype = o2.Hovedtype
				AND o1.Navn = o2.Navn
				AND o1.GyldigFra = o2.MaxGyldigFra
				AND o1.GyldigTil IS NULL
		) o ON v.EnhedsNummer = o.VirksomhesNummer
		JOIN Navn n ON o.DeltagerNummer = n.EnhedsNummer AND n.Parent = 'Navn'
		WHERE v.CVRNummer = ?
		AND n.GyldigTil IS NULL
		ORDER BY n.Navn
	`

	rows, err := db.Query(query, cvr)
	if err != nil {
		return nil, fmt.Errorf("failed to query management: %w", err)
	}
	defer rows.Close()

	var board []ManagementPerson
	var management []ManagementPerson

	for rows.Next() {
		var personName, rolle, orgType string
		err := rows.Scan(&personName, &rolle, &orgType)
		if err != nil {
			return nil, fmt.Errorf("failed to scan management row: %w", err)
		}

		person := ManagementPerson{
			Name: personName,
			Role: rolle,
		}

		// Determine if person is board or management based on orgType
		// "Bestyrelse" = Board, "Direktion" = Management
		if orgType == "Bestyrelse" {
			board = append(board, person)
		} else if orgType == "Direktion" {
			management = append(management, person)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating management rows: %w", err)
	}

	// Initialize empty slices if nil
	if board == nil {
		board = []ManagementPerson{}
	}
	if management == nil {
		management = []ManagementPerson{}
	}

	processingTime := time.Since(startTime).Seconds()

	response := &ManagementResponse{
		Performance: processingTime,
		Result: ManagementResult{
			Board:      board,
			Management: management,
		},
	}

	return response, nil
}
