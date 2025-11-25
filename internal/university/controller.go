package university

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UniversityController struct {
	UniversityService *UniversityService
}

func (uc *UniversityController) GetUniversities(c *gin.Context) {
	data, err := uc.UniversityService.GetUniversities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch universities"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Universities fetched successfully",
		"data":    data,
	})
}

func (uc *UniversityController) GetFields(c *gin.Context) {
	data, err := uc.UniversityService.GetFields()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch fields"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Fields fetched successfully",
		"data":    data,
	})
}

func (uc *UniversityController) GetSpecializations(c *gin.Context) {
	fieldId := c.Query("fieldId")
	data, err := uc.UniversityService.GetSpecializations(fieldId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch specializations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Specializations fetched successfully",
		"data":    data,
	})
}

func (uc *UniversityController) SearchUniversities(c *gin.Context) {
	filters := map[string]string{
		"universityId":     c.Query("universityId"),
		"fieldId":          c.Query("fieldId"),
		"specializationId": c.Query("specializationId"),
		"level":            c.Query("level"),
		"duration":         c.Query("duration"),
		"q":                c.Query("q"),
	}

	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	list, total, pageInt, limitInt, pages, err := uc.UniversityService.SearchUniversities(filters, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":  pageInt,
		"limit": limitInt,
		"total": total,
		"pages": pages,
		"data":  list,
	})
}

func (uc *UniversityController) AddCourse(c *gin.Context) {
	var req Course
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := uc.UniversityService.AddCourse(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish event"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Course update published"})
}

func (uc *UniversityController) UploadCourses(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	tempPath := "./tmp/" + file.Filename
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
		return
	}

	count, err := uc.UniversityService.UploadCourses(tempPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Upload successful",
		"inserted": count,
	})
}
