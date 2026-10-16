package memory

import (
	"sync"

	"yuima-agent/internal/storage/models"
)

type Store struct {
	mu    sync.Mutex
	todos []models.Todo 
	reminders []models.Reminder 
	summary []models.DialogMemory
	nextReminderID int
	nextID int
}

func NewStore() *Store {
	return &Store{
		mu:    sync.Mutex{},
		todos: []models.Todo{},
		reminders: []models.Reminder{},
		nextReminderID: 1,
		nextID: 1,
	}
}