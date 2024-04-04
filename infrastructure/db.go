package infrastructure

import "database/sql"

type Db struct {
	conn *sql.Conn
}

func NewDbConnect(connection *sql.Conn) *Db {
	return &Db{
		conn: connection,
	}
}
