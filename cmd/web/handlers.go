package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/jackc/pgconn"
)

type TableNames struct {
	SourceTables []string `json:"sourceTables"`
	TargetTables []string `json:"targetTables"`
}

type SyncDataRequest struct {
	Table string `json:"table"`
	Type  string `json:"type"`
}

func (c *Config) ConnectToDBHandler(w http.ResponseWriter, r *http.Request) {
	var dbConfigs DBConfigs
	if err := json.NewDecoder(r.Body).Decode(&dbConfigs); err != nil {
		c.errorJSON(w, errors.New("Invalid Request Body"), http.StatusBadRequest)
		return
	}

	sourceDSN := dbConfigs.Source.BuildDSN()
	targetDSN := dbConfigs.Target.BuildDSN()

	sourceConn, err := InitDB(sourceDSN)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "28P01":
				c.errorJSON(w, errors.New("invalid username or password"), http.StatusInternalServerError)
			default:
				c.errorJSON(w, errors.New("couldn't connect to the source database"), http.StatusInternalServerError)
			}
		}
		return
	}

	targetConn, err := InitDB(targetDSN)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "28P01":
				c.errorJSON(w, errors.New("invalid username or password"), http.StatusInternalServerError)
			default:
				c.errorJSON(w, errors.New("couldn't connect to the target database"), http.StatusInternalServerError)
			}
		}
		return
	}

	c.SourceDb = sourceConn
	c.TargetDB = targetConn

	sourceTables, err := FetchTables(c.SourceDb)
	targetTables, err := FetchTables(c.TargetDB)

	tables := TableNames{
		SourceTables: sourceTables,
		TargetTables: targetTables,
	}

	jsonResonse := JsonResponse{
		Message: "Connected Succesfully",
		Data:    tables,
	}
	c.writeJSON(w, 200, jsonResonse)
}

func (c *Config) FetchDBTables(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("type")

	if t == "source" {
		if c.SourceDb == nil {
			c.errorJSON(w, errors.New("Source Database not connected"), http.StatusInternalServerError)
			return
		}
		tables, err := FetchTables(c.SourceDb)
		if err != nil {
			c.errorJSON(w, errors.New("Couldn't fetch tables"), http.StatusInternalServerError)
			return
		}
		c.writeJSON(w, 200, tables)
		return
	} else if t == "target" {
		if c.TargetDB == nil {
			c.errorJSON(w, errors.New("Target Database not connected"), http.StatusInternalServerError)
			return
		}
		tables, err := FetchTables(c.TargetDB)
		if err != nil {
			c.errorJSON(w, errors.New("Couldn't fetch tables"), http.StatusInternalServerError)
			return
		}
		c.writeJSON(w, 200, tables)
		return
	} else {
		c.errorJSON(w, errors.New("You must provide type (source/target)"), http.StatusBadRequest)
	}
}

func (c *Config) FetchRowsForTable(w http.ResponseWriter, r *http.Request) {
	table := r.URL.Query().Get("table")
	t := r.URL.Query().Get("type")

	if t == "source" {
		if c.SourceDb == nil {
			c.errorJSON(w, errors.New("Source Database not connected"), http.StatusInternalServerError)
			return
		}
		rows, err := FetchRecordsFromTable(c.SourceDb, table)
		if err != nil {
			c.errorJSON(w, errors.New("Couldn't fetch rows"), http.StatusInternalServerError)
			return
		}
		c.writeJSON(w, 200, rows)
		return
	} else if t == "target" {
		if c.TargetDB == nil {
			c.errorJSON(w, errors.New("Target Database not connected"), http.StatusInternalServerError)
			return
		}
		rows, err := FetchRecordsFromTable(c.TargetDB, table)
		if err != nil {
			c.errorJSON(w, errors.New("Couldn't fetch rows"), http.StatusInternalServerError)
			return
		}
		c.writeJSON(w, 200, rows)
		return
	} else {
		c.errorJSON(w, errors.New("You must provide type (source/target)"), http.StatusBadRequest)
	}
}

