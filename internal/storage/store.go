package storage

import (
	"time"
	
	"yuima-agent/internal/storage/models"
)

type Store interface {
	TodoStore
	ReminderStore
	AgentStore
}

type TodoStore interface {
	AddTodo(task string) error
	ListTodos() ([]models.Todo, error)
	CompleteTodo(id int) error
	UncompleteTodo(id int) error
	DeleteTodo(id int) 	  error
}

type ReminderStore interface {
	AddReminder(text string, fireAt time.Time) error
	MarkReminderFired(id int) error
	ListReminders() ([]models.Reminder, error)
	EarliestUnfired() (*models.Reminder, error)
}

type AgentStore interface {
	SaveSummaryWithLimit(text string, limit int) error
	LoadSummary() ([]models.DialogMemory, error)
}