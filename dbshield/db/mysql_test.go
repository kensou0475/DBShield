package db_test

import (
	"testing"

	"github.com/qiwihui/DBShield/dbshield/db"
)

func TestInitalMysqlDB(t *testing.T) {

	mysqlDsn := "root:password@tcp(localhost:3306)/dbshield?charset=utf8"
	mysql := new(db.MySQL)
	err := mysql.InitialDB(mysqlDsn, 0, 0)
	if err != nil {
		t.Error("Got error", err)
	}
}
