package ql

import "reflect"

// Eq is a map column -> value pairs which must be matched in a query.
type Eq map[string]interface{}

type whereFragment struct {
	Condition   string
	Values      []interface{}
	EqualityMap map[string]interface{}
}

func newWhereFragment(whereSqlOrMap interface{}, args []interface{}) *whereFragment {
	switch pred := whereSqlOrMap.(type) {
	case string:
		return &whereFragment{Condition: pred, Values: args}
	case map[string]interface{}:
		return &whereFragment{EqualityMap: pred}
	case Eq:
		return &whereFragment{EqualityMap: map[string]interface{}(pred)}
	default:
		panic("Invalid argument passed to Where. Pass a string or an Eq map.")
	}

	return nil
}

// Invariant: only called when len(fragments) > 0.
func writeWhereFragmentsToSql(fragments []*whereFragment, w queryWriter, args *[]interface{}) {
	anyConditions := false
	for _, f := range fragments {
		if f.Condition != "" {
			if anyConditions {
				w.WriteString(" AND (")
			} else {
				w.WriteRune('(')
				anyConditions = true
			}
			w.WriteString(f.Condition)
			w.WriteRune(')')
			if len(f.Values) > 0 {
				*args = append(*args, f.Values...)
			}
		} else if f.EqualityMap != nil {
			anyConditions = writeEqualityMapToSql(f.EqualityMap, w, args, anyConditions)
		} else {
			panic("invalid equality map")
		}
	}
}

func writeEqualityMapToSql(eq map[string]interface{}, w queryWriter, args *[]interface{}, anyConditions bool) bool {
	for k, v := range eq {
		if v == nil {
			anyConditions = writeWhereCondition(w, k, " IS NULL", anyConditions)
		} else {
			vVal := reflect.ValueOf(v)

			if vVal.Kind() == reflect.Array || vVal.Kind() == reflect.Slice {
				vValLen := vVal.Len()
				if vValLen == 0 {
					if vVal.IsNil() {
						anyConditions = writeWhereCondition(w, k, " IS NULL", anyConditions)
					} else {
						if anyConditions {
							w.WriteString(" AND (1=0)")
						} else {
							w.WriteString("(1=0)")
						}
					}
				} else if vValLen == 1 {
					anyConditions = writeWhereCondition(w, k, " = ?", anyConditions)
					*args = append(*args, vVal.Index(0).Interface())
				} else {
					anyConditions = writeWhereCondition(w, k, " IN ?", anyConditions)
					*args = append(*args, v)
				}
			} else {
				anyConditions = writeWhereCondition(w, k, " = ?", anyConditions)
				*args = append(*args, v)
			}
		}
	}

	return anyConditions
}

func writeWhereCondition(w queryWriter, k string, pred string, anyConditions bool) bool {
	if anyConditions {
		w.WriteString(" AND (")
	} else {
		w.WriteRune('(')
		anyConditions = true
	}
	Quoter.writeQuotedColumn(k, w)
	w.WriteString(pred)
	w.WriteRune(')')

	return anyConditions
}
