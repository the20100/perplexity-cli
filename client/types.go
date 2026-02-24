package client

// --- Search API ---

type SearchRequest struct {
	Query                 any      `json:"query"` // string or []string
	MaxTokens             int      `json:"max_tokens,omitempty"`
	MaxTokensPerPage      int      `json:"max_tokens_per_page,omitempty"`
	MaxResults            int      `json:"max_results,omitempty"`
	SearchDomainFilter    []string `json:"search_domain_filter,omitempty"`
	SearchLanguageFilter  []string `json:"search_language_filter,omitempty"`
	SearchRecencyFilter   string   `json:"search_recency_filter,omitempty"` // hour, day, week, month, year
	SearchAfterDateFilter string   `json:"search_after_date_filter,omitempty"`
	SearchBeforeDateFilter string  `json:"search_before_date_filter,omitempty"`
	LastUpdatedBeforeFilter string `json:"last_updated_before_filter,omitempty"`
	LastUpdatedAfterFilter  string `json:"last_updated_after_filter,omitempty"`
	SearchMode            string   `json:"search_mode,omitempty"` // web, academic, sec
	Country               string   `json:"country,omitempty"`
	DisplayServerTime     bool     `json:"display_server_time,omitempty"`
}

type SearchPage struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Snippet     string `json:"snippet"`
	Date        string `json:"date,omitempty"`
	LastUpdated string `json:"last_updated,omitempty"`
}

type SearchResponse struct {
	ID         string       `json:"id"`
	Results    []SearchPage `json:"results"`
	ServerTime string       `json:"server_time,omitempty"`
}

// --- Chat Completions API ---

type Message struct {
	Role    string `json:"role"`    // system, user, assistant
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream,omitempty"`

	// Optional parameters
	MaxTokens         int      `json:"max_tokens,omitempty"`
	Temperature       float64  `json:"temperature,omitempty"`
	TopP              float64  `json:"top_p,omitempty"`
	TopK              int      `json:"top_k,omitempty"`
	PresencePenalty   float64  `json:"presence_penalty,omitempty"`
	FrequencyPenalty  float64  `json:"frequency_penalty,omitempty"`
	SearchDomainFilter []string `json:"search_domain_filter,omitempty"`
	SearchRecencyFilter string  `json:"search_recency_filter,omitempty"`
	ReturnImages      bool     `json:"return_images,omitempty"`
	ReturnRelatedQuestions bool `json:"return_related_questions,omitempty"`
}

type ChatChoice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type ChatDelta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type ChatStreamChoice struct {
	Index        int       `json:"index"`
	Delta        ChatDelta `json:"delta"`
	FinishReason string    `json:"finish_reason,omitempty"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Citation struct {
	URL   string `json:"url"`
	Title string `json:"title,omitempty"`
}

type ChatResponse struct {
	ID      string       `json:"id"`
	Model   string       `json:"model"`
	Choices []ChatChoice `json:"choices"`
	Usage   Usage        `json:"usage"`
	Citations []Citation `json:"citations,omitempty"`
}

type ChatStreamChunk struct {
	ID      string             `json:"id"`
	Model   string             `json:"model"`
	Choices []ChatStreamChoice `json:"choices"`
	Citations []Citation       `json:"citations,omitempty"`
}

// --- Content API ---

type ContentRequest struct {
	URLs []string `json:"urls"`
}

type ContentPage struct {
	URL     string `json:"url"`
	Content string `json:"content"`
	Title   string `json:"title,omitempty"`
}

type ContentResponse struct {
	ID      string        `json:"id"`
	Results []ContentPage `json:"results"`
}

// --- Error ---

type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return "perplexity API error " + itoa(e.StatusCode) + ": " + e.Body
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	b := make([]byte, 0, 10)
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}
