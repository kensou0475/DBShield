package db

import (
	"errors"
	"time"

	"github.com/qiwihui/DBShield/dbshield/sql"
)

var (
	//QueryCounter state
	QueryCounter = uint64(0)

	//AbnormalCounter state
	AbnormalCounter = uint64(0)

	errInvalidPattern = errors.New("Invalid pattern")
)

//BASE interface should get implemented with every added store database(Boltdb, MySQL, Postgre & etc.) structure
type BASE interface {
	InitialDB(string, time.Duration, time.Duration) error
	RecordQueryAction(sql.QueryContext, string) error
	Abnormals() int
	RecordAbnormal(sql.QueryContext) error
	Patterns() int
	CheckPattern([]byte) error
	AddPattern([]byte, sql.QueryContext) error
	PutPattern([]byte, []byte) error
	DeletePattern([]byte) error
	Purge() error
	CheckQuery(sql.QueryContext) bool
	UpdateState() error
	SyncAndClose() error
}

//GenerateLocalDB generate local db
func GenerateLocalDB(dbName string) BASE {
	switch dbName {
	case "mysql":
		return new(MySQL)
	case "boltdb":
		return new(BoltDB)
	default:
		return nil
	}
}
