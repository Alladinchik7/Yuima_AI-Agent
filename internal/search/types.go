package search

type Result struct {
	Title string   `json:"title"`
	URL string     `json:"url"`
	Snippet string `json:"content"`
}

type SearxResponse struct {
	Results []Result `json:"results"`
}