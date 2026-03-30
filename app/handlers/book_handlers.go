package handlers

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
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

	if book.SectionID != "" {
		if isHidden, _ := IsSectionHidden(book.SectionID); isHidden {
			c.JSON(http.StatusNotFound, APIResponse{
				Success: false,
				Message: "Book not found",
			})
			return
		}
	}

	IncrementDownloads(id)
	serveBookFile(c, book, "Book not found")
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

	book, err := GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Book not found",
		})
		return
	}

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
	serveBookFile(c, book, "Invalid or expired link")
}

func serveBookFile(c *gin.Context, book *Book, missingMessage string) {
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
