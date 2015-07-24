package ql

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

func (b *builder) orderDir(expr string, isAsc bool) {
	if isAsc {
		b.OrderBys = append(b.OrderBys, expr+" ASC")
	} else {
		b.OrderBys = append(b.OrderBys, expr+" DESC")
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
