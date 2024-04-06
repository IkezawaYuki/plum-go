package domain

type Gmail struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	FromAddress string `json:"from_address"`
}

func (m *Gmail) Validation() error {
	if m.Content == "" {
		return EmailContentIsEmpty
	}
	return nil
}
