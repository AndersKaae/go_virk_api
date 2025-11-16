package database

import (
	"database/sql"
	"fmt"
	"time"
)

// GetFeed retrieves companies with capital increases for the feed endpoint
func GetFeed(page int, pageSize int) (*FeedResponse, error) {
	startTime := time.Now()
	db := GetDB()
	if db == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	offset := (page - 1) * pageSize

	// Optimized query: Get all data in one go
	// Matches Python query logic: exclude companies with only formation capital
	query := `
		SELECT DISTINCT
			v.CVRNummer,
			v.EnhedsNummer,
			n.Navn as CompanyName,
			b.Branchekode,
			(
				SELECT COUNT(DISTINCT o.DeltagerNummer)
				FROM Organisationer o
				WHERE o.VirksomhesNummer = v.EnhedsNummer
				AND o.Hovedtype = 'REGISTER'
				AND o.Navn = 'EJERREGISTER'
				AND o.GyldigTil IS NULL
			) as OwnerCount
		FROM Vrvirksomhed v
		JOIN Virksomhedsform vf ON v.EnhedsNummer = vf.EnhedsNummer 
			AND vf.VirksomhedsformKode IN (60, 80) 
			AND vf.GyldigTil IS NULL
		LEFT JOIN Navn n ON v.EnhedsNummer = n.EnhedsNummer 
			AND n.Parent = 'Navn' 
			AND n.GyldigTil IS NULL
		LEFT JOIN Branche b ON v.EnhedsNummer = b.EnhedsNummer 
			AND b.Parent = 'Hovedbranche' 
			AND b.GyldigTil IS NULL
		JOIN Attributter a ON v.EnhedsNummer = a.EnhedsNummer 
			AND a.Sekvensnr = 0 
			AND a.Vaerditype = 'decimal'
		JOIN Vaerdier va ON a.Id = va.AttributId
		JOIN (
			SELECT l.EnhedsNummer, MIN(l.GyldigFra) as CompanyStart
			FROM Livsforloeb l
			GROUP BY l.EnhedsNummer
		) lc ON v.EnhedsNummer = lc.EnhedsNummer
		WHERE EXISTS (
			SELECT 1
			FROM Attributter a2
			JOIN Vaerdier va2 ON a2.Id = va2.AttributId
			WHERE a2.EnhedsNummer = v.EnhedsNummer
			AND a2.Sekvensnr = 0
			AND a2.Vaerditype = 'decimal'
			AND va2.GyldigFra != lc.CompanyStart
		)
		GROUP BY v.CVRNummer, v.EnhedsNummer, n.Navn, b.Branchekode
		ORDER BY MAX(va.GyldigFra) DESC
		LIMIT ? OFFSET ?
	`

	rows, err := db.Query(query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query feed: %w", err)
	}
	defer rows.Close()

	var companies []struct {
		CVR          int
		EnhedsNummer int
		Name         sql.NullString
		BusinessCode sql.NullInt64
		OwnerCount   int
	}

	for rows.Next() {
		var company struct {
			CVR          int
			EnhedsNummer int
			Name         sql.NullString
			BusinessCode sql.NullInt64
			OwnerCount   int
		}
		err := rows.Scan(&company.CVR, &company.EnhedsNummer, &company.Name, &company.BusinessCode, &company.OwnerCount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		companies = append(companies, company)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	// Batch fetch all capital data for all companies in ONE query
	capitalMap, err := batchFetchCapital(db, companies)
	if err != nil {
		return nil, fmt.Errorf("failed to batch fetch capital: %w", err)
	}

	// Build feed items
	feed := make([]FeedItem, 0, len(companies))
	for _, company := range companies {
		feedItem := buildFeedItemWithCapital(company.CVR, company.EnhedsNummer, company.Name, company.BusinessCode, company.OwnerCount, capitalMap[company.EnhedsNummer])
		feed = append(feed, feedItem)
	}

	// Calculate processing time
	processingTime := time.Since(startTime).Seconds()

	response := &FeedResponse{
		Stats: Stats{
			URL:         fmt.Sprintf("/api/v1/feed?page=%d", page),
			Performance: processingTime,
		},
		Feed: feed,
	}

	return response, nil
}

// batchFetchCapital fetches all capital increases for multiple companies in a single query
func batchFetchCapital(db *sql.DB, companies []struct {
	CVR          int
	EnhedsNummer int
	Name         sql.NullString
	BusinessCode sql.NullInt64
	OwnerCount   int
}) (map[int][]CapitalIncrease, error) {
	if len(companies) == 0 {
		return make(map[int][]CapitalIncrease), nil
	}

	// Build list of EnhedsNummer for IN clause
	enhedsNummers := make([]int, len(companies))
	for i, company := range companies {
		enhedsNummers[i] = company.EnhedsNummer
	}

	// Create placeholders for IN clause
	placeholders := make([]string, len(enhedsNummers))
	args := make([]interface{}, len(enhedsNummers))
	for i, num := range enhedsNummers {
		placeholders[i] = "?"
		args[i] = num
	}

	// Batch query for all capital data
	query := fmt.Sprintf(`
		SELECT a.EnhedsNummer, va.Vaerdi, va.GyldigFra
		FROM Attributter a
		JOIN Vaerdier va ON a.Id = va.AttributId
		WHERE a.EnhedsNummer IN (%s)
		AND a.Sekvensnr = 0
		AND a.Vaerditype = 'decimal'
		ORDER BY a.EnhedsNummer, va.GyldigFra ASC
	`, string(placeholders[0]))

	// Replace first placeholder with all placeholders
	query = fmt.Sprintf(`
		SELECT a.EnhedsNummer, va.Vaerdi, va.GyldigFra
		FROM Attributter a
		JOIN Vaerdier va ON a.Id = va.AttributId
		WHERE a.EnhedsNummer IN (%s)
		AND a.Sekvensnr = 0
		AND a.Vaerditype = 'decimal'
		ORDER BY a.EnhedsNummer, va.GyldigFra ASC
	`, joinPlaceholders(len(enhedsNummers)))

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to batch query capital: %w", err)
	}
	defer rows.Close()

	// Group capital increases by EnhedsNummer
	capitalMap := make(map[int][]CapitalIncrease)
	for rows.Next() {
		var enhedsNummer int
		var capital float64
		var gyldigFra time.Time
		err := rows.Scan(&enhedsNummer, &capital, &gyldigFra)
		if err != nil {
			return nil, fmt.Errorf("failed to scan capital row: %w", err)
		}
		capitalMap[enhedsNummer] = append(capitalMap[enhedsNummer], CapitalIncrease{
			Capital:   capital,
			ValidFrom: gyldigFra.Format("2006-01-02"),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating capital rows: %w", err)
	}

	return capitalMap, nil
}

// joinPlaceholders creates a comma-separated string of ? placeholders
func joinPlaceholders(count int) string {
	if count == 0 {
		return ""
	}
	result := ""
	for i := 0; i < count; i++ {
		if i > 0 {
			result += ","
		}
		result += "?"
	}
	return result
}

// buildFeedItemWithCapital constructs a FeedItem with all pre-fetched data
func buildFeedItemWithCapital(cvr int, enhedsNummer int, name sql.NullString, businessCode sql.NullInt64, ownerCount int, increases []CapitalIncrease) FeedItem {
	item := FeedItem{
		CVR:            fmt.Sprintf("%d", cvr),
		NumberOfOwners: ownerCount,
		Increases:      increases,
	}

	if name.Valid {
		item.Name = name.String
	}
	if businessCode.Valid {
		item.BusinessCode = int(businessCode.Int64)
	}

	// Ensure increases is not nil
	if item.Increases == nil {
		item.Increases = []CapitalIncrease{}
	}

	return item
}
