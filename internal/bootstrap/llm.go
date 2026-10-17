package bootstrap

import (
	"yuima-agent/internal/agent"
	"yuima-agent/internal/config"
	"yuima-agent/internal/llm"
)

func InitLLM(cfg config.Config) agent.LLMProvider {
	if cfg.UseMockAgent {
		return llm.NewMock()
	}else {
		return llm.NewClient(cfg.AIModel, cfg.AIAPIKey, cfg.AIBaseURL)
	}
}