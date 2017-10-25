package db_test

import (
	"testing"

	"github.com/qiwihui/DBShield/dbshield/db"
)

func TestContains(t *testing.T) {
	testcases := []struct {
		in1 []string
		in2 string
		out bool
	}{
		{
			in1: []string{"SELECT"},
			in2: "SELECT",
			out: true,
		},
		{
			in1: []string{"SELECT"},
			in2: "CREATE",
			out: false,
		},
	}

	for _, tc := range testcases {
		out, err := db.Contains(tc.in1, tc.in2)
		if err != nil {
			t.Error(err)
			continue
		}
		if out != tc.out {
			t.Errorf("%s not Contains('%s')", tc.in1, tc.in2)
		}
	}
}
