package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// GeminiRecommendations asks Gemini for practical suggestions based only on
// aggregate task-completion statistics. Task titles and notes are never sent.
type GeminiRecommendations struct {
	apiKey string
	model  string
	client *http.Client
}

func NewGeminiRecommendations(apiKey, model string) *GeminiRecommendations {
	return &GeminiRecommendations{
		apiKey: strings.TrimSpace(apiKey),
		model:  strings.TrimSpace(model),
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

type recommendationStats struct {
	CompletedThisWeek  int                 `json:"completed_this_week"`
	CompletedThisMonth int                 `json:"completed_this_month"`
	CompletedOnTime    int                 `json:"completed_on_time"`
	Uncompleted        int                 `json:"uncompleted"`
	Overdue            int                 `json:"overdue"`
	WeekDays           []recommendationDay `json:"week_days"`
	MonthDays          []recommendationDay `json:"month_days"`
}

type recommendationDay struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

func (g *GeminiRecommendations) Generate(ctx context.Context, language string, statsJSON []byte) (string, error) {
	if g.apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY is not configured")
	}
	var stats recommendationStats
	if err := json.Unmarshal(statsJSON, &stats); err != nil {
		return "", fmt.Errorf("decode analytics stats: %w", err)
	}
	data, err := json.Marshal(stats)
	if err != nil {
		return "", fmt.Errorf("encode analytics stats: %w", err)
	}
	if language == "" {
		language = "ru"
	}
	prompt := "Ты помощник по планированию учебных и личных задач. На основе приведённых агрегированных чисел составь 3 короткие, конкретные рекомендации на языке " + language +
		". Учитывай просроченные и незавершённые задачи, динамику выполнения по дням и долю выполненных в срок. Не придумывай причин и фактов, если данных мало — скажи об этом. Верни только рекомендации нумерованным списком, без вступления. Данные: " + string(data)
	body, err := json.Marshal(map[string]any{
		"contents":         []any{map[string]any{"parts": []any{map[string]string{"text": prompt}}}},
		"generationConfig": map[string]any{"temperature": 0.4, "maxOutputTokens": 500},
	})
	if err != nil {
		return "", fmt.Errorf("encode Gemini request: %w", err)
	}
	endpoint := "https://generativelanguage.googleapis.com/v1beta/models/" + g.model + ":generateContent"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create Gemini request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", g.apiKey)
	resp, err := g.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("call Gemini API: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("Gemini API returned HTTP %d", resp.StatusCode)
	}
	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&result); err != nil {
		return "", fmt.Errorf("decode Gemini response: %w", err)
	}
	var parts []string
	if len(result.Candidates) > 0 {
		for _, part := range result.Candidates[0].Content.Parts {
			if text := strings.TrimSpace(part.Text); text != "" {
				parts = append(parts, text)
			}
		}
	}
	answer := strings.TrimSpace(strings.Join(parts, "\n"))
	if answer == "" {
		return "", fmt.Errorf("Gemini returned an empty recommendation")
	}
	if len(answer) > 4000 {
		answer = answer[:4000]
	}
	return answer, nil
}
