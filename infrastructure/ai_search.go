package infrastructure

import "plum/domain"

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

func (a *AISearch) CreateIndex() error {
	return nil
}

func (a *AISearch) UploadDocuments(doc domain.MailDoc) error {
	return nil
}

func (a *AISearch) SearchDocuments(search string) error {
	return nil
}
