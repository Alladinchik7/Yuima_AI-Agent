package models

import "time"

type Reminder struct {
	ID int			 `json:"id" gorm:"primaryKey"`
	Text string      `json:"text" gorm:"column:text"`
	FireAt time.Time `json:"fire_at"`
	Fired bool       `json:"fired"`
}