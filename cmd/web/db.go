package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type ColumnSchema struct {
	ColumnName string
	DataType   string
	IsNullable string
}

func InitDB(dsn string) (*sql.DB, error) {
	conn, err := connectToDB(dsn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func connectToDB(dsn string) (*sql.DB, error) {
	// for retrying
	counts := 0
	for {
		conn, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres couldn't connect.. Retrying..")
		} else {
			log.Println("connected to databasej")
			return conn, nil
		}

		if counts > 10 {
			log.Println(err)
			return nil, err
		}

		counts += 1

		time.Sleep(1 * time.Second)
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func FetchTables(db *sql.DB) ([]string, error) {
	query := `
		SELECT tablename FROM pg_tables where schemaname = 'public'
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tables, nil
}

func FetchRecordsFromTable(db *sql.DB, table string) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s", table)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		// Create a slice of interface{} to hold column values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan the row into the value pointers
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// Create a map for the row
		row := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			rawValue := values[i]

			// Convert raw bytes to a usable type
			b, ok := rawValue.([]byte)
			if ok {
				v = string(b)
			} else {
				v = rawValue
			}
			row[col] = v
		}

		// Append the row map to results
		results = append(results, row)
	}

	return results, nil

}

func GetTableSchema(db *sql.DB, table string) ([]ColumnSchema, error) {
	query := fmt.Sprintf("select column_name, data_type, is_nullable FROM information_schema.columns where table_name = '%s'", table)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schema []ColumnSchema
	for rows.Next() {
		var column ColumnSchema
		if err := rows.Scan(&column.ColumnName, &column.DataType, &column.IsNullable); err != nil {
			return nil, err
		}
		schema = append(schema, column)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	fmt.Println("Schema", schema)
	return schema, nil
}

func fetchTableData(db *sql.DB, tableName string) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// Map the values to column names
		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = values[i]
		}
		results = append(results, row)
	}
	return results, nil
}

func FetchChanges(db *sql.DB, lastXmin int) ([]map[string]interface{}, error) {
	rows, err := db.Query("SELECT * from change_log WHERE system_xmin > $1 ORDER BY change_timestamp", lastXmin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var changes []map[string]interface{}
	for rows.Next() {
		var id int
		var tableName, operation string
		var rowData string
		var systemXmin int64
		if err := rows.Scan(&id, &tableName, &operation, &rowData, &systemXmin); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		change := map[string]interface{}{
			"id":          id,
			"table_name":  tableName,
			"operation":   operation,
			"row_data":    rowData,
			"system_xmin": systemXmin,
		}
		changes = append(changes, change)
	}

	return changes, nil
}

func GetPrimaryKey(db *sql.DB, tableName string) (string, error) {
	query := `
		SELECT kcu.column_name
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu
		  ON tc.constraint_name = kcu.constraint_name
		WHERE tc.table_name = $1
		  AND tc.constraint_type = 'PRIMARY KEY';
	`
	row := db.QueryRow(query, tableName)
	var primaryKey string
	if err := row.Scan(&primaryKey); err != nil {
		return "", err
	}
	return primaryKey, nil
}
