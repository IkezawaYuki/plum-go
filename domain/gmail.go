package domain

type Gmail struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
}

type GmailList struct {
	GmailList []Gmail `json:"gmail_list"`
}

func (m *Gmail) Validation() error {
	if m.Content == "" {
		return EmailContentIsEmpty
	}
	return nil
}

func (m *GmailList) Validation() error {
	for _, mail := range m.GmailList {
		if err := mail.Validation(); err != nil {
			return err
		}
	}
	return nil
}
