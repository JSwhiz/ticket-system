package models

import (
	"time"
	"github.com/google/uuid"
)

type TicketHistory struct {
	HistoryID  uuid.UUID              `db:"history_id" json:"history_id"`
	TicketID   uuid.UUID              `db:"ticket_id" json:"ticket_id"`
	ChangedAt  time.Time              `db:"changed_at" json:"changed_at"`
	ChangedBy  *uuid.UUID             `db:"changed_by" json:"changed_by"`
	FieldName  string                 `db:"field_name" json:"field_name"`
	OldValue   *string                `db:"old_value" json:"old_value"`
	NewValue   *string                `db:"new_value" json:"new_value"`
	ChangeData map[string]interface{} `db:"change_data" json:"change_data"`
}