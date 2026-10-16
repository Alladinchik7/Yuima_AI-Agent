package webtools

import (
	"context"
	"encoding/json"
	"yuima-agent/internal/search"
	"yuima-agent/internal/tools/helpers"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

type SearchProvider interface {
	Search(context.Context, string) ([]search.Result, error)
}

type SearchTool struct {
	client SearchProvider
	limit int
}

type checkSearchArgs struct {
	Query string `json:"query"`
}

func NewSearchTool(c SearchProvider, l int) *SearchTool {
	return &SearchTool{
		client: c,
		limit: l,
	}
}

func (t *SearchTool) Definition() openai.ChatCompletionToolUnionParam {
	return helpers.FnDefinition(
		"web_search",
		`Search for information on the internet based on the user's query. 
		Only perform a search if you do not know the answer yourself.`,
		shared.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"query": map[string]any{
					"type": "string",
					"description": "User query for which an answer is required",
				},
			},
			"required": []string{"query"},
		},
	)
}

func (t *SearchTool) Call(
	ctx context.Context,
	tc openai.ChatCompletionMessageToolCallUnion,
) openai.ChatCompletionMessageParamUnion {
	var args checkSearchArgs
	if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
		b, _ := json.Marshal(map[string]string{"error": "failed unmarshal arguments for search"})
		return openai.ToolMessage(string(b), tc.ID)
	}

	if args.Query == "" {
		b, _ := json.Marshal(map[string]string{"error":"empty query"})
		return openai.ToolMessage(string(b), tc.ID)
	}

	res, err := t.client.Search(ctx, args.Query)
	if err != nil {
		b, _ := json.Marshal(map[string]string{"error":err.Error()})
		return openai.ToolMessage(string(b), tc.ID)
	}

	b, _ := json.Marshal(res)
	return openai.ToolMessage(string(b), tc.ID)
}

func (t *SearchTool) MaxCallsPerRun() int {
	return t.limit
}