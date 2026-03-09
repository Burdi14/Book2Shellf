package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/books", GetBooks)
		api.GET("/books/:id", GetBook)
		api.GET("/sections", GetSections)
		api.GET("/sections/:id/books", GetBooksBySection)
	}

	router.POST("/api/login", AdminLogin)

	admin := router.Group("/api/admin")
	admin.Use(AuthMiddleware())
	{
		admin.POST("/books", CreateBook)
		admin.PUT("/books/:id", UpdateBook)
		admin.DELETE("/books/:id", DeleteBook)
		admin.POST("/sections", CreateSection)
		admin.PUT("/sections/:id", UpdateSection)
		admin.DELETE("/sections/:id", DeleteSection)
	}

	return router
}

func TestMain(m *testing.M) {
	os.Remove("./test_book2shelf.db")
	err := initTestDB()
	if err != nil {
		panic("Failed to initialize test database: " + err.Error())
	}
	code := m.Run()
	os.Remove("./test_book2shelf.db")
	os.Exit(code)
}

func initTestDB() error {
	var err error
	db, err = openDB("./test_book2shelf.db")
	if err != nil {
		return err
	}
	return createTables()
}

func TestGetSections(t *testing.T) {
	router := setupTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/sections", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}
}

func TestGetBooks(t *testing.T) {
	router := setupTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/books", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}
}

func TestGetBookNotFound(t *testing.T) {
	router := setupTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/books/nonexistent-id", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestAdminLoginInvalidCredentials(t *testing.T) {
	router := setupTestRouter()
	loginData := AdminCredentials{
		Username: "wrong",
		Password: "wrong",
	}
	jsonData, _ := json.Marshal(loginData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAdminLoginValidCredentials(t *testing.T) {
	router := setupTestRouter()
	loginData := AdminCredentials{
		Username: "admin",
		Password: "B00k2Sh3lf@dm1n!",
	}
	jsonData, _ := json.Marshal(loginData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if response.Token == "" {
		t.Error("Expected token in response")
	}
}

func TestAdminEndpointWithoutAuth(t *testing.T) {
	router := setupTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/books", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestCreateSectionWithAuth(t *testing.T) {
	router := setupTestRouter()

	loginData := AdminCredentials{
		Username: "admin",
		Password: "B00k2Sh3lf@dm1n!",
	}
	jsonData, _ := json.Marshal(loginData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var loginResponse LoginResponse
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse.Token

	sectionData := map[string]string{
		"name":        "Test Section",
		"description": "A test section",
	}
	jsonData, _ = json.Marshal(sectionData)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/admin/sections", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestCreateBookWithAuth(t *testing.T) {
	router := setupTestRouter()

	loginData := AdminCredentials{
		Username: "admin",
		Password: "B00k2Sh3lf@dm1n!",
	}
	jsonData, _ := json.Marshal(loginData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var loginResponse LoginResponse
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse.Token

	bookData := map[string]string{
		"title":       "Test Book",
		"author":      "Test Author",
		"description": "A test book",
	}
	jsonData, _ = json.Marshal(bookData)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/admin/books", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/books", nil)
	router.ServeHTTP(w, req)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	books, ok := response.Data.([]interface{})
	if !ok || len(books) == 0 {
		t.Error("Expected at least one book")
	}
}

func TestInvalidAuthToken(t *testing.T) {
	router := setupTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/books", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}
