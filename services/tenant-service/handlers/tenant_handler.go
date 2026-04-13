package handlers

import (
	"net/http"
	"strings"

	"tenant-service/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TenantHandler struct {
	DB *gorm.DB
}

func NewTenantHandler(db *gorm.DB) *TenantHandler {
	return &TenantHandler{DB: db}
}

type CreateTenantRequest struct {
	Name         string `json:"name" binding:"required"`
	Slug         string `json:"slug" binding:"required"`
	Domain       string `json:"domain" binding:"required"`
	DomainType   string `json:"domain_type" binding:"required"`
	LogoURL      string `json:"logo_url"`
	PrimaryColor string `json:"primary_color"`
}

//////////////////////////////////////////////////////
// CREATE TENANT
//////////////////////////////////////////////////////

func (h *TenantHandler) CreateTenant(c *gin.Context) {
	var req CreateTenantRequest

	// ✅ Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ Clean input
	req.Slug = strings.ToLower(strings.TrimSpace(req.Slug))
	req.Domain = strings.ToLower(strings.TrimSpace(req.Domain))

	// ✅ Validate slug
	if !isValidSlug(req.Slug) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid slug format"})
		return
	}

	// ✅ Check duplicate slug
	var existing models.Tenant
	if err := h.DB.Where("slug = ?", req.Slug).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug already exists"})
		return
	}

	// ✅ Check duplicate domain
	var existingDomain models.TenantDomain
	if err := h.DB.Where("domain = ?", req.Domain).First(&existingDomain).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "domain already exists"})
		return
	}

	// 🔥 Start transaction
	tx := h.DB.Begin()

	// ✅ Create Tenant
	tenant := models.Tenant{
		Name:   req.Name,
		Slug:   req.Slug,
		Status: "active",
	}

	if err := tx.Create(&tenant).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
		"error": err.Error(),
		})
		return
	}

	// ✅ Create Domain
	domain := models.TenantDomain{
		TenantID:  tenant.ID,
		Domain:    req.Domain,
		IsPrimary: true,
		Type:      req.DomainType,
	}

	if err := tx.Create(&domain).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create domain"})
		return
	}

	// ✅ Create Branding
	branding := models.TenantBranding{
		TenantID:     tenant.ID,
		LogoURL:      req.LogoURL,
		PrimaryColor: req.PrimaryColor,
	}

	if err := tx.Create(&branding).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create branding"})
		return
	}

	// ✅ Create Default Settings
	settings := models.TenantSetting{
		TenantID: tenant.ID,
		Timezone: "Asia/Colombo",
		Currency: "LKR",
		Locale:   "en",
	}

	if err := tx.Create(&settings).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create settings"})
		return
	}

	// ✅ Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction commit failed"})
		return
	}

	// 🔥 Reload FULL tenant with relations
	var fullTenant models.Tenant

	if err := h.DB.
		Preload("Domains").
		Preload("Branding").
		First(&fullTenant, "id = ?", tenant.ID).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load tenant"})
		return
	}

	// ✅ Response
	c.JSON(http.StatusCreated, gin.H{
		"message": "tenant created successfully",
		"tenant":  fullTenant,
	})
}

//////////////////////////////////////////////////////
// GET TENANT BY DOMAIN
//////////////////////////////////////////////////////

func (h *TenantHandler) GetTenantByDomain(c *gin.Context) {
	domain := strings.ToLower(strings.TrimSpace(c.Query("domain")))

	if domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "domain is required"})
		return
	}

	var tenant models.Tenant

	err := h.DB.
		Joins("JOIN tenant_domains ON tenant_domains.tenant_id = tenants.id").
		Where("tenant_domains.domain = ?", domain).
		Preload("Domains").
		Preload("Branding").
		First(&tenant).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tenant": tenant,
	})
}