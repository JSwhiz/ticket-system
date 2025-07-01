package auth

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"
)

const (
    CtxUserIDKey = "user_id"
    CtxRoleIDKey = "role_id"
)

func JWTAuthMiddleware(secret string, db *sqlx.DB) gin.HandlerFunc {
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

        userID := claims.UserID

        var roleID string
        if err := db.Get(&roleID,
            "SELECT role_id::text FROM users WHERE user_id = $1", userID,
        ); err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
            return
        }

        c.Set(CtxUserIDKey, userID)
        c.Set(CtxRoleIDKey, roleID)
        c.Next()
    }
}
