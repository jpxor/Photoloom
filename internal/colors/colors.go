package colors

import (
	"fmt"
	"math"
	"strings"
)

func HexToHSL(hex string) (h, s, l float64) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 0, 0, 0
	}

	var r, g, b int
	_, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return 0, 0, 0
	}

	rF := float64(r) / 255.0
	gF := float64(g) / 255.0
	bF := float64(b) / 255.0

	max := math.Max(math.Max(rF, gF), bF)
	min := math.Min(math.Min(rF, gF), bF)
	l = (max + min) / 2.0

	if max == min {
		return 0, 0, l
	}

	d := max - min
	if l > 0.5 {
		s = d / (2.0 - max - min)
	} else {
		s = d / (max + min)
	}

	switch max {
	case rF:
		h = ((gF - bF) / d)
	case gF:
		h = ((bF - rF) / d) + 2.0
	case bF:
		h = ((rF - gF) / d) + 4.0
	}

	h = h * 60
	if h < 0 {
		h += 360
	}

	return h, s, l
}

func HSLToColorName(h, s, l float64) string {
	if l > 0.9 {
		return "white"
	}
	if l < 0.15 {
		return "black"
	}
	if s < 0.15 {
		return "gray"
	}

	switch {
	case h >= 345 || h < 15:
		return "red"
	case h >= 10 && h < 45:
		return "orange"
	case h >= 45 && h < 75:
		if s > 0.5 {
			return "yellow"
		}
		return "brown"
	case h >= 75 && h < 150:
		if s > 0.3 {
			return "green"
		}
		return "brown"
	case h >= 150 && h < 195:
		return "cyan"
	case h >= 195 && h < 255:
		return "blue"
	case h >= 255 && h < 285:
		return "purple"
	case h >= 285 && h < 345:
		if s > 0.5 && l > 0.5 {
			return "pink"
		}
		return "magenta"
	}

	return ""
}

func MapColorToName(hex string) string {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return ""
	}

	h, s, l := HexToHSL(hex)
	return HSLToColorName(h, s, l)
}

func MapColorsToGroups(hexColors []string) []string {
	colorGroups := make(map[string]bool)
	for _, hex := range hexColors {
		colorName := MapColorToName(hex)
		if colorName != "" {
			colorGroups[colorName] = true
		}
	}

	result := make([]string, 0, len(colorGroups))
	for name := range colorGroups {
		result = append(result, name)
	}
	return result
}

func MapColorGroupToColors(hexColors []string, colorGroups []string) map[string]string {
	result := make(map[string]string)
	colorGroupsSet := make(map[string]bool)
	for _, g := range colorGroups {
		colorGroupsSet[g] = true
	}
	for _, hex := range hexColors {
		group := MapColorToName(hex)
		if group != "" && colorGroupsSet[group] {
			if _, exists := result[group]; !exists {
				result[group] = hex
			}
		}
	}
	return result
}
