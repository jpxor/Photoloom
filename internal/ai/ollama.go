package ai

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OllamaProvider struct {
	baseURL string
	model   string
}

func NewOllamaProvider(baseURL, model string) *OllamaProvider {
	return &OllamaProvider{baseURL: baseURL, model: model}
}

func (p *OllamaProvider) Analyze(imageData []byte, prompt string, imagePath string) (*Metadata, error) {
	encoded := base64.StdEncoding.EncodeToString(imageData)

	url := p.baseURL + "/api/generate"

	promptFull := prompt + "\n\nAnalyze this image and respond in the exact format specified."

	payload := map[string]interface{}{
		"model":  p.model,
		"prompt": promptFull,
		"images": []string{encoded},
		"stream": false,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshaling payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %s %s", resp.Status, string(body))
	}

	var result struct {
		Response string `json:"response"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling response: %w", err)
	}

	return ParseResponse(result.Response)
}
