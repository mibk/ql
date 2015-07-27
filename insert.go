package ql

import (
	"bytes"
	"reflect"
)

// InsertBuilder contains the clauses for an INSERT statement.
type InsertBuilder struct {
	executor

	Into string
	Cols []string
	Vals [][]interface{}
	Recs []interface{}
}

func newInsertBuilder(e EventReceiver, r runner, into string) *InsertBuilder {
	b := &InsertBuilder{
		executor: executor{EventReceiver: e, runner: r},
		Into:     into,
	}
	b.executor.builder = b
	return b
}

// InsertInto instantiates a InsertBuilder for the given table.
func (db *Connection) InsertInto(into string) *InsertBuilder {
	return newInsertBuilder(db, db.DB, into)
}

// InsertInto instantiates a InsertBuilder for the given table bound to a transaction.
func (tx *Tx) InsertInto(into string) *InsertBuilder {
	return newInsertBuilder(tx.Connection, tx.Tx, into)
}

// Columns appends columns to insert in the statement.
func (b *InsertBuilder) Columns(columns ...string) *InsertBuilder {
	b.Cols = columns
	return b
}

// Values appends a set of values to the statement.
func (b *InsertBuilder) Values(vals ...interface{}) *InsertBuilder {
	b.Vals = append(b.Vals, vals)
	return b
}

// Record pulls in values to match Columns from the record.
func (b *InsertBuilder) Record(record interface{}) *InsertBuilder {
	b.Recs = append(b.Recs, record)
	return b
}

// Pair adds a key/value pair to the statement.
func (b *InsertBuilder) Pair(column string, value interface{}) *InsertBuilder {
	b.Cols = append(b.Cols, column)
	lenVals := len(b.Vals)
	if lenVals == 0 {
		args := []interface{}{value}
		b.Vals = [][]interface{}{args}
	} else if lenVals == 1 {
		b.Vals[0] = append(b.Vals[0], value)
	} else {
		panic("pair only allows you to specify 1 record to insret")
	}
	return b
}

// ToSql serialized the InsertBuilder to a SQL string.  It returns the string with
// placeholders and a slice of query arguments.
func (b *InsertBuilder) ToSql() (string, []interface{}) {
	if len(b.Into) == 0 {
		panic("no table specified")
	}
	if len(b.Cols) == 0 {
		panic("no columns specified")
	}
	if len(b.Vals) == 0 && len(b.Recs) == 0 {
		panic("no values or records specified")
	}

	sql := new(bytes.Buffer)
	var placeholder bytes.Buffer // Build the placeholder like "(?,?,?)"
	var args []interface{}

	sql.WriteString("INSERT INTO ")
	sql.WriteString(b.Into)
	sql.WriteString(" (")

	// Simulataneously write the cols to the sql buffer, and build a placeholder
	placeholder.WriteRune('(')
	for i, c := range b.Cols {
		if i > 0 {
			sql.WriteRune(',')
			placeholder.WriteRune(',')
		}
		Quoter.writeQuotedColumn(c, sql)
		placeholder.WriteRune('?')
	}
	sql.WriteString(") VALUES ")
	placeholder.WriteRune(')')
	placeholderStr := placeholder.String()

	// Go thru each value we want to insert. Write the placeholders, and collect args
	for i, row := range b.Vals {
		if i > 0 {
			sql.WriteRune(',')
		}
		sql.WriteString(placeholderStr)

		for _, v := range row {
			args = append(args, v)
		}
	}
	anyVals := len(b.Vals) > 0

	// Go thru the records. Write the placeholders, and do reflection on the records to extract args
	for i, rec := range b.Recs {
		if i > 0 || anyVals {
			sql.WriteRune(',')
		}
		sql.WriteString(placeholderStr)

		ind := reflect.Indirect(reflect.ValueOf(rec))
		vals, err := valuesFor(ind.Type(), ind, b.Cols)
		if err != nil {
			panic(err.Error())
		}
		for _, v := range vals {
			args = append(args, v)
		}
	}

	return sql.String(), args
}
