package ql

import (
	"database/sql/driver"
	"testing"
)

func TestInterpolate(t *testing.T) {
	var noArgs []interface{}
	tests := []struct {
		sql    string
		args   []interface{}
		expSql string
		expErr error
	}{
		// NULL
		{"SELECT * FROM x WHERE a = ?", []interface{}{nil},
			"SELECT * FROM x WHERE a = NULL", nil},

		// integers
		{
			`SELECT * FROM x WHERE a = ? AND b = ? AND c = ? AND d = ? AND e = ? AND f = ?
			AND g = ? AND h = ? AND i = ? AND j = ?`,
			[]interface{}{int(1), int8(-2), int16(3), int32(4), int64(5), uint(6), uint8(7),
				uint16(8), uint32(9), uint64(10)},
			`SELECT * FROM x WHERE a = 1 AND b = -2 AND c = 3 AND d = 4 AND e = 5 AND f = 6
			AND g = 7 AND h = 8 AND i = 9 AND j = 10`, nil,
		},

		// boolean
		{"SELECT * FROM x WHERE a = ? AND b = ?", []interface{}{true, false},
			"SELECT * FROM x WHERE a = 1 AND b = 0", nil},

		// floats
		{"SELECT * FROM x WHERE a = ? AND b = ?", []interface{}{float32(0.15625), float64(3.14159)},
			"SELECT * FROM x WHERE a = 0.15625 AND b = 3.14159", nil},

		// strings
		{
			`SELECT * FROM x WHERE a = ?
			AND b = ?`,
			[]interface{}{"hello", "\"hello's \\ world\" \n\r\x00\x1a"},
			`SELECT * FROM x WHERE a = 'hello'
			AND b = '\"hello\'s \\ world\" \n\r\x00\x1a'`, nil,
		},

		// slices
		{"SELECT * FROM x WHERE a = ? AND b = ? AND c = ? AND d = ?",
			[]interface{}{[]int{1}, []int{1, 2, 3}, []uint32{5, 6, 7}, []string{"wat", "ok"}},
			"SELECT * FROM x WHERE a = (1) AND b = (1,2,3) AND c = (5,6,7) AND d = ('wat','ok')", nil},

		// valuers
		{"SELECT * FROM x WHERE a = ? AND b = ?",
			[]interface{}{myString{true, "wat"}, myString{false, "fry"}},
			"SELECT * FROM x WHERE a = 'wat' AND b = NULL", nil},

		// errors
		{"SELECT * FROM x WHERE a = ? AND b = ?", []interface{}{1},
			"", ErrArgumentMismatch},

		{"SELECT * FROM x WHERE", []interface{}{1},
			"", ErrArgumentMismatch},

		{"SELECT * FROM x WHERE a = ?", []interface{}{string([]byte{0x34, 0xFF, 0xFE})},
			"", ErrNotUTF8},

		{"SELECT * FROM x WHERE a = ?", []interface{}{struct{}{}},
			"", ErrInvalidValue},

		{"SELECT * FROM x WHERE a = ?", []interface{}{[]struct{}{struct{}{}, struct{}{}}},
			"", ErrInvalidSliceValue},

		// quoting
		{"SELECT [name] FROM [user]", noArgs, "SELECT `name` FROM `user`", nil},
		{"SELECT [u.name] FROM [user] [u]", noArgs, "SELECT `u`.`name` FROM `user` `u`", nil},
		{"SELECT [u.na`me] FROM [user] [u]", noArgs, "SELECT `u`.`na``me` FROM `user` `u`", nil},
	}

	for _, test := range tests {
		str, err := Interpolate(test.sql, test.args)
		if err != test.expErr {
			t.Errorf("\ngot error: %v\nwant: %v", err, test.expErr)
		}
		if str != test.expSql {
			t.Errorf("\ngot: %v\nwant: %v", str, test.expSql)
		}
	}
}

type myString struct {
	Present bool
	Val     string
}

func (m myString) Value() (driver.Value, error) {
	if m.Present {
		return m.Val, nil
	} else {
		return nil, nil
	}
}
