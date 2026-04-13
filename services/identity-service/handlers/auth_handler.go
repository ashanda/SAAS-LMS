package handlers

import (
	"net/http"
	"strings"
	"time"

	"identity-service/config"
	"identity-service/models"
	"identity-service/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB     *gorm.DB
	Config *config.Config
}

func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{DB: db, Config: cfg}
}

//////////////////////////////////////////////////////
// REGISTER OWNER
//////////////////////////////////////////////////////

func (h *AuthHandler) RegisterOwner(c *gin.Context) {
	var req struct {
		TenantID string `json:"tenant_id" binding:"required"`
		FullName string `json:"full_name" binding:"required"`
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Phone    string `json:"phone"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.TrimSpace(req.Username)
	req.FullName = strings.TrimSpace(req.FullName)
	req.Phone = strings.TrimSpace(req.Phone)

	var existing models.User
	if err := h.DB.Where("email = ? OR username = ?", req.Email, req.Username).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email or username already exists"})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	tid, err := uuid.Parse(req.TenantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant_id"})
		return
	}

	tx := h.DB.Begin()

	user := models.User{
		FullName: req.FullName,
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantUser := models.TenantUser{
		TenantID: tid,
		UserID:   user.ID,
		Status:   "active",
	}

	if err := tx.Create(&tenantUser).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ownerRole models.Role
	if err := tx.Where("code = ?", "owner").First(&ownerRole).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "owner role not found"})
		return
	}

	tur := models.TenantUserRole{
		TenantUserID: tenantUser.ID,
		RoleID:       ownerRole.ID,
	}

	if err := tx.Create(&tur).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "owner registered successfully",
		"user": gin.H{
			"id":        user.ID,
			"full_name": user.FullName,
			"username":  user.Username,
			"email":     user.Email,
			"phone":     user.Phone,
		},
	})
}

//////////////////////////////////////////////////////
// REGISTER STUDENT
//////////////////////////////////////////////////////
func (h *AuthHandler) RegisterStudent(c *gin.Context) {
	var req struct {
		TenantID string `json:"tenant_id" binding:"required"`
		FullName string `json:"full_name" binding:"required"`
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Phone    string `json:"phone"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.TrimSpace(req.Username)
	req.FullName = strings.TrimSpace(req.FullName)
	req.Phone = strings.TrimSpace(req.Phone)

	var existing models.User
	if err := h.DB.Where("email = ? OR username = ?", req.Email, req.Username).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email or username already exists"})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	tid, err := uuid.Parse(req.TenantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant_id"})
		return
	}

	tx := h.DB.Begin()

	user := models.User{
		FullName: req.FullName,
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantUser := models.TenantUser{
		TenantID: tid,
		UserID:   user.ID,
		Status:   "active",
	}

	if err := tx.Create(&tenantUser).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var studentRole models.Role
	if err := tx.Where("code = ?", "student").First(&studentRole).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "student role not found"})
		return
	}

	tur := models.TenantUserRole{
		TenantUserID: tenantUser.ID,
		RoleID:       studentRole.ID,
	}

	if err := tx.Create(&tur).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "student registered successfully",
		"user": gin.H{
			"id":        user.ID,
			"full_name": user.FullName,
			"username":  user.Username,
			"email":     user.Email,
			"phone":     user.Phone,
		},
	})
}

