package ql

import "fmt"

type direction bool

const (
	Asc  direction = false
	Desc direction = true
)

// By is a map of columns and order directions used in builders' Order methods.
// Example of usage:
//	b.Order(ql.By{"col1": ql.Asc,"col2": ql.Desc})
type By map[string]direction

// builder a subset of clauses for the SelectBuilder, InsertBuilder, and DeleteBuilder.
type builder struct {
	WhereFragments []*whereFragment
	OrderBys       []string
	LimitCount     uint64
	LimitValid     bool
	OffsetCount    uint64
	OffsetValid    bool
}

func (b *builder) where(whereSqlOrMap interface{}, args ...interface{}) {
	b.WhereFragments = append(b.WhereFragments, newWhereFragment(whereSqlOrMap, args))
}

func (b *builder) orderBy(expr string) {
	b.OrderBys = append(b.OrderBys, expr)
}

func (b *builder) order(by By) {
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

func (b *builder) limit(v uint64) {
	b.LimitCount = v
	b.LimitValid = true
}

func (b *builder) offset(v uint64) {
	b.OffsetCount = v
	b.OffsetValid = true
}

func (b *builder) buildWhere(w queryWriter, args *[]interface{}) {
	if len(b.WhereFragments) > 0 {
		w.WriteString(" WHERE ")
		writeWhereFragmentsToSql(b.WhereFragments, w, args)
	}
}

func (b *builder) buildOrder(w queryWriter) {
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

func (b *builder) buildLimitAndOffset(w queryWriter) {
	if b.LimitValid {
		w.WriteString(" LIMIT ")
		fmt.Fprint(w, b.LimitCount)
	}
	if b.OffsetValid {
		w.WriteString(" OFFSET ")
		fmt.Fprint(w, b.OffsetCount)
	}
}
