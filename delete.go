package ql

import "bytes"

// DeleteBuilder contains the clauses for a DELETE statement.
type DeleteBuilder struct {
	executor

	From string
	*builder
}

// DeleteFrom creates a new DeleteBuilder for the given table.
func (db *Connection) DeleteFrom(from string) *DeleteBuilder {
	b := &DeleteBuilder{
		executor: executor{Connection: db, runner: db.DB},
		From:     from,
		builder:  new(builder),
	}
	b.executor.builder = b
	return b
}

// DeleteFrom creates a new DeleteBuilder for the given table in the context for
// a transaction.
func (tx *Tx) DeleteFrom(from string) *DeleteBuilder {
	b := &DeleteBuilder{
		executor: executor{Connection: tx.Connection, runner: tx.Tx},
		From:     from,
		builder:  new(builder),
	}
	b.executor.builder = b
	return b
}

// Where appends a WHERE clause to the statement whereSqlOrMap can be a string or map.
// If it's a string, args wil replaces any places holders.
func (b *DeleteBuilder) Where(whereSqlOrMap interface{}, args ...interface{}) *DeleteBuilder {
	b.where(whereSqlOrMap, args...)
	return b
}

// OrderBy appends an ORDER BY clause to the statement.
func (b *DeleteBuilder) OrderBy(expr string) *DeleteBuilder {
	b.orderBy(expr)
	return b
}

// Order accepts By map of columns and directions to ORDER the statement by.
func (b *DeleteBuilder) Order(by By) *DeleteBuilder {
	b.order(by)
	return b
}

// Limit sets a LIMIT clause for the statement; overrides any existing LIMIT.
func (b *DeleteBuilder) Limit(limit uint64) *DeleteBuilder {
	b.limit(limit)
	return b
}

// Offset sets an OFFSET clause for the statement; overrides any existing OFFSET.
func (b *DeleteBuilder) Offset(offset uint64) *DeleteBuilder {
	b.offset(offset)
	return b
}

// ToSql serialized the DeleteBuilder to a SQL string.  It returns the string with
// placeholders and a slice of query arguments.
func (b *DeleteBuilder) ToSql() (string, []interface{}) {
	if len(b.From) == 0 {
		panic("no table specified")
	}

	sql := new(bytes.Buffer)
	var args []interface{}

	sql.WriteString("DELETE FROM ")
	sql.WriteString(b.From)

	b.buildWhere(sql, &args)
	b.buildOrder(sql)
	b.buildLimitAndOffset(sql)

	return sql.String(), args
}
