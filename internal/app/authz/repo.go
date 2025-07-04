package authz

import "github.com/jmoiron/sqlx"

type Repository struct {
    db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
    return &Repository{db: db}
}

func (r *Repository) HasPermission(roleID, permName string) (bool, error) {
    var ok bool
    err := r.db.Get(&ok, `
      SELECT EXISTS(
        SELECT 1
          FROM role_permissions rp
          JOIN permissions p ON p.permission_id = rp.permission_id
         WHERE rp.role_id = $1 AND p.name = $2
      )`, roleID, permName)
    return ok, err
}
