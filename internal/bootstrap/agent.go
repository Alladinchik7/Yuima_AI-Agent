package bootstrap

import (
	"log/slog"
	"time"
	"yuima-agent/internal/agent"
	"yuima-agent/internal/storage"
	"yuima-agent/internal/tools"
)

func InitAgent(
	store storage.Store, 
	llm agent.LLMProvider, 
	registry *tools.Registry,
	logger *slog.Logger,
	timeOut time.Duration,
) *agent.Agent {
	return agent.New(store, llm, registry, logger, timeOut)
}