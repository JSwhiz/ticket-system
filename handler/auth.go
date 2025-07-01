package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/JSwhiz/ticket-system/internal/app/auth"
    "github.com/JSwhiz/ticket-system/pkg/config"
)

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
    Token string `json:"token"`
    User  struct {
        UserID       string `json:"user_id"`
        Username     string `json:"username"`
        RoleID       string `json:"role_id"`
        DepartmentID int    `json:"department_id"`
    } `json:"user"`
}

func RegisterAuthRoutes(rg *gin.RouterGroup, svc *auth.Service, cfg *config.Config) {
    rg.POST("/login", func(c *gin.Context) {
        var req LoginRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        token, user, err := svc.Login(req.Username, req.Password)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
            return
        }

        var resp LoginResponse
        resp.Token = token
        resp.User.UserID = user.ID
        resp.User.Username = user.Username
        resp.User.RoleID = user.RoleID
        resp.User.DepartmentID = user.DepartmentID

        c.JSON(http.StatusOK, resp)
    })
}
