package agent

import (
	"context"

	"github.com/openai/openai-go/v3"
)

// Интерфейс позволяющий использовал реальную ЛЛМ или Мок
type LLMProvider interface {
	AskLLM(ctx context.Context,
		message []openai.ChatCompletionMessageParamUnion,
		tools []openai.ChatCompletionToolUnionParam,
	) (openai.ChatCompletionMessage, error)
}