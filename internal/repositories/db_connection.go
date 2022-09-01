package repositories

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DBConnection struct {
	psqlConfig string
	Session    *sql.DB
}

func NewDBConnection(host string, port string, dbname string, user string, password string, sslmode bool) *DBConnection {
	var sslconfig string
	if sslmode {
		sslconfig = "enable"
	} else {
		sslconfig = "disable"
	}

	conn := &DBConnection{
		fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
			host, port, dbname, user, password, sslconfig),
		nil,
	}

	return conn
}

func (conn *DBConnection) Open() {
	session, err := sql.Open("postgres", conn.psqlConfig)
	if err != nil {
		log.Panic(err)
	}

	conn.Session = session
}

func (conn *DBConnection) Close() {
	conn.Session.Close()
}
