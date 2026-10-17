package bootstrap

import (
	"log/slog"
	"os"
	"path/filepath"
	
	"yuima-agent/internal/config"
	"yuima-agent/internal/storage"
	"yuima-agent/internal/storage/memory"
	"yuima-agent/internal/storage/sqlite"
)

func InitStore(cfg config.Config, logger *slog.Logger) storage.Store {
	if cfg.UseMemory {
		return  memory.NewStore()
	}else {
		path := filepath.Join("..", "..", "data.db")
		logger.Info("db.path", "path", path)
		
		db, err := sqlite.New(path, logger)
		if err != nil {
			logger.Error("db.error", "error", "failed create file db")
			os.Exit(1)
		}
		defer db.Close()
		
		return  sqlite.NewStore(db)
	}
}