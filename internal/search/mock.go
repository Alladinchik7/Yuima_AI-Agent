package search

import "context"

type MockClient struct {
}

func NewMock() *MockClient {
	return &MockClient{}
}

func (c *MockClient) Search(ctx context.Context, query string) ([]Result, error)