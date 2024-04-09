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
	massage := fmt.Sprintf(`A customer asked us the following question: Please compose the body of the email replying to this in Japanese.
「%s」`, content)
	systemMessage := `You are in charge of supporting customer inquiries. Compose your response via email. If you don't know the answer, don't force yourself to answer and escalate the situation.`
	dialog := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemMessage,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: massage,
		},
	}
	f := openai.FunctionDefinition{
		Name:        "escalation",
		Description: "Call this function when you don't know how to answer a question from a customer.",
		Parameters:  nil,
	}
	tool := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &f,
	}
	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       os.Getenv("AZURE_OPENAI_MODEL"),
			Messages:    dialog,
			Temperature: 1.,
			MaxTokens:   400.,
			Tools:       []openai.Tool{tool},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("chatgpt is error: %v", err)
	}
	return &domain.Generated{
		Message:    resp.Choices[0].Message.Content,
		Escalation: resp.Choices[0].Message.FunctionCall != nil,
	}, nil
}
