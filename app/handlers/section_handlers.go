package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
