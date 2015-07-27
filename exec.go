package ql

import (
	"database/sql"
	"time"
)

type executor struct {
	EventReceiver
	runner
	builder queryBuilder
}

// Exec executes the query. It returns the raw database/sql Result and an error if there
// is one.
func (e executor) Exec() (sql.Result, error) {
	fullSql, err := Preprocess(e.builder.ToSql())
	if err != nil {
		return nil, e.EventErrKv("ql.exec.interpolate", err, kvs{"sql": fullSql})
	}

	startTime := time.Now()
	defer func() {
		e.TimingKv("ql.exec", time.Since(startTime).Nanoseconds(), kvs{"sql": fullSql})
	}()

	result, err := e.runner.Exec(fullSql)
	if err != nil {
		return result, e.EventErrKv("ql.exec.exec", err, kvs{"sql": fullSql})
	}

	return result, nil
}

// MustExec is like Exec but panics on error.
func (e executor) MustExec() sql.Result {
	res, err := e.Exec()
	if err != nil {
		panic(err)
	}
	return res
}
