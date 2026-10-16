package search

import (
	"context"
	"fmt"
)

type MockClient struct {
	Default []Result
}

func NewMock() *MockClient {
	return &MockClient{
		Default: []Result{
			{
				Title: "Hello, world!",
				URL: "http://hello.world.com",
				Snippet: "Hello, World!",
			},
		},
	}
}

func (c *MockClient) Search(ctx context.Context, query string) ([]Result, error) {
	if query == "" {
		return nil, fmt.Errorf("empty query")
	}

	switch query{
		case "пирог":
			return []Result{
				{
					Title: "Рецепт пирога",
					URL: "https://recipes.com/pie",
					Snippet: `1. 2 яйца
					2. 2 столовые ложки муки
					...`,
				},
			}, nil
		default:
			return c.Default, nil
	}
}