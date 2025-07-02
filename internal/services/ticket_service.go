package services

import (
	"ticket-system/internal/models"
	"ticket-system/internal/repository"
	"time"

	"github.com/google/uuid"
)

type TicketService struct {
	ticketRepo *repository.TicketRepository
}

func NewTicketService(ticketRepo *repository.TicketRepository) *TicketService {
	return &TicketService{ticketRepo}
}

func (s *TicketService) CreateTicket(title, description string, statusID, priorityID int16, creatorID uuid.UUID, departmentID *int16) (*models.Ticket, error) {
	ticket := &models.Ticket{
		TicketID:     uuid.New(),
		Title:        title,
		Description:  &description, // Указатель на строку
		StatusID:     statusID,
		PriorityID:   priorityID,
		CreatorID:    creatorID,
		DepartmentID: departmentID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		DeletedAt:    nil,
		SearchVector: nil,
	}
	err := s.ticketRepo.CreateTicket(ticket)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}
func (s *TicketService) GetAllTickets() ([]*models.Ticket, error) {
	return s.ticketRepo.GetAllTickets()
}

func (s *TicketService) GetTicket(ticketID string) (*models.Ticket, error) {
	return s.ticketRepo.GetTicketByID(ticketID)
}

func (s *TicketService) UpdateTicket(ticketID string, updates map[string]interface{}) error {
	return s.ticketRepo.UpdateTicket(ticketID, updates)
}

func (s *TicketService) DeleteTicket(ticketID string) error {
	return s.ticketRepo.DeleteTicket(ticketID)
}
