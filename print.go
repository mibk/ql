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

// String returns a string representing a preprocessed, interpolated, query.
func (q *Query) String() string {
	return makeSql(q)
}

// String returns a string representing a preprocessed, interpolated, query.
func (b *DeleteBuilder) String() string {
	return makeSql(b)
}

// String returns a string representing a preprocessed, interpolated, query.
func (b *InsertBuilder) String() string {
	return makeSql(b)
}

// String returns a string representing a preprocessed, interpolated, query.
func (b *SelectBuilder) String() string {
	return makeSql(b)
}

// String returns a string representing a preprocessed, interpolated, query.
func (b *UpdateBuilder) String() string {
	return makeSql(b)
}
