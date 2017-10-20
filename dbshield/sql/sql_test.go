package sql_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/qiwihui/DBShield/dbshield/sql"
)

func TestClassify(t *testing.T) {
	testcases := []struct {
		sql  string
		want int
	}{
		{"select ...", sql.TypeDML},
		{"    select ...", sql.TypeDML},
		{"insert ...", sql.TypeDML},
		{"replace ....", sql.TypeDML},
		{"   update ...", sql.TypeDML},
		{"Update", sql.TypeDML},
		{"UPDATE ...", sql.TypeDML},
		{"\n\t    delete ...", sql.TypeDML},
		{"", sql.TypeUnknown},
		{" ", sql.TypeUnknown},
		{"begin", sql.TypeUnknown},
		{" begin", sql.TypeUnknown},
		{" begin ", sql.TypeUnknown},
		{"\n\t begin ", sql.TypeUnknown},
		{"... begin ", sql.TypeUnknown},
		{"begin ...", sql.TypeUnknown},
		{"start transaction", sql.TypeUnknown},
		{"commit", sql.TypeTCL},
		{"rollback", sql.TypeTCL},
		{"create", sql.TypeDDL},
		{"alter", sql.TypeDDL},
		{"rename", sql.TypeDDL},
		{"drop", sql.TypeDDL},
		{"set", sql.TypeUnknown},
		{"show", sql.TypeUnknown},
		{"use", sql.TypeUnknown},
		{"analyze", sql.TypeUnknown},
		{"describe", sql.TypeUnknown},
		{"desc", sql.TypeUnknown},
		{"explain", sql.TypeUnknown},
		{"repair", sql.TypeUnknown},
		{"optimize", sql.TypeUnknown},
		{"truncate", sql.TypeDDL},
		{"unknown", sql.TypeUnknown},

		{"/* leading comment */ select ...", sql.TypeDML},
		{"/* leading comment */ /* leading comment 2 */ select ...", sql.TypeDML},
		{"-- leading single line comment \n select ...", sql.TypeDML},
		{"-- leading single line comment \n -- leading single line comment 2\n select ...", sql.TypeDML},

		{"/* leading comment no end select ...", sql.TypeUnknown},
		{"-- leading single line comment no end select ...", sql.TypeUnknown},
	}
	for _, tcase := range testcases {
		if got := sql.Classify(tcase.sql); got != tcase.want {
			t.Errorf("Preview(%s): %v, want %v", tcase.sql, got, tcase.want)
		}
	}
}

func TestQueryContext(t *testing.T) {
	c := sql.QueryContext{
		Query:    []byte("select * from test;"),
		Database: []byte("test"),
		User:     []byte("test"),
		Client:   []byte("127.0.0.1"),
		Time:     time.Now(),
	}
	r := c
	b := c.Marshal()
	c.Unmarshal(b)

	if bytes.Compare(c.Query, r.Query) != 0 {
		t.Error("Expected Query:", r.Query, "got", c.Query)
	}

	if bytes.Compare(c.Query, r.Query) != 0 {
		t.Error("Expected Database:", r.Database, "got", c.Database)
	}

	if bytes.Compare(c.User, r.User) != 0 {
		t.Error("Expected User:", r.User, "got", c.User)
	}

	if bytes.Compare(c.Client, r.Client) != 0 {
		t.Error("Expected Client:", r.Client, "got", c.Client)
	}

	if c.Time.Unix() != r.Time.Unix() {
		t.Error("Expected Time:", r.Time, "got", c.Time)
	}
}

func TestPattern(t *testing.T) {
	p := sql.Pattern([]byte("select * from X;"))
	if len(p) < 4 {
		t.Error("Unexpected Pattern")
	}
}

func BenchmarkPattern(b *testing.B) {
	q := []byte("select * from test;")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sql.Pattern(q)
	}
}
