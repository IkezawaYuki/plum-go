package domain

type SupportHistoryRepository interface {
	FindByID(int) (SupportHistory, error)
	FindAll([]SupportHistory, error)
	FindByKeyword(string) (SupportHistory, error)
	Save(*SupportHistory) error
}
