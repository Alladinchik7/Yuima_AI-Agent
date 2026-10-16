package todo

import (
	"context"
	"encoding/json"

	"yuima-agent/internal/storage"
	"yuima-agent/internal/tools/helpers"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

type UncompleteTodoTool struct {
	store storage.TodoStore
	limit int
}

func NewUncompleteTool(store storage.TodoStore, limit int) *UncompleteTodoTool {
	return &UncompleteTodoTool{
		store: store,
		limit: limit,
	}
}

func (u *UncompleteTodoTool) Definition() openai.ChatCompletionToolUnionParam {
	return helpers.FnDefinition(
		"uncomplete_todo",
		`Revert a completed task to an active state. 
		Use only when the user explicitly asks to unmark a task that was previously completed, 
		e.g. "I didn't actually finish it" or "undo that". 
		Requires the task ID — call list_todos first if needed.
		The list_todos result shows tasks with their IDs. When the user says 
		"task 1" or "first task", clarify whether they mean the ID or the 
		position in the list. Prefer to confirm with the user before calling 
		complete_todo or delete_todo.`,
		shared.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"task_id": map[string]any{
					"type": "integer",
					"description": `The numeric ID of the task to unmark, taken from the list_todos result. 
					If the user refers to a task by its text, call list_todos first to resolve the text to the correct ID.`,
				},
				"confirmed": map[string]any{
					"type": "boolean",
					"description": `Set to true only if the user has explicitly confirmed they want to uncomplete the status task 
					(e.g. after you asked "are you sure?" and they said yes). Leave false by default. 
					Never set to true based solely on the user's initial request — always ask for confirmation first.`,
				},
			},
			"required": []string{"task_id"},
		},
	)
}

func (u *UncompleteTodoTool) Call(
	ctx context.Context,
	tc openai.ChatCompletionMessageToolCallUnion,
) openai.ChatCompletionMessageParamUnion {
	var args checkTaskID
	if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
		b, _ := json.Marshal(map[string]string{"error": "bad args: " + err.Error()})
		return openai.ToolMessage(string(b), tc.ID)
	}

	err := u.store.UncompleteTodo(args.TaskID)
	if err != nil {
		b, _ := json.Marshal(map[string]string{"error": "failed to mark task uncomplete: " + err.Error()})
		return openai.ToolMessage(string(b), tc.ID)
	}

	b, _ := json.Marshal(map[string]any{"completed": false, "id": args.TaskID})
	return openai.ToolMessage(string(b), tc.ID)
}

func (u *UncompleteTodoTool) MaxCallsPerRun() int {
	return u.limit
}