package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type apiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func writeJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiResponse{Code: code, Data: data})
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(apiResponse{Code: code, Message: msg})
}

// infoHandler returns server info and current database.
func infoHandler(db *sql.DB, cfg PluginConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var version string
		db.QueryRow("SELECT VERSION()").Scan(&version)

		writeJSON(w, 200, map[string]interface{}{
			"host":     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			"user":     cfg.User,
			"database": cfg.Database,
			"version":  version,
			"readonly": cfg.ReadOnly,
		})
	}
}

// databasesHandler returns list of databases.
func databasesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SHOW DATABASES")
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}
		defer rows.Close()

		var databases []string
		for rows.Next() {
			var name string
			rows.Scan(&name)
			databases = append(databases, name)
		}
		writeJSON(w, 200, databases)
	}
}

// tablesHandler returns list of tables for a database.
// Query param: db (database name)
func tablesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbName := r.URL.Query().Get("db")
		if dbName == "" {
			writeError(w, 400, "db parameter is required")
			return
		}

		rows, err := db.Query(fmt.Sprintf("SHOW TABLES FROM `%s`", sanitizeIdentifier(dbName)))
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}
		defer rows.Close()

		var tables []string
		for rows.Next() {
			var name string
			rows.Scan(&name)
			tables = append(tables, name)
		}
		writeJSON(w, 200, tables)
	}
}

// schemaHandler returns column info for a table.
// Query params: db, table
func schemaHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbName := r.URL.Query().Get("db")
		tableName := r.URL.Query().Get("table")
		if dbName == "" || tableName == "" {
			writeError(w, 400, "db and table parameters are required")
			return
		}

		rows, err := db.Query(fmt.Sprintf("SHOW COLUMNS FROM `%s`.`%s`",
			sanitizeIdentifier(dbName), sanitizeIdentifier(tableName)))
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}
		defer rows.Close()

		type Column struct {
			Field   string  `json:"field"`
			Type    string  `json:"type"`
			Null    string  `json:"null"`
			Key     string  `json:"key"`
			Default *string `json:"default"`
			Extra   string  `json:"extra"`
		}

		var columns []Column
		for rows.Next() {
			var col Column
			rows.Scan(&col.Field, &col.Type, &col.Null, &col.Key, &col.Default, &col.Extra)
			columns = append(columns, col)
		}
		writeJSON(w, 200, columns)
	}
}

// indexesHandler returns index info for a table.
// Query params: db, table
// Uses generic row scanning to handle varying SHOW INDEX columns across MySQL versions.
func indexesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbName := r.URL.Query().Get("db")
		tableName := r.URL.Query().Get("table")
		if dbName == "" || tableName == "" {
			writeError(w, 400, "db and table parameters are required")
			return
		}

		rows, err := db.Query(fmt.Sprintf("SHOW INDEX FROM `%s`.`%s`",
			sanitizeIdentifier(dbName), sanitizeIdentifier(tableName)))
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}
		defer rows.Close()

		columns, _ := rows.Columns()
		var results []map[string]interface{}
		for rows.Next() {
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range values {
				valuePtrs[i] = &values[i]
			}
			rows.Scan(valuePtrs...)
			row := make(map[string]interface{})
			for i, col := range columns {
				val := values[i]
				if b, ok := val.([]byte); ok {
					row[col] = string(b)
				} else {
					row[col] = val
				}
			}
			results = append(results, row)
		}
		writeJSON(w, 200, results)
	}
}

