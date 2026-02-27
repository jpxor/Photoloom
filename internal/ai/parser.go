package ai

import (
	"strings"
)

func ParseResponse(content string) (*Metadata, error) {
	metadata := NewMetadata()

	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "CATEGORY:") {
			metadata.Category = strings.TrimSpace(strings.TrimPrefix(line, "CATEGORY:"))
		} else if strings.HasPrefix(line, "TAGS:") {
			tagsStr := strings.TrimSpace(strings.TrimPrefix(line, "TAGS:"))
			metadata.Tags = strings.Split(tagsStr, ",")
			for i := range metadata.Tags {
				metadata.Tags[i] = strings.TrimSpace(metadata.Tags[i])
			}
		} else if strings.HasPrefix(line, "COLORS:") {
			colorsStr := strings.TrimSpace(strings.TrimPrefix(line, "COLORS:"))
			metadata.Colors = strings.Split(colorsStr, ",")
			for i := range metadata.Colors {
				metadata.Colors[i] = strings.TrimSpace(metadata.Colors[i])
			}
		} else if strings.HasPrefix(line, "DESCRIPTION:") {
			metadata.Description = strings.TrimSpace(strings.TrimPrefix(line, "DESCRIPTION:"))
		}
	}

	if metadata.Category == "" {
		metadata.Category = "other"
	}

	return metadata, nil
}
