package db

import (
	"database/sql"
	"fmt"
	"log/slog"

	//_ "github.com/godror/godror"
	_ "github.com/mattn/go-sqlite3"
)

const (
	dbSqlite   = "SQLite"
	dbOracle   = "Oracle"
	dbPostgres = "PostgreSQL"
)

type connector func(string) (*sql.DB, error)

var registry map[string]connector

func register(name string, fn connector) {
	registry[name] = fn
}

func connect(name, dsn string) (*sql.DB, error) {
	slog.Debug("got " + name + " " + dsn)
	fn, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown db: %s", name)
	}
	return fn(dsn)
}

func InitializeDBRegistry() {
	registry = make(map[string]connector)
	slog.Debug("initializing db registry...")
	slog.Debug("adding sqlite to registry with key: " + dbSqlite)
	register(dbSqlite, func(datasource string) (*sql.DB, error) {
		return sql.Open("sqlite3", datasource)
	})
	slog.Debug("added sqlite to registry")

	//register(dbOracle, func(datasource string) (*sql.DB, error) {
	//	return sql.Open("godror", datasource)
	//})
	//slog.Debug("added oracle to registry")
}
