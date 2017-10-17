package training_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/boltdb/bolt"
	"github.com/qiwihui/DBShield/dbshield/config"
	"github.com/qiwihui/DBShield/dbshield/db"
	"github.com/qiwihui/DBShield/dbshield/sql"
	"github.com/qiwihui/DBShield/dbshield/training"
)

var c0 sql.QueryContext = sql.QueryContext{
	Query:    []byte("select * from test;"),
	Database: []byte("test"),
	User:     []byte("test"),
	Client:   []byte("127.0.0.1"),
	Time:     time.Now(),
}

var mysql = db.MySQL{}
var boltDB = db.BoltDB{}

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard) // Avoid log outputs

	var dbDsn string = "root:password@tcp(localhost:3306)/dbshield?charset=utf8"
	// var mysql = db.MySQL{}
	// config.Config.LocalDB = &mysql
	err := mysql.InitialDB(dbDsn, 0, 0)
	if err != nil {
		fmt.Println("Got error", err)
	}
	pattern := sql.Pattern(c0.Query)
	mysql.AddPattern(pattern, c0)

	tmpfile, err := ioutil.TempFile("", "testdb")
	if err != nil {
		panic(err)
	}
	defer tmpfile.Close()
	path := tmpfile.Name()
	// boltDB := new(db.BoltDB)
	// config.Config.LocalDB = &mysql
	err = boltDB.InitialDB(path, 0, 0)
	if err != nil {
		fmt.Println("Got error", err)
	}

	os.Exit(m.Run())
}

func switchDB(usingDB string) {
	if usingDB == "mysql" {
		config.Config.LocalDB = &mysql
	} else {
		config.Config.LocalDB = &boltDB
	}
}

func TestAddToTrainingSetWithBoltDB(t *testing.T) {
	switchDB("boltdb")
	var err error
	c := sql.QueryContext{
		Query:    []byte("select * from ExpectedTest;"),
		Database: []byte("test"),
		User:     []byte("test"),
		Client:   []byte("127.0.0.1"),
		Time:     time.Now(),
	}
	err = training.AddToTrainingSet(c)
	if err != nil {
		t.Error("Not Expected error", err)
	}

	// Add without buckets
	tmpDB := db.DBCon
	defer func() {
		db.DBCon = tmpDB
	}()
	tmpfile, err := ioutil.TempFile("", "testdbWithoutBuckets")
	if err != nil {
		panic(err)
	}
	defer tmpfile.Close()
	path := tmpfile.Name()
	db.DBCon, err = bolt.Open(path, 0600, nil)
	if err != nil {
		panic(err)
	}
	err = training.AddToTrainingSet(c)
	if err == nil {
		t.Error("Expected error")
	}
}

func TestCheckQueryWithBotDB(t *testing.T) {
	switchDB("boltdb")
	config.Config.CheckUser = true
	config.Config.CheckSource = true
	c1 := sql.QueryContext{
		Query:    []byte("select * from c1test;"),
		Database: []byte("test"),
		User:     []byte("test"),
		Client:   []byte("127.0.0.1"),
		Time:     time.Now(),
	}
	c2 := sql.QueryContext{
		Query:    []byte("select * from c2user;"),
		Database: []byte("test"),
		User:     []byte("test"),
		Client:   []byte("127.0.0.1"),
		Time:     time.Now(),
	}
	training.AddToTrainingSet(c1)
	if !training.CheckQuery(c1) {
		t.Error("Expected false")
	}
	if training.CheckQuery(c2) {
		t.Error("Expected true")
	}

	tmpDB := db.DBCon
	defer func() {
		db.DBCon = tmpDB
	}()
	tmpfile, err := ioutil.TempFile("", "testdb")
	if err != nil {
		panic(err)
	}
	defer tmpfile.Close()
	path := tmpfile.Name()
	db.DBCon, err = bolt.Open(path, 0600, nil)
	if err != nil {
		panic(err)
	}
	db.DBCon.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket([]byte("pattern"))
		return err
	})
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic")
		}
	}()
	training.CheckQuery(c1)
}

func BenchmarkAddToTrainingSetWithMysql(b *testing.B) {
	switchDB("mysql")
	o := orm.NewOrm()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		training.AddToTrainingSet(c0)
		b.StopTimer()
		// clean db to add instead of just checking then passed
		_ = o.Raw("delete from pattern")
		b.StartTimer()
	}
}

func BenchmarkAddToTrainingSetWithBoltDB(b *testing.B) {
	switchDB("boltdb")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		training.AddToTrainingSet(c0)
		b.StopTimer()
		// clean db to add instead of just checking then passed
		db.DBCon.Update(func(tx *bolt.Tx) error {
			err := tx.DeleteBucket([]byte("pattern"))
			if err != nil {
				return fmt.Errorf("delete bucket: %s", err)
			}
			_, err = tx.CreateBucket([]byte("pattern"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		})
		b.StartTimer()
	}
}
