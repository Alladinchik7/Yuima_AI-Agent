package models

import "time"

type Todo struct {
	ID          int        `json:"id"`
	Task        string     `json:"task"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	DoneAt      time.Time  `json:"done_at"`
}