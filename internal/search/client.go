package search

import (
	"context"
	"net/http"
)

type Client struct {
	httpClient http.Client
	baseURL string
	format string
	language string
	safesearch string
}

func NewClient(httpClient http.Client, baseUrl, format, lang, safesearch string) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL: baseUrl,
		format: format,
		language: lang,
		safesearch: safesearch,
	}
}

func (c *Client) Search(ctx context.Context, query string) ([]Result, error)