package tickets

import "github.com/jmoiron/sqlx"

type Service struct {
    repo *Repository
}

func NewService(db *sqlx.DB) *Service {
    return &Service{repo: NewRepository(db)}
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

func (s *Service) Update(id string, u UpdateTicket) (Ticket, error) {
    return s.repo.Update(id, u)
}

func (s *Service) Delete(id string) error {
    return s.repo.Delete(id)
}
