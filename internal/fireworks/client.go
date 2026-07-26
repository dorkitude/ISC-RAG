// Package fireworks implements the small OpenAI-compatible API surface used by ISC-RAG.
package fireworks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Client struct {
	baseURL, apiKey string
	http            *http.Client
}

func NewFromEnv() (*Client, error) {
	key := strings.TrimSpace(os.Getenv("FIREWORKS_API_KEY"))
	if key == "" {
		return nil, fmt.Errorf("FIREWORKS_API_KEY is required")
	}
	base := strings.TrimRight(strings.TrimSpace(os.Getenv("FIREWORKS_API_BASE")), "/")
	if base == "" {
		base = "https://api.fireworks.ai/inference/v1"
	}
	return &Client{baseURL: base, apiKey: key, http: &http.Client{Timeout: 2 * time.Minute}}, nil
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
type chatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Usage usage `json:"usage"`
}
type ChatResult struct {
	Text                                        string
	PromptTokens, CompletionTokens, TotalTokens int
}

func (c *Client) Chat(ctx context.Context, model, system, user string, maxTokens int) (ChatResult, error) {
	return c.chat(ctx, model, system, user, maxTokens, false)
}

// ChatJSON asks the OpenAI-compatible endpoint for a JSON object, while callers
// still validate the resulting schema themselves.
func (c *Client) ChatJSON(ctx context.Context, model, system, user string, maxTokens int) (ChatResult, error) {
	return c.chat(ctx, model, system, user, maxTokens, true)
}

// ChatJSONNoReasoning is for models that otherwise spend the completion budget
// on visible reasoning before returning a structured response. It is opt-in so
// existing experiment conditions are unaffected.
func (c *Client) ChatJSONNoReasoning(ctx context.Context, model, system, user string, maxTokens int) (ChatResult, error) {
	return c.chatWithOptions(ctx, model, system, user, maxTokens, true, true)
}

func (c *Client) chat(ctx context.Context, model, system, user string, maxTokens int, jsonMode bool) (ChatResult, error) {
	return c.chatWithOptions(ctx, model, system, user, maxTokens, jsonMode, false)
}

func (c *Client) chatWithOptions(ctx context.Context, model, system, user string, maxTokens int, jsonMode, noReasoning bool) (ChatResult, error) {
	payload := map[string]any{"model": model, "messages": []Message{{"system", system}, {"user", user}}, "temperature": 0, "max_tokens": maxTokens}
	if jsonMode {
		payload["response_format"] = map[string]string{"type": "json_object"}
	}
	if noReasoning {
		payload["reasoning_effort"] = "none"
	}
	var response chatResponse
	if err := c.request(ctx, "/chat/completions", payload, &response); err != nil {
		return ChatResult{}, err
	}
	if len(response.Choices) == 0 {
		return ChatResult{}, fmt.Errorf("Fireworks returned no completion choices")
	}
	return ChatResult{Text: response.Choices[0].Message.Content, PromptTokens: response.Usage.PromptTokens, CompletionTokens: response.Usage.CompletionTokens, TotalTokens: response.Usage.TotalTokens}, nil
}

type embedResponse struct {
	Data []struct {
		Embedding []float64 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Usage usage `json:"usage"`
}

func (c *Client) Embed(ctx context.Context, model string, input []string) ([][]float64, error) {
	var response embedResponse
	if err := c.request(ctx, "/embeddings", map[string]any{"model": model, "input": input}, &response); err != nil {
		return nil, err
	}
	if len(response.Data) != len(input) {
		return nil, fmt.Errorf("Fireworks returned %d embeddings for %d inputs", len(response.Data), len(input))
	}
	result := make([][]float64, len(input))
	for _, item := range response.Data {
		if item.Index < 0 || item.Index >= len(result) {
			return nil, fmt.Errorf("invalid embedding index %d", item.Index)
		}
		result[item.Index] = item.Embedding
	}
	return result, nil
}

func (c *Client) request(ctx context.Context, path string, payload any, destination any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	var last error
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Content-Type", "application/json")
		res, err := c.http.Do(req)
		if err != nil {
			last = err
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if err = retryPause(ctx, attempt); err != nil {
				return err
			}
			continue
		}
		raw, readErr := io.ReadAll(io.LimitReader(res.Body, 10<<20))
		_ = res.Body.Close()
		if readErr != nil {
			last = readErr
			if err = retryPause(ctx, attempt); err != nil {
				return err
			}
			continue
		}
		if res.StatusCode < 200 || res.StatusCode >= 300 {
			last = fmt.Errorf("Fireworks %s: %s", res.Status, strings.TrimSpace(string(raw)))
			if res.StatusCode == http.StatusTooManyRequests || res.StatusCode >= 500 {
				if err = retryPause(ctx, attempt); err != nil {
					return err
				}
				continue
			}
			return last
		}
		if err = json.Unmarshal(raw, destination); err != nil {
			return fmt.Errorf("decode Fireworks response: %w", err)
		}
		return nil
	}
	return last
}

func retryPause(ctx context.Context, attempt int) error {
	if attempt >= 2 {
		return nil
	}
	timer := time.NewTimer(time.Duration(attempt+1) * time.Second)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

type rerankResponse struct {
	Data []struct {
		Index          int     `json:"index"`
		RelevanceScore float64 `json:"relevance_score"`
	} `json:"data"`
}

// RerankScore pairs an input document index with the reranker's relevance.
type RerankScore struct {
	Index int
	Score float64
}

// Rerank scores documents against a query with a Fireworks reranker model and
// returns the scores sorted most-relevant first, indexed into the input slice.
func (c *Client) Rerank(ctx context.Context, model, query string, documents []string) ([]RerankScore, error) {
	payload := map[string]any{"model": model, "query": query, "documents": documents}
	var response rerankResponse
	if err := c.request(ctx, "/rerank", payload, &response); err != nil {
		return nil, err
	}
	if len(response.Data) != len(documents) {
		return nil, fmt.Errorf("reranker returned %d scores for %d documents", len(response.Data), len(documents))
	}
	out := make([]RerankScore, len(response.Data))
	for i, item := range response.Data {
		if item.Index < 0 || item.Index >= len(documents) {
			return nil, fmt.Errorf("reranker returned out-of-range index %d", item.Index)
		}
		out[i] = RerankScore{Index: item.Index, Score: item.RelevanceScore}
	}
	return out, nil
}
