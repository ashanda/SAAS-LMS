package middleware

import (

	"strings"
	"github.com/gin-gonic/gin"

)

func ForwardIdentityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		if userID, ok := c.Get("user_id"); ok {
			c.Request.Header.Set("X-Internal-User-ID", userID.(string))
		}
		if tenantID, ok := c.Get("tenant_id"); ok {
			c.Request.Header.Set("X-Internal-Tenant-ID", tenantID.(string))
		}
		if roles, ok := c.Get("roles"); ok {
			if roleList, ok := roles.([]string); ok {
				c.Request.Header.Set("X-Internal-Roles", strings.Join(roleList, ","))
			}
		}
		c.Next()
	}
}