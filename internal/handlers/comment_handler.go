package handlers

import (
    "net/http"
    "ticket-system/internal/services"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type CommentHandler struct {
    commentService *services.CommentService
}

func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
    return &CommentHandler{commentService}
}


func (h *CommentHandler) GetComments(c *gin.Context) {
    ticketIDStr := c.Param("id")
    ticketID, err := uuid.Parse(ticketIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID тикета"})
        return
    }

    comments, err := h.commentService.GetCommentsByTicketID(ticketID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении комментариев"})
        return
    }

    c.JSON(http.StatusOK, comments)
}


func (h *CommentHandler) CreateComment(c *gin.Context) {
    ticketIDStr := c.Param("id")
    ticketID, err := uuid.Parse(ticketIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID тикета"})
        return
    }

    var dto struct {
        Content string `json:"content" binding:"required"`
    }
    if err := c.ShouldBindJSON(&dto); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
        return
    }

    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не аутентифицирован"})
        return
    }

    comment, err := h.commentService.CreateComment(ticketID, dto.Content, userID.(uuid.UUID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании комментария"})
        return
    }

    c.JSON(http.StatusCreated, comment)
}