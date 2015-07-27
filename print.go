package ql

type queryBuilder interface {
	ToSql() (string, []interface{})
}

func makeSql(b queryBuilder) string {
	sql, err := Preprocess(b.ToSql())
	if err != nil {
		panic(err)
	}
	return sql
}

func (q *Query) String() string {
	return makeSql(q)
}

func (b *DeleteBuilder) String() string {
	return makeSql(b)
}

func (b *InsertBuilder) String() string {
	return makeSql(b)
}

func (b *SelectBuilder) String() string {
	return makeSql(b)
}

func (b *UpdateBuilder) String() string {
	return makeSql(b)
}
