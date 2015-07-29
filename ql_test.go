package ql

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mibk/ql/null"
)

// Test helpers

// Returns a dbion that's not backed by a database.
func createFakeConnection() *Connection {
	return NewConnection(nil, nil)
}

func createRealConnection() *Connection {
	return NewConnection(realDb(), nil)
}

func createRealConnectionWithFixtures() *Connection {
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

type dbrPerson struct {
	Id    int64
	Name  string
	Email null.String
	Key   null.String
}

func installFixtures(db *sql.DB) {
	createPeopleTable := fmt.Sprintf(`
		CREATE TABLE dbr_people (
			id int(11) DEFAULT NULL auto_increment PRIMARY KEY,
			name varchar(255) NOT NULL,
			email varchar(255),
			%s varchar(255)
		)
	`, "`key`")

	sqlToRun := []string{
		"DROP TABLE IF EXISTS dbr_people",
		createPeopleTable,
		"INSERT INTO dbr_people (name,email) VALUES ('Jonathan', 'jonathan@uservoice.com')",
		"INSERT INTO dbr_people (name,email) VALUES ('Dmitri', 'zavorotni@jadius.com')",
	}

	for _, v := range sqlToRun {
		_, err := db.Exec(v)
		if err != nil {
			log.Fatalln("Failed to execute statement: ", v, " Got error: ", err)
		}
	}
}
