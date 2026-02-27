package pipeline

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"photo-gallery/internal/ai"
	"photo-gallery/internal/config"
	"photo-gallery/internal/hugo"
	"photo-gallery/internal/imageproc"
	"photo-gallery/internal/metadata"
)

type Pipeline struct {
	cfg         *config.Config
	metadataDir string
	contentDir  string
	assetsDir   string
}

func New(cfg *config.Config) *Pipeline {
	return &Pipeline{
		cfg:         cfg,
		metadataDir: filepath.Join(cfg.Paths.CacheDir, "metadata"),
		contentDir:  "hugo/content",
		assetsDir:   "hugo/assets/photos",
	}
}

func (p *Pipeline) RunAIMetadata() error {
	if p.cfg.AI.Provider == "none" {
		log.Println("Skipping AI metadata (provider is 'none')")
		return nil
	}

	log.Println("Running AI metadata extraction...")

	aiPipeline, err := ai.NewPipeline(p.cfg, p.metadataDir)
	if err != nil {
		return err
	}

	if err := aiPipeline.Run(); err != nil {
		return err
	}

	log.Println("AI metadata extraction complete")
	return nil
}

func (p *Pipeline) RunImageProcessing() error {
	log.Println("Running image processing...")

	if err := imageproc.Run(p.cfg); err != nil {
		return err
	}

	log.Println("Image processing complete")
	return nil
}

func (p *Pipeline) GenerateHugoContent() error {
	log.Println("Generating Hugo content...")

	store := metadata.NewStore()

	if err := p.loadMetadataIntoStore(store); err != nil {
		return err
	}

	generator := hugo.NewGenerator(
		p.cfg,
		p.metadataDir,
		p.cfg.Paths.PhotosDir,
		p.contentDir,
	)

	if err := generator.Run(store); err != nil {
		return err
	}

	log.Println("Hugo content generation complete")
	return nil
}

func (p *Pipeline) loadMetadataIntoStore(store *metadata.Store) error {
	images, err := findImages(p.cfg.Paths.PhotosDir)
	if err != nil {
		return err
	}

	for _, imgPath := range images {
		relPath, err := filepath.Rel(p.cfg.Paths.PhotosDir, imgPath)
		if err != nil {
			continue
		}

		metadataPath := filepath.Join(p.metadataDir, relPath+".json")
		data, err := os.ReadFile(metadataPath)
		if err != nil {
			continue
		}

		var meta metadata.PhotoMetadata
		if err := json.Unmarshal(data, &meta); err != nil {
			continue
		}

		store.Set(relPath, &meta)
	}

	return nil
}

func findImages(dir string) ([]string, error) {
	var images []string

	extensions := []string{".jpg", ".jpeg", ".png", ".webp", ".tiff", ".tif", ".gif", ".bmp"}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		for _, e := range extensions {
			if ext == e {
				images = append(images, path)
				break
			}
		}
		return nil
	})

	return images, err
}
