package infrastructure

import (
	"context"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
	"plum/logger"
)

type Db struct {
	conn *sql.DB
}

func NewDb(connection *sql.DB) *Db {
	return &Db{
		conn: connection,
	}
}

func Connect() *sql.DB {
	cfg := mysql.Config{
		User:   os.Getenv("DATABASE_USER"),
		Passwd: os.Getenv("DATABASE_PASSWORD"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "",
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	logger.Logger.Info("DB connected")
	return db
}

func (c *Db) Close() error {
	return c.conn.Close()
}

func (c *Db) UnitOfWork(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := c.conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	var done bool
	defer func() {
		if !done {
			_ = tx.Rollback()
		}
	}()

	if err := fn(ctx); err != nil {
		return err
	}

	done = true
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
