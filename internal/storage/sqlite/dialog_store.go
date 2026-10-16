package sqlite

import (
	"fmt"
	"yuima-agent/internal/storage/models"

	"gorm.io/gorm"
)

func (s *Store) SaveSummaryWithLimit(text string, limit int) error {
	if text == "" {
		return fmt.Errorf("empty text")
	}

	if limit <= 0 {
		limit = 5
	}

	err := s.db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&models.DialogMemory{DialogSummary: text}).Error; err != nil {
			return fmt.Errorf("create summary: %w", err)
		}

		var total int64
		if err := tx.Model(&models.DialogMemory{}).Count(&total).Error; err != nil {
			return fmt.Errorf("count summaries: %w", err)
		}

		if total > int64(limit) {
			err := tx.Where("id NOT IN (?)", 
				tx.Model(&models.DialogMemory{}).
					Select("id").
					Order("id DESC").
					Limit(limit),
			).Delete(&models.DialogMemory{}).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("delete old: %w")
	}

	return nil
}

func (s *Store) LoadSummary() ([]models.DialogMemory, error) {
	var summary []models.DialogMemory
	res := s.db.DB.Order("id ASC").Find(&summary)
	if res.Error != nil {
		return nil, res.Error
	}

	return summary, nil
}