// queryHandler executes a SQL query and returns results.
// POST body: {"sql": "SELECT ...", "db": "mydb"}
func queryHandler(db *sql.DB, readOnly bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, 405, "POST required")
			return
		}

		var req struct {
			SQL      string `json:"sql"`
			Database string `json:"db"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, 400, "invalid request body")
			return
		}

		if strings.TrimSpace(req.SQL) == "" {
			writeError(w, 400, "sql is required")
			return
		}

		if readOnly && isWriteSQL(req.SQL) {
			writeError(w, 403, "blocked: read-only mode is enabled")
			return
		}

		// Switch database if specified
		if req.Database != "" {
			if _, err := db.Exec(fmt.Sprintf("USE `%s`", sanitizeIdentifier(req.Database))); err != nil {
				writeError(w, 500, fmt.Sprintf("failed to switch database: %v", err))
				return
			}
		}

		start := time.Now()

		// Determine if it's a query (SELECT/SHOW/DESCRIBE/EXPLAIN) or exec
		trimmed := strings.TrimSpace(strings.ToUpper(req.SQL))
		isQuery := strings.HasPrefix(trimmed, "SELECT") ||
			strings.HasPrefix(trimmed, "SHOW") ||
			strings.HasPrefix(trimmed, "DESCRIBE") ||
			strings.HasPrefix(trimmed, "DESC") ||
			strings.HasPrefix(trimmed, "EXPLAIN")

		if isQuery {
			rows, err := db.Query(req.SQL)
			if err != nil {
				writeError(w, 500, err.Error())
				return
			}
			defer rows.Close()

			columns, _ := rows.Columns()
			var results []map[string]interface{}

			for rows.Next() {
				values := make([]interface{}, len(columns))
				valuePtrs := make([]interface{}, len(columns))
				for i := range values {
					valuePtrs[i] = &values[i]
				}
				rows.Scan(valuePtrs...)

				row := make(map[string]interface{})
				for i, col := range columns {
					val := values[i]
					// Convert []byte to string for JSON
					if b, ok := val.([]byte); ok {
						row[col] = string(b)
					} else {
						row[col] = val
					}
				}
				results = append(results, row)
			}

			duration := time.Since(start)
			writeJSON(w, 200, map[string]interface{}{
				"columns":  columns,
				"rows":     results,
				"count":    len(results),
				"duration": duration.String(),
			})
		} else {
			result, err := db.Exec(req.SQL)
			if err != nil {
				writeError(w, 500, err.Error())
				return
			}

			affected, _ := result.RowsAffected()
			lastID, _ := result.LastInsertId()
			duration := time.Since(start)

			writeJSON(w, 200, map[string]interface{}{
				"affected_rows":  affected,
				"last_insert_id": lastID,
				"duration":       duration.String(),
			})
		}
	}
}

// isWriteSQL checks if a SQL statement is a write operation.
func isWriteSQL(sql string) bool {
	t := strings.TrimSpace(strings.ToUpper(sql))
	writeKeywords := []string{"INSERT", "UPDATE", "DELETE", "DROP", "ALTER", "CREATE", "TRUNCATE", "RENAME", "REPLACE", "GRANT", "REVOKE"}
	for _, k := range writeKeywords {
		if strings.HasPrefix(t, k) {
			return true
		}
	}
	return false
}

// sanitizeIdentifier removes backticks to prevent SQL injection in identifiers.
func sanitizeIdentifier(s string) string {
	return strings.ReplaceAll(s, "`", "")
}

// exportSQLHandler exports table structure and optionally data as SQL.
// Query params: db, table (optional — omit for whole database), mode (structure|all)
func exportSQLHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbName := r.URL.Query().Get("db")
		table := r.URL.Query().Get("table")
		mode := r.URL.Query().Get("mode")
		if dbName == "" {
			writeError(w, 400, "db parameter is required")
			return
		}
		if mode == "" {
			mode = "structure"
		}

		var sb strings.Builder

		if table != "" {
			// ── Single table export ──
			sb.WriteString("-- Table: `" + sanitizeIdentifier(dbName) + "`.`" + sanitizeIdentifier(table) + "`\n")
			sb.WriteString("-- Generated by Shield MySQL\n\n")
			if err := exportTableSQL(db, dbName, table, mode, &sb); err != nil {
				writeError(w, 500, err.Error())
				return
			}
		} else {
			// ── Whole database export ──
			sb.WriteString("-- Database: `" + sanitizeIdentifier(dbName) + "`\n")
			sb.WriteString("-- Generated by Shield MySQL\n")
			sb.WriteString("-- Mode: " + mode + "\n\n")

			rows, err := db.Query(fmt.Sprintf("SHOW TABLES FROM `%s`", sanitizeIdentifier(dbName)))
			if err != nil {
				writeError(w, 500, err.Error())
				return
			}
			defer rows.Close()

			var tables []string
			for rows.Next() {
				var t string
				rows.Scan(&t)
				tables = append(tables, t)
			}

			if len(tables) == 0 {
				writeError(w, 404, "no tables found in database")
				return
			}

			// First pass: CREATE TABLEs (structure only)
			for _, t := range tables {
				if err := exportTableSQL(db, dbName, t, "structure", &sb); err != nil {
					writeError(w, 500, err.Error())
					return
				}
				sb.WriteString("\n")
			}

			// Second pass: INSERT data (after all tables created, so FK references are valid)
			if mode == "all" {
				sb.WriteString("\n-- ══════════════════════════════════════════════\n")
				sb.WriteString("--  Data\n")
				sb.WriteString("-- ══════════════════════════════════════════════\n")
				for _, t := range tables {
					if err := exportTableData(db, dbName, t, &sb); err != nil {
						writeError(w, 500, err.Error())
						return
					}
				}
			}
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(sb.String()))
	}
}

// exportTableSQL exports a single table's DDL (and optionally data) into sb.
func exportTableSQL(db *sql.DB, dbName, table, mode string, sb *strings.Builder) error {
	sb.WriteString("\n-- Table: `" + sanitizeIdentifier(dbName) + "`.`" + sanitizeIdentifier(table) + "`\n")

	// Use SHOW CREATE TABLE for accurate MySQL DDL
	var tableName, createSQL string
	err := db.QueryRow(fmt.Sprintf("SHOW CREATE TABLE `%s`.`%s`",
		sanitizeIdentifier(dbName), sanitizeIdentifier(table))).Scan(&tableName, &createSQL)
	if err != nil {
		return err
	}
	sb.WriteString(createSQL + ";\n")

	if mode == "all" {
		if err := exportTableData(db, dbName, table, sb); err != nil {
			return err
		}
	}
	return nil
}

