package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"lms-service/models"
)

type SubjectHandler struct {
	DB *gorm.DB
}

func NewSubjectHandler(db *gorm.DB) *SubjectHandler {
	return &SubjectHandler{DB: db}
}

//////////////////////////////////////////////////////
// ➕ CREATE SUBJECT
//////////////////////////////////////////////////////

func (h *SubjectHandler) CreateSubject(c *gin.Context) {

	tenantID := c.GetString("tenant_id")

	var req struct {
		Name    string `json:"name"`
		GradeID string `json:"grade_id"`
	}

	// validate input
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// parse IDs
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid tenant id"})
		return
	}

	gid, err := uuid.Parse(req.GradeID)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid grade id"})
		return
	}

	// 🔒 CHECK: grade belongs to same tenant
	var grade models.Grade
	if err := h.DB.
		Where("id = ? AND tenant_id = ?", gid, tid).
		First(&grade).Error; err != nil {

		c.JSON(404, gin.H{"error": "grade not found in this tenant"})
		return
	}

	// create subject
	subject := models.Subject{
		ID:       uuid.New(),
		Name:     req.Name,
		TenantID: tid,
		GradeID:  gid,
	}

	if err := h.DB.Create(&subject).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, subject)
}

//////////////////////////////////////////////////////
// 📄 LIST SUBJECTS (TENANT SAFE)
//////////////////////////////////////////////////////

func (h *SubjectHandler) ListSubjects(c *gin.Context) {

	tenantID := c.GetString("tenant_id")

	var subjects []models.Subject

	if err := h.DB.
		Where("tenant_id = ?", tenantID).
		Find(&subjects).Error; err != nil {

		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, subjects)
}