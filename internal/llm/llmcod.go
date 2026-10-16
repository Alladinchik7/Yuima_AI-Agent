package llm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type Client struct {
	model string
	sdk openai.Client
}

func NewClient(aiProvider, apiKey, baseURL string) *Client {
	return &Client{
		model: aiProvider,
		sdk: openai.NewClient(
			option.WithAPIKey(apiKey),
			option.WithBaseURL(baseURL),
		),
	}
}

func (l *Client) AskLLM(
	ctx context.Context,
	messages []openai.ChatCompletionMessageParamUnion,
	tools []openai.ChatCompletionToolUnionParam,
) (openai.ChatCompletionMessage, error) {
	if len(messages) == 0 {
		return openai.ChatCompletionMessage{}, fmt.Errorf("no messages")
	}

	completion, err := l.sdk.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:   l.model,
		Messages: messages,
		Tools: tools,
		ParallelToolCalls: openai.Bool(true),
	})
	if err != nil {
		return openai.ChatCompletionMessage{}, err
	}
	
	if len(completion.Choices) == 0 {
		raw, _ := json.MarshalIndent(completion, "", "  ")
		return openai.ChatCompletionMessage{}, fmt.Errorf("empty choices: %s", raw)
	}
	
	return completion.Choices[0].Message, nil
}