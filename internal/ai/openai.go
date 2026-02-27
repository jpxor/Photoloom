package ai

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

type OpenAIProvider struct {
	apiKey       string
	model        string
	previewDir   string
	previewWidth int
}

func NewOpenAIProvider(apiKey, model, previewDir string, previewWidth int) *OpenAIProvider {
	return &OpenAIProvider{
		apiKey:       apiKey,
		model:        model,
		previewDir:   previewDir,
		previewWidth: previewWidth,
	}
}

func (p *OpenAIProvider) Analyze(imageData []byte, prompt string, imagePath string) (*Metadata, error) {
	img, err := imaging.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}

	img = imaging.Resize(img, p.previewWidth, p.previewWidth, imaging.Lanczos)

	var buf bytes.Buffer
	if err := imaging.Encode(&buf, img, imaging.JPEG); err != nil {
		return nil, fmt.Errorf("encoding image: %w", err)
	}

	if imagePath != "" && p.previewDir != "" {
		savePath := filepath.Join(p.previewDir, imagePath+".jpg")
		if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
			return nil, fmt.Errorf("creating preview dir: %w", err)
		}
		if err := os.WriteFile(savePath, buf.Bytes(), 0644); err != nil {
			return nil, fmt.Errorf("writing preview: %w", err)
		}
	}

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	url := "https://api.openai.com/v1/chat/completions"

	systemPrompt := "You are an expert at analyzing photographs. " + prompt

	payload := map[string]interface{}{
		"model": p.model,
		"messages": []map[string]interface{}{
			{
				"role":    "system",
				"content": systemPrompt,
			},
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "image_url",
						"image_url": map[string]interface{}{
							"url":    "data:image/jpeg;base64," + encoded,
							"detail": "low",
						},
					},
				},
			},
		},
		"max_tokens": 1000,
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
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("OpenAI API response status: %s", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %s %s", resp.Status, string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling response: %w", err)
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no response from API")
	}

	return ParseResponse(result.Choices[0].Message.Content)
}
