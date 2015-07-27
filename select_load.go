package ql

import (
	"reflect"
	"time"
)

// Unvetted thots:
// Given a query and given a structure (field list), there's 2 sets of fields.
// Take the intersection. We can fill those in. great.
// For fields in the structure that aren't in the query, we'll let that slide if db:"-".
// For fields in the structure that aren't in the query but without db:"-", return error.
// For fields in the query that aren't in the structure, we'll ignore them.

type loader struct {
	*Connection
	runner
	builder queryBuilder
}

// All executes the query and loads the resulting data into the dest, which can be a slice of
// either structs, or primitive values. It returns n found items (which is not necessarily the
// number of items set).
func (l loader) All(dest interface{}) (n int, err error) {
	valOfDest := reflect.ValueOf(dest)
	if valOfDest.Kind() != reflect.Ptr {
		panic("dest must be a pointer to a slice")
	}

	valOfIndirect := reflect.Indirect(valOfDest)
	if valOfIndirect.Kind() != reflect.Slice {
		panic("dest must be a pointer to a slice")
	}

	originType := valOfIndirect.Type().Elem()
	elemType := originType

	canBeStruct := true
	if originType.Kind() != reflect.Ptr {
		canBeStruct = false
	} else {
		elemType = originType.Elem()
	}

	switch elemType.Kind() {
	case reflect.Struct:
		if !canBeStruct {
			panic("elements of the dest slice must be pointers to structs")
		}
		return l.loadStructs(dest, valOfIndirect, elemType)
	default:
		return l.loadValues(dest, valOfIndirect, originType)
	}
}

// One executes the query and loads the resulting data into the dest, which can be either
// a struct, or a primitive value. Returns ErrNotFound if no item was found, and it was
// therefore not set.
func (l loader) One(dest interface{}) error {
	valOfDest := reflect.ValueOf(dest)
	if valOfDest.Kind() != reflect.Ptr {
		panic("dest must be a pointer")
	}

	valOfIndirect := reflect.Indirect(valOfDest)
	switch valOfIndirect.Kind() {
	case reflect.Struct:
		return l.loadStruct(dest, valOfIndirect)
	default:
		return l.loadValue(dest)
	}
}

// loadStructs executes the query and loads the resulting data into a slice of structs,
// dest must be a pointer to a slice of pointers to structs. It returns the number of items
// found (which is not necessarily the number of items set).
func (l loader) loadStructs(dest interface{}, valueOfDest reflect.Value, elemType reflect.Type) (int, error) {
	fullSql, err := Preprocess(l.builder.ToSql())
	if err != nil {
		return 0, l.EventErr("dbr.select.load_all.interpolate", err)
	}

	numberOfRowsReturned := 0

	startTime := time.Now()
	defer func() { l.TimingKv("dbr.select", time.Since(startTime).Nanoseconds(), kvs{"sql": fullSql}) }()

	rows, err := l.runner.Query(fullSql)
	if err != nil {
		return 0, l.EventErrKv("dbr.select.load_all.query", err, kvs{"sql": fullSql})
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return numberOfRowsReturned, l.EventErrKv("dbr.select.load_one.rows.Columns", err, kvs{"sql": fullSql})
	}

	fieldMap, err := l.calculateFieldMap(elemType, columns, false)
	if err != nil {
		return numberOfRowsReturned, l.EventErrKv("dbr.select.load_all.calculateFieldMap", err, kvs{"sql": fullSql})
	}

	// Build a 'holder', which is an []interface{}. Each value will be the set to address of the field corresponding to our newly made records:
	holder := make([]interface{}, len(fieldMap))

	// Iterate over rows and scan their data into the structs
	sliceValue := valueOfDest
	for rows.Next() {
		// Create a new record to store our row:
		pointerToNewRecord := reflect.New(elemType)
		newRecord := reflect.Indirect(pointerToNewRecord)

		// Prepare the holder for this record
		scannable, err := l.prepareHolderFor(newRecord, fieldMap, holder)
		if err != nil {
			return numberOfRowsReturned, l.EventErrKv("dbr.select.load_all.holderFor", err, kvs{"sql": fullSql})
		}

		// Load up our new structure with the row's values
		err = rows.Scan(scannable...)
		if err != nil {
			return numberOfRowsReturned, l.EventErrKv("dbr.select.load_all.scan", err, kvs{"sql": fullSql})
		}

		// Append our new record to the slice:
		sliceValue = reflect.Append(sliceValue, pointerToNewRecord)

		numberOfRowsReturned++
	}
	valueOfDest.Set(sliceValue)

	// Check for errors at the end. Supposedly these are error that can happen during iteration.
	if err = rows.Err(); err != nil {
		return numberOfRowsReturned, l.EventErrKv("dbr.select.load_all.rows_err", err, kvs{"sql": fullSql})
	}

	return numberOfRowsReturned, nil
}

