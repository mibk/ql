package ql

import (
	"bytes"
	"database/sql"
	"fmt"
)

// DeleteBuilder contains the clauses for a DELETE statement.
type DeleteBuilder struct {
	*Connection
	runner

	From string
	*builder
}

// DeleteFrom creates a new DeleteBuilder for the given table.
func (db *Connection) DeleteFrom(from string) *DeleteBuilder {
	return &DeleteBuilder{
		Connection: db,
		runner:     db.Db,
		From:       from,
		builder:    new(builder),
	}
}

// DeleteFrom creates a new DeleteBuilder for the given table in the context for
// a transaction.
func (tx *Tx) DeleteFrom(from string) *DeleteBuilder {
	return &DeleteBuilder{
		Connection: tx.Connection,
		runner:     tx.Tx,
		From:       from,
		builder:    new(builder),
	}
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

	var sql bytes.Buffer
	var args []interface{}

	sql.WriteString("DELETE FROM ")
	sql.WriteString(b.From)

	// Write WHERE clause if we have any fragments
	if len(b.WhereFragments) > 0 {
		sql.WriteString(" WHERE ")
		writeWhereFragmentsToSql(b.WhereFragments, &sql, &args)
	}

	// Ordering and limiting
	if len(b.OrderBys) > 0 {
		sql.WriteString(" ORDER BY ")
		for i, s := range b.OrderBys {
			if i > 0 {
				sql.WriteString(", ")
			}
			sql.WriteString(s)
		}
	}

	if b.LimitValid {
		sql.WriteString(" LIMIT ")
		fmt.Fprint(&sql, b.LimitCount)
	}

	if b.OffsetValid {
		sql.WriteString(" OFFSET ")
		fmt.Fprint(&sql, b.OffsetCount)
	}

	return sql.String(), args
}

// Exec executes the statement represented by the DeleteBuilder. It returns the raw
// database/sql Result and an error if there was one.
func (b *DeleteBuilder) Exec() (sql.Result, error) {
	return exec(b.runner, b, b, "delete")
}
