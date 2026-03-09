package handlers

import "time"

// Book represents a book in the library
type Book struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	Year        int       `json:"year"`
	CoverURL    string    `json:"cover_url"`
	FileURL     string    `json:"file_url"`
	FileName    string    `json:"file_name"`
	FileSize    int64     `json:"file_size"`
	SectionID   string    `json:"section_id"`
	SectionName string    `json:"section_name,omitempty"`
	Downloads   int       `json:"downloads"`
	ShareToken  string    `json:"share_token,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Section represents a category/section of books
type Section struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Hidden      bool      `json:"hidden"`
	BookCount   int       `json:"book_count,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// AdminCredentials for login
type AdminCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse with token
type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

// APIResponse generic response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
