package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/openai/openai-go/v3"
)

type MockLLM struct {
	Default  string
}

func NewMock() *MockLLM {
	return &MockLLM{
		Default:  "Hello, World!",
	}
}

func (m *MockLLM) AskLLM(
	ctx context.Context,
	history []openai.ChatCompletionMessageParamUnion,
	tools []openai.ChatCompletionToolUnionParam,
) (openai.ChatCompletionMessage, error) {
	if len(history) == 0 {
		return openai.ChatCompletionMessage{Content: m.Default}, nil
	}

	last, checkUsrMsg := lastUserMessage(history)
	if !checkUsrMsg {
		return openai.ChatCompletionMessage{}, fmt.Errorf("Не получилось достать строку пользователя")
	}

	if strings.HasPrefix(last, "Инструмент: ") {
		command := strings.TrimPrefix(last, "Инструмент: ")
		switch command {
		case "Погода":
			if toolMsg := checkLastMessageOfToolCall(history); !toolMsg {
				return openai.ChatCompletionMessage{
					Content: "",
					ToolCalls: []openai.ChatCompletionMessageToolCallUnion{
						{
							ID: "call_1",
							Function: openai.ChatCompletionMessageFunctionToolCallFunction{
								Name: "weather_check",
								Arguments: `{"city":"Москва"}`,
							},
						},
					},
				}, nil
			} else {
				return openai.ChatCompletionMessage{Content: "Погода в Москве +14"}, nil
			}
		case "Добавь задачу":
			if toolMsg := checkLastMessageOfToolCall(history); !toolMsg {
				return openai.ChatCompletionMessage{
					Content: "",
					ToolCalls: []openai.ChatCompletionMessageToolCallUnion{
						{
							ID: "call_2",
							Function: openai.ChatCompletionMessageFunctionToolCallFunction{
								Name: "add_todo",
								Arguments: `{"task":"купить хлеб"}`,
							},
						},
					},
				}, nil
			} else {
				return openai.ChatCompletionMessage{Content: "Задача добавлена"}, nil
			}
		default:
			return openai.ChatCompletionMessage{Content: "Инструмент не указан"}, nil
		}
	} else {
		switch last{
			case "Долго":
				select {
					case <-ctx.Done():
						return openai.ChatCompletionMessage{}, ctx.Err()
					case <-time.After(3 * time.Second):
						return openai.ChatCompletionMessage{Content:"Ожидание прошло"}, nil 
				}
			case "Эхо":
				return openai.ChatCompletionMessage{Content: last}, nil
			case "История":
				bytes, _ := json.Marshal(history)
				singleLineString := string(bytes)
				return openai.ChatCompletionMessage{Content: singleLineString}, nil	
			default:
				return openai.ChatCompletionMessage{Content: m.Default}, nil	
		}
	}
}

func checkLastMessageOfToolCall(history []openai.ChatCompletionMessageParamUnion) bool {
	lastUstIdx := -1
	for i := len(history)-1; i >=0; i-- {
		if history[i].OfUser != nil {
			lastUstIdx = i
			break
		}
	}

	if lastUstIdx == -1 {
		return false
	}

	for j := lastUstIdx + 1; j < len(history); j++ {
		if history[j].OfTool != nil {
			return true
		}
	}
	
	return false
}

func lastUserMessage(history []openai.ChatCompletionMessageParamUnion) (string, bool) {
	for i := len(history)-1; i >= 0; i-- {
		if history[i].OfUser != nil {
			return history[i].OfUser.Content.OfString.String(), true
		}
	}
	return "", false
}