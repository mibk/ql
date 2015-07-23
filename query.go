package ql

import "database/sql"

type Query struct {
	// methods for loading structs and values
	loader

	rawSql string
	args   []interface{}
}

// Query creates Query by the raw SQL query and args
func (db *Connection) Query(sql string, args ...interface{}) *Query {
	q := &Query{
		loader: loader{Connection: db, runner: db.Db},
		rawSql: sql,
		args:   args,
	}
	q.loader.builder = q
	return q
}

// ToSql returns the raw SQL query and args
func (q *Query) ToSql() (string, []interface{}) {
	return q.rawSql, q.args
}

// Exec executes the query
func (q *Query) Exec() (sql.Result, error) {
	return exec(q.runner, q, q, "query")
}
