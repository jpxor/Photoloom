# Photoloom

Weave your photos into stunning galleries with AI-powered automation.

## What is Photoloom?

Photoloom is a free, open-source tool that transforms your messy photo folders into a beautiful, searchable static website. Just point it at your photos—it uses AI to auto-tag, describe, and add copyright watermarks, then generates a fast static gallery you can host anywhere for free.

**No subscriptions. No privacy concerns. No hosting fees. Just your photos—beautifully presented.**

## Why Photoloom?

- **AI-first preprocessing** - Automatically tags, categorizes, and describes your photos using AI
- **Truly free hosting** - Generates static sites that can be hosted on Netlify, GitHub Pages, Cloudflare Pages, or any static host
- **Privacy-first** - Runs locally on your machine. Your photos never leave your computer until you deploy
- **Incremental builds** - Only reprocess new photos, saving time and AI API tokens
- **Copyright protection** - Built-in watermarking keeps your work protected
- **Hugo-powered** - Leverages the proven Hugo static site generator ecosystem

## Quick Start

```bash
# 1. Clone and configure
git clone https://github.com/yourusername/photoloom.git
cd photoloom
cp .env.example .env

# 2. Add your photos
cp your_photos/* ./photos/

# 3. Add your API key (or use 'none' for local-only)
# Edit .env and set OPENAI_API_KEY=your-key-here

# 4. Build
go run ./cmd/build

# 5. Preview locally
cd hugo && hugo server
```

## Features

### AI-Powered Metadata Extraction

Photoloom analyzes every photo using AI to extract:

- **Categories** - Portrait, Landscape, Wildlife, Nature, Architecture, Events, Travel, Family, and more
- **Tags** - 2-5 relevant keywords per photo
- **Descriptions** - AI-generated 1-2 sentence descriptions
- **Colors** - Dominant colors in hex format

Supported AI providers:
- **OpenAI** - Uses GPT-4o vision (configurable model)
- **Ollama** - Local AI using llava or other vision models
- **None** - Skip AI and use EXIF data only

### EXIF Data Extraction

Automatically extracts camera metadata:

- Date taken
- Camera make & model
- Lens information
- ISO, aperture, shutter speed
- Focal length
- Image dimensions

### Image Processing

- **Web optimization** - Resizes photos for optimal web display (default: 1600px)
- **Thumbnail generation** - Creates thumbnails for gallery views (default: 400px)
- **Watermarking** - Adds customizable copyright watermark to large photos
- **Format conversion** - Converts to optimized JPEG format

### Album Organization

- **Folder-based albums** - Automatically creates albums from folder structure
- **AI-category grouping** - Optional grouping by AI-detected category
- **Date-based sorting** - Photos sorted by EXIF date or file modification time
- **Color grouping** - Groups photos by dominant color palette

### Hugo Integration

- Generates clean Hugo content files with full metadata
- Built-in gallery theme with lightbox viewer
- Taxonomy support for categories, tags, and colors
- Responsive design for mobile and desktop

## Installation

### Prerequisites

- **Go 1.21+** - For building the CLI tool
- **Hugo 0.110+** - For generating the static site
- **OpenAI API key** or **Ollama** (optional) - For AI metadata extraction

### Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/photoloom.git
   cd photoloom
   ```

2. Copy the environment file:
   ```bash
   cp .env.example .env
   ```

3. Add your photos to the `photos` directory (or configure a different path in `config.yaml`)

4. Configure your AI provider (optional):
   - For OpenAI: Add your API key to `.env`
   - For Ollama: Ensure Ollama is running locally
   - For no AI: Set `ai.provider: "none"` in config.yaml

5. Build:
   ```bash
   go run ./cmd/build
   ```

6. Preview:
   ```bash
   cd hugo && hugo server
   ```

## Configuration

All configuration is in `config.yaml`:

```yaml
paths:
  photosDir: "./photos"      # Your source photos
  cacheDir: "./cache"        # Processed metadata cache
  outputDir: "./output"      # Final static site

images:
  thumbnail:
    width: 400               # Thumbnail size
    height: 400
    quality: 80
  web:
    width: 1600              # Web image size (0 = auto)
    height: 0
    quality: 85
  watermarkThreshold: 1080   # Only watermark images > this width
  watermark:
    text: "© Your Name"     # Watermark text
    opacity: 0.7            # 0.0 - 1.0
    fontPath: "fonts/AlexBrush-Regular.ttf"
    position: "bottom-right" # bottom-left, top-right, top-left

ai:
  provider: "openai"        # "openai", "ollama", or "none"
  openai:
    apiKey: "$OPENAI_API_KEY"
    model: "gpt-4o-mini"
    concurrency: 4
  ollama:
    baseURL: "http://localhost:11434"
    model: "llava"

albums:
  groupBy: ["folder"]        # How to group photos into albums
  minPhotos: 1               # Minimum photos per album

hugo:
  theme: "gallery"
  flags:
    - "--minify"
    - "--gc"

build:
  workers: 4                 # Parallel processing workers
```

### CLI Flags

```bash
go run ./cmd/build              # Full build
go run ./cmd/build -skip-ai     # Skip AI processing (use cached)
go run ./cmd/build -skip-images # Skip image processing
go run ./cmd/build -skip-hugo   # Skip Hugo build
go run ./cmd/build -clean       # Clear cache before building
```

## Project Structure

```
photoloom/
├── cmd/build/                  # Main CLI entry point
├── internal/
│   ├── ai/                    # AI provider (OpenAI, Ollama)
│   ├── config/                # Configuration loading
│   ├── exif/                  # EXIF metadata extraction
│   ├── hugo/                  # Hugo content generator
│   ├── imageproc/             # Image processing & watermarking
│   ├── metadata/              # Metadata storage
│   └── colors/                # Color grouping logic
├── hugo/
│   ├── config.yaml            # Hugo configuration
│   └── themes/gallery/        # Custom gallery theme
├── fonts/                     # Watermark fonts
├── config.yaml                # Photoloom configuration
├── .env                       # API keys (not committed)
└── .env.example               # Environment template
```

## Hosting

Photoloom generates a completely static site. Deploy anywhere for free:

### GitHub Pages

```bash
cd hugo
hugo --theme=gallery --destination=../docs
# Enable GitHub Pages in repo settings, point to /docs
```

### Netlify

1. Connect your repo to Netlify
2. Build command: `go run ./cmd/build`
3. Publish directory: `output`

### Cloudflare Pages

1. Connect your repo to Cloudflare
2. Build command: `go run ./cmd/build`
3. Publish directory: `output`

### VPS / Any Web Server

Simply upload the `output/` directory to any web server.

## Customization

### Changing the Theme

Edit `hugo/config.yaml` to use a different Hugo theme, or create your own in `hugo/themes/`.

### Custom Watermark

1. Add your font file to the `fonts/` directory
2. Update `config.yaml` with your font path and watermark text
3. Rebuild

### Custom AI Prompt

Modify the `ai.prompt` section in `config.yaml` to customize how AI analyzes your photos.

## Contributing

Contributions welcome! Please open an issue or submit a PR.

## License

MIT License - see LICENSE for details.

## Acknowledgments

- Built with [Hugo](https://gohugo.io/)
- Image processing with [imaging](https://github.com/disintegration/imaging)
- AI vision powered by [OpenAI](https://openai.com/) and [Ollama](https://ollama.com/)