// loadStruct executes the query and loads the resulting data into a struct,
// dest must be a pointer to a struct. Returns ErrNotFound if nothing was found.
func (l loader) loadStruct(dest interface{}, valueOfDest reflect.Value) error {
	fullSql, err := Preprocess(l.builder.ToSql())
	if err != nil {
		return err
	}

	startTime := time.Now()
	defer func() {
		l.TimingKv("dbr.select", time.Since(startTime).Nanoseconds(), kvs{"sql": fullSql})
	}()

	rows, err := l.runner.Query(fullSql)
	if err != nil {
		return l.EventErrKv("dbr.select.load_one.query", err, kvs{"sql": fullSql})
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return l.EventErrKv("dbr.select.load_one.rows.Columns", err, kvs{"sql": fullSql})
	}

	fieldMap, err := l.calculateFieldMap(valueOfDest.Type(), columns, false)
	if err != nil {
		return l.EventErrKv("dbr.select.load_one.calculateFieldMap", err, kvs{"sql": fullSql})
	}

	// Build a 'holder', which is an []interface{}. Each value will be the set to address of the field corresponding to our newly made records:
	holder := make([]interface{}, len(fieldMap))

	if rows.Next() {
		// Build a 'holder', which is an []interface{}. Each value will be the address of the field corresponding to our newly made record:
		scannable, err := l.prepareHolderFor(valueOfDest, fieldMap, holder)
		if err != nil {
			return l.EventErrKv("dbr.select.load_one.holderFor", err, kvs{"sql": fullSql})
		}

		// Load up our new structure with the row's values
		err = rows.Scan(scannable...)
		if err != nil {
			return l.EventErrKv("dbr.select.load_one.scan", err, kvs{"sql": fullSql})
		}
		return nil
	}

	if err := rows.Err(); err != nil {
		return l.EventErrKv("dbr.select.load_one.rows_err", err, kvs{"sql": fullSql})
	}

	return ErrNotFound
}

// loadValues executes the query and loads the resulting data into a slice of
// primitive values. Returns ErrNotFound if no value was found, and it was therefore not set.
func (l loader) loadValues(dest interface{}, valueOfDest reflect.Value, elemType reflect.Type) (int, error) {
	fullSql, err := Preprocess(l.builder.ToSql())
	if err != nil {
		return 0, err
	}

	numberOfRowsReturned := 0

	startTime := time.Now()
	defer func() { l.TimingKv("dbr.select", time.Since(startTime).Nanoseconds(), kvs{"sql": fullSql}) }()

	rows, err := l.runner.Query(fullSql)
	if err != nil {
		return numberOfRowsReturned, l.EventErrKv("dbr.select.load_all_values.query", err, kvs{"sql": fullSql})
	}
	defer rows.Close()

	sliceValue := valueOfDest
	for rows.Next() {
		// Create a new value to store our row:
		pointerToNewValue := reflect.New(elemType)
		newValue := reflect.Indirect(pointerToNewValue)

		err = rows.Scan(pointerToNewValue.Interface())
		if err != nil {
			return numberOfRowsReturned, l.EventErrKv("dbr.select.load_all_values.scan", err, kvs{"sql": fullSql})
		}

		// Append our new value to the slice:
		sliceValue = reflect.Append(sliceValue, newValue)

		numberOfRowsReturned++
	}
	valueOfDest.Set(sliceValue)

	if err := rows.Err(); err != nil {
		return numberOfRowsReturned, l.EventErrKv("dbr.select.load_all_values.rows_err", err, kvs{"sql": fullSql})
	}

	return numberOfRowsReturned, nil
}

// loadValue executes the query and loads the resulting data into a primitive value.
// Returns ErrNotFound if no value was found, and it was therefore not set.
func (l loader) loadValue(dest interface{}) error {
	fullSql, err := Preprocess(l.builder.ToSql())
	if err != nil {
		return err
	}

	startTime := time.Now()
	defer func() {
		l.TimingKv("dbr.select", time.Since(startTime).Nanoseconds(), kvs{"sql": fullSql})
	}()

	// Run the query:
	rows, err := l.runner.Query(fullSql)
	if err != nil {
		return l.EventErrKv("dbr.select.load_value.query", err, kvs{"sql": fullSql})
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(dest)
		if err != nil {
			return l.EventErrKv("dbr.select.load_value.scan", err, kvs{"sql": fullSql})
		}
		return nil
	}

	if err := rows.Err(); err != nil {
		return l.EventErrKv("dbr.select.load_value.rows_err", err, kvs{"sql": fullSql})
	}

	return ErrNotFound
}
