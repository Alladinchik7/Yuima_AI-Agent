package todo

import (
	"context"
	"encoding/json"

	"yuima-agent/internal/storage"
	"yuima-agent/internal/tools/helpers"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

type CompleteTodoTool struct {
	store storage.TodoStore
	limit int
}

func NewCompleteTool(store storage.TodoStore, limit int) *CompleteTodoTool {
	return &CompleteTodoTool{
		store: store, 
		limit: limit,
	}
}

func (c *CompleteTodoTool) Definition() openai.ChatCompletionToolUnionParam {
	return helpers.FnDefinition(
		"complete_todo",
		`Mark a task as completed.
	 	Use when the user says they finished a task, or asks to mark it as done. 
		Requires the task ID — if you don't have it, call list_todos first to find the correct ID.
		The list_todos result shows tasks with their IDs. When the user says 
		"task 1" or "first task", clarify whether they mean the ID or the 
		position in the list. Prefer to confirm with the user before calling 
		complete_todo or delete_todo.`,
		shared.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"task_id": map[string]any {
					"type": "integer",
					"description": "The ID of the task to mark as completed.",
				},
			},
			"required": []string{"task_id"},
		},
	)
}

func (c *CompleteTodoTool) Call(
	ctx context.Context,
	tc openai.ChatCompletionMessageToolCallUnion,
) openai.ChatCompletionMessageParamUnion {
	var args checkTaskID
	if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
		b, _ := json.Marshal(map[string]string{"error": "bad args: " + err.Error()})
		return openai.ToolMessage(string(b), tc.ID)
	}

	err := c.store.CompleteTodo(args.TaskID)
	if err != nil {
		b, _ := json.Marshal(map[string]string{"error": "failed to mark task complete: " + err.Error()})
		return openai.ToolMessage(string(b), tc.ID)
	}

	b, _ := json.Marshal(map[string]any{"completed": true, "id": args.TaskID})
	return openai.ToolMessage(string(b), tc.ID)
}

func (c *CompleteTodoTool) MaxCallsPerRun() int {
	return c.limit
}