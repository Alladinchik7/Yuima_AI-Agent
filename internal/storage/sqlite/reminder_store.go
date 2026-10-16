package sqlite

import (
	"fmt"
	"time"

	"yuima-agent/internal/storage/models"
)

func (s *Store) AddReminder(text string, fireAt time.Time) error {
	if text == "" {
		return fmt.Errorf("emty text")
	}

	if fireAt.Before(time.Now()) {
		return fmt.Errorf("fire_at in the past")
	}
	
	res := s.db.DB.Create(&models.Reminder{
		Text: text,
		FireAt: fireAt,
		Fired: false,
	})

	if res.Error != nil {
		return fmt.Errorf("failed add reminder from database: %v", res.Error)
	}

	return nil
}

func (s *Store) MarkReminderFired(id int) error {
	res := s.db.DB.Model(&models.Reminder{}).Where("id = ?", id).Update("Fired", true)
	if res.Error != nil {
		return fmt.Errorf("failed update mark remined with id %d fired: %v", id, res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("fired did not update")
	}

	return nil
}

func (s *Store) ListReminders() ([]models.Reminder, error) {
	var reminders []models.Reminder
	res := s.db.DB.Find(&reminders)
	if res.Error != nil {
		return nil , res.Error
	}

	return reminders, nil
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