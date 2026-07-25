package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
)

type QueryResultMsg struct {
	Columns []table.Column
	Rows    []table.Row
	Err     error
}

type ActiveConnection struct {
	ConnId int
	db     *sql.DB
}

func (a *ActiveConnection) ExecuteQueryInDatabase(query string) tea.Cmd {
	return func() tea.Msg {
		q := strings.TrimSpace(strings.ToLower(query))
		if strings.HasPrefix(q, "select") || strings.HasPrefix(q, "with") {
			return a.executeSelectQuery(query)
		}

		return a.executeUpdateQuery(query)
	}
}

func (a *ActiveConnection) executeSelectQuery(query string) QueryResultMsg {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFunc()

	queryResult, err := a.db.QueryContext(ctx, query)
	if err != nil {
		// return nil on success, or an error
		return QueryResultMsg{Err: err}
	}
	defer queryResult.Close()

	columns, columnCount, err := extractHeaders(queryResult)
	if err != nil {
		// return nil on success, or an error
		return QueryResultMsg{Err: err}
	}

	rows, err := extractRows(queryResult, columnCount)
	if err != nil {
		// return nil on success, or an error
		return QueryResultMsg{Err: err}
	}

	return QueryResultMsg{Columns: columns, Rows: rows}
}

func (a *ActiveConnection) executeUpdateQuery(query string) QueryResultMsg {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	result, err := a.db.ExecContext(ctx, query)
	if err != nil {
		// return nil on success, or an error
		return QueryResultMsg{Err: err}
	}

	affected, _ := result.RowsAffected()
	id, _ := result.LastInsertId()

	headers := []table.Column{
		{
			Title: "Result",
			Width: 13,
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
		for i := range columnCount {
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
