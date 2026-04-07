package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
)

type QueryResultMsg struct {
	Columns []table.Column
	Rows    []table.Row
	Err     error
}

func executeSelectQuery(query string) QueryResultMsg {
	queryResult, err := runQuery(query)
	if err != nil {
		// return nil on success, or an error
		return QueryResultMsg{Columns: nil, Rows: nil, Err: err}
	}
	defer queryResult.Close()

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

func executeUpdateQuery(query string) QueryResultMsg {
	result, err := runExec(query)
	if err != nil {
		// return nil on success, or an error
		return QueryResultMsg{Err: err}
	}

	affected, _ := result.RowsAffected()
	id, _ := result.LastInsertId()

	headers := []table.Column{
		{
			Title: "Result",
			Width: 8,
		},
		{
			Title: "Affected Rows",
			Width: 13,
		},
		{
			Title: "Last Inserted Id",
			Width: 16,
		},
	}

	rows := []table.Row{
		{"InputQuery OK", fmt.Sprintf("%v", affected), fmt.Sprintf("%v", id)},
	}

	return QueryResultMsg{Columns: headers, Rows: rows}
}

func runExec(query string) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return database.ExecContext(ctx, query)
}

func runQuery(query string) (*sql.Rows, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()

	return database.QueryContext(ctx, query)
}

func extractHeaders(rows *sql.Rows) ([]table.Column, int, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, 0, err
	}

	headers := make([]table.Column, len(columns))
	for idx, column := range columns {
		headers[idx] = table.Column{
			Title: column,
			Width: len(column),
		}
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
