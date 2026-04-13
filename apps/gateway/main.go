package main

import (
	"log"
	"strings"

	"gateway/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// =========================================
	// 🌐 CORS
	// =========================================
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Tenant-Domain")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// =========================================
	// 🩺 HEALTH
	// =========================================
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "gateway running"})
	})

	// =========================================
	// 🏢 TENANT SERVICE (PUBLIC)
	// =========================================
	r.Any("/tenants", gin.WrapH(proxy("http://tenant-service:8081", "")))
	r.Any("/tenants/*path", gin.WrapH(proxy("http://tenant-service:8081", "")))

	// =========================================
	// 🔥 TENANT GROUP
	// =========================================
	tenantGroup := r.Group("/")
	tenantGroup.Use(middleware.TenantMiddleware())
	{
		// =====================================
		// 🔐 AUTH ROUTES (FIXED)
		// =====================================
		tenantGroup.Any("/auth/*path", func(c *gin.Context) {
			tenantID := c.MustGet("tenant_id").(string)

			log.Println("Gateway Tenant ID:", tenantID)

			// set header BEFORE proxy
			c.Request.Header.Set("X-Tenant-ID", tenantID)

			proxy("http://identity-service:8082", "").ServeHTTP(c.Writer, c.Request)
		})

		// =====================================
		// 🔍 DEBUG ENDPOINT
		// =====================================
		tenantGroup.GET("/tenant-health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"tenant_id":   c.MustGet("tenant_id"),
				"tenant_slug": c.MustGet("tenant_slug"),
			})
		})

		// =====================================
		// 🔒 PROTECTED ROUTES
		// =====================================
		protected := tenantGroup.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/me", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"user_id":   c.MustGet("user_id"),
					"roles":     c.MustGet("roles"),
					"tenant_id": c.MustGet("tenant_id"),
				})
			})

			protected.Any("/lms/*path", func(c *gin.Context) {
				// inject user + tenant
				c.Request.Header.Set("X-Internal-User-ID", c.MustGet("user_id").(string))
				c.Request.Header.Set("X-Internal-Tenant-ID", c.MustGet("tenant_id").(string))

				if roles, ok := c.MustGet("roles").([]string); ok {
					c.Request.Header.Set("X-Internal-Roles", strings.Join(roles, ","))
				}

				proxy("http://lms-service:8083", "/lms").ServeHTTP(c.Writer, c.Request)
			})
		}
	}

	log.Println("🚀 Gateway running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}