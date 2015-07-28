// Package null defines Null variants of the primitive types. It is supposed
// that users may define their own types to satisfy their needs. The purpose
// of these types is to give the basic idea on how to define custom types.
//
// The advantage of these types over original database/sql Null* types is
// that (un)marshaling to JSON works as expected. So, instead of:
// 	{
// 		"title": {
// 			"Valid": true,
// 			"String": "Octopus"
// 		},
// 		"description": {
// 			"Valid": false,
// 			"String": ""
// 		}
// 	}
// it is now marshaled like this:
// 	{
// 		"title": "Octopus",
// 		"description": null
// 	}
//
// Other kinds of marshaling (e.g. XML) are not implemented yet. Anyway, it is not
// possible (nor desired) to cover everyone's needs for Null* types, so it is recommended
// to define own types.
package null

import (
	"database/sql"
	"encoding/json"
)

var nullString = []byte("null")

// String is a type that can be null or a string.
type String struct {
	sql.NullString
}

// MarshalJSON correctly serializes a String to JSON.
func (n *String) MarshalJSON() ([]byte, error) {
	if n.Valid {
		j, e := json.Marshal(n.String)
		return j, e
	}
	return nullString, nil
}

// UnmarshalJSON correctly deserializes a String from JSON.
func (n *String) UnmarshalJSON(b []byte) error {
	var s interface{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return n.Scan(s)
}

// Float64 is a type that can be null or a float64.
type Float64 struct {
	sql.NullFloat64
}

// MarshalJSON correctly serializes a Float64 to JSON.
func (n *Float64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		j, e := json.Marshal(n.Float64)
		return j, e
	}
	return nullString, nil
}

// UnmarshalJSON correctly deserializes a Float64 from JSON.
func (n *Float64) UnmarshalJSON(b []byte) error {
	var s interface{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return n.Scan(s)
}

// Int64 is a type that can be null or an int.
type Int64 struct {
	sql.NullInt64
}

// MarshalJSON correctly serializes a Int64 to JSON.
func (n *Int64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		j, e := json.Marshal(n.Int64)
		return j, e
	}
	return nullString, nil
}

// UnmarshalJSON correctly deserializes a Int64 from JSON.
func (n *Int64) UnmarshalJSON(b []byte) error {
	var s interface{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return n.Scan(s)
}

// Bool is a type that can be null or a bool.
type Bool struct {
	sql.NullBool
}

// MarshalJSON correctly serializes a Bool to JSON.
func (n *Bool) MarshalJSON() ([]byte, error) {
	if n.Valid {
		j, e := json.Marshal(n.Bool)
		return j, e
	}
	return nullString, nil
}

// UnmarshalJSON correctly deserializes a Bool from JSON.
func (n *Bool) UnmarshalJSON(b []byte) error {
	var s interface{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return n.Scan(s)
}
