package search

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	httpClient *http.Client
	baseURL string
	format string
	language string
	safesearch string
}

func NewClient(httpClient *http.Client, baseUrl, format, lang, safesearch string) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL: baseUrl,
		format: format,
		language: lang,
		safesearch: safesearch,
	}
}

func (c *Client) Search(ctx context.Context, query string) ([]Result, error) {
	baseURL := c.baseURL + "/search"
	params := url.Values{}
	params.Set("q", query)
	params.Set("format", c.format)
	params.Set("language", c.language)
	params.Set("safesearch", c.safesearch)

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("http request creation failed: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search upstream failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search upstream failed: status=%d body %s", resp.StatusCode, body)
	}

	var res SearxResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("invalid upstream payload: %v", err)
	}

	if len(res.Results) == 0 {
		return res.Results, nil
	}

	return res.Results, nil
}