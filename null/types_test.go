package null

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/mibk/ql"
	"github.com/stretchr/testify/assert"
)

type nullTypedRecord struct {
	Id         int64
	StringVal  String
	Int64Val   Int64
	Float64Val Float64
	TimeVal    Time
	BoolVal    Bool
}

func installFixtures(db *sql.DB) {
	createNullTypesTable := `
		CREATE TABLE null_types (
			id int(11) DEFAULT NULL auto_increment PRIMARY KEY,
			string_val varchar(255) NULL,
			int64_val int(11) NULL,
			float64_val float NULL,
			time_val datetime NULL,
			bool_val bool NULL
		)
	`

	sqlToRun := []string{
		"DROP TABLE IF EXISTS null_types",
		createNullTypesTable,
	}

	for _, v := range sqlToRun {
		_, err := db.Exec(v)
		if err != nil {
			log.Fatalln("Failed to execute statement: ", v, " Got error: ", err)
		}
	}
}

func createRealConnection() *ql.Connection {
	return ql.NewConnection(realDb(), nil)
}

func createRealConnectionWithFixtures() *ql.Connection {
	db := createRealConnection()
	installFixtures(db.DB)
	return db
}

func realDb() *sql.DB {
	driver := os.Getenv("DBR_TEST_DRIVER")
	if driver == "" {
		driver = "mysql"
	}

	dsn := os.Getenv("DBR_TEST_DSN")
	if dsn == "" {
		dsn = "root@/ql_dev?charset=utf8&parseTime=true"
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Fatalln("Mysql error ", err)
	}

	return db
}

func TestNullTypeScanning(t *testing.T) {
	s := createRealConnectionWithFixtures()

	type nullTypeScanningTest struct {
		record *nullTypedRecord
		valid  bool
	}

	tests := []nullTypeScanningTest{
		nullTypeScanningTest{
			record: &nullTypedRecord{},
			valid:  false,
		},
		nullTypeScanningTest{
			record: newNullTypedRecordWithData(),
			valid:  true,
		},
	}

	for _, test := range tests {
		// Create the record in the db
		res, err := s.InsertInto("null_types").Columns("string_val", "int64_val", "float64_val", "time_val", "bool_val").Record(test.record).Exec()
		assert.NoError(t, err)
		id, err := res.LastInsertId()
		assert.NoError(t, err)

		// Scan it back and check that all fields are of the correct validity and are
		// equal to the reference record
		nullTypeSet := &nullTypedRecord{}
		err = s.Select("*").From("null_types").Where("id = ?", id).One(nullTypeSet)
		assert.NoError(t, err)

		// autoincrement ID is not set
		nullTypeSet.Id = 0

		assert.Equal(t, test.record, nullTypeSet)
		assert.Equal(t, test.valid, nullTypeSet.StringVal.Valid)
		assert.Equal(t, test.valid, nullTypeSet.Int64Val.Valid)
		assert.Equal(t, test.valid, nullTypeSet.Float64Val.Valid)
		assert.Equal(t, test.valid, nullTypeSet.TimeVal.Valid)
		assert.Equal(t, test.valid, nullTypeSet.BoolVal.Valid)

		nullTypeSet.StringVal.String = "newStringVal"
		assert.NotEqual(t, test.record, nullTypeSet)
	}
}

func TestNullTypeJSONMarshal(t *testing.T) {
	type nullTypeJSONTest struct {
		record       *nullTypedRecord
		expectedJSON []byte
	}

	tests := []nullTypeJSONTest{
		nullTypeJSONTest{
			record:       &nullTypedRecord{},
			expectedJSON: []byte(`{"Id":0,"StringVal":null,"Int64Val":null,"Float64Val":null,"TimeVal":null,"BoolVal":null}`),
		},
		nullTypeJSONTest{
			record:       newNullTypedRecordWithData(),
			expectedJSON: []byte(`{"Id":0,"StringVal":"wow","Int64Val":42,"Float64Val":1.618,"TimeVal":"2009-01-03T18:15:05Z","BoolVal":true}`),
		},
	}

	for _, test := range tests {
		// Marshal the record
		rawJSON, err := json.Marshal(test.record)
		assert.NoError(t, err)
		assert.Equal(t, test.expectedJSON, rawJSON)

		// Unmarshal it back
		newRecord := &nullTypedRecord{}
		err = json.Unmarshal([]byte(rawJSON), newRecord)
		assert.NoError(t, err)
		assert.Equal(t, test.record, newRecord)
	}
}

func newNullTypedRecordWithData() *nullTypedRecord {
	return &nullTypedRecord{
		StringVal:  String{sql.NullString{String: "wow", Valid: true}},
		Int64Val:   Int64{sql.NullInt64{Int64: 42, Valid: true}},
		Float64Val: Float64{sql.NullFloat64{Float64: 1.618, Valid: true}},
		TimeVal:    Time{mysql.NullTime{Time: time.Date(2009, 1, 3, 18, 15, 5, 0, time.UTC), Valid: true}},
		BoolVal:    Bool{sql.NullBool{Bool: true, Valid: true}},
	}
}
