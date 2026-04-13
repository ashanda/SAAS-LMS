package main

import (
	"log"
	"tenant-service/config"
	"tenant-service/database"
	"tenant-service/handlers"
	"tenant-service/models"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	db := database.Connect(cfg)

	// enable UUID
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`).Error; err != nil {
		log.Fatal("failed to enable pgcrypto: ", err)
	}

	// migrate
	if err := db.AutoMigrate(
		&models.Tenant{},
		&models.TenantDomain{},
		&models.TenantBranding{},
		&models.TenantSetting{},
	); err != nil {
		log.Fatal("failed to migrate tenant-service tables: ", err)
	}

	h := handlers.NewTenantHandler(db)

	r := gin.Default()

	// disable redirect
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false

	// routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "tenant service running"})
	})

	r.POST("/tenants", h.CreateTenant)
	r.GET("/tenants/by-domain", h.GetTenantByDomain)

	log.Println("tenant-service running on :" + cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}