//////////////////////////////////////////////////////
// REGISTER TEACHER
//////////////////////////////////////////////////////
func (h *AuthHandler) RegisterTeacher(c *gin.Context) {
	var req struct {
		TenantID string `json:"tenant_id" binding:"required"`
		FullName string `json:"full_name" binding:"required"`
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Phone    string `json:"phone"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.TrimSpace(req.Username)

	// check existing
	var existing models.User
	if err := h.DB.Where("email = ? OR username = ?", req.Email, req.Username).
		First(&existing).Error; err == nil {

		c.JSON(400, gin.H{"error": "user already exists"})
		return
	}

	hashed, _ := utils.HashPassword(req.Password)
	tid, _ := uuid.Parse(req.TenantID)

	tx := h.DB.Begin()

	user := models.User{
		FullName: req.FullName,
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashed,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tenantUser := models.TenantUser{
		TenantID: tid,
		UserID:   user.ID,
		Status:   "active",
	}

	tx.Create(&tenantUser)

	// 🔥 assign teacher role
	var role models.Role
	if err := tx.Where("code = ?", "teacher").First(&role).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "teacher role not found"})
		return
	}

	tx.Create(&models.TenantUserRole{
		TenantUserID: tenantUser.ID,
		RoleID:       role.ID,
	})

	tx.Commit()

	c.JSON(201, gin.H{
		"message": "teacher created",
		"user_id": user.ID,
	})
}


//////////////////////////////////////////////////////
// LOGIN
//////////////////////////////////////////////////////

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		TenantID string `json:"tenant_id"` // optional fallback
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	// parse body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// =========================================
	// 🔥 FIX: GET TENANT ID FROM HEADER
	// =========================================
	tenantID := c.GetHeader("X-Tenant-ID")

	// fallback (optional)
	if tenantID == "" {
		tenantID = req.TenantID
	}

	// DEBUG
	println("IDENTITY RECEIVED TENANT ID:", tenantID)

	tid, err := uuid.Parse(tenantID)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid tenant_id"})
		return
	}

	// =========================================
	// 🔍 FIND USER
	// =========================================
	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	// =========================================
	// 🔐 CHECK PASSWORD
	// =========================================
	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	// =========================================
	// 🔗 CHECK TENANT USER
	// =========================================
	var tenantUser models.TenantUser
	if err := h.DB.
		Where("tenant_id = ? AND user_id = ? AND status = ?", tid, user.ID, "active").
		First(&tenantUser).Error; err != nil {

		c.JSON(401, gin.H{"error": "user is not active in this tenant"})
		return
	}

	// =========================================
	// 🎭 LOAD ROLES
	// =========================================
	var roles []models.Role
	if err := h.DB.
		Table("roles").
		Select("roles.*").
		Joins("JOIN tenant_user_roles ON tenant_user_roles.role_id = roles.id").
		Where("tenant_user_roles.tenant_user_id = ?", tenantUser.ID).
		Find(&roles).Error; err != nil {

		c.JSON(500, gin.H{"error": "failed to load roles"})
		return
	}

	roleCodes := make([]string, 0, len(roles))
	for _, role := range roles {
		roleCodes = append(roleCodes, role.Code)
	}

	// =========================================
	// 🔐 GENERATE TOKEN
	// =========================================
	token, err := utils.GenerateToken(
		h.Config.JWTSecret,
		user.ID.String(),
		tenantUser.TenantID.String(),
		roleCodes,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}

	refresh := uuid.New().String()

	h.DB.Create(&models.RefreshToken{
		UserID:    user.ID,
		Token:     refresh,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	})

	// =========================================
	// ✅ RESPONSE
	// =========================================
	c.JSON(200, gin.H{
		"token":         token,
		"refresh_token": refresh,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"roles": roleCodes,
		},
	})
}

//////////////////////////////////////////////////////
// LOGOUT
//////////////////////////////////////////////////////

func (h *AuthHandler) Logout(c *gin.Context) {
	token := strings.TrimSpace(c.GetHeader("X-Refresh-Token"))

	// ❌ Missing token
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing refresh token"})
		return
	}

	// ❌ JWT passed instead of refresh token
	if strings.Contains(token, ".") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "use refresh token, not access token"})
		return
	}

	

	// 🗑 Delete token
	result := h.DB.Where("token = ?", token).Delete(&models.RefreshToken{})

	// ❌ Token not found
	if result.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or already logged out token"})
		return
	}

	// ❌ DB error
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
		return
	}

	// ✅ Success
	c.JSON(http.StatusOK, gin.H{
		"message": "logged out successfully",
	})
}