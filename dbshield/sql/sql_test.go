package sql_test

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/qiwihui/DBShield/dbshield/sql"
)

func TestClassify(t *testing.T) {
	testcases := []struct {
		sql  string
		want string
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

func TestGetTableName(t *testing.T) {
	testcases := []struct {
		in, out string
	}{{
		in:  "select * from t",
		out: "t",
	}, {
		in:  "select * from t.t",
		out: "",
	}, {
		in:  "select * from (select * from t) as tt",
		out: "",
	}, {
		in:  "insert into t (num) values (1)",
		out: "",
	}}

	for _, tc := range testcases {
		out, err := sql.GetTableName(tc.in)
		if err != nil {
			t.Error(err)
			continue
		}
		if out != tc.out {
			t.Errorf("GetTableName('%s'): %s, want %s", tc.in, out, tc.out)
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

func TestExtractTableNames(t *testing.T) {
	testcases := []struct {
		in  string
		out []string
	}{
		{
			in:  "select * from test",
			out: []string{"test"},
		},
		{
			in:  "SELECT * FROM test",
			out: []string{"test"},
		},
		{
			in:  "select * from test as t where id=3;",
			out: []string{"test"},
		},
		{
			in:  "delete * from test",
			out: []string{"test"},
		},
		{
			in:  "select * from (select * from test) as tt",
			out: []string{"test"},
		},
		{
			in:  "insert into test (num) values (1)",
			out: []string{"test"},
		},
		{
			in:  "update test set name='test' where id=1",
			out: []string{"test"},
		},
		{
			in:  "UPDATE test SET name='test' WHERE id=1",
			out: []string{"test"},
		},
		{
			in: `select * from table1 l
				left join table2 p on l.patid=p.patientid
				join table3 c on   l.patid=c.patid inner join table4 ph on l.patid=ph.patid
			 	from table5 p where p.outvisitid=l.outvisitid) unit all;`,
			out: []string{"table1", "table2", "table3", "table4", "table5"},
		},
		{
			in:  "select count(*) as count from table1 as a left join table2 as b on a.user_id = b.user_id where a.title like '%Cloth%'",
			out: []string{"table1", "table2"},
		},
		{
			in: `CREATE TABLE test (
				PersonID int,
				LastName varchar(255),
				FirstName varchar(255),
				Address varchar(255),
				City varchar(255) 
			);`,
			out: []string{"test"},
		},
		{
			in:  "DROP TABLE test;",
			out: []string{"test"},
		},
		{
			in:  "ALTER TABLE test ADD column_name datatype;",
			out: []string{"test"},
		},
		{
			in:  "CREATE DATABASE testDB;",
			out: []string{},
		},
		{
			in:  "DROP DATABASE testDB;",
			out: []string{},
		},
		{
			in: `SELECT column_name(s)
			FROM test1
			LEFT JOIN test2 ON table1.column_name = table2.column_name;`,
			out: []string{"test1", "test2"},
		},
		{
			in: `SELECT column_name(s)
			FROM test1
			RIGHT JOIN test2 ON table1.column_name = table2.column_name;`,
			out: []string{"test1", "test2"},
		},
		{
			in: `SELECT column_name(s)
			FROM test1
			FULL OUTER JOIN test2 ON table1.column_name = table2.column_name;`,
			out: []string{"test1", "test2"},
		},
		// {
		// 	in: `SELECT A.CustomerName AS CustomerName1, B.CustomerName AS CustomerName2, A.City
		// 	FROM Customers A, Customers B
		// 	WHERE A.CustomerID <> B.CustomerID
		// 	AND A.City = B.City
		// 	ORDER BY A.City;`,
		// 	out: []string{},
		// },
		{
			in: `SELECT 'Customer' As Type, ContactName, City, Country
			FROM Customers
			UNION
			SELECT 'Supplier', ContactName, City, Country
			FROM Suppliers;`,
			out: []string{"Customers", "Suppliers"},
		},
		// {
		// 	in: `SELECT Customers.CustomerName, Orders.OrderID
		// 	INTO CustomersOrderBackup2017
		// 	FROM Customers
		// 	LEFT JOIN Orders ON Customers.CustomerID = Orders.CustomerID;`,
		// 	out: []string{"CustomersOrderBackup2017", "Customers", "Orders"},
		// },
	}

	for _, tc := range testcases {
		out, err := sql.ExtractTableNames(tc.in)
		if err != nil {
			t.Error(err)
			continue
		}
		if !reflect.DeepEqual(out, tc.out) {
			t.Errorf("ExtractTableNames('%s'): %s, want %s", tc.in, out, tc.out)
		}
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
