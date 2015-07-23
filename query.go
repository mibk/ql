package ql

import "database/sql"

type Query struct {
	*Connection
	runner
	rawSql string
	args   []interface{}
}

// Query creates Query by the raw SQL query and args
func (db *Connection) Query(sql string, args ...interface{}) *Query {
	return &Query{
		Connection: db,
		runner:     db.Db,
		rawSql:     sql,
		args:       args,
	}
}

// ToSql returns the raw SQL query and args
func (q *Query) ToSql() (string, []interface{}) {
	return q.rawSql, q.args
}

// Exec executes the query
func (q *Query) Exec() (sql.Result, error) {
	return exec(q.runner, q, q, "query")
}
