package db_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/qiwihui/DBShield/dbshield/db"
	"github.com/qiwihui/DBShield/dbshield/sql"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

var mysql = db.MySQL{}
var c1 sql.QueryContext = sql.QueryContext{
	Query:    []byte("select * from test;"),
	Database: []byte("test"),
	User:     []byte("test"),
	Client:   []byte("127.0.0.1"),
	Time:     time.Now(),
}
var dbDsn string = "root:password@tcp(localhost:3306)/dbshield?charset=utf8"

func init() {
	var mysql = db.MySQL{}
	err := mysql.InitialDB(dbDsn, 0, 0)
	if err != nil {
		fmt.Println("Got error", err)
	}
	pattern := sql.Pattern(c1.Query)
	mysql.AddPattern(pattern, c1)
}

func TestCheckQuery(t *testing.T) {
	// test exists
	cExist := sql.QueryContext{
		Query:    []byte("select * from testExist;"),
		Database: []byte("test"),
		User:     []byte("test"),
		Client:   []byte("127.0.0.1"),
		Time:     time.Now(),
	}
	pattern := sql.Pattern(cExist.Query)
	mysql.AddPattern(pattern, cExist)
	if !mysql.CheckQuery(cExist, true, true) {
		t.Error("Expected false")
	}

	// test not exists
	cNotExist := sql.QueryContext{
		Query:    []byte("select * from testNotExist;"),
		Database: []byte("test"),
		User:     []byte("test"),
		Client:   []byte("127.0.0.1"),
		Time:     time.Now(),
	}
	if mysql.CheckQuery(cNotExist, true, true) {
		t.Error("Expected true")
	}
}

func BenchmarkCheckQuery(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mysql.CheckQuery(c1, true, true)
	}
}
