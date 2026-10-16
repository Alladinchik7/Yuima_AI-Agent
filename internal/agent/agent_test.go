package agent

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"yuima-agent/internal/storage/models"
	"yuima-agent/internal/tools"

	"github.com/openai/openai-go/v3"
)

type fakeStore struct {
	
}

func (s *fakeStore) SaveSummaryWithLimit(text string, limit int) error {
	return nil
}

func (s *fakeStore) LoadSummary() ([]models.DialogMemory, error) {
	return []models.DialogMemory{{DialogSummary: "empty summary"}}, nil
}

type fakeLLM struct {
	resp openai.ChatCompletionMessage
	lastMessages []openai.ChatCompletionMessageParamUnion
	calls int
}

func (f *fakeLLM) AskLLM(ctx context.Context,
		message []openai.ChatCompletionMessageParamUnion,
		tools []openai.ChatCompletionToolUnionParam,
) (openai.ChatCompletionMessage, error) {
	f.calls++
	f.lastMessages = make([]openai.ChatCompletionMessageParamUnion, len(message))
	copy(f.lastMessages, message)
	return f.resp, nil
}

func TestAgent_Run_EmptyInput_UseFakeLLM(t *testing.T) {
	cases := []struct{
		name string
		input string
	} {
		{"empty", ""},
		{"spaces", "   "},
		{"newlines", "\n\n\n"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fakeHistory := &fakeStore{}
			fake := &fakeLLM{}
			a := New(fakeHistory, fake, &tools.Registry{}, slog.Default(), 5*time.Minute)

			resp, err := a.Run(context.Background(), c.input)

			if err == nil {
				t.Errorf("input=%q: ожидалась ошибка", c.input)
			}

			if fake.calls != 0 {
                t.Errorf("input=%q: LLM вызвана %d раз", c.input, fake.calls)
            }

            if resp != "" {
            	t.Errorf("ожидался пустой ответ, получили %q", resp)
            }

            if len(fake.lastMessages) > 0 {
                t.Errorf("ожидалось 0 сообщений, получили %d", len(fake.lastMessages))
            }

            for i, m := range fake.lastMessages {
                if m.OfUser != nil {
                    t.Errorf("сообщение %d — user, а не должно быть", i)
                }
            }
		})
	}
}