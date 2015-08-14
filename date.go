package ql

import (
	"database/sql/driver"
	"time"
)

type date time.Time

// Date turns time into date type. It is supposed to be used in queries to use
// 2006-01-02 time format.
func Date(t time.Time) date {
	return date(t)
}

func (d date) Value() (driver.Value, error) {
	return time.Time(d).Format("2006-01-02"), nil
}
