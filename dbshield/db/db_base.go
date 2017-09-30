package db

import (
	"encoding/binary"
	"errors"
	"net"
	"time"

	"github.com/qiwihui/DBShield/dbshield/sql"
)

var (
	//QueryCounter state
	QueryCounter = uint64(0)

	//AbnormalCounter state
	AbnormalCounter = uint64(0)

	errInvalidPattern = errors.New("Invalid pattern")
	errInvalidUser    = errors.New("Invalid user")
	errInvalidClient  = errors.New("Invalid client")
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
	CheckQuery(sql.QueryContext, bool, bool) bool
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

func fourByteBigEndianToIP(data []byte) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, binary.BigEndian.Uint32(data))
	return ip.String()
}
