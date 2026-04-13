package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type TenantLookupResponse struct {
	Tenant struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Slug   string `json:"slug"`
		Status string `json:"status"`
	} `json:"tenant"`
}

func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		domain := strings.ToLower(strings.TrimSpace(c.GetHeader("X-Tenant-Domain")))

		if domain == "" {
			host := c.Request.Host
			if strings.Contains(host, ":") {
				host = strings.Split(host, ":")[0]
			}
			domain = strings.ToLower(strings.TrimSpace(host))
		}

		if domain == "" || domain == "localhost" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unable to resolve tenant domain"})
			return
		}

		url := fmt.Sprintf("http://tenant-service:8081/tenants/by-domain?domain=%s", domain)
		resp, err := http.Get(url)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "tenant service unavailable"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
			return
		}

		var tenantResp TenantLookupResponse
		if err := json.NewDecoder(resp.Body).Decode(&tenantResp); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to parse tenant response"})
			return
		}

		if tenantResp.Tenant.Status != "active" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "tenant is not active"})
			return
		}

		c.Set("tenant_id", tenantResp.Tenant.ID)
		c.Set("tenant_slug", tenantResp.Tenant.Slug)
		c.Next()
	}
}