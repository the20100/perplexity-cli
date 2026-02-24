package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const baseURL = "https://api.perplexity.ai"

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func New(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (c *Client) do(method, path string, body any, out any) error {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{StatusCode: resp.StatusCode, Body: string(respBody)}
	}

	if out != nil {
		if err := json.Unmarshal(respBody, out); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}
	return nil
}

// Search performs a web search.
func (c *Client) Search(req *SearchRequest) (*SearchResponse, error) {
	var resp SearchResponse
	if err := c.do("POST", "/search", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ChatComplete sends a chat request and returns the full response.
func (c *Client) ChatComplete(req *ChatRequest) (*ChatResponse, error) {
	req.Stream = false
	var resp ChatResponse
	if err := c.do("POST", "/chat/completions", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ChatStream sends a chat request and streams chunks to the provided callback.
// The callback receives each content delta as it arrives.
// The final full response is returned when streaming completes.
func (c *Client) ChatStream(req *ChatRequest, onChunk func(delta string)) (*ChatStreamChunk, error) {
	req.Stream = true

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", baseURL+"/chat/completions", bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, &APIError{StatusCode: resp.StatusCode, Body: string(body)}
	}

	var lastChunk *ChatStreamChunk
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}
		var chunk ChatStreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}
		if len(chunk.Choices) > 0 {
			delta := chunk.Choices[0].Delta.Content
			if delta != "" && onChunk != nil {
				onChunk(delta)
			}
		}
		lastChunk = &chunk
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("stream read: %w", err)
	}
	return lastChunk, nil
}

// GetContent extracts content from URLs.
func (c *Client) GetContent(urls []string) (*ContentResponse, error) {
	req := &ContentRequest{URLs: urls}
	var resp ContentResponse
	if err := c.do("POST", "/content", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
