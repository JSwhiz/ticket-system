package handlers

import (
	"net/http"
	"ticket-system/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TicketHandler struct {
	ticketService *services.TicketService
	authService   *services.AuthService // Для проверки токена
}

func NewTicketHandler(ticketService *services.TicketService, authService *services.AuthService) *TicketHandler {
	return &TicketHandler{ticketService, authService}
}

func (h *TicketHandler) CreateTicket(c *gin.Context) {
	var req struct {
		Title       string  `json:"title" binding:"required"`
		Description string  `json:"description" binding:"required"`
		StatusID    int16   `json:"status_id" binding:"required"`
		PriorityID  int16   `json:"priority_id" binding:"required"`
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

	// Извлечение user_id из JWT (временная заглушка)
	userID, _ := uuid.Parse("a2657c63-77cb-4c3c-8577-27c19f017a65") // Замени на валидацию токена

	ticket, err := h.ticketService.CreateTicket(req.Title, req.Description, req.StatusID, req.PriorityID, userID, req.DepartmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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