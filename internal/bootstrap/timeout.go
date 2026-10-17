package bootstrap

import (
	"time"
	"yuima-agent/internal/config"
)

const (
	localTimeout = 5*time.Minute // Таймаут для локальной ЛЛМ
	cloudTimeout = 180*time.Second// Таймаут для облачной ЛЛМ
)

func SetTimeOut(cfg config.Config) time.Duration {
	if cfg.UseCloud {
		return cloudTimeout
	}else {
		return localTimeout
	}
}