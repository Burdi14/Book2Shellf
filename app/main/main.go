package main

import (
	"log"
	"os"

	"book2shelf/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	os.MkdirAll("./data", 0755)
	os.MkdirAll("./uploads/books", 0755)
	os.MkdirAll("./uploads/covers", 0755)

	// Initialize database
	if err := handlers.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	handlers.SyncBookFileSizes()

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	router.Static("/uploads", "./uploads")

	// Public API routes
	api := router.Group("/api")
	{
		api.GET("/books", handlers.GetBooks)
		api.GET("/books/:id", handlers.GetBook)
		api.GET("/books/:id/download", handlers.DownloadBook)
		api.GET("/sections", handlers.GetSections)
		api.GET("/sections/:id/books", handlers.GetBooksBySection)
		api.GET("/share/:token", handlers.SharedDownload)
	}

	// Admin routes - protected
	admin := router.Group("/api/admin")
	admin.Use(handlers.AuthMiddleware())
	{
		admin.GET("/books", handlers.GetBooksAdmin)
		admin.GET("/sections", handlers.GetSectionsAdmin)
		admin.POST("/login", handlers.AdminLogin)
		admin.POST("/books", handlers.CreateBook)
		admin.PUT("/books/:id", handlers.UpdateBook)
		admin.DELETE("/books/:id", handlers.DeleteBook)
		admin.POST("/sections", handlers.CreateSection)
		admin.PUT("/sections/:id", handlers.UpdateSection)
		admin.DELETE("/sections/:id", handlers.DeleteSection)
		admin.POST("/upload/book", handlers.UploadBook)
		admin.POST("/upload/cover", handlers.UploadCover)
		admin.POST("/cover/crop", handlers.CropCover)
	}

	// Login route (unprotected)
	router.POST("/api/login", handlers.AdminLogin)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API server starting on port %s", port)
	router.Run(":" + port)
}
