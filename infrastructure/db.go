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

func Connect() *sql.DB {
	cfg := mysql.Config{
		User:   os.Getenv("DATABASE_USER"),
		Passwd: os.Getenv("DATABASE_PASSWORD"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "plum",
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

const selectChatgptSettingStmt = "select prompt, system_message chatgpt_setting where id = 1"

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

const updatePromptStmt = "update chatgpt_setting set prompt = ? where id = 1"

func (c *Db) UpdatePrompt(prompt string) error {
	result, err := c.conn.Exec(updatePromptStmt, prompt)
	if err != nil {
		return err
	}
	fmt.Println(result.LastInsertId())
	return nil
}

const updateSystemMsgStmt = "update chatgpt_setting set system_message = ? where id = 1"

func (c *Db) UpdateSystemMessage(systemMsg string) error {
	return nil
}
