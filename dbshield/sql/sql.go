package sql

import (
	"bytes"
	"encoding/binary"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/xwb1989/sqlparser"
)

//Query
const (
	TypeDDL     = "DDL"
	TypeDML     = "DML"
	TypeDCL     = "DCL"
	TypeTCL     = "TCL"
	TypeUnknown = "Unknown"
)

//Statement
const (
	StmtSelect = iota
	StmtInsert
	StmtReplace
	StmtMerge
	StmtCall
	StmtUpdate
	StmtDelete
	StmtDDL
	StmtBegin
	StmtCommit
	StmtRollback
	StmtSavePoint
	StmtSet
	StmtShow
	StmtUse
	StmtCreate
	StmtAlert
	StmtDrop
	StmtTruncate
	StmtComment
	StmtRename
	StmtGrant
	StmtRevoke
	StmtExplainPlan
	StmtLockTable
	StmtSetTransaction
	StmtOther
	StmtUnknown
)

//Classify 分类sql语句
func Classify(sql string) string {
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
		return "DML"
	case "commit", "rollback", "savepoint":
		return "TCL"
	case "create", "alter", "drop", "truncate", "comment", "rename":
		return "DDL"
	case "grant", "revoke":
		return "DCL"
	}
	switch strings.ToLower(trimmed) {
	case "explain plan", "lock table":
		return "DML"
	case "set transaction":
		return "TCL"
	}
	return "Unknown"
}

//GetType 获取sql语句类型
func GetType(sql string) int {
	trimmed := sqlparser.StripLeadingComments(sql)

	firstWord := trimmed
	if end := strings.IndexFunc(trimmed, unicode.IsSpace); end != -1 {
		firstWord = trimmed[:end]
	}

	// Comparison is done in order of priority.
	loweredFirstWord := strings.ToLower(firstWord)
	switch loweredFirstWord {
	case "select":
		return StmtSelect
	case "insert":
		return StmtInsert
	case "replace":
		return StmtReplace
	case "merge":
		return StmtMerge
	case "call":
		return StmtCall
	case "update":
		return StmtUpdate
	case "delete":
		return StmtDelete
	case "commit":
		return StmtCommit
	case "rollback":
		return StmtRollback
	case "savepoint":
		return StmtSavePoint
	case "create":
		return StmtCreate
	case "alert":
		return StmtAlert
	case "drop":
		return StmtDrop
	case "truncate":
		return StmtTruncate
	case "comment":
		return StmtComment
	case "rename":
		return StmtRename
	case "grant":
		return StmtGrant
	case "revoke":
		return StmtRevoke
	}
	switch strings.ToLower(trimmed) {
	case "explain plan":
		return StmtExplainPlan
	case "lock table":
		return StmtLockTable
	case "set transaction":
		return StmtSetTransaction
	}
	return StmtUnknown
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

// GetTableName from sql select statement
func GetTableName(query string) (string, error) {
	tree, err := sqlparser.Parse(query)
	if err != nil {
		return "", err
	}
	out := sqlparser.GetTableName(tree.(*sqlparser.Select).From[0].(*sqlparser.AliasedTableExpr).Expr)
	return out.String(), nil
}

// ExtractTableNames from sql statement
func ExtractTableNames(fromSQL string) (rets []string, err error) {
	rets = []string{}
	// TODO 增加对sql union时候的检测
	reg := "(?i)((alter|drop|create)\\s+table\\s+(?P<table0>\\w+)(\\s+(as\\s+)?(?P<alias0>\\w+))?)|((update|\\s+(from|into))\\s+(?P<table1>\\w+)(\\s+(as\\s+)?(?P<alias1>\\w+))?(\\s+(union|where|left|right|outer|inner))?)|(\\s+join\\s+(?P<table2>\\w+)\\s+((as\\s+)?(\\w+)\\s+)?on)"
	r, err := regexp.Compile(reg)
	if err != nil {
		return []string{}, err
	}
	n1 := r.SubexpNames()
	finds := r.FindAllStringSubmatch(fromSQL, -1)
	for _, v := range finds {
		for ii, vv := range v {
			if (n1[ii] == "table0" || n1[ii] == "table1" || n1[ii] == "table2") && vv != "" {
				// fmt.Println("=>", n1[ii], vv)
				rets = append(rets, string(vv))
			}

		}
	}
	return
}
