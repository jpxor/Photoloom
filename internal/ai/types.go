package ai

type Metadata struct {
	SourcePath   string   `json:"sourcePath"`
	Category     string   `json:"category"`
	Tags         []string `json:"tags"`
	Colors       []string `json:"colors"`
	Description  string   `json:"description"`
	DateTaken    string   `json:"dateTaken"`
	CameraMake   string   `json:"cameraMake"`
	CameraModel  string   `json:"cameraModel"`
	Lens         string   `json:"lens"`
	ISO          int      `json:"iso"`
	Aperture     string   `json:"aperture"`
	ShutterSpeed string   `json:"shutterSpeed"`
	FocalLength  string   `json:"focalLength"`
	Width        int      `json:"width"`
	Height       int      `json:"height"`
}

func NewMetadata() *Metadata {
	return &Metadata{
		Tags:   []string{},
		Colors: []string{},
	}
}
