package elasticsearch

type (
	// SearchRequest represents an incoming search from the app
	SearchRequest struct {
		Query     string   `json:"query"`
		Filters   []string `json:"filters"`
		Offset    int      `json:"offset"`
		Limit     int      `json:"limit"`
		Facets    []string `json:"facets"`
		Stats     []string `json:"stats"`
		Fields    []string `json:"fields"`
		Highlight []string `json:"highlight"`
	}

	// SearchResult represents an outgoing result to the app
	SearchResult struct {
		Total     int64                               `json:"total"`
		QueryTime int64                               `json:"qtime"`
		Documents []map[string]interface{}            `json:"docs"`
		Stats     map[string]interface{}              `json:"stats"`
		Facets    map[string][]map[string]interface{} `json:"facets"`
	}

	IndexMetadata struct {
		Name        string  `json:"name"`
		Size        float32 `json:"indexSize"`
		SizeContent float32 `json:"contentSize"`
	}
)

var (
	titleExclude = []string{
		"Back Matter",
		"Front Matter",
	}
)
