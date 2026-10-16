package weathertools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"yuima-agent/internal/tools/helpers"
	"yuima-agent/internal/weather"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

// WeatherProvider — то, что нужно тулу от погодного слоя.
// Интерфейс живёт здесь, потому что потребляется здесь.
type WeatherProvider interface {
	GetWeather(ctx context.Context, city string) (*weather.WeatherData, error)
}

// Tool — инструмент weather_check.
type Tool struct {
	client WeatherProvider
	limit int
}

type checkWeatherArgs struct {
	City string `json:"city"`
}

// NewWeatherTool принимает любой WeatherProvider (реальный или фейк в тестах).
func NewCheckTool(c WeatherProvider, l int) *Tool {
	return &Tool{
		client: c,
		limit: l,
	}
}

// Definition описывает инструмент для LLM.
func (t *Tool) Definition() openai.ChatCompletionToolUnionParam {
	return helpers.FnDefinition(
		"weather_check",
		"Поиск информации о погоде в выбранном городе",
		shared.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"city": map[string]any{
					"type":        "string",
					"description": "The city to check for weather in",
				},
			},
			"required": []string{"city"},
		},
	)
}

// Call выполняет один tool_call и возвращает готовое ToolMessage.
func (t *Tool) Call(
	ctx context.Context,
	tc openai.ChatCompletionMessageToolCallUnion,
) openai.ChatCompletionMessageParamUnion {
	var args checkWeatherArgs
	if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
		b, _ := json.Marshal(map[string]string{"error": "bad args: " + err.Error()})
		return openai.ToolMessage(string(b), tc.ID)
	}

	data, err := t.client.GetWeather(ctx, args.City)
	if err != nil {
		var msg []byte
		switch {
		case errors.Is(err, weather.ErrCityNotFound):
			// LLM полезно знать — попробует уточнить у пользователя.
			msg, _ = json.Marshal(map[string]string{"error": fmt.Sprintf(`city %q not found, ask user to clarify`, args.City)})
		case errors.Is(err, weather.ErrUpstreamFailed):
			msg = []byte(`{"error":"weather service unavailable, try again later"}`)
		default:
			msg, _ = json.Marshal(map[string]string{"error": err.Error()})
		}
		return openai.ToolMessage(string(msg), tc.ID)
	}

	b, _ := json.Marshal(data)
	return openai.ToolMessage(string(b), tc.ID)
}

func (t *Tool) MaxCallsPerRun() int {
	return t.limit
}