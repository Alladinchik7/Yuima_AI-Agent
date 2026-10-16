package memory

import (
	"fmt"
	"time"
	
	"yuima-agent/internal/storage/models"
)

func (s *Store) AddTodo(task string) error {
	if task == "" {
		return fmt.Errorf("empty task")
	}
	
	s.mu.Lock()
	defer s.mu.Unlock()
	s.todos = append(s.todos, models.Todo{
		ID:          s.nextID,
		Task:        task,
		Completed:   false,
		CreatedAt:   time.Now(),
		DoneAt:      time.Time{},
	})
	s.nextID++

	return nil
}

func (s *Store) ListTodos() ([]models.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	copyTodos := make([]models.Todo, len(s.todos))
	copy(copyTodos, s.todos)
	return copyTodos, nil
}

func (s *Store) CompleteTodo(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := 0; i < len(s.todos); i++ {
		if s.todos[i].ID == id {
			s.todos[i].Completed = true
			s.todos[i].DoneAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("todo with id %d not found", id)
}

func (s *Store) UncompleteTodo(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := 0; i < len(s.todos); i++ {
		if s.todos[i].ID == id {
			s.todos[i].Completed = false
			s.todos[i].DoneAt = time.Time{}
			return nil
		}
	}
	return fmt.Errorf("todo %d not found", id)
}	

func (s *Store) DeleteTodo(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := 0; i < len(s.todos); i++ {
		if s.todos[i].ID == id {
			s.todos = append(s.todos[:i], s.todos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("todo %d not found", id)
}