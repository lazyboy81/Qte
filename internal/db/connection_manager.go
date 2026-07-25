package db

import (
	"database/sql"
	"query-terminal/internal/store"

	tea "charm.land/bubbletea/v2"
	_ "github.com/godror/godror"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
)

const (
	Sqlite   = "SQLite"
	Oracle   = "Oracle"
	Postgres = "PostgreSQL"
)

type ConnectResultMsg struct {
	Conn *ActiveConnection
	Err  error
}

func ConnectToDatabase(selectedDriver, databaseUrl, dbName string, connId int) tea.Cmd {
	return func() tea.Msg {
		connection, err := openConnection(selectedDriver, databaseUrl)
		if err == nil {
			var exists bool
			exists, err = store.ConnectionsCfg().Exists(dbName, selectedDriver, databaseUrl)
			if err != nil {
				return err
			}

			if !exists {
				err = store.ConnectionsCfg().SaveConnection(dbName, selectedDriver, databaseUrl)
			}

			connection.ConnId = connId
		}

		// return nil on success, or an error
		return &ConnectResultMsg{Conn: connection, Err: err}
	}
}

func openConnection(selectedDriver, datasource string) (*ActiveConnection, error) {
	open, err := sql.Open(selectedDriver, datasource)
	if err != nil {
		return nil, err
	}

	err = open.Ping()
	if err != nil {
		return nil, err
	}

	return &ActiveConnection{db: open}, nil
}
