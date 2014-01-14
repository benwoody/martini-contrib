package psql

import (
  "github.com/3d0c/martini-contrib/config"
  "testing"
)

func TestPsql(t *testing.T) {
  config.Init("./psql_test.json")

  con := Get()

  if !con.Default {
    t.Fatalf("Expected con.Default = true, got %v\n", con.Default)
  }
}
