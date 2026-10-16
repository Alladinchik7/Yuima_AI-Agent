package remindertools

import (
	"context"
	"encoding/json"
	"yuima-agent/internal/storage/models"
	"yuima-agent/internal/tools/helpers"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

type ReminderLister interface {
	ListReminders() ([]models.Reminder, error)
}

type ListRemindersTool struct {
	store ReminderLister
	limit int
}

type reminderItem struct {
	Pos    int    `json:"position"`
	ID     int    `json:"id"`
	Text   string `json:"text"`
	FireAt string `json:"fire_at"`
	Fired  bool   `json:"fired"`
}

func NewListTool(store ReminderLister, limit int) *ListRemindersTool {
	return &ListRemindersTool{
		store: store,
		limit: limit,
	}
}

func (l *ListRemindersTool) Definition() openai.ChatCompletionToolUnionParam {
	return helpers.FnDefinition(
		"list_reminders",
		"Show the user's current reminder list. Use when the user asks what's on their list or wants to see reminders.",
		shared.FunctionParameters{},
	)
}

func (l *ListRemindersTool) Call(
	ctx context.Context,
	tc openai.ChatCompletionMessageToolCallUnion,
) openai.ChatCompletionMessageParamUnion {
	reminders, err := l.store.ListReminders()
	if err != nil {
		b, _ := json.Marshal(map[string]string{"error":"failed list reminders: "+err.Error()})
		return openai.ToolMessage(string(b), tc.ID)
	}

	if len(reminders) == 0 {
		return openai.ToolMessage(`{"reminders":[]}`, tc.ID)
	}

	items := make([]reminderItem, len(reminders))
	for i, r := range reminders {
		items[i] = reminderItem{
			Pos: i+1,
			ID: r.ID,
			Text: r.Text,
			FireAt: r.FireAt.Format("2006-01-02 15:04:05"),
			Fired: r.Fired,
		}
	}

	b, err := json.Marshal(map[string]any{"reminders":items})
	if err != nil {
    	return openai.ToolMessage(`{"error":"marshal failed"}`, tc.ID)
	}
	return openai.ToolMessage(string(b), tc.ID)
}

func (l *ListRemindersTool) MaxCallsPerRun() int {
	return l.limit
}