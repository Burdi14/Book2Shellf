package handlers

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
)

// generateCoverFromPDF renders the first PDF page to a PNG using pdftoppm (poppler)
func generateCoverFromPDF(pdfPath string) (string, error) {
	baseName := strings.TrimSuffix(filepath.Base(pdfPath), filepath.Ext(pdfPath))
	outputDir := "./uploads/covers"
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return "", fmt.Errorf("create cover dir: %w", err)
	}

	outputBase := filepath.Join(outputDir, baseName)
	cmd := exec.Command("pdftoppm", "-f", "1", "-l", "1", "-png", pdfPath, outputBase)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.Printf("pdftoppm failed, falling back to placeholder: %v: %s", err, strings.TrimSpace(stderr.String()))
		return generatePlaceholderCover(outputDir, baseName)
	}

	coverFile := outputBase + "-1.png"
	if _, err := os.Stat(coverFile); err != nil {
		altPadded := outputBase + "-01.png"
		if _, errPad := os.Stat(altPadded); errPad == nil {
			coverFile = altPadded
		} else {
			matches, _ := filepath.Glob(outputBase + "-*.png")
			if len(matches) > 0 {
				coverFile = matches[0]
			} else {
				log.Printf("pdftoppm output missing, using placeholder: %v", err)
				return generatePlaceholderCover(outputDir, baseName)
			}
		}
	}

	return "/uploads/covers/" + filepath.Base(coverFile), nil
}

func generatePlaceholderCover(outputDir, baseName string) (string, error) {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return "", fmt.Errorf("create placeholder dir: %w", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, 640, 900))

	bg := color.RGBA{R: 18, G: 18, B: 18, A: 255}
	accent := color.RGBA{R: 255, G: 176, B: 0, A: 255}
	stripe := color.RGBA{R: 215, G: 38, B: 56, A: 255}

	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.Set(x, y, bg)
		}
	}

	for x := 0; x < img.Bounds().Dx(); x++ {
		img.Set(x, 0, accent)
		img.Set(x, img.Bounds().Dy()-1, accent)
	}
	for y := 0; y < img.Bounds().Dy(); y++ {
		img.Set(0, y, accent)
		img.Set(img.Bounds().Dx()-1, y, accent)
	}

	for y := 0; y < img.Bounds().Dy(); y += 24 {
		for x := 0; x < img.Bounds().Dx(); x++ {
			if (x+y)%48 < 12 {
				img.Set(x, y, stripe)
			}
		}
	}

	name := baseName + "-placeholder.png"
	path := filepath.Join(outputDir, name)
	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("create placeholder file: %w", err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return "", fmt.Errorf("encode placeholder: %w", err)
	}

	return "/uploads/covers/" + name, nil
}
