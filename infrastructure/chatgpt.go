package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"io"
	"net/http"
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

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AuthenticationParams struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

type DataSourceParams struct {
	Endpoint              string               `json:"endpoint"`
	IndexName             string               `json:"index_name"`
	SemanticConfiguration string               `json:"semantic_configuration"`
	QueryType             string               `json:"query_type"`
	FieldsMapping         map[string]string    `json:"fields_mapping"`
	InScope               bool                 `json:"in_scope"`
	RoleInformation       string               `json:"role_information"`
	Filter                interface{}          `json:"filter"`
	Strictness            int                  `json:"strictness"`
	TopNDocuments         int                  `json:"top_n_documents"`
	Authentication        AuthenticationParams `json:"authentication"`
	Key                   string               `json:"key"`
	IndexNameAlias        string               `json:"indexName"`
}

type DataSource struct {
	Type       string           `json:"type"`
	Parameters DataSourceParams `json:"parameters"`
}

type CompletionRequest struct {
	Messages     []Message    `json:"messages"`
	DeploymentID string       `json:"deployment_id"`
	DataSources  []DataSource `json:"data_sources"`
	Enhancements interface{}  `json:"enhancements"`
	Temperature  float64      `json:"temperature"`
	TopP         float64      `json:"top_p"`
	MaxTokens    int          `json:"max_tokens"`
	Stop         interface{}  `json:"stop"`
	Stream       bool         `json:"stream"`
}

func (c *ChatGPT) Generate(string) (*domain.Generated, error) {
	apiBase := "https://plum-strategy-drive.openai.azure.com"
	apiVersion := "2024-02-15-preview"
	deploymentID := "plum-ai"

	client := &http.Client{}

	searchEndpoint := "https://plum-search.search.windows.net"
	searchKey := os.Getenv("AI_SEARCH_API_KEY")
	fmt.Println(searchKey)

	messageText := []Message{
		{
			Role:    "user",
			Content: "What are the differences between Azure Machine Learning and Azure AI services?",
		},
	}

	dataSources := []DataSource{
		{
			Type: "azure_search",
			Parameters: DataSourceParams{
				Endpoint:              searchEndpoint,
				IndexName:             "email-index",
				SemanticConfiguration: "default",
				QueryType:             "simple",
				FieldsMapping:         make(map[string]string),
				InScope:               true,
				RoleInformation:       "You are an AI assistant that helps people find information.",
				Filter:                nil,
				Strictness:            3,
				TopNDocuments:         5,
				Authentication: AuthenticationParams{
					Type: "api_key",
					Key:  searchKey,
				},
				Key:            searchKey,
				IndexNameAlias: "email-index",
			},
		},
	}

	completionRequest := CompletionRequest{
		Messages:     messageText,
		DeploymentID: deploymentID,
		DataSources:  dataSources,
		Enhancements: nil,
		Temperature:  0,
		TopP:         1,
		MaxTokens:    800,
		Stop:         nil,
		Stream:       true,
	}

	requestBody, err := json.Marshal(completionRequest)
	if err != nil {
		fmt.Println("Error marshaling request:", err)
		return nil, err
	}

	url := fmt.Sprintf("%s/openai/deployments/%s/extensions/chat/completions?api-version=%s", apiBase, deploymentID, apiVersion)
	fmt.Println(url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AZURE_OPENAI_KEY")))
	fmt.Println(os.Getenv("AZURE_OPENAI_KEY"))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	}

	fmt.Println(string(body))
	return nil, nil
}
