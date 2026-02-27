# Photo Gallery Static Site Builder - Specification

I need to build a static photo gallery website using Hugo. Please plan and implement a complete solution with the following requirements:

## Source & Output

- **Input**: Point to a photos directory (configurable path, e.g., `/path/to/photos` or `./photos`)
- **Output**: Static Hugo site in `./output/` or similar
- **Original photos**: Must NEVER be modified or moved - they stay in place

## Image Processing Pipeline

1. **Thumbnails**: Generate small preview versions (e.g., 400x400 max, smart crop for portraits)
2. **Web-optimized**: Create optimized full-size images (max 1600px, quality 85%, WebP + JPEG fallback)
   - **Watermark**: Add small copyright watermark in corner for full-size images only
   - Lower resolution/quality photos (below configurable threshold, e.g., < 2048px) skip watermark
3. Use **bimg/Imagor** (Go) or similar Go-native tool for image processing
4. Processed images go to a separate cache/build directory, not overwriting originals

## AI-Powered Metadata Extraction

For each photo, extract and store:
- **Content category**: (landscape, portrait, urban, nature, event, etc.)
- **Tags**: 5-10 relevant tags
- **Dominant colors**: Top 3-5 colors as hex codes
- **Description**: Brief AI-generated caption

Support **multiple AI providers**:
- OpenAI GPT-4o (with vision)
- Ollama (local) with vision-capable model (llava, etc.)

Make it configurable which provider to use

## Album Auto-Generation

- Automatically group photos into albums based on:
  - Folder structure in source (each subfolder = album)
  - Date taken (EXIF data)
  - AI-detected content categories
- Album pages show cover image, photo count, date range

## Hugo Theme Requirements

- Create a **custom Hugo theme** (not use existing theme)
- Clean, minimalist design with:
  - Homepage with album grid
  - Album page with photo grid
  - Photo detail page with full image + metadata
  - Lightbox for viewing
  - Responsive (mobile-friendly)
- Use Hugo templating and static assets

## Technical Requirements

- Configuration via `config.yaml` or `config.toml`
- Build script that orchestrates: image processing → AI metadata → Hugo build
- Support incremental builds (only process new/modified photos)
- No external database needed

## Deliverables

1. Project structure documentation
2. Image processing script(s)
3. AI metadata extraction script(s) with provider abstraction
4. Hugo custom theme files
5. Main build orchestration script
6. Configuration example
