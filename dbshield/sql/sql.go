package sql

import (
	"bytes"
	"encoding/binary"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/xwb1989/sqlparser"
)

//Query
const (
	TypeDDL = iota
	TypeDML
	TypeDCL
	TypeTCL
	TypeUnknown
)

//Classify sql commmands
func Classify(sql string) int {
	// DDL包括：CREATE，ALTER，DROP，TRUNCATE，COMMENT，RENAME；
	// DML包括：SELECT，INSERT，UPDATE，DELETE，MERGE，CALL，EXPLAIN PLAN，LOCK TABLE；
	// DCL包括：GRANT，REVOKE
	// TCL包括：COMMIT，ROLLBACK，SAVEPOINT，SET TRANSACTION

	trimmed := sqlparser.StripLeadingComments(sql)

	firstWord := trimmed
	if end := strings.IndexFunc(trimmed, unicode.IsSpace); end != -1 {
		firstWord = trimmed[:end]
	}

	// Comparison is done in order of priority.
	loweredFirstWord := strings.ToLower(firstWord)
	switch loweredFirstWord {
	case "select", "insert", "update", "delete", "replace", "merge", "call":
		return TypeDML
	case "commit", "rollback", "SAVEPOINT":
		return TypeTCL
	case "create", "alter", "drop", "truncate", "comment", "rename":
		return TypeDDL
	case "grant", "revoke":
		return TypeDCL
	}
	switch strings.ToLower(trimmed) {
	case "explain plan", "lock table":
		return TypeDML
	case "set transaction":
		return TypeTCL
	}
	return TypeUnknown
}

//QueryContext holds information around query
type QueryContext struct {
	Query    []byte
	User     []byte
	Client   []byte
	Database []byte
	Time     time.Time
}

// QueryAction action and duration
type QueryAction struct {
	QueryContext
	Action   string
	Duration time.Duration
}

//Unmarshal []byte into QueryContext
func (c *QueryContext) Unmarshal(b []byte) (size uint32) {
	n := binary.BigEndian.Uint32(b)
	b = b[4:]

	c.Query = b[:n]
	size = n

	b = b[n:]
	n = binary.BigEndian.Uint32(b)
	b = b[4:]
	c.User = b[:n]
	size += n

	b = b[n:]
	n = binary.BigEndian.Uint32(b)
	b = b[4:]
	c.Client = b[:n]
	size += n

	b = b[n:]
	n = binary.BigEndian.Uint32(b)
	b = b[4:]
	c.Database = b[:n]
	size += n

	c.Time.UnmarshalBinary(b[n:])
	size += 8
	return
}

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

//Marshal load []byte into QueryContext
func (c *QueryContext) Marshal() []byte {
	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)
	buf.Reset()
	l := make([]byte, 4)
	binary.BigEndian.PutUint32(l, uint32(len(c.Query)))
	buf.Write(l)
	buf.Write(c.Query)

	binary.BigEndian.PutUint32(l, uint32(len(c.User)))
	buf.Write(l)
	buf.Write(c.User)

	binary.BigEndian.PutUint32(l, uint32(len(c.Client)))
	buf.Write(l)
	buf.Write(c.Client)

	binary.BigEndian.PutUint32(l, uint32(len(c.Database)))
	buf.Write(l)
	buf.Write(c.Database)

	t, _ := c.Time.MarshalBinary()
	buf.Write(t)
	return buf.Bytes()
}

//Pattern returns pattern of given query
func Pattern(query []byte) []byte {
	tokenizer := sqlparser.NewStringTokenizer(string(query))
	buf := bytes.Buffer{}
	l := make([]byte, 4)
	for {
		typ, val := tokenizer.Scan()
		switch typ {
		case sqlparser.ID: //table, database, variable & ... names
			buf.Write(val)
		case 0: //End of query
			return buf.Bytes()
		default:
			binary.BigEndian.PutUint32(l, uint32(typ))
			buf.Write(l)
		}
	}
}

// GetTableName from sql statement
func GetTableName(query string) (string, error) {
	tree, err := sqlparser.Parse(query)
	if err != nil {
		return "", err
	}
	out := sqlparser.GetTableName(tree.(*sqlparser.Select).From[0].(*sqlparser.AliasedTableExpr).Expr)
	return out.String(), nil
}
