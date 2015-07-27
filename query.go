package ql

type Query struct {
	loader // methods for loading structs and values
	executor

	rawSql string
	args   []interface{}
}

// Query creates Query by the raw SQL query and args.
func (db *Connection) Query(sql string, args ...interface{}) *Query {
	q := &Query{
		loader:   loader{EventReceiver: db, runner: db.DB},
		executor: executor{EventReceiver: db, runner: db.DB},
		rawSql:   sql,
		args:     args,
	}
	q.loader.builder = q
	q.executor.builder = q
	return q
}

// ToSql returns the raw SQL query and args.
func (q *Query) ToSql() (string, []interface{}) {
	return q.rawSql, q.args
}