// exportTableData exports INSERT statements for a single table.
func exportTableData(db *sql.DB, dbName, table string, sb *strings.Builder) error {
	dataRows, err := db.Query(fmt.Sprintf("SELECT * FROM `%s`.`%s`",
		sanitizeIdentifier(dbName), sanitizeIdentifier(table)))
	if err != nil {
		return err
	}
	defer dataRows.Close()

	dataCols, _ := dataRows.Columns()
	colTypes, _ := dataRows.ColumnTypes()
	hasData := false

	for dataRows.Next() {
		if !hasData {
			sb.WriteString("\n-- Data for: `" + sanitizeIdentifier(dbName) + "`.`" + sanitizeIdentifier(table) + "`\n")
			hasData = true
		}
		values := make([]interface{}, len(dataCols))
		valuePtrs := make([]interface{}, len(dataCols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		dataRows.Scan(valuePtrs...)

		sb.WriteString("INSERT INTO `" + sanitizeIdentifier(dbName) + "`.`" + sanitizeIdentifier(table) + "` (")
		for i, c := range dataCols {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString("`" + sanitizeIdentifier(c) + "`")
		}
		sb.WriteString(") VALUES (")
		for i, val := range values {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(formatSQLValue(val, colTypes[i]))
		}
		sb.WriteString(");\n")
	}
	return nil
}

// formatSQLValue converts a Go value to a SQL literal.
func formatSQLValue(val interface{}, ct *sql.ColumnType) string {
	if val == nil {
		return "NULL"
	}
	switch v := val.(type) {
	case []byte:
		s := string(v)
		return "'" + strings.ReplaceAll(s, "'", "''") + "'"
	case string:
		return "'" + strings.ReplaceAll(v, "'", "''") + "'"
	case bool:
		if v {
			return "TRUE"
		}
		return "FALSE"
	case int64, int32, int16, int, float64, float32:
		return fmt.Sprintf("%v", v)
	case time.Time:
		return "'" + v.Format("2006-01-02 15:04:05") + "'"
	default:
		s := fmt.Sprintf("%v", v)
		return "'" + strings.ReplaceAll(s, "'", "''") + "'"
	}
}

// erHandler returns all tables with columns and foreign key relationships for ER diagram.
// Query param: db (database name)
func erHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbName := r.URL.Query().Get("db")
		if dbName == "" {
			writeError(w, 400, "db parameter is required")
			return
		}

		// Get all tables with their columns
		colRows, err := db.Query(`
			SELECT c.TABLE_NAME, c.COLUMN_NAME, c.COLUMN_TYPE,
				CASE WHEN c.COLUMN_KEY = 'PRI' THEN 1 ELSE 0 END AS is_pk
			FROM INFORMATION_SCHEMA.COLUMNS c
			WHERE c.TABLE_SCHEMA = ?
			ORDER BY c.TABLE_NAME, c.ORDINAL_POSITION`, dbName)
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}
		defer colRows.Close()

		type ERColumn struct {
			Name string `json:"name"`
			Type string `json:"type"`
			PK   bool   `json:"pk"`
		}
		type ERTable struct {
			Name    string     `json:"name"`
			Columns []ERColumn `json:"columns"`
		}

		tableMap := make(map[string]*ERTable)
		var tableOrder []string

		for colRows.Next() {
			var tbl, col, dtype string
			var pk bool
			colRows.Scan(&tbl, &col, &dtype, &pk)

			t, ok := tableMap[tbl]
			if !ok {
				t = &ERTable{Name: tbl}
				tableMap[tbl] = t
				tableOrder = append(tableOrder, tbl)
			}
			t.Columns = append(t.Columns, ERColumn{Name: col, Type: dtype, PK: pk})
		}

		var tables []ERTable
		for _, name := range tableOrder {
			tables = append(tables, *tableMap[name])
		}

		// Get foreign key relationships
		fkRows, err := db.Query(`
			SELECT
				kcu.CONSTRAINT_NAME,
				kcu.TABLE_NAME AS from_table,
				kcu.COLUMN_NAME AS from_column,
				kcu.REFERENCED_TABLE_NAME AS to_table,
				kcu.REFERENCED_COLUMN_NAME AS to_column
			FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
			WHERE kcu.TABLE_SCHEMA = ?
				AND kcu.REFERENCED_TABLE_NAME IS NOT NULL
			ORDER BY kcu.TABLE_NAME, kcu.COLUMN_NAME`, dbName)
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}
		defer fkRows.Close()

		type ERRelation struct {
			Constraint string `json:"constraint"`
			FromTable  string `json:"from_table"`
			FromColumn string `json:"from_column"`
			ToTable    string `json:"to_table"`
			ToColumn   string `json:"to_column"`
		}

		var relations []ERRelation
		for fkRows.Next() {
			var rel ERRelation
			fkRows.Scan(&rel.Constraint, &rel.FromTable, &rel.FromColumn, &rel.ToTable, &rel.ToColumn)
			relations = append(relations, rel)
		}

		writeJSON(w, 200, map[string]interface{}{
			"tables":    tables,
			"relations": relations,
		})
	}
}
