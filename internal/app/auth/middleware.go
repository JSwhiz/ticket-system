package auth

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
)

const (
    CtxUserIDKey = "user_id"
    CtxRoleIDKey = "role_id"
)

func JWTAuthMiddleware(secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid Authorization header"})
            return
        }
        tokenStr := parts[1]

        claims, err := ParseToken(secret, tokenStr)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }

        c.Set(CtxUserIDKey, claims.UserID)
        c.Set(CtxRoleIDKey, claims.RoleID)
        c.Next()
    }
}
