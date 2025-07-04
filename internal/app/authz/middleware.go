package authz

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/JSwhiz/ticket-system/internal/app/auth"
)

func Authorize(repo *Repository, permName string) gin.HandlerFunc {
    return func(c *gin.Context) {
        v, exists := c.Get(auth.CtxRoleIDKey)
        if !exists {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"missing role in token"})
            return
        }
        roleID := v.(string)

        ok, err := repo.HasPermission(roleID, permName)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        if !ok {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error":"forbidden"})
            return
        }
        c.Next()
    }
}
