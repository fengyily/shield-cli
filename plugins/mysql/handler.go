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

// queryHandler executes a SQL query and returns results.
// POST body: {"sql": "SELECT ...", "db": "mydb"}
func queryHandler(db *sql.DB) http.HandlerFunc {
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

// sanitizeIdentifier removes backticks to prevent SQL injection in identifiers.
func sanitizeIdentifier(s string) string {
	return strings.ReplaceAll(s, "`", "")
}
