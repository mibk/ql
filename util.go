package ql

import "io"

// NameMapping is the routine to use when mapping column names to struct properties
var NameMapping = camelCaseToSnakeCase

func camelCaseToSnakeCase(name string) string {
	var newstr []rune
	firstTime := true

	for _, chr := range name {
		if isUpper := 'A' <= chr && chr <= 'Z'; isUpper {
			if firstTime == true {
				firstTime = false
			} else {
				newstr = append(newstr, '_')
			}
			chr -= ('A' - 'a')
		}
		newstr = append(newstr, chr)
	}

	return string(newstr)
}

// queryWriter is used to write a query.
type queryWriter interface {
	io.Writer
	WriteString(s string) (n int, err error)
	WriteRune(r rune) (n int, err error)
}
