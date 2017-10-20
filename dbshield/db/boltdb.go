package db

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/boltdb/bolt"
	"github.com/qiwihui/DBShield/dbshield/logger"
	"github.com/qiwihui/DBShield/dbshield/sql"
)

var (
	// DBCon boltdb
	DBCon *bolt.DB
)

//BoltDB local db
type BoltDB struct {
	name string
}

// RecordQueryAction record query and action
func (m *BoltDB) RecordQueryAction(context sql.QueryContext, action string, elapsed time.Duration) error {
	logger.Debugf("action: %s", action)
	// DBCon.Update(func(tx *bolt.Tx) error {
	// 	b := tx.Bucket([]byte("abnormal"))
	// 	if b == nil {
	// 		panic("Invalid DB")
	// 	}
	// 	return nil
	// })
	return errors.New("Not Impletement")
}

// RecordAbnormal record abnormal query
func (m *BoltDB) RecordAbnormal(context sql.QueryContext) error {
	// pattern := sql.Pattern(context.Query)
	atomic.AddUint64(&AbnormalCounter, 1)
	return DBCon.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("abnormal"))
		if b == nil {
			panic("Invalid DB")
		}
		id, _ := b.NextSequence()
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(id))
		return b.Put(buf, context.Marshal())
	})
}

//UpdateState update
func (m *BoltDB) UpdateState() error {
	DBCon.Update(func(tx *bolt.Tx) error {
		//Supplied value must remain valid for the life of the transaction
		qCount := make([]byte, 8)
		abCount := make([]byte, 8)

		b := tx.Bucket([]byte("state"))
		binary.BigEndian.PutUint64(qCount, QueryCounter)
		b.Put([]byte("QueryCounter"), qCount)

		binary.BigEndian.PutUint64(abCount, AbnormalCounter)
		b.Put([]byte("AbnormalCounter"), abCount)

		return nil
	})
	return nil
}

// SyncAndClose local databases
func (m *BoltDB) SyncAndClose() error {
	DBCon.Sync()
	DBCon.Close()
	return nil
}

// CheckPattern check if pattern exist
func (m *BoltDB) CheckPattern(pattern []byte) error {
	return errors.New("Not Impletement")
}

// PutPattern put pattern
func (m *BoltDB) PutPattern(pattern []byte, query []byte) error {
	return errors.New("Not Impletement")
}

// DeletePattern delete pattern
func (m *BoltDB) DeletePattern(pattern []byte) error {
	return DBCon.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("pattern"))
		if b != nil {
			return b.Delete(pattern)
		}
		return nil
	})
}

// Abnormals list Abnormals
func (m *BoltDB) Abnormals() (count int) {

	DBCon.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("abnormal"))
		if b != nil {
			return b.ForEach(func(k, v []byte) error {
				var c sql.QueryContext
				c.Unmarshal(v)
				fmt.Printf("[%s] [User: %s] [Database: %s] %s\n",
					c.Time.Format(time.RFC1123),
					c.User,
					c.Database,
					c.Query)
				count++
				return nil
			})
		}
		return nil
	})
	return count
}

// Patterns list Patterns
func (m *BoltDB) Patterns() (count int) {
	DBCon.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("pattern"))
		if b != nil {
			return b.ForEach(func(k, v []byte) error {
				if strings.Index(string(k), "_client_") == -1 && strings.Index(string(k), "_user_") == -1 {
					fmt.Printf(
						`-----Pattern: 0x%x
Sample: %s
`,
						k,
						v,
					)
					count++
				}
				return nil
			})
		}
		return nil
	})
	return
}

// Purge local database
func (m *BoltDB) Purge() error {
	return errors.New("Not Impletement")
	// return os.Remove(path.Join(config.Config.DBDir,
	// 	config.Config.TargetIP+"_"+config.Config.DBType) + ".db")
}

// AddPattern add
func (m *BoltDB) AddPattern(pattern []byte, context sql.QueryContext) error {
	atomic.AddUint64(&QueryCounter, 1)
	if err := DBCon.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("pattern"))
		if b == nil {
			return errors.New("Invalid DB")
		}
		if b.Get(pattern) == nil {
			b.Put(pattern, context.Query)
		}

		uKey := bytes.Buffer{}
		uKey.Write(pattern)
		uKey.WriteString("_user_")
		uKey.Write(context.User)
		b.Put(uKey.Bytes(), []byte{0x11})

		cKey := bytes.Buffer{}
		cKey.Write(pattern)
		cKey.WriteString("_client_")
		cKey.Write(context.Client)
		b.Put(cKey.Bytes(), []byte{0x11})
		return nil
	}); err != nil {
		logger.Warning(err)
		return err
	}
	return nil
}

// CheckQuery check if Query exist
func (m *BoltDB) CheckQuery(context sql.QueryContext, checkUser bool, checkSource bool) bool {
	atomic.AddUint64(&QueryCounter, 1)
	pattern := sql.Pattern(context.Query)
	if err := DBCon.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("pattern"))
		if b == nil {
			panic("Invalid DB")
		}
		if b.Get(pattern) == nil {
			return errInvalidPattern
		}

		key := bytes.Buffer{}
		if checkUser {
			key.Write(pattern)
			key.WriteString("_user_")
			key.Write(context.User)
			if b.Get(key.Bytes()) == nil {
				return errInvalidUser
			}
		}
		if checkSource {
			key.Reset()
			key.Write(pattern)
			key.WriteString("_client_")
			key.Write(context.Client)
			if b.Get(key.Bytes()) == nil {
				return errInvalidClient
			}
		}
		return nil
	}); err != nil {
		logger.Warning(err)
		//Record abnormal
		// RecordAbnormal(pattern, context)
		return false
	}
	return true
}

//InitialDB local databases
func (m *BoltDB) InitialDB(str string, syncInterval time.Duration, timeout time.Duration) error {
	// logger.Infof("Internal DB: %s", path)
	if DBCon == nil {
		DBCon, _ = bolt.Open(str, 0600, nil)
		DBCon.Update(func(tx *bolt.Tx) error {
			tx.CreateBucketIfNotExists([]byte("pattern"))
			tx.CreateBucketIfNotExists([]byte("abnormal"))
			tx.CreateBucketIfNotExists([]byte("query_action"))
			b, _ := tx.CreateBucketIfNotExists([]byte("state"))
			v := b.Get([]byte("QueryCounter"))
			if v != nil {
				QueryCounter = binary.BigEndian.Uint64(v)
			}
			v = b.Get([]byte("AbnormalCounter"))
			if v != nil {
				AbnormalCounter = binary.BigEndian.Uint64(v)
			}
			return nil
		})
	}

	if syncInterval != 0 {
		DBCon.NoSync = true
		ticker := time.NewTicker(syncInterval)
		go func() {
			for range ticker.C {
				DBCon.Sync()
			}
		}()
	}
	return nil
}
