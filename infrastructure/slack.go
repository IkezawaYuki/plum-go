package infrastructure

import (
	"bytes"
	"fmt"
	"net/http"
	"plum/domain"
	"plum/logger"
	"time"
)

type Slack struct {
	webhookURL string
}

func NewSlack(webhookURL string) *Slack {
	return &Slack{
		webhookURL: webhookURL,
	}
}

func getMsg(contact domain.Contact) string {
	return fmt.Sprintf("サポートサイトより問い合わせがありました。 対応をお願いします。\n")
}

func (s *Slack) Report(contact domain.Contact) error {
	return nil
}

func (s *Slack) Escalation(contact domain.Contact) error {
	jsonStr := []byte(`{"text":"` + getMsg(contact) + `"}`)

	req, err := http.NewRequest("POST", s.webhookURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

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
	
	if resp.StatusCode == http.StatusOK {
		logger.Logger.Info("message sent successfully")
	} else {
		logger.Logger.Info("failed to send message")
	}
	return nil
}
