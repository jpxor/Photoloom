package hugo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"photo-gallery/internal/colors"
	"photo-gallery/internal/config"
	"photo-gallery/internal/metadata"
)

type Generator struct {
	cfg         *config.Config
	metadataDir string
	photosDir   string
	contentDir  string
}

func NewGenerator(cfg *config.Config, metadataDir, photosDir, contentDir string) *Generator {
	return &Generator{
		cfg:         cfg,
		metadataDir: metadataDir,
		photosDir:   photosDir,
		contentDir:  contentDir,
	}
}

func (g *Generator) Run(store *metadata.Store) error {
	log.Println("Generating Hugo content...")

	if err := os.MkdirAll(g.contentDir, 0755); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Join(g.contentDir, "albums"), 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(g.contentDir, "photos"), 0755); err != nil {
		return err
	}

	if err := g.generateIndex(); err != nil {
		return err
	}

	albums, err := g.discoverAlbums(store)
	if err != nil {
		return err
	}

	photoCount := 0
	for albumName, photos := range albums {
		albumSlug := albumName

		coverImage := ""
		for _, photo := range photos {
			if coverImage == "" {
				coverImage = strings.ReplaceAll("photos/"+photo, "\\", "/")
			}

			if err := g.createPhotoContent(albumSlug, photo); err != nil {
				log.Printf("Warning: failed to create content for %s: %v", photo, err)
				continue
			}
			photoCount++
		}

		if err := g.createAlbumContent(albumName, albumSlug, coverImage, len(photos)); err != nil {
			log.Printf("Warning: failed to create album %s: %v", albumName, err)
		}
	}

	log.Printf("Generated %d photos across %d albums", photoCount, len(albums))
	return nil
}

func (g *Generator) generateIndex() error {
	content := `---
title: "Home"
type: page
---
`
	return os.WriteFile(filepath.Join(g.contentDir, "_index.md"), []byte(content), 0644)
}

func (g *Generator) discoverAlbums(store *metadata.Store) (map[string][]string, error) {
	albums := make(map[string][]string)

	allMetadata := store.GetAll()
	for relPath := range allMetadata {
		ext := strings.ToLower(filepath.Ext(relPath))
		if !isImage(ext) {
			continue
		}

		albumName := filepath.Dir(relPath)
		if albumName == "." {
			albumName = "uncategorized"
		}
		firstComponent := albumName
		if idx := strings.Index(firstComponent, "/"); idx != -1 {
			firstComponent = firstComponent[:idx]
		}
		albums[firstComponent] = append(albums[firstComponent], relPath)
	}

	return albums, nil
}

