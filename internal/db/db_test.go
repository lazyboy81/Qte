package db

import (
	"context"
	"database/sql"
	"slices"
	"testing"
	"time"

	"charm.land/bubbles/v2/table"
)

func TestConnectToDatabase(t *testing.T) {
	connection, err := openConnection("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("encountered an error when connecting to database: %v", err)
	}

	if connection == nil {
		t.Fatal("connection struct is nil")
	}

	if connection.db == nil {
		t.Fatal("sql.DB is nil")
	}
	defer connection.db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Basic working test
	_, err = connection.db.ExecContext(ctx, `CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT);`)
	if err != nil {
		t.Fatalf("create table: %v", err)
	}

	_, err = connection.db.ExecContext(ctx, `INSERT INTO test (name) VALUES (?)`, "alice")
	if err != nil {
		t.Fatalf("insert: %v", err)
	}

	var name string
	err = connection.db.QueryRowContext(ctx, `SELECT name FROM test WHERE id = 1`).Scan(&name)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if name != "alice" {
		t.Fatalf("got %q, want %q", name, "alice")
	}
}

func TestExtractColumns(t *testing.T) {
	expectedColumns := []table.Column{
		{Title: "id", Width: 2},
		{Title: "name", Width: 4},
	}

	open, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("encountered an error when connecting to database: %v", err)
	}
	defer open.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = open.ExecContext(ctx, `CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT);`)
	if err != nil {
		t.Fatalf("encountered an error when creating a table: %v", err)
	}

	_, err = open.ExecContext(ctx, `INSERT INTO test (name) VALUES (?)`, "alice")
	if err != nil {
		t.Fatalf("encountered an error when inserting a record in table: %v", err)
	}

	rows, err := open.QueryContext(ctx, `select * from test`)
	if err != nil {
		t.Fatalf("encountered an error when selecting records from table: %v", err)
	}

	columns, columnCount, err := extractHeaders(rows)
	if err != nil {
		t.Fatalf("encountered an error when extracting table column names: %v", err)
	}

	if columnCount != 2 {
		t.Fatal("columns count should be 2")
	}

	cmp := slices.CompareFunc(columns, expectedColumns, func(first table.Column, second table.Column) int {
		if first.Title != second.Title || first.Width != second.Width {
			return 1
		}

		return 0
	})

	if cmp != 0 {
		t.Fatalf("the result columns are equal to expected columns:\n expected: %v got: %v", expectedColumns, columns)
	}
}

func TestExtractRows(t *testing.T) {
	expectedRows := []table.Row{{"1", "alice"}}

	open, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("encountered an error when connecting to database: %v", err)
	}
	defer open.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = open.ExecContext(ctx, `CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT);`)
	if err != nil {
		t.Fatalf("encountered an error when creating a table: %v", err)
	}

	_, err = open.ExecContext(ctx, `INSERT INTO test (name) VALUES (?)`, "alice")
	if err != nil {
		t.Fatalf("encountered an error when inserting a record in table: %v", err)
	}

	rows, err := open.QueryContext(ctx, `select * from test`)
	if err != nil {
		t.Fatalf("encountered an error when selecting records from table: %v", err)
	}

	tableRows, err := extractRows(rows, 2)
	if err != nil {
		t.Fatalf("encountered an error when extracting table column names: %v", err)
	}

	if len(tableRows) != 1 {
		t.Fatal("row count should be 1")
	}

	cmp := slices.CompareFunc(tableRows, expectedRows, func(first table.Row, second table.Row) int {
		return slices.Compare(first, second)
	})

	if cmp != 0 {
		t.Fatalf("the result columns are equal to expected columns:\n expected: %v got: %v", expectedRows, tableRows)
	}
}

func TestActiveConnection_ExecuteUpdateQuery(t *testing.T) {
	connection, err := openConnection("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("encountered an error when connecting to database: %v", err)
	}
	defer connection.db.Close()

	rst := connection.executeUpdateQuery(`CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT);`)
	if rst.Err != nil {
		t.Fatalf("create table: %v", err)
	}

	if rst.Rows[0][0] != "InputQuery OK" {
		t.Fatalf("create - first row result is not: %v", "InputQuery OK")
	}

	rst = connection.executeUpdateQuery(`INSERT INTO test (name) VALUES ('alice')`)
	if rst.Err != nil {
		t.Fatalf("insert table: %v", err)
	}

	if rst.Rows[0][0] != "InputQuery OK" {
		t.Fatalf("insert - first row result is not: %v", "InputQuery OK")
	}
}

func TestActiveConnection_ExecuteSelectQuery(t *testing.T) {
	connection, err := openConnection("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("encountered an error when connecting to database: %v", err)
	}
	defer connection.db.Close()

	rst := connection.executeUpdateQuery(`CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT);`)
	if rst.Err != nil {
		t.Fatalf("create table: %v", err)
	}

	rst = connection.executeUpdateQuery(`INSERT INTO test (name) VALUES ('alice')`)
	if rst.Err != nil {
		t.Fatalf("insert table: %v", err)
	}

	result := connection.executeSelectQuery(`SELECT name FROM test WHERE id = 1`)
	if rst.Err != nil {
		t.Fatalf("select table: %v", err)
	}

	if result.Columns[0].Title != "name" {
		t.Fatal("title is not 'name'")
	}

	if len(result.Rows) != 1 {
		t.Fatal("must have one row to show")
	}

	if result.Rows[0][0] != "alice" {
		t.Fatal("the given row does not have a value")
	}
}
