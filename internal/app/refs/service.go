package refs

type Service struct {
    repo *Repository
}

func NewService(repo *Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) GetDepartments() ([]Department, error) {
    return s.repo.Departments()
}

func (s *Service) GetStatuses() ([]Status, error) {
    return s.repo.Statuses()
}

func (s *Service) GetPriorities() ([]Priority, error) {
    return s.repo.Priorities()
}

func (s *Service) GetRoles() ([]Role, error) {
    return s.repo.Roles()
}

func (s *Service) GetPermissions() ([]Permission, error) {
    return s.repo.Permissions()
}
