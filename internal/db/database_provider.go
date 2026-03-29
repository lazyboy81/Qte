package db

import (
	"database/sql"
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/godror/godror"
	_ "github.com/mattn/go-sqlite3"
)

const (
	DbSqlite   = "SQLite"
	DbOracle   = "Oracle"
	DbPostgres = "PostgreSQL"
)

var (
	connectorRegistry map[string]connector
	driverRegistry    map[string]string
	database          *sql.DB
)

type ConnectResultMsg struct {
	Err error
}

type QueryResultMsg struct {
	Columns []table.Column
	Rows    []table.Row
	Err     error
}

type connector func(string, string) (*sql.DB, error)

func InitializeDBRegistry() {
	connectorRegistry = make(map[string]connector)
	driverRegistry = make(map[string]string)

	registerDriverAndConnector(DbSqlite, "sqlite3", openConnection)
}

func registerDriverAndConnector(name, driver string, fn connector) {
	connectorRegistry[name] = fn
	driverRegistry[name] = driver
}

func openConnection(driver, datasource string) (*sql.DB, error) {
	open, err := sql.Open(driver, datasource)
	if err != nil {
		return nil, err
	}

	err = open.Ping()
	if err != nil {
		return nil, err
	}

	return open, nil
}

func ConnectToDatabase(selectedProvider, databaseUrl string) tea.Cmd {
	return func() tea.Msg {
		var err error
		database, err = connect(selectedProvider, databaseUrl)

		// return nil on success, or an error
		return ConnectResultMsg{Err: err}
	}
}

func connect(name, dsn string) (*sql.DB, error) {
	fn, connectorOk := connectorRegistry[name]
	if !connectorOk {
		return nil, fmt.Errorf("unknown db: %s", name)
	}

	driver, driverOk := driverRegistry[name]
	if !driverOk {
		return nil, fmt.Errorf("unknown db driver: %s", driver)
	}

	return fn(driver, dsn)
}

func ExecuteQueryInDatabase(query string, args ...any) tea.Cmd {
	return func() tea.Msg {
		queryResult, err := executeSelectQuery(query, args)
		defer queryResult.Close()

		if err != nil {
			// return nil on success, or an error
			return QueryResultMsg{Columns: nil, Rows: nil, Err: err}
		}
		columns, columnCount, err := extractHeaders(queryResult)
		if err != nil {
			// return nil on success, or an error
			return QueryResultMsg{Columns: nil, Rows: nil, Err: err}
		}
		rows, err := extractRows(queryResult, columnCount)
		if err != nil {
			// return nil on success, or an error
			return QueryResultMsg{Columns: nil, Rows: nil, Err: err}
		}

		return QueryResultMsg{Columns: columns, Rows: rows, Err: nil}
	}
}

func executeUpdateQuery(query string, arg ...any) (sql.Result, error) {
	tnx, err := database.Begin()
	if err != nil {
		return nil, err
	}

	result, err := tnx.Exec(query, arg)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func executeSelectQuery(query string, arg ...any) (*sql.Rows, error) {
	rows, err := database.Query(query)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func extractHeaders(rows *sql.Rows) ([]table.Column, int, error) {
	headers := make([]table.Column, 0)

	columns, err := rows.Columns()
	if err != nil {
		return headers, 0, err
	}

	for _, column := range columns {
		headers = append(headers, table.Column{
			Title: column,
			Width: len(column),
		})
	}

	return headers, len(headers), nil
}

func extractRows(rows *sql.Rows, columnCount int) ([]table.Row, error) {
	var result []table.Row

	for rows.Next() {
		values := make([]any, columnCount)

		valuesPtr := make([]any, columnCount)
		for i := 0; i < columnCount; i++ {
			valuesPtr[i] = &values[i]
		}

		err := rows.Scan(valuesPtr...)
		if err != nil {
			return result, err
		}

		row := make(table.Row, columnCount)
		for i, value := range values {
			// we use indexes here to avoid a wierd bug which by using append(), table.renderRow() would think we have
			// more columns that what we actually wrote into table.Row, probably due to underlying slice
			if value == nil {
				row[i] = "NULL"
			} else {
				row[i] = fmt.Sprintf("%v", value)
			}
		}

		result = append(result, row)
	}

	return result, nil
}
