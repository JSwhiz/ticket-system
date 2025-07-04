package tickets
import (
	"time"
	"github.com/jmoiron/sqlx"
)

type HistoryEntry struct {
    ID        string    `db:"history_id"  json:"history_id"`
    TicketID  string    `db:"ticket_id"   json:"ticket_id"`
    ChangedBy *string   `db:"changed_by"  json:"changed_by,omitempty"`
    FieldName string    `db:"field_name"  json:"field_name"`
    OldValue  *string   `db:"old_value"   json:"old_value,omitempty"`
    NewValue  *string   `db:"new_value"   json:"new_value,omitempty"`
    When      time.Time `db:"changed_at"  json:"changed_at"`
}

type HistoryRepo struct { db *sqlx.DB }
func NewHistoryRepo(db *sqlx.DB) *HistoryRepo { return &HistoryRepo{db} }

func (r *HistoryRepo) ListByTicket(ticketID string) ([]HistoryEntry, error) {
    var out []HistoryEntry
    err := r.db.Select(&out, `
      SELECT history_id, ticket_id, changed_by, field_name,
             old_value, new_value, changed_at
      FROM ticket_history
      WHERE ticket_id = $1
      ORDER BY changed_at ASC`, ticketID)
    return out, err
}
