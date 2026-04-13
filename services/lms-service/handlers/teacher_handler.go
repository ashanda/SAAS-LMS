package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"lms-service/models"
)

type TeacherHandler struct {
	DB *gorm.DB
}

func NewTeacherHandler(db *gorm.DB) *TeacherHandler {
	return &TeacherHandler{DB: db}
}

//////////////////////////////////////////////////////
// ➕ CREATE TEACHER PROFILE
//////////////////////////////////////////////////////

func (h *TeacherHandler) CreateTeacher(c *gin.Context) {

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

	// prevent duplicate
	var existing models.Teacher
	if err := h.DB.
		Where("tenant_id = ? AND user_id = ?", tid, uid).
		First(&existing).Error; err == nil {

		c.JSON(400, gin.H{"error": "teacher already exists"})
		return
	}

	teacher := models.Teacher{
		ID:       uuid.New(),
		TenantID: tid,
		UserID:   uid,
	}

	h.DB.Create(&teacher)

	c.JSON(201, teacher)
}


//////////////////////////////////////////////////////
// 📄 LIST TEACHER
//////////////////////////////////////////////////////

func (h *TeacherHandler) ListTeachers(c *gin.Context) {

	tenantID := c.GetString("tenant_id")

	var teachers []models.Teacher

	h.DB.Where("tenant_id = ?", tenantID).Find(&teachers)

	c.JSON(200, teachers)
}

//////////////////////////////////////////////////////
// 📄 GET TEACHER SUBJECTS
//////////////////////////////////////////////////////

func (h *TeacherHandler) MySubjects(c *gin.Context) {

	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")

	tid, _ := uuid.Parse(tenantID)
	uid, _ := uuid.Parse(userID)

	// 🔥 find teacher
	var teacher models.Teacher
	if err := h.DB.
		Where("user_id = ? AND tenant_id = ?", uid, tid).
		First(&teacher).Error; err != nil {

		c.JSON(404, gin.H{"error": "teacher not found"})
		return
	}

	// 🔥 join subjects
	var subjects []models.Subject

	h.DB.
		Table("subjects").
		Select("subjects.*").
		Joins("JOIN subject_teachers st ON st.subject_id = subjects.id").
		Where("st.teacher_id = ?", teacher.ID).
		Where("subjects.tenant_id = ?", tid).
		Find(&subjects)

	c.JSON(200, subjects)
}

//////////////////////////////////////////////////////
// 📄 GET STUDENTS PER SUBJECT
//////////////////////////////////////////////////////
func (h *TeacherHandler) SubjectStudents(c *gin.Context) {

	tenantID := c.GetString("tenant_id")
	subjectID := c.Param("id")

	tid, _ := uuid.Parse(tenantID)
	subid, _ := uuid.Parse(subjectID)

	// 🔥 get students via enrollments
	var students []models.Student

	h.DB.
		Table("students").
		Select("students.*").
		Joins("JOIN enrollments e ON e.student_id = students.id").
		Where("e.subject_id = ?", subid).
		Where("students.tenant_id = ?", tid).
		Find(&students)

	c.JSON(200, students)
}