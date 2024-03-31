package infrastructure

type GmailService struct {
}

func NewGmailService() *GmailService {
	return &GmailService{}
}

func (g *GmailService) CreateDraft() error {
	return nil
}

func (g *GmailService) Crawling() error {
	return nil
}
