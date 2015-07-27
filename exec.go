package ql

import (
	"database/sql"
	"time"
)

func exec(db runner, b queryBuilder, r EventReceiver, logAction string) (sql.Result, error) {
	fullSql, err := Preprocess(b.ToSql())
	if err != nil {
		return nil, r.EventErrKv("ql."+logAction+".exec.interpolate", err, kvs{"sql": fullSql})
	}

	startTime := time.Now()
	defer func() {
		r.TimingKv("ql."+logAction, time.Since(startTime).Nanoseconds(), kvs{"sql": fullSql})
	}()

	result, err := db.Exec(fullSql)
	if err != nil {
		return result, r.EventErrKv("ql."+logAction+".exec.exec", err, kvs{"sql": fullSql})
	}

	return result, nil
}
