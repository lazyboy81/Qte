package db

import (
	"database/sql"
	"query-terminal/internal/store"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/godror/godror"
	_ "github.com/mattn/go-sqlite3"
)

const (
	Sqlite   = "SQLite"
	Oracle   = "Oracle"
	Postgres = "PostgreSQL"
)

var (
	database *sql.DB
)

type ConnectResultMsg struct {
	Err error
}

func ConnectToDatabase(selectedDriver, databaseUrl, dbName string) tea.Cmd {
	return func() tea.Msg {
		var err error
		database, err = openConnection(selectedDriver, databaseUrl)
		if err == nil {
			var exists bool
			exists, err = store.ConnectionsCfg().Exists(dbName, selectedDriver, databaseUrl)
			if !exists {
				err = store.ConnectionsCfg().SaveConnection(dbName, selectedDriver, databaseUrl)
			}
		}

		// return nil on success, or an error
		return ConnectResultMsg{Err: err}
	}
}

func openConnection(selectedDriver, datasource string) (*sql.DB, error) {
	open, err := sql.Open(selectedDriver, datasource)
	if err != nil {
		return nil, err
	}

	err = open.Ping()
	if err != nil {
		return nil, err
	}

	return open, nil
}

func ExecuteQueryInDatabase(query string) tea.Cmd {
	return func() tea.Msg {
		if strings.HasPrefix(strings.ToLower(query), "select") {
			return executeSelectQuery(query)
		} else {
			return executeUpdateQuery(query)
		}
	}
}
