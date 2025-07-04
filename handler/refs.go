package handler

import (
    "net/http"

    "github.com/JSwhiz/ticket-system/internal/app/refs"
    "github.com/gin-gonic/gin"
)

func RegisterRefRoutes(rg *gin.RouterGroup, svc *refs.Service) {
    r := rg.Group("/refs")

    r.GET("/departments", func(c *gin.Context) {
        out, err := svc.GetDepartments()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, out)
    })

    r.GET("/statuses", func(c *gin.Context) {
        out, err := svc.GetStatuses()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, out)
    })

    r.GET("/priorities", func(c *gin.Context) {
        out, err := svc.GetPriorities()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, out)
    })

    r.GET("/roles", func(c *gin.Context) {
        out, err := svc.GetRoles()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, out)
    })

    r.GET("/permissions", func(c *gin.Context) {
        out, err := svc.GetPermissions()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, out)
    })
}
