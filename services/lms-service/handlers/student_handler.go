package handlers

import (
	

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"lms-service/models"
)

type StudentHandler struct {
	DB *gorm.DB
}

func NewStudentHandler(db *gorm.DB) *StudentHandler {
	return &StudentHandler{DB: db}
}

//////////////////////////////////////////////////////
// ➕ CREATE STUDENT
//////////////////////////////////////////////////////

func (h *StudentHandler) CreateStudent(c *gin.Context) {

	tenantID := c.GetString("tenant_id")

	var req struct {
		UserID string `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tid, _ := uuid.Parse(tenantID)
	uid, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user id"})
		return
	}

	// 🔒 prevent duplicate
	var existing models.Student
	if err := h.DB.
		Where("tenant_id = ? AND user_id = ?", tid, uid).
		First(&existing).Error; err == nil {

		c.JSON(400, gin.H{"error": "student already exists"})
		return
	}

	student := models.Student{
		ID:       uuid.New(),
		TenantID: tid,
		UserID:   uid,
	}

	h.DB.Create(&student)

	c.JSON(201, student)
}

//////////////////////////////////////////////////////
// 📄 LIST STUDENTS
//////////////////////////////////////////////////////

func (h *StudentHandler) ListStudents(c *gin.Context) {

	tenantID := c.GetString("tenant_id")

	var students []models.Student

	h.DB.Where("tenant_id = ?", tenantID).Find(&students)

	c.JSON(200, students)
}


//////////////////////////////////////////////////////
// 📄 My Subjects
//////////////////////////////////////////////////////
func (h *StudentHandler) MySubjects(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")

	tid, _ := uuid.Parse(tenantID)
	uid, _ := uuid.Parse(userID)

	var student models.Student
	if err := h.DB.
		Where("user_id = ? AND tenant_id = ?", uid, tid).
		First(&student).Error; err != nil {
		c.JSON(404, gin.H{"error": "student not found"})
		return
	}

	var subjects []models.Subject
	h.DB.
		Table("subjects").
		Select("subjects.*").
		Joins("JOIN enrollments e ON e.subject_id = subjects.id").
		Where("e.student_id = ?", student.ID).
		Where("subjects.tenant_id = ?", tid).
		Find(&subjects)

	if subjects == nil {
		subjects = []models.Subject{}
	}

	c.JSON(200, subjects)
}