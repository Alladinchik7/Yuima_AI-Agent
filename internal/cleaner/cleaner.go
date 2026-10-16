package cleaner

import (
	"context"
	"log/slog"
	"time"

	"yuima-agent/internal/storage/models"
)

type CleanerStore interface {
	ListTodos() ([]models.Todo, error)
	DeleteTodo(id int) error
}

type Cleaner struct {
	store CleanerStore
	ttl   time.Duration
	interval time.Duration
	logger   *slog.Logger
}

func New(
	store CleanerStore,
	logger *slog.Logger,
) *Cleaner {
	return &Cleaner{
		store:    store,
		ttl:      24 * time.Hour,
		interval: 5 * time.Minute,
		logger:   logger,
	}
}

func (c *Cleaner) Start(ctx context.Context) {
	c.cleanup(ctx)
	
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.cleanup(ctx)
		}
	}
}

func (c *Cleaner) cleanup(ctx context.Context) {
	todos, err := c.store.ListTodos()
	if err != nil {
		c.logger.Error("failed to list todos", "error", err)
		return
	}
	
	var deleteID []int
	for _, t := range todos {
		if t.Completed && !t.DoneAt.IsZero() && time.Since(t.DoneAt) > c.ttl {
			deleteID = append(deleteID, t.ID)
		}
	}

	for _, id := range deleteID {
		if err := c.store.DeleteTodo(id); err != nil {
			c.logger.Error("failed delete todo", "error", err)
		}
	}

	if len(deleteID) > 0 {
		c.logger.Info("deleted task", "count", len(deleteID))
	}
}
