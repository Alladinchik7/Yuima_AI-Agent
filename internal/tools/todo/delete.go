package todo

import (
	"context"
	"encoding/json"

	"yuima-agent/internal/storage"
	"yuima-agent/internal/tools/helpers"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

type DeleteTodoTool struct {
	store storage.TodoStore
	limit int
}

func NewDeleteTool(store storage.TodoStore, limit int) *DeleteTodoTool {
	return &DeleteTodoTool{
		store: store,
		limit: limit,
	}
}

func (d *DeleteTodoTool) Definition() openai.ChatCompletionToolUnionParam {
	return helpers.FnDefinition(
		"delete_todo",
		`Permanently remove a task from the list. 
		Use when the user asks to delete, remove, or get rid of a task. 
		Do NOT use this to mark a task as completed — use complete_todo for that. 
		Requires the task ID — call list_todos first if needed.
		The list_todos result shows tasks with their IDs. When the user says 
		"task 1" or "first task", clarify whether they mean the ID or the 
		position in the list.`,
		shared.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"task_id": map[string]any{
					"type": "integer",
					"description": `The numeric ID of the task to delete, taken from the list_todos result. 
						Do NOT guess or invent the ID — if you don't know it, 
						call list_todos first to find the correct one.`,
				},
				"confirmed": map[string]any{
					"type": "boolean",
					"description": `Set to true only if the user has explicitly confirmed they want to delete the task 
					(e.g. after you asked "are you sure?" and they said yes). Leave false by default. 
					Never set to true based solely on the user's initial request — always ask for confirmation first.`,
				},
			},
			"required": []string{"task_id"},
		},
	)
}

func (d *DeleteTodoTool) Call(
	ctx context.Context,
	tc openai.ChatCompletionMessageToolCallUnion,
) openai.ChatCompletionMessageParamUnion {
	var args checkTaskID
	if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
		b, _ := json.Marshal(map[string]string{"error": "bad args: "+err.Error()})
		return openai.ToolMessage(string(b), tc.ID)
	}

	err := d.store.DeleteTodo(args.TaskID)
	if err != nil {
		b, _ := json.Marshal(map[string]string{"error": err.Error()+" - it may have been already deleted. Do not retry."})
		return openai.ToolMessage(string(b), tc.ID)
	}

	b, _ := json.Marshal(map[string]int{"delete id": args.TaskID})
	return openai.ToolMessage(string(b), tc.ID)
}

func (d *DeleteTodoTool) MaxCallsPerRun() int {
	return d.limit
}