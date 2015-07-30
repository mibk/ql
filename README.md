# ql [![GoDoc](https://godoc.org/github.com/mibk/ql?status.png)](https://godoc.org/github.com/mibk/ql)

This is a fork of the wonderful [gocraft/dbr](https://github.com/gocraft/dbr). I have made some changes to the
package interface to better fit my needs. Before you look at this README be sure to have read the
[original README](https://github.com/gocraft/dbr) for **gocraft/dbr**.

## Instalation

```bash
go get github.com/mibk/ql
```

## Differences from *gocraft/dbr*

### No Session

All methods are available directly on `*ql.Connection`, there is no need to do `conn.NewSession(nil)`. I've found
it not necessary as the `*dbr.Session` is  only a wrapper around `*dbr.Connection` whitch enables setting a
different `EventReceiver`. If it is so, I would suggest a different approach for setting different `EventReceiver`
(currently not provided).

### Query

There is only `Query` method indstead of `SelectBySql` and `UpdateBySql`. For me, they were 2 methods doing the same
thing. `Query` is a superior of these 2. It handles arbitrary SQL statements (not only SELECT and UPDATE, although
the previous 2 were actualy capable of it).

There is `Exec` method for running INSERT, UPDATE, DELETE, ..., and there are methods for loading values returned
by SELECT statement (will be explained later).

### All and One methods

`SelectBuilder` and `Query` have new methods for loading values: `All` and `One`. They replace old methods
for loading (`LoadStructs`, `LoadStruct`, `LoadValues`, `LoadValue`). `All` handles pointer to a slice of
pointers to an arbitrary value (primitive value, or a struct), and `One` works likewise for just one pointer
to an arbitrary value.

Methods for quick returning returning primitive types (`ReturnInt64`, `ReturnStrings`, ...) were remained.

### String methods

For all builders and `Query` there is the String method, which returns an interpolated (and preprocessed) SQL
statement. Useful for debugging (it is possible to just `fmt.Println` a builder).

### Functions for opening DB

There are shortcut functions for opening a DB and creating new `*Connection` (`Open`, `MustOpen`, and
`MustOpenAndVerify`, which performs a *ping* to the database).

### Quoting identifiers (columns)

It is possible to escape columns using brackets.

```go
var u User
db.Query(`SELECT [name], [age], [from] AS [birth_place]
	FROM [user]
	WHERE [id] = ?`, 15).One(&u)
// It executes:
// SELECT `name`, `age`, `from` AS `birth_place` FROM `user` WHERE `id` = 15
```

I have found it convenient to write SQL statements on multiple lines for readability reasons. But that is only
possible using *backticks* that are also used for column escaping.

### Null types in `null` package

Null* types where moved to the `ql/null` package. Example of usage:

```go
import "github.com/mibk/ql/null"

type Nonsence struct {
	Number      null.Int64
	Description null.String
}
```

### Updated `Where` (and `Having`) method

`Where` now accepts only string or `ql.And` type as the first argument. `ql.And` is the replacement for
the original `Eq` type. It's because `ql.And` goes beyond an equality map. When an ql.And map is the
input, it is only a shorthand for calling the `Where` (`Having`) method multiple times with 2 arguments:
an expression string and one parameter. Example:

```go
b.Where(ql.And{"name": "Igor", "age >": 50})
// is equivalent to:
b.Where("name", "Igor")
b.Where("age >", 50)
```

As you can see in these simple examples it is not required to type an operator, or a placeholder. That
is so called "short notation" of the format:
```
column_name [operator] [?]
```
Default operator is `=` and typing a placeholder is just optional. If the expression doesn't match
this format, it works as before.

### Order method

`OrderBy` methods works as before, but there is new method `Order` which replaces the `OrderDir`.
It expects the `ql.By` type as a parameter which is a map of column names and order directions.
Example:

```go
b.Order(ql.By{"col1:": ql.Asc, "col2:": ql.Desc})
// ORDER BY `col1` ASC, `col2` DESC
```

### Removed objects
* `Paginate` on `*SelectBuilder` was removed. It should probably be in another layer.
* No `NullTime` as it was dependent on the *mysql driver*.

## Quickstart

```go
package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mibk/ql"
	"github.com/mibk/ql/null"
)

type User struct {
	Id      int64
	Name    string
	Comment null.String
}

var conn *ql.Connection

func main() {
	// It panics on error and after successful opening it tries to ping the database to
	// test the connection.
	conn = ql.MustOpenAndVerify("mysql", "root@/your_database")

	var u User
	b := conn.Select("id, [title]").From("user").Where("id = ?", 13)

	// This will print the interpolated query:
	//     SELECT id, `title` FROM user WHERE (id = 13)
	fmt.Println(b)

	// Method One will execute the query and load the result to the u struct.
	// If it was not set before, it also sets LIMIT to 1 as there is no need
	// for returning multiple rows.
	//     SELECT id, `title` FROM user WHERE (id = 13) LIMIT 1
	if err := b.One(&u); err != nil {
		panic(err)
	}
	fmt.Printf("User: %+v", u)
}
```

## Driver support

Currently only MySQL is supported. I am planning to move all driver dependent parts to `gl/dialect` package.

## Authors:

* Jonathan Novak ([github](https://github.com/cypriss))
* Tyler Smith ([github](https://github.com/tyler-smith))
* Michal Bohuslávek ([github](https://github.com/mibk))
* Sponsored by [UserVoice](https://eng.uservoice.com)




