package query

import "io"

// Writer is used to write a query.
type Writer interface {
	io.Writer
	WriteString(s string) (n int, err error)
	WriteRune(r rune) (n int, err error)
}
