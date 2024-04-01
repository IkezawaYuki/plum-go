package infrastructure

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

type Slack struct {
}

func NewSlack() *Slack {
	return &Slack{}
}

func (s *Slack) SendMessage(webhookURL, msg string) error {
	jsonStr := []byte(`{"text":"` + msg + `"}`)

	// HTTP POSTリクエストを作成
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// HTTPクライアントを作成してリクエストを送信
	client := &http.Client{
		Timeout: time.Minute,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// ステータスコードをチェック
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Message sent successfully")
	} else {
		fmt.Println("Failed to send message")
	}
	return nil
}
