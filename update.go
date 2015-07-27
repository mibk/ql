package ql

import "bytes"

// UpdateBuilder contains the clauses for an UPDATE statement.
type UpdateBuilder struct {
	executor

	Table      string
	SetClauses []*setClause
	*builder
}

type setClause struct {
	column string
	value  interface{}
}

// Update creates a new UpdateBuilder for the given table.
func (db *Connection) Update(table string) *UpdateBuilder {
	b := &UpdateBuilder{
		executor: executor{EventReceiver: db, runner: db.DB},
		Table:    table,
		builder:  new(builder),
	}
	b.executor.builder = b
	return b
}

// Update creates a new UpdateBuilder for the given table bound to a transaction.
func (tx *Tx) Update(table string) *UpdateBuilder {
	b := &UpdateBuilder{
		executor: executor{EventReceiver: tx.Connection, runner: tx.Tx},
		Table:    table,
		builder:  new(builder),
	}
	b.executor.builder = b
	return b
}

// Set appends a column/value pair for the statement.
func (b *UpdateBuilder) Set(column string, value interface{}) *UpdateBuilder {
	b.SetClauses = append(b.SetClauses, &setClause{column: column, value: value})
	return b
}

// SetMap appends the elements of the map as column/value pairs for the statement.
func (b *UpdateBuilder) SetMap(clauses map[string]interface{}) *UpdateBuilder {
	for col, val := range clauses {
		b = b.Set(col, val)
	}
	return b
}

// Where appends a WHERE clause to the statement.
func (b *UpdateBuilder) Where(whereSqlOrMap interface{}, args ...interface{}) *UpdateBuilder {
	b.where(whereSqlOrMap, args...)
	return b
}

// OrderBy appends a column to ORDER the statement by.
func (b *UpdateBuilder) OrderBy(expr string) *UpdateBuilder {
	b.orderBy(expr)
	return b
}

// Order accepts By map of columns and directions to ORDER the statement by.
func (b *UpdateBuilder) Order(by By) *UpdateBuilder {
	b.order(by)
	return b
}

// Limit sets a limit for the statement; overrides any existing LIMIT.
func (b *UpdateBuilder) Limit(limit uint64) *UpdateBuilder {
	b.limit(limit)
	return b
}

// Offset sets an offset for the statement; overrides any existing OFFSET.
func (b *UpdateBuilder) Offset(offset uint64) *UpdateBuilder {
	b.offset(offset)
	return b
}

// ToSql serialized the UpdateBuilder to a SQL string. It returns the string with
// placeholders and a slice of query arguments.
func (b *UpdateBuilder) ToSql() (string, []interface{}) {
	if len(b.Table) == 0 {
		panic("no table specified")
	}
	if len(b.SetClauses) == 0 {
		panic("no set clauses specified")
	}

	sql := new(bytes.Buffer)
	var args []interface{}

	sql.WriteString("UPDATE ")
	sql.WriteString(b.Table)
	sql.WriteString(" SET ")

	// Build SET clause SQL with placeholders and add values to args
	for i, c := range b.SetClauses {
		if i > 0 {
			sql.WriteString(", ")
		}
		Quoter.writeQuotedColumn(c.column, sql)
		if e, ok := c.value.(*expr); ok {
			sql.WriteString(" = ")
			sql.WriteString(e.Sql)
			args = append(args, e.Values...)
		} else {
			sql.WriteString(" = ?")
			args = append(args, c.value)
		}
	}

	b.buildWhere(sql, &args)
	b.buildOrder(sql)
	b.buildLimitAndOffset(sql)

	return sql.String(), args
}
