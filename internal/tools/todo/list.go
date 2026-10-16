package todo

import (
	"context"
	"encoding/json"
	"yuima-agent/internal/storage"
	"yuima-agent/internal/tools/helpers"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

type ListTodoTool struct {
	store storage.TodoStore
	limit int
}

type todoItem struct {
    Position  int    `json:"position"`
    ID        int    `json:"id"`
    Task      string `json:"task"`
    Completed bool   `json:"completed"`
}

func NewListTool(store storage.TodoStore, limit int) *ListTodoTool {
	return &ListTodoTool{
		store: store,
		limit: limit,
	}
}

func (t *ListTodoTool) Definition() openai.ChatCompletionToolUnionParam {
	return helpers.FnDefinition(
		"list_todos",
		"Show the user's current todo list. Use when the user asks what's on their list or wants to see tasks.",
		shared.FunctionParameters{},
	)
}

func (t *ListTodoTool) Call(
	ctx context.Context,
	tc openai.ChatCompletionMessageToolCallUnion,
) openai.ChatCompletionMessageParamUnion {
	todos, err := t.store.ListTodos()
	if err != nil {
		b, _ := json.Marshal(map[string]string{"error": "list todos failed: " + err.Error()})
		return openai.ToolMessage(string(b), tc.ID)
	}

	if len(todos) == 0 {
		return openai.ToolMessage(`{"todos": "[]"}`, tc.ID)
	}

	items := make([]todoItem, len(todos))
	for i, t := range todos {
		items[i] = todoItem{
			Position: i+1,
			ID: t.ID,
			Task: t.Task,
			Completed: t.Completed,
		}
	}

	b, err := json.Marshal(map[string]any{"todos":items})
	if err != nil {
    	return openai.ToolMessage(`{"error":"marshal failed"}`, tc.ID)
	}
	return openai.ToolMessage(string(b), tc.ID)
}

func (t *ListTodoTool) MaxCallsPerRun() int {
	return t.limit
}