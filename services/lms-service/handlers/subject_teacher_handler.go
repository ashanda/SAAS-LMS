package handlers

import (


	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"lms-service/models"
)

type SubjectTeacherHandler struct {
	DB *gorm.DB
}

func NewSubjectTeacherHandler(db *gorm.DB) *SubjectTeacherHandler {
	return &SubjectTeacherHandler{DB: db}
}

//////////////////////////////////////////////////////
// 🔗 ASSIGN TEACHER → SUBJECT
//////////////////////////////////////////////////////

func (h *SubjectTeacherHandler) AssignTeacher(c *gin.Context) {

	tenantID := c.GetString("tenant_id")

	var req struct {
		TeacherID string `json:"teacher_id"`
		SubjectID string `json:"subject_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tid, _ := uuid.Parse(tenantID)
	tchid, _ := uuid.Parse(req.TeacherID)
	subid, _ := uuid.Parse(req.SubjectID)

	// validate teacher
	var teacher models.Teacher
	if err := h.DB.Where("id = ? AND tenant_id = ?", tchid, tid).
		First(&teacher).Error; err != nil {

		c.JSON(404, gin.H{"error": "teacher not found"})
		return
	}

	// validate subject
	var subject models.Subject
	if err := h.DB.Where("id = ? AND tenant_id = ?", subid, tid).
		First(&subject).Error; err != nil {

		c.JSON(404, gin.H{"error": "subject not found"})
		return
	}

	// prevent duplicate
	var existing models.SubjectTeacher
	if err := h.DB.
		Where("teacher_id = ? AND subject_id = ?", tchid, subid).
		First(&existing).Error; err == nil {

		c.JSON(400, gin.H{"error": "already assigned"})
		return
	}

	assign := models.SubjectTeacher{
		ID:        uuid.New(),
		TenantID:  tid,
		TeacherID: tchid,
		SubjectID: subid,
	}

	h.DB.Create(&assign)

	c.JSON(201, assign)
}