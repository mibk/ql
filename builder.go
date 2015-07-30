package ql

import "github.com/mibk/ql/query"

type direction bool

// These constant are used to indicate a direction of an ORDER clause. They
// are used as a value in the column/direction map in Order method.
const (
	Asc  direction = false
	Desc direction = true
)

// By is a map of columns and order directions used in builders' Order methods.
// Example of usage:
//	b.Order(ql.By{"col1": ql.Asc,"col2": ql.Desc})
type By map[string]direction

// builder a subset of clauses for the SelectBuilder, InsertBuilder, and DeleteBuilder.
type baseBuilder struct {
	WhereFragments []*whereFragment
	OrderBys       []string
	LimitCount     uint64
	LimitValid     bool
	OffsetCount    uint64
	OffsetValid    bool
}

func (b *baseBuilder) where(exprOrMap interface{}, args ...interface{}) {
	handleExprType(exprOrMap, args, func(expr string, args ...interface{}) {
		expr, args = handleShortNotation(expr, args)
		b.WhereFragments = append(b.WhereFragments, &whereFragment{expr, args})
	})
}

func (b *baseBuilder) orderBy(expr string) {
	b.OrderBys = append(b.OrderBys, expr)
}

func (b *baseBuilder) order(by By) {
	for col, dir := range by {
		expr := "[" + col + "]"
		if dir == Desc {
			expr += " DESC"
		} else {
			expr += " ASC"
		}
		b.orderBy(expr)
	}
}

func (b *baseBuilder) limit(v uint64) {
	b.LimitCount = v
	b.LimitValid = true
}

func (b *baseBuilder) offset(v uint64) {
	b.OffsetCount = v
	b.OffsetValid = true
}

func (b *baseBuilder) buildWhere(w query.Writer, args *[]interface{}) {
	if len(b.WhereFragments) > 0 {
		w.WriteString(" WHERE ")
		writeWhereFragmentsToSql(w, b.WhereFragments, args)
	}
}

func (b *baseBuilder) buildOrder(w query.Writer) {
	if len(b.OrderBys) > 0 {
		w.WriteString(" ORDER BY ")
		for i, s := range b.OrderBys {
			if i > 0 {
				w.WriteString(", ")
			}
			w.WriteString(s)
		}
	}
}

func (b *baseBuilder) buildLimitAndOffset(w query.Writer) {
	if b.LimitValid || b.OffsetValid {
		D.ApplyLimitAndOffset(w, b.LimitCount, b.OffsetCount)
	}
}