func (g *Generator) createPhotoContent(albumSlug, photoRel string) error {
	title := strings.TrimSuffix(filepath.Base(photoRel), filepath.Ext(photoRel))

	imagePath := strings.ReplaceAll("photos/"+photoRel, "\\", "/")

	relPath := photoRel
	store := g.loadMetadataForPhoto(relPath)

	meta, _ := store.Get(relPath)

	categories := []string{}
	tags := []string{}
	hexColors := []string{}
	description := ""
	dateTaken := ""
	cameraMake := ""
	cameraModel := ""
	lens := ""
	iso := 0
	aperture := ""
	shutterSpeed := ""
	focalLength := ""
	width := 0
	height := 0

	if meta != nil {
		categories = parseCategories(meta.Category)
		seen := make(map[string]bool)
		var uniqueCategories []string
		for _, c := range categories {
			if !seen[c] {
				seen[c] = true
				uniqueCategories = append(uniqueCategories, c)
			}
		}
		categories = uniqueCategories
		tags = meta.Tags
		hasBirdTag := false
		for _, t := range tags {
			if t == "bird" || t == "birds" {
				hasBirdTag = true
				break
			}
		}
		if hasBirdTag {
			hasBirdsCategory := false
			for _, c := range categories {
				if c == "birds" {
					hasBirdsCategory = true
					break
				}
			}
			if !hasBirdsCategory {
				categories = append(categories, "birds")
			}
		}
		hexColors = meta.Colors
		description = meta.Description
		dateTaken = meta.DateTaken
		cameraMake = meta.CameraMake
		cameraModel = meta.CameraModel
		lens = meta.Lens
		iso = meta.ISO
		aperture = meta.Aperture
		shutterSpeed = meta.ShutterSpeed
		focalLength = meta.FocalLength
		width = meta.Width
		height = meta.Height
	}

	photoDate := dateTaken
	if photoDate != "" {
		if t, err := time.Parse("2006:01:02 15:04:05", dateTaken); err == nil {
			photoDate = t.Format("2006-01-02T15:04:05Z07:00")
		}
	}
	if photoDate == "" {
		photoPath := filepath.Join(g.photosDir, photoRel)
		if info, err := os.Stat(photoPath); err == nil {
			photoDate = info.ModTime().Format("2006-01-02T15:04:05Z07:00")
		}
	}

	var frontmatter bytes.Buffer
	frontmatter.WriteString("---\n")
	frontmatter.WriteString(fmt.Sprintf("title: %q\n", title))
	frontmatter.WriteString("type: photo\n")
	frontmatter.WriteString(fmt.Sprintf("album: %q\n", albumSlug))
	frontmatter.WriteString(fmt.Sprintf("image: %q\n", imagePath))
	if photoDate != "" {
		frontmatter.WriteString(fmt.Sprintf("date: %q\n", photoDate))
	}
	if cameraMake != "" || cameraModel != "" {
		frontmatter.WriteString(fmt.Sprintf("camera_make: %q\n", cameraMake))
		frontmatter.WriteString(fmt.Sprintf("camera_model: %q\n", cameraModel))
	}
	if lens != "" {
		frontmatter.WriteString(fmt.Sprintf("lens: %q\n", lens))
	}
	if iso > 0 {
		frontmatter.WriteString(fmt.Sprintf("iso: %d\n", iso))
	}
	if aperture != "" {
		frontmatter.WriteString(fmt.Sprintf("aperture: %q\n", aperture))
	}
	if shutterSpeed != "" {
		frontmatter.WriteString(fmt.Sprintf("shutter_speed: %q\n", shutterSpeed))
	}
	if focalLength != "" {
		frontmatter.WriteString(fmt.Sprintf("focal_length: %q\n", focalLength))
	}
	if width > 0 && height > 0 {
		frontmatter.WriteString(fmt.Sprintf("width: %d\n", width))
		frontmatter.WriteString(fmt.Sprintf("height: %d\n", height))
	}
	if len(categories) > 0 {
		frontmatter.WriteString(fmt.Sprintf("categories: [%s]\n", strings.Join(quoteStrings(categories), ", ")))
	}
	if len(tags) > 0 {
		frontmatter.WriteString(fmt.Sprintf("tags: [%s]\n", strings.Join(quoteStrings(tags), ", ")))
	}
	if len(hexColors) > 0 {
		frontmatter.WriteString(fmt.Sprintf("specificColors: [%s]\n", strings.Join(quoteStrings(hexColors), ", ")))

		colorGroups := colors.MapColorsToGroups(hexColors)
		if len(colorGroups) > 0 {
			frontmatter.WriteString(fmt.Sprintf("colors: [%s]\n", strings.Join(quoteStrings(colorGroups), ", ")))

			colorGroupColors := colors.MapColorGroupToColors(hexColors, colorGroups)
			if len(colorGroupColors) > 0 {
				frontmatter.WriteString("colorGroupColors:\n")
				for group, color := range colorGroupColors {
					frontmatter.WriteString(fmt.Sprintf("  %s: %q\n", group, color))
				}
			}
		}
	}
	if description != "" {
		frontmatter.WriteString(fmt.Sprintf("description: %q\n", description))
	}
	frontmatter.WriteString("---\n")

	slug := strings.TrimSuffix(filepath.Base(photoRel), filepath.Ext(photoRel))
	outPath := filepath.Join(g.contentDir, "photos", albumSlug, slug+".md")

	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(outPath, frontmatter.Bytes(), 0644)
}

func (g *Generator) loadMetadataForPhoto(relPath string) *metadata.Store {
	store := metadata.NewStore()

	metadataPath := filepath.Join(g.metadataDir, relPath+".json")
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return store
	}

	var meta metadata.PhotoMetadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return store
	}

	store.Set(relPath, &meta)
	return store
}

func (g *Generator) createAlbumContent(albumName, albumSlug, cover string, photoCount int) error {
	title := albumName
	if title == "uncategorized" {
		title = "Uncategorized"
	}

	var frontmatter bytes.Buffer
	frontmatter.WriteString("---\n")
	frontmatter.WriteString(fmt.Sprintf("title: %q\n", title))
	frontmatter.WriteString("type: album\n")
	if cover != "" {
		frontmatter.WriteString(fmt.Sprintf("cover: %q\n", cover))
	}
	frontmatter.WriteString(fmt.Sprintf("description: %d photos\n", photoCount))
	frontmatter.WriteString("---\n")

	outPath := filepath.Join(g.contentDir, "albums", albumSlug+".md")
	return os.WriteFile(outPath, frontmatter.Bytes(), 0644)
}

func isImage(ext string) bool {
	images := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true,
		".webp": true, ".tiff": true, ".tif": true,
		".gif": true, ".bmp": true,
	}
	return images[ext]
}

func parseCategories(categoryStr string) []string {
	if categoryStr == "" {
		return nil
	}

	categoryStr = strings.ReplaceAll(categoryStr, "/", ",")
	categoryStr = strings.ReplaceAll(categoryStr, "|", ",")

	parts := strings.Split(categoryStr, ",")
	var categories []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		p = strings.ToLower(p)
		if p != "" {
			categories = append(categories, p)
		}
	}
	return categories
}

func quoteStrings(ss []string) []string {
	result := make([]string, len(ss))
	for i, s := range ss {
		result[i] = fmt.Sprintf("%q", s)
	}
	return result
}
