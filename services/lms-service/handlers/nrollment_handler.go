package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"lms-service/models"
)

type EnrollmentHandler struct {
	DB *gorm.DB
}

func NewEnrollmentHandler(db *gorm.DB) *EnrollmentHandler {
	return &EnrollmentHandler{DB: db}
}

//////////////////////////////////////////////////////
// ➕ ENROLL STUDENT
//////////////////////////////////////////////////////

func (h *EnrollmentHandler) EnrollStudent(c *gin.Context) {

	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id") // 🔥 from auth

	var req struct {
		SubjectID string `json:"subject_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tid, _ := uuid.Parse(tenantID)
	uid, _ := uuid.Parse(userID)

	// 🔥 find student by user_id
	var student models.Student
	if err := h.DB.
		Where("user_id = ? AND tenant_id = ?", uid, tid).
		First(&student).Error; err != nil {

		c.JSON(404, gin.H{"error": "student profile not found"})
		return
	}

	subid, err := uuid.Parse(req.SubjectID)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid subject id"})
		return
	}

	// validate subject
	var subject models.Subject
	if err := h.DB.
		Where("id = ? AND tenant_id = ?", subid, tid).
		First(&subject).Error; err != nil {

		c.JSON(404, gin.H{"error": "subject not found"})
		return
	}

	// prevent duplicate
	var existing models.Enrollment
	if err := h.DB.
		Where("student_id = ? AND subject_id = ?", student.ID, subid).
		First(&existing).Error; err == nil {

		c.JSON(400, gin.H{"error": "already enrolled"})
		return
	}

	enrollment := models.Enrollment{
		ID:        uuid.New(),
		TenantID:  tid,
		StudentID: student.ID,
		SubjectID: subid,
	}

	h.DB.Create(&enrollment)

	c.JSON(201, enrollment)
}

//////////////////////////////////////////////////////
// 📄 LIST ENROLLMENTS
//////////////////////////////////////////////////////

func (h *EnrollmentHandler) ListEnrollments(c *gin.Context) {

	tenantID := c.GetString("tenant_id")

	var enrollments []models.Enrollment

	h.DB.Where("tenant_id = ?", tenantID).Find(&enrollments)

	c.JSON(200, enrollments)
}