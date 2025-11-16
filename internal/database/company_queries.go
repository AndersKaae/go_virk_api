package database

import (
	"database/sql"
	"fmt"
	"time"
)

// GetCompanyInfo retrieves basic company information by CVR number
func GetCompanyInfo(cvr string) (*CompanyInfoResponse, error) {
	startTime := time.Now()
	db := GetDB()
	if db == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	query := `
		SELECT 
			v.CVRNummer,
			COALESCE(n.Navn, '') as Name,
			MIN(l.GyldigFra) as Start,
			COALESCE(b.Branchetekst, '') as Branchekode,
			COALESCE(ba.Vejnavn, '') as Vejnavn,
			COALESCE(ba.HusnummerFra, 0) as HusnummerFra,
			COALESCE(ba.Etage, '') as Etage,
			COALESCE(ba.Sidedoer, '') as Sidedoer,
			COALESCE(ba.Postnummer, 0) as Postnummer,
			COALESCE(ba.Postdistrikt, '') as Postdistrikt,
			k.Kontaktoplysning as Website
		FROM Vrvirksomhed v
		LEFT JOIN Navn n ON v.EnhedsNummer = n.EnhedsNummer 
			AND n.Parent = 'Navn' 
			AND n.GyldigTil IS NULL
		LEFT JOIN Livsforloeb l ON v.EnhedsNummer = l.EnhedsNummer
		LEFT JOIN Branche b ON v.EnhedsNummer = b.EnhedsNummer 
			AND b.Parent = 'Hovedbranche' 
			AND b.GyldigTil IS NULL
		LEFT JOIN Beliggenhedsadresse ba ON v.EnhedsNummer = ba.EnhedsNummer 
			AND ba.Parent = 'Beliggenhedsadresse'
			AND ba.GyldigTil IS NULL
		LEFT JOIN Kontaktoplysning k ON v.EnhedsNummer = k.EnhedsNummer
			AND k.Parent = 'Webadresse'
			AND k.GyldigTil IS NULL
		WHERE v.CVRNummer = ?
		GROUP BY v.CVRNummer, n.Navn, b.Branchetekst, ba.Vejnavn, ba.HusnummerFra, ba.Etage, ba.Sidedoer, ba.Postnummer, ba.Postdistrikt, k.Kontaktoplysning
		LIMIT 1
	`

	var cvrNum, name, brancheText, vejnavn, etage, sidedoer, postdistrikt string
	var websiteStr sql.NullString
	var start time.Time
	var husnummerFra, postnummer int

	err := db.QueryRow(query, cvr).Scan(
		&cvrNum, &name, &start, &brancheText,
		&vejnavn, &husnummerFra, &etage, &sidedoer, &postnummer, &postdistrikt, &websiteStr,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("company not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query company info: %w", err)
	}

	// Build address string
	address := buildAddress(vejnavn, husnummerFra, etage, sidedoer, postnummer, postdistrikt)

	// Get website if available
	var website *string
	if websiteStr.Valid && websiteStr.String != "" {
		website = &websiteStr.String
	}

	// Calculate processing time
	processingTime := time.Since(startTime).Seconds()

	response := &CompanyInfoResponse{
		Performance: fmt.Sprintf("%v", processingTime),
		CVR:         cvrNum,
		Name:        name,
		Start:       start.Format("2006-01-02 15:04:05"),
		Branchekode: brancheText,
		Website:     website,
		Address:     address,
	}

	return response, nil
}

// buildAddress constructs an address string from components
func buildAddress(vejnavn string, husnummerFra int, etage string, sidedoer string, postnummer int, postdistrikt string) string {
	if vejnavn == "" {
		return ""
	}

	address := vejnavn
	if husnummerFra > 0 {
		address += fmt.Sprintf(" %d", husnummerFra)
	}
	if etage != "" {
		address += fmt.Sprintf(" %s", etage)
	}
	if sidedoer != "" {
		address += fmt.Sprintf(" %s", sidedoer)
	}
	if postnummer > 0 && postdistrikt != "" {
		address += fmt.Sprintf(", %d %s", postnummer, postdistrikt)
	}

	return address
}
