package infrastructure

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"os"
	"plum/domain"
)

type ChatGPT struct {
	client *openai.Client
}

func NewChatGPT(apiKey, baseURL string) *ChatGPT {
	config := openai.DefaultAzureConfig(apiKey, baseURL)
	client := openai.NewClientWithConfig(config)
	return &ChatGPT{client: client}
}

func (c *ChatGPT) Create(content string) (*domain.Generated, error) {
	massage := fmt.Sprintf(`顧客から次のような問い合わせがありました。
「%s」
これに返信するメール本文を作成してください。
`, content)
	systemMessage := `you are nice guy`
	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: os.Getenv("AZURE_OPENAI_MODEL"),
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemMessage,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: massage,
				},
			},
			Temperature: 1.,
			MaxTokens:   400.,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("chatgpt is error: %v", err)
	}
	//logger.Logger.Info("message is generated", resp.Choices[0].Message.Content)
	return &domain.Generated{
		Message:    resp.Choices[0].Message.Content,
		Escalation: resp.Choices[0].Message.FunctionCall != nil,
	}, nil
}
