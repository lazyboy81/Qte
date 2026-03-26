package db

import (
	"database/sql"
	"fmt"

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
	Rows *sql.Rows
	Err  error
}

type connector func(string, string) (*sql.DB, error)

func InitializeDBRegistry() {
	connectorRegistry = make(map[string]connector)
	driverRegistry = make(map[string]string)

	registerConnector(DbSqlite, openConnection)
	registerDriver(DbSqlite, "sqlite3")

	//registerConnector(DbOracle, func(datasource string) (*sql.DB, error) {
	//	return sql.Open("godror", datasource)
	//})
	//slog.Debug("added oracle to connectorRegistry")
}

func registerConnector(name string, fn connector) {
	connectorRegistry[name] = fn
}

func registerDriver(name, driver string) {
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

	return open, err
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
		rows, err := executeQuery(query, args)

		// return nil on success, or an error
		return QueryResultMsg{Rows: rows, Err: err}
	}
}

func executeQuery(query string, arg ...any) (*sql.Rows, error) {
	rows, err := database.Query(query, arg)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
