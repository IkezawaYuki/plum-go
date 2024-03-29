package infrastructure

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
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
	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	fmt.Println(resp.Choices[0].Message.Content)
	return resp.Choices[0].Message.Content, nil
}
