package tickets

import "github.com/jmoiron/sqlx"

type Service struct {
    repo        *Repository
    historyRepo *HistoryRepo
    commentRepo *CommentRepo
    attachRepo  *AttachmentRepo
}

func NewService(db *sqlx.DB) *Service {
    return &Service{
        repo:        NewRepository(db),
        historyRepo: NewHistoryRepo(db),
        commentRepo: NewCommentRepo(db),
        attachRepo:  NewAttachmentRepo(db),
    }
}

func (s *Service) List(filters map[string]interface{}) ([]Ticket, error) {
    return s.repo.List(filters)
}

func (s *Service) Create(n NewTicket, creatorID string) (Ticket, error) {
    return s.repo.Create(n, creatorID)
}

func (s *Service) Get(id string) (Ticket, error) {
    return s.repo.GetByID(id)
}

func (s *Service) Update(id string, u UpdateTicket, userID string) (Ticket, error) {
    return s.repo.UpdateWithUser(id, u, userID)
}

func (s *Service) Delete(id string) error {
    return s.repo.Delete(id)
}

func (s *Service) History(ticketID string) ([]HistoryEntry, error) {
    return s.historyRepo.ListByTicket(ticketID)
}

func (s *Service) Comments(ticketID string) ([]Comment, error) {
    return s.commentRepo.ListByTicket(ticketID)
}
func (s *Service) AddComment(ticketID, authorID, content string) (Comment, error) {
    return s.commentRepo.Create(ticketID, authorID, content)
}

func (s *Service) Attachments(ticketID string) ([]Attachment, error) {
    return s.attachRepo.ListByTicket(ticketID)
}

func (s *Service) UploadAttachment(ticketID, userID, filename string, data []byte) (Attachment, error) {
    return s.attachRepo.Create(ticketID, userID, filename, data)
}

func (s *Service) DownloadAttachment(attachmentID string) (string, []byte, error) {
    return s.attachRepo.GetFile(attachmentID)
}