package ai

import (
	"fmt"

	"photo-gallery/internal/config"
)

type Provider interface {
	Analyze(imageData []byte, prompt string, imagePath string) (*Metadata, error)
}

func NewProvider(cfg *config.Config) (Provider, error) {
	switch cfg.AI.Provider {
	case "openai":
		previewWidth := cfg.AI.PreviewWidth
		if previewWidth == 0 {
			previewWidth = 512
		}
		return NewOpenAIProvider(
			cfg.AI.OpenAI.APIKey,
			cfg.AI.OpenAI.Model,
			cfg.AI.PreviewDir,
			previewWidth,
		), nil
	case "ollama":
		return NewOllamaProvider(
			cfg.AI.Ollama.BaseURL,
			cfg.AI.Ollama.Model,
		), nil
	case "none":
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown AI provider: %s", cfg.AI.Provider)
	}
}
