package infrastructure

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"plum/logger"
)

type ChatGPT struct {
	client *openai.Client
}

func NewChatGPT(token string) *ChatGPT {
	return &ChatGPT{
		client: openai.NewClient(token),
	}
}

func (c *ChatGPT) Create(content string) (string, error) {
	massage := fmt.Sprintf(`顧客から次のような問い合わせがありました。
「%s」
これに返信するメール本文を作成してください。
`, content)
	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: massage,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	logger.Logger.Info("message is generated", resp.Choices[0].Message.Content)
	return resp.Choices[0].Message.Content, nil
}
