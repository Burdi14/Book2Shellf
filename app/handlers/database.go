package handlers

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// openDB opens a connection to the SQLite database
func openDB(path string) (*sql.DB, error) {
	return sql.Open("sqlite3", path)
}

// createTables creates the database schema
func createTables() error {
	createTablesSQL := `
	CREATE TABLE IF NOT EXISTS sections (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		hidden INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS books (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		author TEXT,
		description TEXT,
		year INTEGER DEFAULT 0,
		cover_url TEXT,
		file_url TEXT,
		file_name TEXT,
		file_size INTEGER DEFAULT 0,
		section_id TEXT,
		downloads INTEGER DEFAULT 0,
		share_token TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (section_id) REFERENCES sections(id)
	);

	CREATE INDEX IF NOT EXISTS idx_books_section ON books(section_id);
	`

	_, err := db.Exec(createTablesSQL)
	if err != nil {
		return err
	}

	// Run migrations for existing databases
	runMigrations()

	// Create index on share_token after migrations ensure the column exists
	db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_books_share_token ON books(share_token)`)

	return nil
}

// runMigrations adds new columns to existing databases
func runMigrations() {
	// Add year column if it doesn't exist
	db.Exec(`ALTER TABLE books ADD COLUMN year INTEGER DEFAULT 0`)
	// Add hidden column to sections
	db.Exec(`ALTER TABLE sections ADD COLUMN hidden INTEGER DEFAULT 0`)
	// Add share_token column to books
	db.Exec(`ALTER TABLE books ADD COLUMN share_token TEXT`)
	// Back-fill share tokens for existing books that don't have one
	rows, err := db.Query(`SELECT id FROM books WHERE share_token IS NULL OR share_token = ''`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id string
			if rows.Scan(&id) == nil {
				db.Exec(`UPDATE books SET share_token = ? WHERE id = ?`, uuid.New().String(), id)
			}
		}
	}

	backfillMissingFileSizes()
}

func resolveStoredFileSize(fileURL string) (int64, error) {
	if strings.TrimSpace(fileURL) == "" {
		return 0, os.ErrNotExist
	}

	cleanPath := filepath.Clean(strings.TrimPrefix(fileURL, "/"))
	info, err := os.Stat(filepath.Join(".", cleanPath))
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}

func backfillMissingFileSizes() {
	rows, err := db.Query(`SELECT id, file_url FROM books WHERE COALESCE(file_size, 0) <= 0 AND COALESCE(file_url, '') != ''`)
	if err != nil {
		return
	}
	defer rows.Close()

	var repaired int
	for rows.Next() {
		var id, fileURL string
		if err := rows.Scan(&id, &fileURL); err != nil {
			continue
		}

		size, err := resolveStoredFileSize(fileURL)
		if err != nil || size <= 0 {
			continue
		}

		if _, err := db.Exec(`UPDATE books SET file_size = ?, updated_at = ? WHERE id = ?`, size, time.Now(), id); err == nil {
			repaired++
		}
	}

	if repaired > 0 {
		log.Printf("Repaired file_size for %d books", repaired)
	}
}

// InitDB initializes the SQLite database
func InitDB() error {
	var err error

	// Ensure data directory exists
	os.MkdirAll("./data", 0755)

	db, err = openDB("./data/book2shelf.db")
	if err != nil {
		return err
	}

	err = createTables()
	if err != nil {
		return err
	}

	// Insert default section if none exists
	var count int
	db.QueryRow("SELECT COUNT(*) FROM sections").Scan(&count)
	if count == 0 {
		_, err = db.Exec(`INSERT INTO sections (id, name, description) VALUES (?, ?, ?)`,
			uuid.New().String(), "General", "General collection of books")
		if err != nil {
			log.Println("Warning: Could not create default section:", err)
		}
	}

	return nil
}

// GetAllSections returns only visible (non-hidden) sections
func GetAllSections() ([]Section, error) {
	rows, err := db.Query(`
		SELECT s.id, s.name, s.description, COALESCE(s.hidden, 0), s.created_at, 
			   COALESCE((SELECT COUNT(*) FROM books WHERE section_id = s.id), 0) as book_count
		FROM sections s
		WHERE COALESCE(s.hidden, 0) = 0
		ORDER BY s.name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []Section
	for rows.Next() {
		var s Section
		var hidden int
		err := rows.Scan(&s.ID, &s.Name, &s.Description, &hidden, &s.CreatedAt, &s.BookCount)
		if err != nil {
			return nil, err
		}
		s.Hidden = hidden != 0
		sections = append(sections, s)
	}
	return sections, nil
}

