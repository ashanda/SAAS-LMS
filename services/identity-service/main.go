package main

import (
	"log"

	"identity-service/config"
	"identity-service/database"
	"identity-service/handlers"
	"identity-service/models"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	db := database.Connect(cfg)

	if err := db.AutoMigrate(
		&models.User{},
		&models.TenantUser{},
		&models.Role{},
		&models.TenantUserRole{},
		&models.RefreshToken{},
		&models.Invite{},
	); err != nil {
		log.Fatal("failed to migrate identity-service tables: ", err)
	}

	defaultRoles := []models.Role{
		{Code: "owner", Name: "Owner"},
		{Code: "admin", Name: "Admin"},
		{Code: "teacher", Name: "Teacher"},
		{Code: "student", Name: "Student"},
		{Code: "staff", Name: "Staff"},
	}

	for _, role := range defaultRoles {
		var existing models.Role
		if err := db.Where("code = ?", role.Code).First(&existing).Error; err != nil {
			if err := db.Create(&role).Error; err != nil {
				log.Fatal("failed to seed role: ", role.Code, err)
			}
		}
	}

	h := handlers.NewAuthHandler(db, cfg)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "identity-service running"})
	})

	r.POST("/auth/register-owner", h.RegisterOwner)
	r.POST("/auth/register-student", h.RegisterStudent)
	r.POST("/auth/register-teacher", h.RegisterTeacher)
	r.POST("/auth/login", h.Login)
	r.POST("/auth/logout", h.Logout)

	log.Println("identity-service running on :8082")
	if err := r.Run(":8082"); err != nil {
		log.Fatal(err)
	}
}