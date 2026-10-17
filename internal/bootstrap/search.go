package bootstrap

import (
	"log/slog"
	"net/http"
	"time"
	
	"yuima-agent/internal/config"
	"yuima-agent/internal/search"
	"yuima-agent/internal/tools/webtools"
)

func InitSearchClient(cfg config.Config, logger *slog.Logger) webtools.SearchProvider {
	if cfg.UseMockSearch {
		 return search.NewMock()
	}else {
		httpClient := http.Client{Timeout: 15*time.Second}
		return search.NewClient(&httpClient, cfg.SearxNGURL, cfg.FormatSearxNG, cfg.LangSerxNG, cfg.Safesearch)
	}
}