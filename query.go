package ql

// Query represent an arbitrary SQL statement.
type Query struct {
	loader // methods for loading structs and values
	executor

	rawSql string
	args   []interface{}
}

func newQuery(e EventReceiver, r runner, sql string, args ...interface{}) *Query {
	q := &Query{
		loader:   loader{EventReceiver: e, runner: r},
		executor: executor{EventReceiver: e, runner: r},
		rawSql:   sql,
		args:     args,
	}
	q.loader.builder = q
	q.executor.builder = q
	return q
}

// Query creates Query by the raw SQL query and args.
func (db *Connection) Query(sql string, args ...interface{}) *Query {
	return newQuery(db, db.DB, sql, args...)
}

// Query creates Query by the raw SQL query and args. Query is bound to the transaction.
func (tx *Tx) Query(sql string, args ...interface{}) *Query {
	return newQuery(tx.Connection, tx.Tx, sql, args...)
}

// ToSql returns the raw SQL query and args.
func (q *Query) ToSql() (string, []interface{}) {
	return q.rawSql, q.args
}
