package sqlite

import (
	"fmt"
	"log/slog"
	"yuima-agent/internal/storage/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Conn struct {
	DB *gorm.DB
	logger *slog.Logger
}

func New(path string, log *slog.Logger) (*Conn, error) {
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	}
	db, err := gorm.Open(sqlite.Open(path), config)
	if err != nil {
		log.Error("DB.error", "result_connect", "unseccessfuly", "error", err)
		return nil, fmt.Errorf("failed connect from database: %v", err)
	}
	log.Info("DB", "result_connect", "seccessfully")

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

	err = db.AutoMigrate(&models.Reminder{}, &models.Todo{}, &models.DialogMemory{})
	if err != nil {
		log.Error("DB", "result_db_migrate", "unsuccessfully", "error", err)
		return nil, fmt.Errorf("failed automigrate: %v", err)
	}
	log.Info("DB", "result_db_migrate", "successfully")
	return &Conn{
		DB: db,
		logger: log,
	}, nil
}

func (d *Conn) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	
	if err := sqlDB.Close(); err != nil {
		d.logger.Error("db.close", "result", "error:"+err.Error())
		return fmt.Errorf("failed close database: %v", err)
	}
	
	return nil
}