// GetAllSectionsAdmin returns all sections including hidden ones
func GetAllSectionsAdmin() ([]Section, error) {
	rows, err := db.Query(`
		SELECT s.id, s.name, s.description, COALESCE(s.hidden, 0), s.created_at, 
			   COALESCE((SELECT COUNT(*) FROM books WHERE section_id = s.id), 0) as book_count
		FROM sections s
		ORDER BY s.name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []Section
	for rows.Next() {
		var s Section
		var hidden int
		err := rows.Scan(&s.ID, &s.Name, &s.Description, &hidden, &s.CreatedAt, &s.BookCount)
		if err != nil {
			return nil, err
		}
		s.Hidden = hidden != 0
		sections = append(sections, s)
	}
	return sections, nil
}

// CreateSectionDB creates a new section
func CreateSectionDB(name, description string, hidden bool) (*Section, error) {
	id := uuid.New().String()
	now := time.Now()

	hiddenInt := 0
	if hidden {
		hiddenInt = 1
	}

	_, err := db.Exec(`INSERT INTO sections (id, name, description, hidden, created_at) VALUES (?, ?, ?, ?, ?)`,
		id, name, description, hiddenInt, now)
	if err != nil {
		return nil, err
	}

	return &Section{
		ID:          id,
		Name:        name,
		Description: description,
		Hidden:      hidden,
		CreatedAt:   now,
	}, nil
}

// UpdateSectionDB updates a section
func UpdateSectionDB(id, name, description string, hidden bool) error {
	hiddenInt := 0
	if hidden {
		hiddenInt = 1
	}
	_, err := db.Exec(`UPDATE sections SET name = ?, description = ?, hidden = ? WHERE id = ?`,
		name, description, hiddenInt, id)
	return err
}

// DeleteSectionDB deletes a section
func DeleteSectionDB(id string) error {
	// First, set books in this section to have no section
	db.Exec(`UPDATE books SET section_id = NULL WHERE section_id = ?`, id)

	_, err := db.Exec(`DELETE FROM sections WHERE id = ?`, id)
	return err
}

// GetAllBooks returns only books in visible (non-hidden) sections
func GetAllBooks() ([]Book, error) {
	rows, err := db.Query(`
		SELECT b.id, b.title, b.author, b.description, COALESCE(b.year, 0), b.cover_url, b.file_url, 
			   b.file_name, b.file_size, b.section_id, COALESCE(s.name, 'Uncategorized') as section_name,
			   b.downloads, COALESCE(b.share_token, ''), b.created_at, b.updated_at
		FROM books b
		LEFT JOIN sections s ON b.section_id = s.id
		WHERE COALESCE(s.hidden, 0) = 0 OR b.section_id IS NULL
		ORDER BY b.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		var sectionID sql.NullString
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Description, &b.Year, &b.CoverURL, &b.FileURL,
			&b.FileName, &b.FileSize, &sectionID, &b.SectionName, &b.Downloads, &b.ShareToken, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if sectionID.Valid {
			b.SectionID = sectionID.String
		}
		books = append(books, b)
	}
	return books, nil
}

