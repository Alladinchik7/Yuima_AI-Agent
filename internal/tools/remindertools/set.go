package remindertools

import (
	"context"
	"encoding/json"
	"time"
	"yuima-agent/internal/tools/helpers"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

type ReminderCreator interface {
    AddReminder(text string, fireAt time.Time) error
}

type SetReminderTool struct {
    store ReminderCreator
    wakeUp chan struct{}
    limit int
}

type setReminderArgs struct {
    Text   string `json:"text"`
    FireAt string `json:"fire_at"`
}

func NewSetTool(store ReminderCreator, wakeUp chan struct{}, limit int) *SetReminderTool {
	return &SetReminderTool{
		store: store, 
		wakeUp: wakeUp,
		limit: limit,
	}
}

func (s *SetReminderTool) Definition() openai.ChatCompletionToolUnionParam {
	return helpers.FnDefinition(
		"set_reminder",
		`Create a reminder that will fire at a specific time. 
		Use when the user asks to be reminded about something at a particular moment, 
		e.g. "remind me to call mom tomorrow at 12" or "remind me in an hour".
		The fire_at parameter must be an ISO 8601 datetime string in the user's local timezone, 
		e.g. "2026-10-09T12:00:00+03:00".`,
		shared.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"text": map[string]any{
					"type": "string",
					"description": `The reminder text, e.g. "call mom".`},
				"fire_at": map[string]any{
					"type": "string",
					"description": `ISO 8601 datetime in the user's local timezone, 
					e.g. "2026-10-09T12:00:00+03:00". 
					Do NOT use relative phrases like "tomorrow" — convert them to an absolute datetime.`,
				},
			},
			"required": []string{"text", "fire_at"},
		},
	)
}

func (s *SetReminderTool) Call(
	ctx context.Context,
	tc openai.ChatCompletionMessageToolCallUnion,
) openai.ChatCompletionMessageParamUnion {
	var args setReminderArgs
	if  err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
		b, _ := json.Marshal(map[string]string{"error": "bad args: " + err.Error()})
		return openai.ToolMessage(string(b), tc.ID)
	}

	if args.Text == "" {
		b, _ := json.Marshal(map[string]string{"error": "text missing"})
		return openai.ToolMessage(string(b), tc.ID)
	}
	
	fireAt, err := time.Parse(time.RFC3339, args.FireAt)
	if err != nil {
		b, _ := json.Marshal(map[string]string{"error": "invalid fire_at format, expected ISO 8601, got: " + err.Error()})
		return openai.ToolMessage(string(b), tc.ID)
	}

	if fireAt.Before(time.Now()) {
		b, _ := json.Marshal(map[string]string{"error": "fire_at is in the past"})
		return openai.ToolMessage(string(b), tc.ID)
	}

	if err := s.store.AddReminder(args.Text, fireAt); err != nil {
		b, _ := json.Marshal(map[string]string{"error": "failed add reminder: " + err.Error()})
		return openai.ToolMessage(string(b), tc.ID)
	}

	select{
		case s.wakeUp <- struct{}{}:
		default:
		
	}

	b, _ := json.Marshal(map[string]any{"created": true, "text": args.Text, "fire_at": args.FireAt})
	return openai.ToolMessage(string(b), tc.ID)
}

func (s *SetReminderTool) MaxCallsPerRun() int {
	return s.limit
}