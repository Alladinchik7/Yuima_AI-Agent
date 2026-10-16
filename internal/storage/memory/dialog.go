package memory

import "yuima-agent/internal/storage/models"

func (s *Store) SaveSummaryWithLimit(text string, limit int) error {
	s.summary = append(s.summary, models.DialogMemory{ID: 1, DialogSummary: text})
	return nil
}

func (s *Store) LoadSummary() ([]models.DialogMemory, error) {
	return s.summary, nil
}