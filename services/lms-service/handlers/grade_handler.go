package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"lms-service/models"
)

type GradeHandler struct {
	DB *gorm.DB
}

func NewGradeHandler(db *gorm.DB) *GradeHandler {
	return &GradeHandler{DB: db}
}

//////////////////////////////////////////////////////
// ➕ CREATE GRADE
//////////////////////////////////////////////////////

func (h *GradeHandler) CreateGrade(c *gin.Context) {

	tenantID := c.GetString("tenant_id")

	var req struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	grade := models.Grade{
		ID:       uuid.New(),
		Name:     req.Name,
		TenantID: uuid.MustParse(tenantID),
	}

	h.DB.Create(&grade)

	c.JSON(201, grade)
}

//////////////////////////////////////////////////////
// 📄 LIST GRADES
//////////////////////////////////////////////////////

func (h *GradeHandler) ListGrades(c *gin.Context) {

	tenantID := c.GetString("tenant_id")

	var grades []models.Grade

	h.DB.Where("tenant_id = ?", tenantID).Find(&grades)

	c.JSON(200, grades)
}