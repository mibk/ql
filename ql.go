package ql

import (
	"database/sql"
)

// Connection is a connection to the database with an EventReceiver to send events,
// errors, and timings to.
type Connection struct {
	DB *sql.DB
	EventReceiver
}

// NewConnection instantiates a Connection for a given database/sql connection
// and event receiver.
func NewConnection(db *sql.DB, log EventReceiver) *Connection {
	if log == nil {
		log = nullReceiver
	}

	return &Connection{DB: db, EventReceiver: log}
}

// Close closes the database, releasing any open resources.
func (c *Connection) Close() error {
	return c.DB.Close()
}

type runner interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

// Open opens a database by calling sql.Open. It returns new Connection with
// nil EventReceiver.
func Open(driverName, dataSourceName string) (*Connection, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return NewConnection(db, nil), nil
}

// MustOpen is like Open but panics on error.
func MustOpen(driverName, dataSourceName string) *Connection {
	conn, err := Open(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}
	return conn
}
