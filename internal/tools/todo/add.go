package todo

import (
	"context"
	"encoding/json"

	"yuima-agent/internal/storage"
	"yuima-agent/internal/tools/helpers"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

type AddTodoTool struct {
	store storage.TodoStore
	limit int
}

func NewAddTool(store storage.TodoStore, limit int) *AddTodoTool {
	return &AddTodoTool{
		store: store,
		limit: limit,
	}
}

func (t *AddTodoTool) Definition() openai.ChatCompletionToolUnionParam {
	return helpers.FnDefinition(
		"add_todo",
		`When the user asks to add multiple tasks, call add_todo once per task.
		Example: user says "add milk, eggs, bread" → you call add_todo 3 times.
		Add ONE task to the user's todo list. If the user wants multiple 
		tasks, call this function ONCE PER TASK — do not combine them into 
		a single string.`,
		shared.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"task": map[string]any{
					"type":        "string",
					"description": "Название задачи",
				},
			},
			"required": []string{"task"},
		},
	)
}

func (t *AddTodoTool) Call(
	ctx context.Context,
	tc openai.ChatCompletionMessageToolCallUnion,
) openai.ChatCompletionMessageParamUnion {
	var args checkAddTodoArgs
	if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
		b, _ := json.Marshal(map[string]string{"error": "bad args: " + err.Error()})
		return openai.ToolMessage(string(b), tc.ID)
	}

	if args.Task == "" {
		return openai.ToolMessage(`{"error":"task is required"}`, tc.ID)
	}

	err := t.store.AddTodo(args.Task)
	if err != nil {
		b, _ := json.Marshal(map[string]string{"error": err.Error()})
		return openai.ToolMessage(string(b), tc.ID)
	}

	b, _ := json.Marshal(map[string]string{"added": args.Task})
	return openai.ToolMessage(string(b), tc.ID)
}

func (t *AddTodoTool) MaxCallsPerRun() int {
	return t.limit
}