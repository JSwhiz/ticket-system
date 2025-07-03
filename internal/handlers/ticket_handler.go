package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"ticket-system/internal/models"
	"ticket-system/internal/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TicketHandler struct {
	ticketService *services.TicketService
	authService   *services.AuthService
}

func NewTicketHandler(ticketService *services.TicketService, authService *services.AuthService) *TicketHandler {
	return &TicketHandler{ticketService, authService}
}

func (h *TicketHandler) CreateTicket(c *gin.Context) {
	var req struct {
		Title        string `json:"title" binding:"required"`
		Description  string `json:"description" binding:"required"`
		StatusID     int16  `json:"status_id" binding:"required"`
		PriorityID   int16  `json:"priority_id" binding:"required"`
		DepartmentID *int16 `json:"department_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    http.StatusBadRequest,
				"message": "Validation failed",
				"details": err.Error(),
			},
		})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		log.Println("JWT token missing user_id")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized: JWT token invalid or missing",
			},
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		log.Printf("Failed to parse user_id from JWT: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Invalid user_id in JWT token",
			},
		})
		return
	}

	ticket, err := h.ticketService.CreateTicket(req.Title, req.Description, req.StatusID, req.PriorityID, userID, req.DepartmentID)
	if err != nil {
		log.Printf("Failed to create ticket: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Failed to create ticket",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

func (h *TicketHandler) GetAllTickets(c *gin.Context) {
	tickets, err := h.ticketService.GetAllTickets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tickets)
}

func (h *TicketHandler) GetTicket(c *gin.Context) {
	ticketID := c.Param("id")
	ticket, err := h.ticketService.GetTicket(ticketID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}
	c.JSON(http.StatusOK, ticket)
}

func (h *TicketHandler) UpdateTicket(c *gin.Context) {
	ticketID := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.ticketService.UpdateTicket(ticketID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket updated successfully"})
}

func (h *TicketHandler) DeleteTicket(c *gin.Context) {
	ticketID := c.Param("id")

	err := h.ticketService.DeleteTicket(ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted successfully"})
}

func (h *TicketHandler) GetAttachments(c *gin.Context) {
	ticketID := c.Param("id")
	if ticketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket ID is required"})
		return
	}

	attachments, err := h.ticketService.GetAttachmentsByTicketID(ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attachments)
}

func (h *TicketHandler) UploadAttachment(c *gin.Context) {
	ticketID := c.Param("id")
	if ticketID == "" {
		log.Printf("Ticket ID is missing")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket ID is required"})
		return
	}

	log.Printf("Processing upload for ticket ID: %s", ticketID)

	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("Failed to get file from form: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	log.Printf("File received: %s, size: %d bytes", file.Filename, file.Size)

	if file.Size > 10*1024*1024 { // Ограничение 10MB
		log.Printf("File too large: %s, size: %d bytes", file.Filename, file.Size)
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large"})
		return
	}

	filename := fmt.Sprintf("%s_%s", uuid.New().String(), file.Filename)
	uploadDir := "./Uploads"
	log.Printf("Attempting to create upload directory: %s", uploadDir)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Printf("Failed to create upload directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}
	log.Printf("Upload directory created successfully: %s", uploadDir)

	filePath := filepath.Join(uploadDir, filename)
	log.Printf("Saving file to: %s", filePath)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		log.Printf("Failed to save file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	log.Printf("File saved successfully: %s", filePath)

	userIDStr := c.GetString("user_id") // Полагаемся на middleware
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("Failed to parse user_id from context: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user_id format"})
		return
	}

	attachment := &models.TicketAttachment{
		AttachmentID: uuid.New(),
		TicketID:     uuid.MustParse(ticketID),
		Filename:     filename,
		FilePath:     filePath,
		UploadedBy:   userID,
		UploadedAt:   time.Now(),
	}
	err = h.ticketService.CreateAttachment(attachment)
	if err != nil {
		log.Printf("Failed to create attachment in database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, attachment)
}

func (h *TicketHandler) DownloadAttachment(c *gin.Context) {
	attachmentID := c.Param("att_id")
	if attachmentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Attachment ID is required"})
		return
	}

	attachment, err := h.ticketService.GetAttachmentByID(attachmentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attachment not found"})
		return
	}

	c.FileAttachment(attachment.FilePath, attachment.Filename)
}

func (h *TicketHandler) DeleteAttachment(c *gin.Context) {
	attachmentID := c.Param("att_id")
	if attachmentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Attachment ID is required"})
		return
	}

	err := h.ticketService.DeleteAttachment(attachmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
