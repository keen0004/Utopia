package database

import (
	"errors"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type Database struct {
	path       string
	connection *sqlite.Conn
}

// new database instance (sqlite3)
func NewDatabase(path string) *Database {
	return &Database{
		path:       path,
		connection: nil,
	}
}

// open database for read and write
func (db *Database) Open() error {
	if db.connection != nil {
		db.Close()
	}

	var err error
	db.connection, err = sqlite.OpenConn(db.path, sqlite.OpenCreate|sqlite.OpenReadWrite)
	if err != nil {
		return nil
	}

	return nil
}

func (db *Database) Close() {
	if db.connection != nil {
		db.connection.Close()
	}

	db.connection = nil
}

// exec sql on database
func (db *Database) ExecSql(sql string, args ...interface{}) (int, error) {
	if db.connection == nil {
		return 0, errors.New("Database not opened")
	}

	err := sqlitex.Execute(db.connection, sql, &sqlitex.ExecOptions{
		Args: args,
	})
	if err != nil {
		return 0, err
	}

	return db.connection.Changes(), nil
}

func (db *Database) Query(sql string, args ...interface{}) ([][]interface{}, error) {
	if db.connection == nil {
		return make([][]interface{}, 0), errors.New("Database not opened")
	}

	result := make([][]interface{}, 0)
	err := sqlitex.Execute(db.connection, sql, &sqlitex.ExecOptions{
		Args: args,
		ResultFunc: func(stmt *sqlite.Stmt) error {
			col := stmt.ColumnCount()
			rows := make([]interface{}, 0, col)

			// parse result, support int, float, text, bytes
			for i := 0; i < col; i++ {
				switch stmt.ColumnType(i) {
				case sqlite.TypeInteger:
					rows = append(rows, stmt.ColumnInt(i))
				case sqlite.TypeFloat:
					rows = append(rows, stmt.ColumnFloat(i))
				case sqlite.TypeText:
					rows = append(rows, stmt.ColumnText(i))
				case sqlite.TypeBlob:
				default:
					buf := make([]byte, 0, stmt.ColumnLen(i))
					stmt.ColumnBytes(i, buf)
					rows = append(rows, buf)
				}
			}

			result = append(result, rows)
			return nil
		},
	})

	if err != nil {
		return make([][]interface{}, 0), err
	}

	return result, nil
}
