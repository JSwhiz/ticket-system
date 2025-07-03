package services

import (
    "ticket-system/internal/models"
    "ticket-system/internal/repository"
    "github.com/google/uuid"
    "time"
)

type CommentService struct {
    commentRepo *repository.CommentRepository
}

func NewCommentService(commentRepo *repository.CommentRepository) *CommentService {
    return &CommentService{commentRepo}
}
func (s *CommentService) GetCommentsByTicketID(ticketID uuid.UUID) ([]*models.Comment, error) {
    return s.commentRepo.GetCommentsByTicketID(ticketID)
}


func (s *CommentService) CreateComment(ticketID uuid.UUID, content string, authorID uuid.UUID) (*models.Comment, error) {
    comment := &models.Comment{
        ID:        uuid.New(),
        TicketID:  ticketID,
        AuthorID:  authorID,
        Content:   content,
        CreatedAt: time.Now(),
    }
    err := s.commentRepo.CreateComment(comment)
    if err != nil {
        return nil, err
    }
    return comment, nil
}