package ai

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"photo-gallery/internal/config"
	"photo-gallery/internal/exif"
)

type Pipeline struct {
	cfg         *config.Config
	provider    Provider
	metadataDir string
}

func NewPipeline(cfg *config.Config, metadataDir string) (*Pipeline, error) {
	provider, err := NewProvider(cfg)
	if err != nil {
		return nil, err
	}

	return &Pipeline{
		cfg:         cfg,
		provider:    provider,
		metadataDir: metadataDir,
	}, nil
}

func (p *Pipeline) Run() error {
	if p.cfg.AI.Provider == "none" {
		log.Println("AI provider is 'none', skipping AI metadata extraction")
		return nil
	}

	photosDir := p.cfg.Paths.PhotosDir

	if err := os.MkdirAll(p.metadataDir, 0755); err != nil {
		return fmt.Errorf("creating metadata dir: %w", err)
	}

	images, err := findImages(photosDir)
	if err != nil {
		return fmt.Errorf("finding images: %w", err)
	}

	fmt.Printf("Extracting metadata from %d images using %s\n", len(images), p.cfg.AI.Provider)

	processed := 0
	errors := 0

	concurrency := p.getConcurrency()
	sem := make(chan struct{}, concurrency)

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, imgPath := range images {
		wg.Add(1)
		go func(imgPath string) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			relPath, err := filepath.Rel(photosDir, imgPath)
			if err != nil {
				log.Printf("Error getting relative path for %s: %v", imgPath, err)
				mu.Lock()
				errors++
				mu.Unlock()
				return
			}

			metadataPath := filepath.Join(p.metadataDir, relPath+".json")
			metadataDirPath := filepath.Dir(metadataPath)

			if err := os.MkdirAll(metadataDirPath, 0755); err != nil {
				log.Printf("Error creating metadata dir for %s: %v", relPath, err)
				mu.Lock()
				errors++
				mu.Unlock()
				return
			}

			if shouldSkipMetadata(imgPath, metadataPath) {
				mu.Lock()
				processed++
				mu.Unlock()
				return
			}

			data, err := os.ReadFile(imgPath)
			if err != nil {
				log.Printf("Error reading file %s: %v", imgPath, err)
				mu.Lock()
				errors++
				mu.Unlock()
				return
			}

			result := NewMetadata()

			exifData, err := exif.Extract(imgPath)
			if err == nil {
				result.DateTaken = exifData.DateTaken
				result.CameraMake = exifData.CameraMake
				result.CameraModel = exifData.CameraModel
				result.Lens = exifData.Lens
				result.ISO = exifData.ISO
				result.Aperture = exifData.Aperture
				result.ShutterSpeed = exifData.ShutterSpeed
				result.FocalLength = exifData.FocalLength
				result.Width = exifData.Width
				result.Height = exifData.Height
			}

			if p.provider != nil {
				aiResult, err := p.provider.Analyze(data, p.cfg.AI.Prompt, relPath)
				if err != nil {
					log.Printf("Error analyzing %s: %v", relPath, err)
				} else {
					result.Category = aiResult.Category
					result.Tags = aiResult.Tags
					result.Colors = aiResult.Colors
					result.Description = aiResult.Description
				}
			}

			result.SourcePath = relPath

			jsonData, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				log.Printf("Error marshaling JSON for %s: %v", relPath, err)
				mu.Lock()
				errors++
				mu.Unlock()
				return
			}

			if err := os.WriteFile(metadataPath, jsonData, 0644); err != nil {
				log.Printf("Error writing metadata for %s: %v", relPath, err)
				mu.Lock()
				errors++
				mu.Unlock()
				return
			}

			mu.Lock()
			processed++
			mu.Unlock()
		}(imgPath)
	}

	wg.Wait()

	fmt.Printf("Processed: %d, Errors: %d\n", processed, errors)

	return nil
}

func (p *Pipeline) getConcurrency() int {
	switch p.cfg.AI.Provider {
	case "openai":
		if p.cfg.AI.OpenAI.Concurrency > 0 {
			return p.cfg.AI.OpenAI.Concurrency
		}
		return 4
	case "ollama":
		return 2
	default:
		return 1
	}
}

func findImages(dir string) ([]string, error) {
	var images []string

	extensions := []string{".jpg", ".jpeg", ".png", ".webp"}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		for _, e := range extensions {
			if ext == e || ext == ".JPG" || ext == ".JPEG" {
				images = append(images, path)
				break
			}
		}
		return nil
	})

	return images, err
}

func shouldSkipMetadata(srcPath, metadataPath string) bool {
	metaInfo, err := os.Stat(metadataPath)
	if err != nil {
		return false
	}

	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return false
	}

	return !metaInfo.ModTime().Before(srcInfo.ModTime())
}