func (c *Config) SynchronizeData(w http.ResponseWriter, r *http.Request) {
	// get the request body
	var syncDataRequest []SyncDataRequest

	if err := json.NewDecoder(r.Body).Decode(&syncDataRequest); err != nil {
		c.errorJSON(w, errors.New("Invalid Request Body"), http.StatusBadRequest)
		return
	}

	if len(syncDataRequest) != 2 {
		c.errorJSON(w, errors.New("There should be precise two requests"), http.StatusBadRequest)
		return
	}

	source := SyncDataRequest{}
	target := SyncDataRequest{}

	for _, v := range syncDataRequest {
		if v.Type == "source" {
			source = v
		} else if v.Type == "target" {
			target = v
		} else {
			c.errorJSON(w, errors.New("You must provide type (source/target)"), http.StatusBadRequest)
		}
	}

	if source.Table == "" || target.Table == "" {
		c.errorJSON(w, errors.New("You must provide table name"), http.StatusBadRequest)
		return
	}

	isCompatible, err := c.compareTableSchemas(source, target)

	if err != nil {
		c.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	if !isCompatible {
		c.errorJSON(w, errors.New("Data is not compatible"), http.StatusBadRequest)
		return
	}

	sourcePrimaryKey, err := GetPrimaryKey(c.SourceDb, source.Table)
	if err != nil {
		c.errorJSON(w, errors.New("Failed to fetch source primary key"), http.StatusBadRequest)
		return
	}

	targetPrimaryKey, err := GetPrimaryKey(c.TargetDB, target.Table)
	if err != nil {
		c.errorJSON(w, errors.New("Failed to fetch target primary key"), http.StatusBadRequest)
		return
	}

	if targetPrimaryKey != sourcePrimaryKey {
		c.errorJSON(w, errors.New("Primary keys are not matching"), http.StatusBadRequest)
		return
	}

	// sync the data
	sourceData, err := fetchTableData(c.SourceDb, source.Table)
	targetData, err := fetchTableData(c.TargetDB, target.Table)

	insertRows, updateRows, deleteRows := compareRows(sourceData, targetData, sourcePrimaryKey)

	err = synchronizeTables(c.TargetDB, target.Table, targetPrimaryKey, insertRows, updateRows, deleteRows)
	if err != nil {
		c.errorJSON(w, errors.New("Failed to synchronize data"), http.StatusBadRequest)
		return
	}

	c.writeJSON(w, 201, "Data Synced Successfully")
}

func (c *Config) CheckDbHealth(w http.ResponseWriter, r *http.Request) {
	if c.SourceDb == nil && c.TargetDB == nil {
		c.errorJSON(w, errors.New("Database not connected"), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Service is healthy"))
	return
}

func (c *Config) compareTableSchemas(source, target SyncDataRequest) (bool, error) {

	targetSchema, err := GetTableSchema(c.TargetDB, target.Table)
	if err != nil {
		return false, errors.New("Couldn't fetch target table schema")
	}

	sourceSchema, err := GetTableSchema(c.SourceDb, source.Table)
	if err != nil {
		return false, errors.New("Couldn't fetch source table schema")
	}

	if len(sourceSchema) != len(targetSchema) {
		return false, nil
	}

	for i := 0; i < len(sourceSchema); i++ {
		if sourceSchema[i].ColumnName != targetSchema[i].ColumnName || sourceSchema[i].DataType != targetSchema[i].DataType || sourceSchema[i].IsNullable != targetSchema[i].IsNullable {
			return false, nil
		}
	}
	return true, nil
}

func compareRows(sourceRows, targetRows []map[string]interface{}, primaryKey string) (insertRows, updateRows, deleteRows []map[string]interface{}) {
	targetMap := make(map[interface{}]map[string]interface{})
	for _, row := range targetRows {
		targetMap[row[primaryKey]] = row
	}

	for _, srcRow := range sourceRows {
		if targetRow, exists := targetMap[srcRow[primaryKey]]; exists {
			if !reflect.DeepEqual(srcRow, targetRow) {
				updateRows = append(updateRows, srcRow)
			}
			delete(targetMap, srcRow[primaryKey])
		} else {
			insertRows = append(insertRows, srcRow)
		}
	}

	for _, tgtRow := range targetMap {
		deleteRows = append(deleteRows, tgtRow)
	}

	return insertRows, updateRows, deleteRows
}

func synchronizeTables(db *sql.DB, tableName string, primaryKey string, insertRows, updateRows, deleteRows []map[string]interface{}) error {
	for _, row := range insertRows {
		columns, values := prepareInsertQuery(row)
		query := fmt.Sprintf("INSERT INTO  %s (%s) OVERRIDING SYSTEM VALUE VALUES (%s)", tableName, columns, values)
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("insert failed: %w", err)
		}
	}

	for _, row := range updateRows {
		setClause, whereClause := prepareUpdateQuery(row, primaryKey)
		query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, setClause, whereClause)
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("update failed: %w", err)
		}
	}

	for _, row := range deleteRows {
		whereClause := prepareDeleteQuery(row, primaryKey)
		query := fmt.Sprintf("DELETE FROM %s WHERE %s", tableName, whereClause)
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("delete failed: %w", err)
		}
	}

	return nil
}

func prepareInsertQuery(row map[string]interface{}) (string, string) {
	columns := []string{}
	values := []string{}
	for col, val := range row {
		columns = append(columns, col)
		values = append(values, fmt.Sprintf("'%v'", val))
	}
	return strings.Join(columns, ", "), strings.Join(values, ", ")
}

func prepareUpdateQuery(row map[string]interface{}, primaryKey string) (string, string) {
	setClauses := []string{}
	var whereClause string
	for col, val := range row {
		if col == primaryKey {
			whereClause = fmt.Sprintf("%s = '%v'", col, val)
		} else {
			setClauses = append(setClauses, fmt.Sprintf("%s = '%v'", col, val))
		}
	}
	return strings.Join(setClauses, ", "), whereClause
}

func prepareDeleteQuery(row map[string]interface{}, primaryKey string) string {
	return fmt.Sprintf("%s = '%v'", primaryKey, row[primaryKey])
}
