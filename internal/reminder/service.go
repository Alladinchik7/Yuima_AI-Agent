package reminder

import (
	"context"
	"log/slog"
	"time"

	"yuima-agent/internal/storage/models"
)

type ReminderStore interface {
	EarliestUnfired() (*models.Reminder, error)
    MarkReminderFired(id int) error
}

type Service struct {
	store ReminderStore
	notify func(text string)
	wakeUp chan struct{}
	logger *slog.Logger
}

func New(
	store ReminderStore,
	logger *slog.Logger,
	wakeUp chan struct{},
	notify func(text string), 
) *Service {
	return &Service{
		store: store,
		wakeUp: wakeUp,
		logger: logger,
		notify: notify,
	}
}

func (s *Service) Start(ctx context.Context) {
	for {
		next, err := s.store.EarliestUnfired()
		if err != nil {
			s.logger.Error("earliest unfired failed", "err", err)
        	// всё равно ставим таймер на час — попробуем позже
         	next = nil
		}

		duration := time.Hour // По умолчанию - час
		if next != nil {
			duration = time.Until(next.FireAt)
			if duration < 0 {
				duration = 0
			}
		}

		timer := time.NewTimer(duration)
		select{
			case <-timer.C:
				if next != nil {
					s.notify(next.Text)
					if err := s.store.MarkReminderFired(next.ID); err != nil {
						s.logger.Warn("mark fired failed", "id", next.ID, "err", err)
					}
				}
			case <-s.wakeUp:
				timer.Stop()
			case <-ctx.Done():
				timer.Stop()
				return
		}
	}
}