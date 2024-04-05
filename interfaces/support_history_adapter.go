package interfaces

import "plum/domain"

type SupportHistoryAdapter struct {
}

func (s *SupportHistoryAdapter) FindByID(int) (domain.SupportHistory, error) {
	return domain.SupportHistory{}, nil
}

func (s *SupportHistoryAdapter) FindAll([]domain.SupportHistory, error) {
	return
}

func (s *SupportHistoryAdapter) FindByKeyword(string) (domain.SupportHistory, error) {
	return domain.SupportHistory{}, nil
}

func (s *SupportHistoryAdapter) Save(support *domain.SupportHistory) error {
	return nil
}
