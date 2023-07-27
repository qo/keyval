package logger

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteLogger struct {
	db *sql.DB
	// lastEventId IdType ; let db handle this
	events chan<- Event
	errors <-chan error
}

// More options can be found here:
// https://github.com/mattn/go-sqlite3#connection-string
type SQLiteConfig struct {
	Path string
}

func getConnectionString(cfg SQLiteConfig) string {
	// Example of  more complex connection string
	// can be found here:
	// https://github.com/mattn/go-sqlite3#dsn-examples
	return fmt.Sprintf("file:%v", cfg.Path)
}

const tableName = "events"

func tableExists(db *sql.DB) (bool, error) {
	const statementTemplate = "SELECT name FROM sqlite_master WHERE type='table' AND name='%s';"
	statement := fmt.Sprintf(
		statementTemplate,
		tableName,
	)
	rows, err := db.Query(statement)
	if err != nil {
		return false, fmt.Errorf(
			"failed to execute sql query to check if table exists: %w",
			err,
		)
	}

	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
	}
	exists := count == 1
	return exists, nil
}

type tableColumn struct {
	name       string
	properties string
	insertable bool
}

var tableColumns = []tableColumn{
	{
		name:       "event_id",
		properties: "INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT",
		insertable: false,
	},
	{
		name:       "event_type",
		properties: "INTEGER NOT NULL",
		insertable: true,
	},
	{
		name:       "event_key",
		properties: "TEXT NOT NULL",
		insertable: true,
	},
	{
		name:       "event_value",
		properties: "TEXT NOT NULL",
		insertable: true,
	},
}

func getColumnsStringForCreatingTable() (string, error) {
	last := len(tableColumns) - 1
	cols := ""
	for idx, col := range tableColumns {
		cols += col.name + " " + col.properties
		if idx != last {
			cols += ", "
		}
	}
	return cols, nil
}

func getColumnsStringForInserting() (string, error) {
	last := len(tableColumns) - 1
	cols := ""
	for idx, col := range tableColumns {
		if col.insertable {
			cols += col.name
			if idx != last {
				cols += ", "
			}
		}
	}
	return cols, nil
}

func getColumnsStringForSelecting() (string, error) {
	last := len(tableColumns) - 1
	cols := ""
	for idx, col := range tableColumns {
		cols += col.name
		if idx != last {
			cols += ", "
		}
	}
	return cols, nil
}

func getAmountOfValuesForInserting() (int, error) {
	res := 0
	for _, col := range tableColumns {
		if col.insertable {
			res++
		}
	}
	return res, nil
}

func createTable(db *sql.DB) error {
	columns, err := getColumnsStringForCreatingTable()
	if err != nil {
		return fmt.Errorf("failed to get table columns string for creating table: %w", err)
	}

	statementTemplate := "CREATE TABLE %s( %s );"
	statement := fmt.Sprintf(
		statementTemplate,
		tableName,
		columns,
	)

	_, err = db.Exec(statement)
	if err != nil {
		return fmt.Errorf(
			"failed to execute sql query to create table: %w",
			err,
		)
	}

	return nil
}

func NewSQLiteLogger(cfg SQLiteConfig) (Logger, error) {
	connStr := getConnectionString(cfg)

	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	exists, err := tableExists(db)
	if err != nil {
		return nil, fmt.Errorf("failed to check if table exists: %w", err)
	}
	if !exists {
		if err = createTable(db); err != nil {
			return nil, fmt.Errorf("failed to create table: %w", err)
		}
	}

	logger := &SQLiteLogger{
		db: db,
	}

	return logger, nil
}

func (l *SQLiteLogger) WritePut(key, val string) {
	l.events <- Event{
		Type: EventPut,
		Key:  key,
		Val:  val,
	}
}

func (l *SQLiteLogger) WriteDelete(key string) {
	l.events <- Event{
		Type: EventDelete,
		Key:  key,
	}
}

func (l *SQLiteLogger) Err() <-chan error {
	return l.errors
}

func getValuesTemplateString() (string, error) {
	const valueTemplate = "?"
	res := ""
	insertableColsAmt, err := getAmountOfValuesForInserting()
	if err != nil {
		return "", err
	}
	for i := 0; i < insertableColsAmt; i++ {
		if i != insertableColsAmt-1 {
			res += valueTemplate + ", "
		} else {
			res += valueTemplate
		}
	}
	return res, nil
}

func (l *SQLiteLogger) Run() {
	events := make(chan Event, eventChanCapacity)
	l.events = events

	errors := make(chan error, errChanCapacity)
	l.errors = errors

	go func() {

		insertableCols, err := getColumnsStringForInserting()
		if err != nil {
			fmt.Errorf(
				"failed to get insertable table columns string: %w",
				err,
			)
		}

		valuesTemplate, err := getValuesTemplateString()
		if err != nil {
			fmt.Errorf(
				"failed to get values template string: %w",
				err,
			)
		}

		query := fmt.Sprintf(
			"INSERT INTO %s ( %s ) VALUES ( %s )",
			tableName,
			insertableCols,
			valuesTemplate,
		)

		for e := range events {
			_, err := l.db.Exec(
				query,
				e.Type,
				e.Key,
				e.Val,
			)
			if err != nil {
				errors <- err
			}
		}
	}()
}

func (l *SQLiteLogger) ReadEvents() (<-chan Event, <-chan error) {
	outEvent := make(chan Event)
	outErr := make(chan error, errChanCapacity)

	go func() {
		defer close(outEvent)
		defer close(outErr)

		cols, err := getColumnsStringForSelecting()
		if err != nil {
			outErr <- fmt.Errorf(
				"failed to get columns string: %w",
				err,
			)
		}

		const indexOfColumnToOrderBy = 0 // corresponds to event_id column
		minIndex, maxIndex := 0, len(tableColumns)-1
		if indexOfColumnToOrderBy < minIndex || indexOfColumnToOrderBy > maxIndex {
			outErr <- fmt.Errorf(
				"index of column to order by (%d) is out of range ([%d, %d])",
				indexOfColumnToOrderBy,
				minIndex,
				maxIndex,
			)
		}
		columnToOrderBy := tableColumns[indexOfColumnToOrderBy].name

		query := fmt.Sprintf(
			"SELECT %s FROM %s ORDER BY %s",
			cols,
			tableName,
			columnToOrderBy,
		)

		rows, err := l.db.Query(query)
		if err != nil {
			outErr <- fmt.Errorf(
				"failed to execute sql query: %w",
				err,
			)
		}

		// if no rows are present
		if rows == nil {
			outErr <- fmt.Errorf(
				"failed to read rows: %w",
				errors.New("no rows present"),
			)
		}

		defer rows.Close()

		e, err := Event{}, nil

		for rows.Next() {
			err = rows.Scan(
				&e.Id,
				&e.Type,
				&e.Key,
				&e.Val,
			)

			if err != nil {
				outErr <- fmt.Errorf(
					"failed to read row: %w",
					err,
				)
				return
			}

			outEvent <- e
		}

		err = rows.Err()
		if err != nil {
			outErr <- fmt.Errorf(
				"failed to read log: %w",
				err,
			)
		}
	}()

	return outEvent, outErr
}
