package database

import (
	"database/sql"
	"fmt"
)

// TableStats represents statistics for a single table
type TableStats struct {
	TableName   string  `json:"tableName"`
	RowCount    int64   `json:"rowCount"`
	DataSize    int64   `json:"dataSize"`    // in bytes
	IndexSize   int64   `json:"indexSize"`   // in bytes
	TotalSize   int64   `json:"totalSize"`   // in bytes
	DataSizeMB  float64 `json:"dataSizeMB"`
	IndexSizeMB float64 `json:"indexSizeMB"`
	TotalSizeMB float64 `json:"totalSizeMB"`
}

// DatabaseStats represents overall database statistics
type DatabaseStats struct {
	DatabaseName string       `json:"databaseName"`
	Tables       []TableStats `json:"tables"`
	TotalRows    int64        `json:"totalRows"`
	TotalDataSize int64       `json:"totalDataSize"`     // in bytes
	TotalIndexSize int64      `json:"totalIndexSize"`    // in bytes
	TotalSize     int64       `json:"totalSize"`         // in bytes
	TotalSizeMB   float64     `json:"totalSizeMB"`
	TotalSizeGB   float64     `json:"totalSizeGB"`
	TableCount    int         `json:"tableCount"`
}

// GetDatabaseStats retrieves statistics for all tables in the database
func GetDatabaseStats() (*DatabaseStats, error) {
	db := GetDB()
	if db == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	// Get database name from environment
	dbName := getEnv("DB_NAME", "virk2")

	query := `
		SELECT 
			TABLE_NAME,
			TABLE_ROWS,
			DATA_LENGTH,
			INDEX_LENGTH,
			DATA_LENGTH + INDEX_LENGTH AS TOTAL_SIZE
		FROM information_schema.TABLES
		WHERE TABLE_SCHEMA = ?
		ORDER BY TOTAL_SIZE DESC
	`

	rows, err := db.Query(query, dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to query table stats: %w", err)
	}
	defer rows.Close()

	stats := &DatabaseStats{
		DatabaseName: dbName,
		Tables:       []TableStats{},
	}

	for rows.Next() {
		var table TableStats
		err := rows.Scan(
			&table.TableName,
			&table.RowCount,
			&table.DataSize,
			&table.IndexSize,
			&table.TotalSize,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Convert to MB for readability
		table.DataSizeMB = float64(table.DataSize) / (1024 * 1024)
		table.IndexSizeMB = float64(table.IndexSize) / (1024 * 1024)
		table.TotalSizeMB = float64(table.TotalSize) / (1024 * 1024)

		stats.Tables = append(stats.Tables, table)
		stats.TotalRows += table.RowCount
		stats.TotalDataSize += table.DataSize
		stats.TotalIndexSize += table.IndexSize
		stats.TotalSize += table.TotalSize
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	// Calculate totals in MB and GB
	stats.TotalSizeMB = float64(stats.TotalSize) / (1024 * 1024)
	stats.TotalSizeGB = float64(stats.TotalSize) / (1024 * 1024 * 1024)
	stats.TableCount = len(stats.Tables)

	return stats, nil
}

// GetTableStats retrieves statistics for a specific table
func GetTableStats(tableName string) (*TableStats, error) {
	db := GetDB()
	if db == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	dbName := getEnv("DB_NAME", "virk2")

	query := `
		SELECT 
			TABLE_NAME,
			TABLE_ROWS,
			DATA_LENGTH,
			INDEX_LENGTH,
			DATA_LENGTH + INDEX_LENGTH AS TOTAL_SIZE
		FROM information_schema.TABLES
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
	`

	var table TableStats
	err := db.QueryRow(query, dbName, tableName).Scan(
		&table.TableName,
		&table.RowCount,
		&table.DataSize,
		&table.IndexSize,
		&table.TotalSize,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("table not found: %s", tableName)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query table stats: %w", err)
	}

	// Convert to MB
	table.DataSizeMB = float64(table.DataSize) / (1024 * 1024)
	table.IndexSizeMB = float64(table.IndexSize) / (1024 * 1024)
	table.TotalSizeMB = float64(table.TotalSize) / (1024 * 1024)

	return &table, nil
}
