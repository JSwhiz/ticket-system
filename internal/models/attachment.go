package models

import (
	"time"
	"github.com/google/uuid"
)

type TicketAttachment struct {
	AttachmentID uuid.UUID  `json:"attachment_id" db:"attachment_id"`
	TicketID     uuid.UUID  `json:"ticket_id" db:"ticket_id"`
	Filename     string     `json:"filename" db:"filename"`
	FilePath     string     `json:"file_path" db:"file_path"`
	UploadedBy   uuid.UUID  `json:"uploaded_by" db:"uploaded_by"`
	UploadedAt   time.Time  `json:"uploaded_at" db:"uploaded_at"`
}