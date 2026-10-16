package search

import "context"

type SearchClient interface {
	Search(context.Context, string) ([]Result, error)
}

type Result struct {
	Title string
	URL string
	Snippet string
}