package tools

import (
	"context"
	"sort"

	"github.com/openai/openai-go/v3"
)

// Tool — контракт, который знает агент. Конкретные инструменты
// (погода, время и т.д.) реализуют его в отдельных файлах.
type Tool interface {
	Definition() openai.ChatCompletionToolUnionParam
	Call(ctx context.Context, tc openai.ChatCompletionMessageToolCallUnion) openai.ChatCompletionMessageParamUnion
	MaxCallsPerRun() int
}

// Registry — реестр тулов по имени функции.
type Registry struct {
	tools map[string]Tool
	ordered []Tool
}

// NewRegistry собирает реестр. Имя берётся из Definition().GetFunction().Name.
// GetFunction() возвращает *shared.FunctionDefinitionParam, поэтому имя — fn.Name.
func NewRegistry(ts ...Tool) *Registry {
	r := &Registry{tools: make(map[string]Tool, len(ts))}
	for _, t := range ts {
		def := t.Definition()
		fn := def.GetFunction()
		if fn == nil {
			// Не function-tool (например, custom) — пропускаем.
			continue
		}
		r.tools[fn.Name] = t
		r.ordered = append(r.ordered, t)
	}
	sort.Slice(r.ordered, func(i, j int) bool {
		return r.ordered[i].Definition().GetFunction().Name < r.ordered[j].Definition().GetFunction().Name
	})
	return r
}

// Definitions — для поля Tools в ChatCompletionNewParams.
func (r *Registry) Definitions() []openai.ChatCompletionToolUnionParam {
	out := make([]openai.ChatCompletionToolUnionParam, 0, len(r.ordered))
	for _, t := range r.ordered {
		out = append(out, t.Definition())
	}
	return out
}

// Get ищет тул по имени из tool_call.
func (r *Registry) Get(name string) (Tool, bool) {
	t, ok := r.tools[name]
	return t, ok
}