package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardHandler struct {
	DB *gorm.DB
}

func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{DB: db}
}

func (h *DashboardHandler) AdminStats(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	var gradesCount int64
	var subjectsCount int64
	var studentsCount int64
	var teachersCount int64
	var enrollmentsCount int64

	h.DB.Table("grades").Where("tenant_id = ?", tenantID).Count(&gradesCount)
	h.DB.Table("subjects").Where("tenant_id = ?", tenantID).Count(&subjectsCount)
	h.DB.Table("students").Where("tenant_id = ?", tenantID).Count(&studentsCount)
	h.DB.Table("teachers").Where("tenant_id = ?", tenantID).Count(&teachersCount)
	h.DB.Table("enrollments").Where("tenant_id = ?", tenantID).Count(&enrollmentsCount)

	c.JSON(200, gin.H{
		"grades_count":      gradesCount,
		"subjects_count":    subjectsCount,
		"students_count":    studentsCount,
		"teachers_count":    teachersCount,
		"enrollments_count": enrollmentsCount,
	})
}