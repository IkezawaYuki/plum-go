package infrastructure

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
	"plum/domain"
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

func DbConnect() *sql.DB {
	cfg := mysql.Config{
		User:   os.Getenv("DATABASE_USER"),
		Passwd: os.Getenv("DATABASE_PASSWORD"),
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:3306", os.Getenv("DATABASE_HOST")),
		DBName: os.Getenv("DATABASE_SCHEME"),
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

const selectChatgptSettingStmt = "select prompt, system_message from chatgpt_setting where id = 1"

func (c *Db) GetChatgptSetting() (domain.ChatgptSetting, error) {
	result := domain.ChatgptSetting{}
	row := c.conn.QueryRow(selectChatgptSettingStmt)
	if err := row.Scan(&result.Prompt, &result.SystemMessage); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, nil
		}
		return result, err
	}
	return result, nil
}

const updatePromptStmt = "update chatgpt_setting set prompt = ?, system_message = ? where id = 1"

func (c *Db) UpdateChatgptSetting(setting domain.ChatgptSetting) error {
	result, err := c.conn.Exec(updatePromptStmt, setting.Prompt, setting.SystemMessage)
	if err != nil {
		return err
	}
	fmt.Println(result.LastInsertId())
	return nil
}
