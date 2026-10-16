package memory

import (
	"fmt"
	"time"
	
	"yuima-agent/internal/storage/models"
)

func (s *Store) AddReminder(text string, fireAt time.Time) error {
	if text == "" {
		return fmt.Errorf("empty text")
	}

	if fireAt.Before(time.Now()) {
		return fmt.Errorf("fire_at in the past")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.reminders = append(s.reminders, models.Reminder{
		ID: s.nextReminderID,
		Text: text,
		FireAt: fireAt,
		Fired: false,
	})
	s.nextReminderID++
	
	return nil
}

func (s *Store) MarkReminderFired(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	for i := range s.reminders {
		if s.reminders[i].ID == id {
			s.reminders[i].Fired = true
			return nil
		}
	}

	return fmt.Errorf("reminder %d not found", id)
}
func (s *Store) ListReminders() ([]models.Reminder, error) {
	s.mu.Lock()
	copyRemindes := make([]models.Reminder, len(s.reminders))
	copy(copyRemindes, s.reminders)
	s.mu.Unlock()

	return copyRemindes, nil
}

func (s *Store) EarliestUnfired() (*models.Reminder, error) {
	reminders, err := s.ListReminders()
	if err != nil {
		return nil, fmt.Errorf("list reminders: %w", err)
	}

	var earliest *models.Reminder
	for i := range reminders {
		r := &reminders[i]

		if reminders[i].Fired {
			continue
		}

		if earliest == nil || r.FireAt.Before(earliest.FireAt) {
			earliest = r
		}
	}

	if earliest == nil {
		return nil, nil
	}

	return earliest, nil
}