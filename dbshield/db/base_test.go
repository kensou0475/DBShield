package db_test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

// Setup
func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard) // Avoid log outputs
	os.Exit(m.Run())
}
