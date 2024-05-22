package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"plum/domain"
)

type AISearch struct {
	baseURL string
	apiKey  string
}

func NewAISearch(baseURL, apiKey string) *AISearch {
	return &AISearch{
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

func (a *AISearch) UploadDocuments(doc domain.MailDoc) error {
	url := fmt.Sprintf("%s/indexes/%s/docs/index?api-version=2023-11-01", a.baseURL, "email_index")
	requestBody, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("error marshaling request: %v", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", a.apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error: status code %d, %s", resp.StatusCode, body)
	}
	return nil
}

func (a *AISearch) SearchDocuments(search string) (string, error) {
	url := fmt.Sprintf("%s/indexes/%s/docs?api-version=2023-11-01&search=%s", a.baseURL, "email-index", search)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", a.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: status code %d, %s", resp.StatusCode, body)
	}

	fmt.Println(string(body))
	return string(body), nil
}
