package psql

import (
  "github.com/3d0c/martini-contrib/config"
  _ "github.com/lib/pq"
  "database/sql"
  "log"
  "fmt"
)

type Connection struct {
  Host    string      `json:"host"`
  Port    string      `json:"port"`
  DbName  string      `json:"db_name"`
  Default bool        `json:"default"`
  SslMode string      `json:"ssl_mode"`
  session interface{} `json:"-"`
}

type Connections map[string]Connection

type Config struct {
  Connections `json:"connections"`
}

var pool *Config

func init() {
  log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

func Get(args ...string) Connection {
  var name string

  if pool == nil {
    pool = &Config{}

    config.LoadInto(pool)

    if len(pool.Connections) == 0 {
      log.Println("Warning! 'connections' is not configured.") 
    }
  }
  
  if len(args) > 0 {
    name = args[0]
  } else {
    name = getDefaultConnection(pool.Connections)
  }

  con, ok := pool.Connections[name]
  if !ok {
    log.Println("Connection with '%s' not found.", name)
    return Connection{}
  }

  if con.session == nil {
    var err error
    
    connect := fmt.Sprintf("%s:%s/%s" + "?sslmode=%s", con.Host, con.Port, con.DbName, con.SslMode)
    con.session, err = sql.Open("postgres", connect)
    if err != nil {
      log.Fatal(err)
    }
  }
  return con
}

func (this Connection) DB() interface{} {
  return this.session
}

func (this Connection) Session() *sql.DB {
  return this.session.(*sql.DB)
}

func getDefaultConnection(c Connections) string {
  for name, con := range c {
    if con.Default {
      return name
    }
  }
  panic("No default connection found.")
}
