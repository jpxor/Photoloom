package exif

import (
	"fmt"
	"os"

	"github.com/rwcarlsen/goexif/exif"
)

type ImageMetadata struct {
	DateTaken    string `json:"dateTaken"`
	CameraMake   string `json:"cameraMake"`
	CameraModel  string `json:"cameraModel"`
	Lens         string `json:"lens"`
	ISO          int    `json:"iso"`
	Aperture     string `json:"aperture"`
	ShutterSpeed string `json:"shutterSpeed"`
	FocalLength  string `json:"focalLength"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}

func Extract(imgPath string) (*ImageMetadata, error) {
	metadata := &ImageMetadata{}

	f, err := os.Open(imgPath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decoding exif: %w", err)
	}

	if dateStr, err := x.Get(exif.DateTimeOriginal); err == nil {
		if dateVal, err := dateStr.StringVal(); err == nil {
			metadata.DateTaken = dateVal
		}
	} else if dateStr, err := x.Get(exif.DateTimeDigitized); err == nil {
		if dateVal, err := dateStr.StringVal(); err == nil {
			metadata.DateTaken = dateVal
		}
	}

	if make, err := x.Get(exif.Make); err == nil {
		if makeVal, err := make.StringVal(); err == nil {
			metadata.CameraMake = makeVal
		}
	}

	if model, err := x.Get(exif.Model); err == nil {
		if modelVal, err := model.StringVal(); err == nil {
			metadata.CameraModel = modelVal
		}
	}

	if lens, err := x.Get(exif.LensModel); err == nil {
		if lensVal, err := lens.StringVal(); err == nil {
			metadata.Lens = lensVal
		}
	}

	if iso, err := x.Get(exif.ISOSpeedRatings); err == nil {
		if isoVal, err := iso.Int(0); err == nil {
			metadata.ISO = int(isoVal)
		}
	}

	if aperture, err := x.Get(exif.FNumber); err == nil {
		if num, den, err := aperture.Rat2(0); err == nil {
			if den == 1 {
				metadata.Aperture = fmt.Sprintf("%.1f", float64(num))
			} else {
				metadata.Aperture = fmt.Sprintf("%.1f", float64(num)/float64(den))
			}
		}
	}

	if shutter, err := x.Get(exif.ExposureTime); err == nil {
		if num, den, err := shutter.Rat2(0); err == nil {
			if num == 1 {
				metadata.ShutterSpeed = fmt.Sprintf("1/%d", den)
			} else if den == 1 {
				metadata.ShutterSpeed = fmt.Sprintf("%d", num)
			} else {
				metadata.ShutterSpeed = fmt.Sprintf("1/%d", den/num)
			}
		}
	}

	if focal, err := x.Get(exif.FocalLength); err == nil {
		if num, den, err := focal.Rat2(0); err == nil {
			if den == 1 {
				metadata.FocalLength = fmt.Sprintf("%d", num)
			} else {
				metadata.FocalLength = fmt.Sprintf("%.0f", float64(num)/float64(den))
			}
		}
	}

	if width, err := x.Get(exif.PixelXDimension); err == nil {
		if w, err := width.Int(0); err == nil {
			metadata.Width = int(w)
		}
	}

	if height, err := x.Get(exif.PixelYDimension); err == nil {
		if h, err := height.Int(0); err == nil {
			metadata.Height = int(h)
		}
	}

	return metadata, nil
}
