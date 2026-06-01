package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	// ContextKeyUserID is the key for user ID in gin context.
	ContextKeyUserID = "user_id"
	// ContextKeyUserRole is the key for user role in gin context.
	ContextKeyUserRole = "user_role"
	// ContextKeySessionID is the key for session ID in gin context.
	ContextKeySessionID = "session_id"
)

// AuthMiddleware validates JWT access tokens.
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "missing authorization header",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "invalid authorization header format",
			})
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "invalid or expired token",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "invalid token claims",
			})
			return
		}

		userID, _ := claims["sub"].(string)
		role, _ := claims["role"].(string)
		sessionID, _ := claims["sid"].(string)

		c.Set(ContextKeyUserID, userID)
		c.Set(ContextKeyUserRole, role)
		c.Set(ContextKeySessionID, sessionID)

		c.Next()
	}
}

// RoleMiddleware restricts access to specific roles.
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextKeyUserRole)
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "access denied",
			})
			return
		}

		userRole := role.(string)
		for _, allowed := range allowedRoles {
			if userRole == allowed {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "insufficient permissions",
		})
	}
}

// OptionalAuth extracts user info from token if present, but doesn't require it.
func OptionalAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next()
			return
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				userID, _ := claims["sub"].(string)
				role, _ := claims["role"].(string)
				c.Set(ContextKeyUserID, userID)
				c.Set(ContextKeyUserRole, role)
			}
		}

		c.Next()
	}
}

// GetUserID extracts user ID from gin context.
func GetUserID(c *gin.Context) string {
	id, _ := c.Get(ContextKeyUserID)
	if id == nil {
		return ""
	}
	return id.(string)
}

// GetUserRole extracts user role from gin context.
func GetUserRole(c *gin.Context) string {
	role, _ := c.Get(ContextKeyUserRole)
	if role == nil {
		return ""
	}
	return role.(string)
}
