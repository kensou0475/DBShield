package db_test

import (
	"testing"
	"time"

	"github.com/qiwihui/DBShield/dbshield/db"
	"github.com/qiwihui/DBShield/dbshield/sql"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

var dbDsn string = "root:password@tcp(localhost:3306)/dbshield?charset=utf8"

func TestInitalMysqlDB(t *testing.T) {
	mysqlDsn := dbDsn
	mysql := new(db.MySQL)
	err := mysql.InitialDB(mysqlDsn, 0, 0)
	if err != nil {
		t.Error("Got error", err)
	}
}

func TestCheckQuery(t *testing.T) {
	var mysql = db.MySQL{}
	err := mysql.InitialDB(dbDsn, 0, 0)
	if err != nil {
		t.Error("Got error", err)
	}

	c1 := sql.QueryContext{
		Query:    []byte("select * from test;"),
		Database: []byte("test"),
		User:     []byte("test"),
		Client:   []byte("127.0.0.1"),
		Time:     time.Now(),
	}
	c2 := sql.QueryContext{
		Query:    []byte("select * from user;"),
		Database: []byte("test"),
		User:     []byte("test"),
		Client:   []byte("127.0.0.1"),
		Time:     time.Now(),
	}
	pattern := sql.Pattern(c1.Query)
	mysql.AddPattern(pattern, c1)
	if !mysql.CheckQuery(c1, true, true) {
		t.Error("Expected false")
	}
	if mysql.CheckQuery(c2, true, true) {
		t.Error("Expected true")
	}

}

func BenchmarkCheckQuery(b *testing.B) {
	var mysql = db.MySQL{}
	err := mysql.InitialDB(dbDsn, 0, 0)
	if err != nil {
		b.Error("Got error", err)
	}

	c1 := sql.QueryContext{
		Query:    []byte("select * from test;"),
		Database: []byte("test"),
		User:     []byte("test"),
		Client:   []byte("127.0.0.1"),
		Time:     time.Now(),
	}
	pattern := sql.Pattern(c1.Query)
	mysql.AddPattern(pattern, c1)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mysql.CheckQuery(c1, true, true)
	}
}
