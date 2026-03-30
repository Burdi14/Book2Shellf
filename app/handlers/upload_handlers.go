package handlers

import (
	"image"
	"image/draw"
	"image/png"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CropCover lets admin choose a portion of the first-page cover image
func CropCover(c *gin.Context) {
	var req struct {
		CoverURL string  `json:"cover_url"`
		X        float64 `json:"x"`
		Y        float64 `json:"y"`
		Width    float64 `json:"width"`
		Height   float64 `json:"height"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Success: false, Message: "Invalid crop payload"})
		return
	}

	if req.CoverURL == "" {
		c.JSON(http.StatusBadRequest, APIResponse{Success: false, Message: "cover_url required"})
		return
	}

	x := math.Max(0, math.Min(1, req.X))
	y := math.Max(0, math.Min(1, req.Y))
	w := math.Max(0, math.Min(1, req.Width))
	h := math.Max(0, math.Min(1, req.Height))

	if w == 0 || h == 0 {
		c.JSON(http.StatusBadRequest, APIResponse{Success: false, Message: "Crop size must be greater than zero"})
		return
	}
	if x+w > 1 {
		w = 1 - x
	}
	if y+h > 1 {
		h = 1 - y
	}

	if w <= 0 || h <= 0 {
		c.JSON(http.StatusBadRequest, APIResponse{Success: false, Message: "Crop area invalid"})
		return
	}

	coverFile := strings.TrimPrefix(req.CoverURL, "/uploads/covers/")
	if strings.Contains(coverFile, "..") || coverFile == "" {
		c.JSON(http.StatusBadRequest, APIResponse{Success: false, Message: "Invalid cover path"})
		return
	}
	coverPath := filepath.Join("./uploads/covers", coverFile)

	f, err := os.Open(coverPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Success: false, Message: "Cover not found"})
		return
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Success: false, Message: "Cannot decode cover"})
		return
	}

	bounds := img.Bounds()
	rect := image.Rect(
		int(x*float64(bounds.Dx())),
		int(y*float64(bounds.Dy())),
		int((x+w)*float64(bounds.Dx())),
		int((y+h)*float64(bounds.Dy())),
	)
	rect = rect.Intersect(bounds)
	if rect.Dx() < 10 || rect.Dy() < 10 {
		c.JSON(http.StatusBadRequest, APIResponse{Success: false, Message: "Crop too small"})
		return
	}

	cropped := image.NewRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
	draw.Draw(cropped, cropped.Bounds(), img, rect.Min, draw.Src)

	if err := os.MkdirAll("./uploads/covers", 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Success: false, Message: "Failed to prepare cover folder"})
		return
	}

	newName := uuid.New().String() + "-crop.png"
	outPath := filepath.Join("./uploads/covers", newName)
	out, err := os.Create(outPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Success: false, Message: "Failed to save cropped cover"})
		return
	}
	defer out.Close()

	if err := png.Encode(out, cropped); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Success: false, Message: "Failed to encode cropped cover"})
		return
	}

	newURL := "/uploads/covers/" + newName

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Cover cropped",
		Data: map[string]interface{}{
			"cover_url": newURL,
		},
	})
}

// UploadBook handles book file upload (PDF, EPUB, DJVU, etc.)
func UploadBook(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "No file uploaded",
		})
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	filename := uuid.New().String() + ext
	filePath := filepath.Join("./uploads/books", filename)

	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to prepare upload folder",
		})
		return
	}

	out, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to save file",
		})
		return
	}
	defer out.Close()

	written, err := io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to save file",
		})
		return
	}

	var coverURL string
	if strings.ToLower(ext) == ".pdf" {
		coverURL, _ = generateCoverFromPDF(filePath)
	}

	if coverURL == "" {
		fallbackName := strings.TrimSuffix(filename, ext)
		placeholderURL, phErr := generatePlaceholderCover("./uploads/covers", fallbackName)
		if phErr == nil {
			coverURL = placeholderURL
		} else {
			log.Printf("placeholder cover failed for %s: %v", header.Filename, phErr)
		}
	}

	bookObj := map[string]interface{}{
		"url":           "/uploads/books/" + filename,
		"original_name": header.Filename,
		"size":          written,
	}
	coverObj := map[string]interface{}{
		"url":       coverURL,
		"generated": true,
	}

	responseData := map[string]interface{}{
		"url":           bookObj["url"],
		"original_name": header.Filename,
		"size":          written,
		"cover_url":     coverURL,
		"book":          bookObj,
		"cover":         coverObj,
		"lib_unit": map[string]interface{}{
			"book":  bookObj,
			"cover": coverObj,
		},
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "File uploaded successfully",
		Data:    responseData,
	})
}

// UploadCover handles cover image upload
func UploadCover(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "No file uploaded",
		})
		return
	}
	defer file.Close()

	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	valid := false
	for _, allowed := range allowedExts {
		if ext == allowed {
			valid = true
			break
		}
	}
	if !valid {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Only image files (jpg, png, gif, webp) are allowed",
		})
		return
	}

	filename := uuid.New().String() + ext
	filePath := filepath.Join("./uploads/covers", filename)

	out, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to save file",
		})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to save file",
		})
		return
	}

	coverObj := map[string]interface{}{
		"url":       "/uploads/covers/" + filename,
		"generated": false,
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Cover uploaded successfully",
		Data: map[string]interface{}{
			"url":       coverObj["url"],
			"cover_url": coverObj["url"],
			"cover":     coverObj,
			"lib_unit": map[string]interface{}{
				"cover": coverObj,
			},
		},
	})
}