// GetAllBooksAdmin returns all books including those in hidden sections
func GetAllBooksAdmin() ([]Book, error) {
	rows, err := db.Query(`
		SELECT b.id, b.title, b.author, b.description, COALESCE(b.year, 0), b.cover_url, b.file_url, 
			   b.file_name, b.file_size, b.section_id, COALESCE(s.name, 'Uncategorized') as section_name,
			   b.downloads, COALESCE(b.share_token, ''), b.created_at, b.updated_at
		FROM books b
		LEFT JOIN sections s ON b.section_id = s.id
		ORDER BY b.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		var sectionID sql.NullString
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Description, &b.Year, &b.CoverURL, &b.FileURL,
			&b.FileName, &b.FileSize, &sectionID, &b.SectionName, &b.Downloads, &b.ShareToken, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if sectionID.Valid {
			b.SectionID = sectionID.String
		}
		books = append(books, b)
	}
	return books, nil
}

// GetBooksBySection returns books in a specific section
func GetBooksBySectionDB(sectionID string) ([]Book, error) {
	rows, err := db.Query(`
		SELECT b.id, b.title, b.author, b.description, COALESCE(b.year, 0), b.cover_url, b.file_url, 
			   b.file_name, b.file_size, b.section_id, COALESCE(s.name, 'Uncategorized') as section_name,
			   b.downloads, COALESCE(b.share_token, ''), b.created_at, b.updated_at
		FROM books b
		LEFT JOIN sections s ON b.section_id = s.id
		WHERE b.section_id = ?
		ORDER BY b.created_at DESC
	`, sectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		var secID sql.NullString
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Description, &b.Year, &b.CoverURL, &b.FileURL,
			&b.FileName, &b.FileSize, &secID, &b.SectionName, &b.Downloads, &b.ShareToken, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if secID.Valid {
			b.SectionID = secID.String
		}
		books = append(books, b)
	}
	return books, nil
}

// GetBookByID returns a single book
func GetBookByID(id string) (*Book, error) {
	var b Book
	var sectionID sql.NullString
	err := db.QueryRow(`
		SELECT b.id, b.title, b.author, b.description, COALESCE(b.year, 0), b.cover_url, b.file_url, 
			   b.file_name, b.file_size, b.section_id, COALESCE(s.name, 'Uncategorized') as section_name,
			   b.downloads, COALESCE(b.share_token, ''), b.created_at, b.updated_at
		FROM books b
		LEFT JOIN sections s ON b.section_id = s.id
		WHERE b.id = ?
	`, id).Scan(&b.ID, &b.Title, &b.Author, &b.Description, &b.Year, &b.CoverURL, &b.FileURL,
		&b.FileName, &b.FileSize, &sectionID, &b.SectionName, &b.Downloads, &b.ShareToken, &b.CreatedAt, &b.UpdatedAt)

	if err != nil {
		return nil, err
	}
	if sectionID.Valid {
		b.SectionID = sectionID.String
	}
	return &b, nil
}

// GetBookByShareToken returns a book by its share token
func GetBookByShareToken(token string) (*Book, error) {
	var b Book
	var sectionID sql.NullString
	err := db.QueryRow(`
		SELECT b.id, b.title, b.author, b.description, COALESCE(b.year, 0), b.cover_url, b.file_url, 
			   b.file_name, b.file_size, b.section_id, COALESCE(s.name, 'Uncategorized') as section_name,
			   b.downloads, COALESCE(b.share_token, ''), b.created_at, b.updated_at
		FROM books b
		LEFT JOIN sections s ON b.section_id = s.id
		WHERE b.share_token = ?
	`, token).Scan(&b.ID, &b.Title, &b.Author, &b.Description, &b.Year, &b.CoverURL, &b.FileURL,
		&b.FileName, &b.FileSize, &sectionID, &b.SectionName, &b.Downloads, &b.ShareToken, &b.CreatedAt, &b.UpdatedAt)

	if err != nil {
		return nil, err
	}
	if sectionID.Valid {
		b.SectionID = sectionID.String
	}
	return &b, nil
}

// CreateBookDB creates a new book
func CreateBookDB(book *Book) error {
	book.ID = uuid.New().String()
	book.ShareToken = uuid.New().String()
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()

	if book.FileSize <= 0 && book.FileURL != "" {
		if size, err := resolveStoredFileSize(book.FileURL); err == nil {
			book.FileSize = size
		}
	}

	var sectionID interface{} = book.SectionID
	if book.SectionID == "" {
		sectionID = nil
	}

	_, err := db.Exec(`
		INSERT INTO books (id, title, author, description, year, cover_url, file_url, file_name, file_size, section_id, downloads, share_token, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, book.ID, book.Title, book.Author, book.Description, book.Year, book.CoverURL, book.FileURL,
		book.FileName, book.FileSize, sectionID, 0, book.ShareToken, book.CreatedAt, book.UpdatedAt)

	return err
}

// UpdateBookDB updates a book
func UpdateBookDB(book *Book) error {
	book.UpdatedAt = time.Now()

	if book.FileSize <= 0 && book.FileURL != "" {
		if size, err := resolveStoredFileSize(book.FileURL); err == nil {
			book.FileSize = size
		}
	}

	var sectionID interface{} = book.SectionID
	if book.SectionID == "" {
		sectionID = nil
	}

	_, err := db.Exec(`
		UPDATE books 
		SET title = ?, author = ?, description = ?, year = ?, cover_url = ?, file_url = ?, 
			file_name = ?, file_size = ?, section_id = ?, updated_at = ?
		WHERE id = ?
	`, book.Title, book.Author, book.Description, book.Year, book.CoverURL, book.FileURL,
		book.FileName, book.FileSize, sectionID, book.UpdatedAt, book.ID)

	return err
}

// DeleteBookDB deletes a book
func DeleteBookDB(id string) error {
	_, err := db.Exec(`DELETE FROM books WHERE id = ?`, id)
	return err
}

// IncrementDownloads increments the download count
func IncrementDownloads(id string) error {
	_, err := db.Exec(`UPDATE books SET downloads = downloads + 1 WHERE id = ?`, id)
	return err
}

// IsSectionHidden checks if a section is marked hidden
func IsSectionHidden(sectionID string) (bool, error) {
	var hidden int
	err := db.QueryRow(`SELECT COALESCE(hidden, 0) FROM sections WHERE id = ?`, sectionID).Scan(&hidden)
	if err != nil {
		return false, err
	}
	return hidden != 0, nil
}
