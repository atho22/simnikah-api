package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"simnikah/pkg/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token and extracts user info
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
				"error":   "Token otorisasi tidak disediakan",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
				"error":   "Token tidak valid atau kedaluwarsa",
			})
			return
		}

		// Validate role
		validRoles := map[string]bool{
			"user_biasa": true,
			"penghulu":   true,
			"staff":      true,
			"kepala_kua": true,
		}
		if !validRoles[claims.Role] {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
				"error":   "Role tidak valid",
			})
			return
		}

		// Set claims to context
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("email", claims.Email)
		c.Set("nama", claims.Nama)

		c.Next()
	}
}

// RoleMiddleware validates specific role
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
				"error":   "Role tidak ditemukan",
			})
			c.Abort()
			return
		}

		if role.(string) != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Forbidden",
				"error":   fmt.Sprintf("Akses ditolak. Role %s diperlukan", requiredRole),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// MultiRoleMiddleware validates multiple roles
func MultiRoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
				"error":   "Role tidak ditemukan",
			})
			c.Abort()
			return
		}

		userRole := role.(string)
		hasAccess := false
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Forbidden",
				"error":   fmt.Sprintf("Akses ditolak. Role yang diizinkan: %s", strings.Join(allowedRoles, ", ")),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

