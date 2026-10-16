package web

import (
	"context"
	"yuima-agent/internal/search"

	"github.com/openai/openai-go/v3"
)

type SearchProvider interface {
	Search(context.Context, string) ([]search.Result, error)
}

type SearchTool struct {
	client SearchProvider
	limit int
}

func NewSearchTool(c SearchProvider, l int) *SearchTool {
	return &SearchTool{
		client: c,
		limit: l,
	}
}

func (t *SearchTool) Definition() openai.ChatCompletionToolUnionParam 

func (t *SearchTool) Call(
	ctx context.Context,
	tc openai.ChatCompletionMessageToolCallUnion,
) openai.ChatCompletionMessageParamUnion

func (t *SearchTool) MaxCallsPerRun() int {
	return t.limit
}