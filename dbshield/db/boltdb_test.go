package db_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/qiwihui/DBShield/dbshield/db"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard) // Avoid log outputs
	m.Run()
}

func TestInitalDB(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "testdb")
	if err != nil {
		panic(err)
	}
	defer tmpfile.Close()
	path := tmpfile.Name()
	boltDB := new(db.BoltDB)
	err = boltDB.InitialDB(path, 0, 0)
	if err != nil {
		t.Error("Got error", err)
	}
}
