package handlers

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"log"
	"math"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	_ "image/gif"
	_ "image/jpeg"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetBooks(c *gin.Context) {
	books, err := GetAllBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to fetch books",
		})
		return
	}

	if books == nil {
		books = []Book{}
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    books,
	})
}

// GetBook returns a single book (public - blocks hidden section books)
func GetBook(c *gin.Context) {
	id := c.Param("id")
	book, err := GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Book not found",
		})
		return
	}

	// Check if the book's section is hidden
	if book.SectionID != "" {
		if isHidden, _ := IsSectionHidden(book.SectionID); isHidden {
			c.JSON(http.StatusNotFound, APIResponse{
				Success: false,
				Message: "Book not found",
			})
			return
		}
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    book,
	})
}

// DownloadBook serves the book file and increments download count (public - blocks hidden)
func DownloadBook(c *gin.Context) {
	id := c.Param("id")
	book, err := GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Book not found",
		})
		return
	}

	// Block download of books in hidden sections via public route
	if book.SectionID != "" {
		if isHidden, _ := IsSectionHidden(book.SectionID); isHidden {
			c.JSON(http.StatusNotFound, APIResponse{
				Success: false,
				Message: "Book not found",
			})
			return
		}
	}

	// Increment download counter
	IncrementDownloads(id)

	// Serve the file
	filePath := "." + book.FileURL
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "File not found",
		})
		return
	}

	contentType := mime.TypeByExtension(filepath.Ext(book.FileName))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", book.FileName))
	c.Header("Content-Type", contentType)
	c.File(filePath)
}

func GetSections(c *gin.Context) {
	sections, err := GetAllSections()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to fetch sections",
		})
		return
	}

	if sections == nil {
		sections = []Section{}
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    sections,
	})
}

// GetBooksBySection returns books in a specific section (public - blocks hidden sections)
func GetBooksBySection(c *gin.Context) {
	sectionID := c.Param("id")

	// Block access to hidden sections
	if isHidden, _ := IsSectionHidden(sectionID); isHidden {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Section not found",
		})
		return
	}

	books, err := GetBooksBySectionDB(sectionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to fetch books",
		})
		return
	}

	if books == nil {
		books = []Book{}
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    books,
	})
}

// CreateBook creates a new book (admin only)
func CreateBook(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if book.Title == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Title is required",
		})
		return
	}

	if err := CreateBookDB(&book); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to create book",
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "Book created successfully",
		Data:    book,
	})
}

// UpdateBook updates a book (admin only)
func UpdateBook(c *gin.Context) {
	id := c.Param("id")

	// Check if book exists
	existingBook, err := GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Book not found",
		})
		return
	}

	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	book.ID = id
	// Preserve file info if not updated
	if book.FileURL == "" {
		book.FileURL = existingBook.FileURL
		book.FileName = existingBook.FileName
		book.FileSize = existingBook.FileSize
	}
	if book.CoverURL == "" {
		book.CoverURL = existingBook.CoverURL
	}

	if err := UpdateBookDB(&book); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to update book",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Book updated successfully",
		Data:    book,
	})
}

// DeleteBook deletes a book (admin only)
func DeleteBook(c *gin.Context) {
	id := c.Param("id")

	// Get book to delete associated files
	book, err := GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Book not found",
		})
		return
	}

	// Delete associated files
	if book.FileURL != "" {
		os.Remove("." + book.FileURL)
	}
	if book.CoverURL != "" {
		os.Remove("." + book.CoverURL)
	}

	if err := DeleteBookDB(id); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to delete book",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Book deleted successfully",
	})
}

// CreateSection creates a new section (admin only)
func CreateSection(c *gin.Context) {
	var section struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Hidden      bool   `json:"hidden"`
	}
	if err := c.ShouldBindJSON(&section); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if section.Name == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Name is required",
		})
		return
	}

	newSection, err := CreateSectionDB(section.Name, section.Description, section.Hidden)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to create section",
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "Section created successfully",
		Data:    newSection,
	})
}

// UpdateSection updates a section (admin only)
func UpdateSection(c *gin.Context) {
	id := c.Param("id")
	var section struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Hidden      bool   `json:"hidden"`
	}
	if err := c.ShouldBindJSON(&section); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if err := UpdateSectionDB(id, section.Name, section.Description, section.Hidden); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to update section",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Section updated successfully",
	})
}

// DeleteSection deletes a section (admin only)
func DeleteSection(c *gin.Context) {
	id := c.Param("id")

	if err := DeleteSectionDB(id); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to delete section",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Section deleted successfully",
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

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := uuid.New().String() + ext
	filePath := filepath.Join("./uploads/books", filename)

	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to prepare upload folder",
		})
		return
	}

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to save file",
		})
		return
	}
	defer out.Close()

	// Copy the file
	written, err := io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to save file",
		})
		return
	}

	var coverURL string

	// Only attempt PDF cover extraction for .pdf files
	if strings.ToLower(ext) == ".pdf" {
		coverURL, _ = generateCoverFromPDF(filePath)
	}

	if coverURL == "" {
		// Generate a placeholder cover for non-PDF files or when PDF extraction fails
		fallbackName := strings.TrimSuffix(filename, ext)
		if placeholderURL, phErr := generatePlaceholderCover("./uploads/covers", fallbackName); phErr == nil {
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

// generateCoverFromPDF renders the first PDF page to a PNG using pdftoppm (poppler)
func generateCoverFromPDF(pdfPath string) (string, error) {
	// Build unique base name from PDF filename (already UUID)
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

	// pdftoppm writes <outputBase>-1.png
	coverFile := outputBase + "-1.png"
	if _, err := os.Stat(coverFile); err != nil {
		// Some pdftoppm builds pad page numbers (e.g. -01.png). Try common variants before falling back.
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

	// Background fill
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.Set(x, y, bg)
		}
	}

	// Accent border
	for x := 0; x < img.Bounds().Dx(); x++ {
		img.Set(x, 0, accent)
		img.Set(x, img.Bounds().Dy()-1, accent)
	}
	for y := 0; y < img.Bounds().Dy(); y++ {
		img.Set(0, y, accent)
		img.Set(img.Bounds().Dx()-1, y, accent)
	}

	// Diagonal stripes for texture
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

	// Clamp values to [0,1]
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

	// Sanitize path
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

	// Validate file type
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

	// Generate unique filename
	filename := uuid.New().String() + ext
	filePath := filepath.Join("./uploads/covers", filename)

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to save file",
		})
		return
	}
	defer out.Close()

	// Copy the file
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

// GetBooksAdmin returns all books including hidden section books (admin only)
func GetBooksAdmin(c *gin.Context) {
	books, err := GetAllBooksAdmin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to fetch books",
		})
		return
	}

	if books == nil {
		books = []Book{}
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    books,
	})
}

// GetSectionsAdmin returns all sections including hidden ones (admin only)
func GetSectionsAdmin(c *gin.Context) {
	sections, err := GetAllSectionsAdmin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to fetch sections",
		})
		return
	}

	if sections == nil {
		sections = []Section{}
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    sections,
	})
}

// SharedDownload serves a book file via share token (no auth required)
func SharedDownload(c *gin.Context) {
	token := c.Param("token")
	book, err := GetBookByShareToken(token)
	if err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Invalid or expired link",
		})
		return
	}

	IncrementDownloads(book.ID)

	filePath := "." + book.FileURL
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "File not found",
		})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", book.FileName))
	contentType := mime.TypeByExtension(filepath.Ext(book.FileName))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	c.Header("Content-Type", contentType)
	c.File(filePath)
}
