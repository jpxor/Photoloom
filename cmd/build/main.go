package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"photo-gallery/internal/config"
	"photo-gallery/internal/pipeline"
)

var (
	configPath = flag.String("config", "config.yaml", "Path to config file")
	skipImages = flag.Bool("skip-images", false, "Skip image processing")
	skipAI     = flag.Bool("skip-ai", false, "Skip AI metadata extraction")
	skipHugo   = flag.Bool("skip-hugo", false, "Skip Hugo build")
	clean      = flag.Bool("clean", false, "Clean cache before build")
)

func main() {
	flag.Parse()

	log.Println("Starting photo gallery build...")

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if *clean {
		log.Println("Cleaning cache directory...")
		if err := os.RemoveAll(cfg.Paths.CacheDir); err != nil {
			log.Printf("Warning: failed to clean cache: %v", err)
		}
	}

	start := time.Now()

	p := pipeline.New(cfg)

	if !*skipAI {
		if err := p.RunAIMetadata(); err != nil {
			log.Printf("Warning: AI metadata extraction failed: %v", err)
		}
	}

	if !*skipImages {
		if err := p.RunImageProcessing(); err != nil {
			log.Fatalf("Image processing failed: %v", err)
		}
	}

	if err := p.GenerateHugoContent(); err != nil {
		log.Fatalf("Failed to generate Hugo content: %v", err)
	}

	if !*skipHugo {
		if err := runHugoBuild(cfg); err != nil {
			log.Fatalf("Hugo build failed: %v", err)
		}
	}

	log.Printf("Build completed in %v", time.Since(start))
}

func runHugoBuild(cfg *config.Config) error {
	log.Println("Running Hugo build...")

	hugoCmd := "hugo"
	hugoPath := "./hugo.exe"
	if _, err := os.Stat(hugoPath); err == nil {
		hugoCmd, _ = filepath.Abs(hugoPath)
	}

	outputDir := "output"

	cmd := exec.Command(hugoCmd,
		"--source", "hugo",
		"--destination", outputDir,
		"--theme", "gallery",
		"--cleanDestinationDir",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("hugo failed: %w", err)
	}

	log.Printf("Hugo build complete. Output: %s", outputDir)
	return nil
}
