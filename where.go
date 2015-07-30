package ql

import (
	"reflect"
	"regexp"

	"github.com/mibk/ql/query"
)

// And is a map column -> value pairs which must be matched in a query.
type And map[string]interface{}

type whereFragment struct {
	Condition string
	Values    []interface{}
}

type addFragmentFn func(expr string, args ...interface{})

func handleExprType(exprOrMap interface{}, args []interface{}, fn addFragmentFn) {
	switch v := exprOrMap.(type) {
	case string:
		fn(v, args...)
	case And:
		if len(args) > 0 {
			panic("args are not expected when passing an And map")
		}
		for ex, arg := range v {
			fn(ex, arg)
		}
	default:
		panic("invalid argument passed to Where, only a string or an And map is allowed")
	}
}

var shortNotation = regexp.MustCompile(`^\s*([a-zA-Z._]+)\s*([a-zA-Z=<>!]+)?\s*\??\s*$`)

func handleShortNotation(expr string, args []interface{}) (string, []interface{}) {
	if len(args) == 1 {
		if m := shortNotation.FindStringSubmatch(expr); m != nil {
			col, op := m[1], m[2]
			if op == "" {
				op = "="
			}
			expr = "[" + col + "]"

			arg := args[0]
			if arg == nil {
				expr += " IS NULL"
				args = args[:0]
			} else {
				v := reflect.ValueOf(arg)
				if v.Kind() == reflect.Array || v.Kind() == reflect.Slice {
					if v.Len() == 0 {
						if v.IsNil() {
							expr += " IS NULL"
						} else {
							expr = "1=0"
						}
						args = args[:0]
					} else {
						expr += " IN ?"
					}
				} else {
					expr += " " + op + " ?"
				}
			}
		}
	}
	return expr, args
}

// Invariant: only called when len(fragments) > 0.
func writeWhereFragmentsToSql(w query.Writer, fragments []*whereFragment, args *[]interface{}) {
	anyConditions := false
	for _, f := range fragments {
		if f.Condition == "" {
			panic("invalid condition expression")
		}
		if anyConditions {
			w.WriteString(" AND ")
		}
		anyConditions = true
		w.WriteString("(" + f.Condition + ")")
		if len(f.Values) > 0 {
			*args = append(*args, f.Values...)
		}
	}
}
