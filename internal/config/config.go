package config

import (
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Paths  PathsConfig  `yaml:"paths"`
	Images ImagesConfig `yaml:"images"`
	AI     AIConfig     `yaml:"ai"`
	Albums AlbumsConfig `yaml:"albums"`
	Hugo   HugoConfig   `yaml:"hugo"`
	Build  BuildConfig  `yaml:"build"`
}

type PathsConfig struct {
	PhotosDir string `yaml:"photosDir"`
	CacheDir  string `yaml:"cacheDir"`
	OutputDir string `yaml:"outputDir"`
}

type ImagesConfig struct {
	Thumbnail          ImageConfig     `yaml:"thumbnail"`
	Web                ImageConfig     `yaml:"web"`
	WatermarkThreshold int             `yaml:"watermarkThreshold"`
	Watermark          WatermarkConfig `yaml:"watermark"`
}

type ImageConfig struct {
	Width   int    `yaml:"width"`
	Height  int    `yaml:"height"`
	Quality int    `yaml:"quality"`
	Format  string `yaml:"format"`
	Fit     string `yaml:"fit"`
}

type WatermarkConfig struct {
	Text     string  `yaml:"text"`
	Opacity  float64 `yaml:"opacity"`
	FontPath string  `yaml:"fontPath"`
	Position string  `yaml:"position"`
}

type AIConfig struct {
	Provider     string       `yaml:"provider"`
	OpenAI       OpenAIConfig `yaml:"openai"`
	Ollama       OllamaConfig `yaml:"ollama"`
	Prompt       string       `yaml:"prompt"`
	PreviewDir   string       `yaml:"previewDir"`
	PreviewWidth int          `yaml:"previewWidth"`
}

type OpenAIConfig struct {
	APIKey      string `yaml:"apiKey"`
	Model       string `yaml:"model"`
	Concurrency int    `yaml:"concurrency"`
}

type OllamaConfig struct {
	BaseURL string `yaml:"baseURL"`
	Model   string `yaml:"model"`
}

type AlbumsConfig struct {
	GroupBy   []string `yaml:"groupBy"`
	MinPhotos int      `yaml:"minPhotos"`
}

type HugoConfig struct {
	Theme string   `yaml:"theme"`
	Flags []string `yaml:"flags"`
}

type BuildConfig struct {
	CleanCache bool `yaml:"cleanCache"`
	Workers    int  `yaml:"workers"`
}

func Load(path string) (*Config, error) {
	godotenv.Load()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.AI.OpenAI.APIKey != "" && cfg.AI.OpenAI.APIKey[0] == '$' {
		envVar := cfg.AI.OpenAI.APIKey[1:]
		cfg.AI.OpenAI.APIKey = os.Getenv(envVar)
	}

	if cfg.Build.Workers == 0 {
		cfg.Build.Workers = 4
	}

	if cfg.Images.WatermarkThreshold == 0 {
		cfg.Images.WatermarkThreshold = 2048
	}

	return &cfg, nil
}
