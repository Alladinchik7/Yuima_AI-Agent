package sqlite

import (
	"fmt"
	"time"

	"yuima-agent/internal/storage/models"
)

func (s *Store) AddTodo(task string) error {
	if task == "" {
		return fmt.Errorf("empty task")
	}
	
	res := s.db.DB.Create(&models.Todo{
		Task: task,
		Completed: false,
		CreatedAt: time.Now(),
		DoneAt: time.Time{},
	})
	if res.Error != nil {
		return fmt.Errorf("failed add task from database: %v", res.Error)
	}
	
	return nil
}

func (s *Store) ListTodos() ([]models.Todo, error) {
	var tasks []models.Todo
	res := s.db.DB.Find(&tasks)
	if res.Error != nil {
		return nil, res.Error
	}

	return tasks, nil
}

func (s *Store) CompleteTodo(id int) error {
	res := s.db.DB.Model(&models.Todo{}).Where("id = ?", id).Updates(map[string]any{"Completed": true, "DoneAt":time.Now()})
	if res.Error != nil {
		return fmt.Errorf("failed update task with id %d to complete: %v", id, res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("task status complete did not update")
	}

	return nil
}

func (s *Store) UncompleteTodo(id int) error {
	res := s.db.DB.Model(&models.Todo{}).Where("id = ?", id).Updates(map[string]any{"Completed":false, "DoneAt":time.Time{}})
	if res.Error != nil {
		return fmt.Errorf("failed update task with id %d to uncomplete: %v", id, res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("task status uncomplete did not update")
	}

	return nil
}

func (s *Store) DeleteTodo(id int) error {
	res := s.db.DB.Delete(&models.Todo{}, id)
	if res.Error != nil {
		return fmt.Errorf("failed deleted task with id %d: %v", id, res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("task not deleted")
	}

	return nil
}