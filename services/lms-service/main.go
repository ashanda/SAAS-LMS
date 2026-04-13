package main

import (
	"log"
	"os"
	"strings"

	"lms-service/config"
	"lms-service/database"
	"lms-service/handlers"
	"lms-service/models"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-Internal-User-ID")
		tenantID := c.GetHeader("X-Internal-Tenant-ID")
		roles := c.GetHeader("X-Internal-Roles")

		if userID == "" || tenantID == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing internal auth headers"})
			return
		}

		c.Set("user_id", userID)
		c.Set("tenant_id", tenantID)
		c.Set("roles", roles)
		c.Next()
	}
}

func RequireRole(allowed ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := c.GetString("roles")
		if raw == "" {
			c.AbortWithStatusJSON(403, gin.H{"error": "roles not found"})
			return
		}

		userRoles := strings.Split(raw, ",")
		for _, ur := range userRoles {
			for _, ar := range allowed {
				if strings.TrimSpace(ur) == ar {
					c.Next()
					return
				}
			}
		}

		c.AbortWithStatusJSON(403, gin.H{"error": "permission denied"})
	}
}

func main() {
	cfg := config.Load()
	db := database.Connect(cfg)

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`).Error; err != nil {
		log.Fatal("failed to enable pgcrypto: ", err)
	}

	if err := db.AutoMigrate(
		&models.Grade{},
		&models.Subject{},
		&models.Student{},
		&models.Teacher{},
		&models.SubjectTeacher{},
		&models.Enrollment{},
	); err != nil {
		log.Fatal("failed to migrate lms-service tables: ", err)
	}

	gradeHandler := handlers.NewGradeHandler(db)
	subjectHandler := handlers.NewSubjectHandler(db)
	studentHandler := handlers.NewStudentHandler(db)
	teacherHandler := handlers.NewTeacherHandler(db)
	enrollmentHandler := handlers.NewEnrollmentHandler(db)
	subjectTeacherHandler := handlers.NewSubjectTeacherHandler(db)
	dashboardHandler := handlers.NewDashboardHandler(db)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "lms-service running"})
	})

	api := r.Group("/")
	api.Use(AuthMiddleware())
	{
		adminOnly := api.Group("/")
		adminOnly.Use(RequireRole("owner", "admin", "staff"))
		{
			adminOnly.POST("/grades", gradeHandler.CreateGrade)
			adminOnly.GET("/grades", gradeHandler.ListGrades)

			adminOnly.POST("/subjects", subjectHandler.CreateSubject)
			adminOnly.GET("/subjects", subjectHandler.ListSubjects)

			adminOnly.POST("/students", studentHandler.CreateStudent)
			adminOnly.GET("/students", studentHandler.ListStudents)

			adminOnly.POST("/teachers", teacherHandler.CreateTeacher)
			adminOnly.POST("/assign-teacher", subjectTeacherHandler.AssignTeacher)

			adminOnly.GET("/dashboard/admin", dashboardHandler.AdminStats)
		}

		studentOnly := api.Group("/")
		studentOnly.Use(RequireRole("student"))
		{
			studentOnly.POST("/enrollments", enrollmentHandler.EnrollStudent)
			studentOnly.GET("/dashboard/student/subjects", studentHandler.MySubjects)
		}

		teacherOnly := api.Group("/")
		teacherOnly.Use(RequireRole("teacher"))
		{
			teacherOnly.GET("/teacher/subjects", teacherHandler.MySubjects)
			teacherOnly.GET("/teacher/subjects/:id/students", teacherHandler.SubjectStudents)
		}
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8083"
	}

	log.Println("lms-service running on :" + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}