package domain

type Mail struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	FromAddress string `json:"from_address"`
}

func (m *Mail) Validation() error {
	if m.Content == "" {
		return EmailContentIsEmpty
	}
	return nil
}
