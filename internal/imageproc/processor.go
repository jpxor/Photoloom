package imageproc

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"photo-gallery/internal/config"

	"github.com/disintegration/imaging"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func Run(cfg *config.Config) error {
	photosDir := cfg.Paths.PhotosDir
	assetsDir := "hugo/assets/photos"

	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		return fmt.Errorf("creating assets dir: %w", err)
	}

	images, err := findImages(photosDir)
	if err != nil {
		return fmt.Errorf("finding images: %w", err)
	}

	fmt.Printf("Found %d images to process\n", len(images))

	processed := 0
	skipped := 0
	errors := 0

	sem := make(chan struct{}, cfg.Build.Workers)

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, img := range images {
		wg.Add(1)
		go func(img string) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			relPath, err := filepath.Rel(photosDir, img)
			if err != nil {
				mu.Lock()
				errors++
				mu.Unlock()
				return
			}

			destPath := filepath.Join(assetsDir, relPath)

			dir := filepath.Dir(destPath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				mu.Lock()
				errors++
				mu.Unlock()
				return
			}

			if shouldSkip(img, destPath) {
				mu.Lock()
				skipped++
				mu.Unlock()
				return
			}

			if err := ProcessImage(img, destPath, cfg); err != nil {
				mu.Lock()
				errors++
				mu.Unlock()
				return
			}

			mu.Lock()
			processed++
			mu.Unlock()
		}(img)
	}

	wg.Wait()

	fmt.Printf("Processed: %d, Skipped: %d, Errors: %d\n", processed, skipped, errors)

	return nil
}

func shouldSkip(srcPath, destPath string) bool {
	destInfo, err := os.Stat(destPath)
	if err != nil {
		return false
	}

	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return false
	}

	return !destInfo.ModTime().Before(srcInfo.ModTime())
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

func ProcessImage(srcPath, destPath string, cfg *config.Config) error {
	img, err := imaging.Open(srcPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", srcPath, err)
	}

	bounds := img.Bounds()
	origSize := bounds.Dx()
	if bounds.Dy() > origSize {
		origSize = bounds.Dy()
	}

	if origSize >= cfg.Images.WatermarkThreshold {
		img, err = addWatermark(img, cfg.Images.Watermark)
		if err != nil {
			return fmt.Errorf("adding watermark: %w", err)
		}
	}

	if err := saveImage(destPath, img, cfg.Images.Web.Quality); err != nil {
		return fmt.Errorf("saving image: %w", err)
	}

	return nil
}

func saveImage(path string, img image.Image, quality int) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()

	ext := filepath.Ext(path)
	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Encode(outFile, img, &jpeg.Options{Quality: quality})
	case ".png":
		return png.Encode(outFile, img)
	default:
		return jpeg.Encode(outFile, img, &jpeg.Options{Quality: quality})
	}
}

func addWatermark(img image.Image, wm config.WatermarkConfig) (image.Image, error) {
	fontBytes, err := ioutil.ReadFile(wm.FontPath)
	if err != nil {
		return nil, fmt.Errorf("reading font file %s: %w", wm.FontPath, err)
	}

	f, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("parsing font: %w", err)
	}

	bounds := img.Bounds()
	margin := bounds.Dx() / 50
	if margin < 10 {
		margin = 10
	}

	watermarkHeight := bounds.Dy() / 30
	if watermarkHeight < 16 {
		watermarkHeight = 16
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    float64(watermarkHeight),
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		return nil, fmt.Errorf("creating font face: %w", err)
	}

	watermarked := image.NewRGBA(bounds)
	draw.Draw(watermarked, bounds, img, image.Point{}, draw.Src)

	textWidth := font.MeasureString(face, wm.Text).Ceil()

	x := bounds.Dx() - textWidth - margin
	y := bounds.Dy() - watermarkHeight - margin

	switch wm.Position {
	case "bottom-left":
		x = margin
	case "top-right":
		y = margin
	case "top-left":
		x = margin
		y = margin
	}

	opacity := uint8(wm.Opacity * 255)
	if opacity < 50 {
		opacity = 128
	}

	outlineOffsets := [][2]int{
		{-1, -1}, {0, -1}, {1, -1},
		{-1, 0}, {1, 0},
		{-1, 1}, {0, 1}, {1, 1},
	}

	black := &font.Drawer{
		Dst:  watermarked,
		Src:  &image.Uniform{color.RGBA{0, 0, 0, 255}},
		Face: face,
		Dot:  fixed.Point26_6{},
	}

	for _, offset := range outlineOffsets {
		black.Dot = fixed.Point26_6{
			X: fixed.I(x + offset[0]),
			Y: fixed.I(y + offset[1]),
		}
		black.DrawString(wm.Text)
	}

	white := &font.Drawer{
		Dst:  watermarked,
		Src:  &image.Uniform{color.RGBA{255, 255, 255, opacity}},
		Face: face,
		Dot: fixed.Point26_6{
			X: fixed.I(x),
			Y: fixed.I(y),
		},
	}
	white.DrawString(wm.Text)

	return watermarked, nil
}

func GetImageDimensions(path string) (int, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	cfg, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}

	return cfg.Width, cfg.Height, nil
}
