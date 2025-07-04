package tickets

import (
    "database/sql"
    "time"

    "github.com/jmoiron/sqlx"
)

type Attachment struct {
    AttachmentID string         `db:"attachment_id" json:"attachment_id"`
    TicketID     string         `db:"ticket_id"     json:"ticket_id"`
    Filename     string         `db:"filename"      json:"filename"`
    // если храните бинарник в БД:
    FileData     []byte         `db:"file_data"     json:"-"`
    // если храните на диске — используем nullable string
    FilePath     sql.NullString `db:"file_path"     json:"file_path,omitempty"`
    UploadedBy   string         `db:"uploaded_by"   json:"uploaded_by"`
    UploadedAt   time.Time      `db:"uploaded_at"   json:"uploaded_at"`
}

type AttachmentRepo struct {
    db *sqlx.DB
}

func NewAttachmentRepo(db *sqlx.DB) *AttachmentRepo {
    return &AttachmentRepo{db: db}
}

func (r *AttachmentRepo) ListByTicket(ticketID string) ([]Attachment, error) {
    const q = `
      SELECT attachment_id, ticket_id, filename, file_path, uploaded_by, uploaded_at
        FROM ticket_attachments
       WHERE ticket_id = $1
       ORDER BY uploaded_at ASC
    `
    out := make([]Attachment, 0)
    if err := r.db.Select(&out, q, ticketID); err != nil {
        return nil, err
    }
    return out, nil
}

func (r *AttachmentRepo) Create(ticketID, userID, filename string, data []byte) (Attachment, error) {
    const q = `
      INSERT INTO ticket_attachments
        (ticket_id, filename, file_data, uploaded_by)
      VALUES ($1, $2, $3, $4)
      RETURNING attachment_id, ticket_id, filename, file_path, uploaded_by, uploaded_at
    `
    var a Attachment
    if err := r.db.Get(&a, q, ticketID, filename, data, userID); err != nil {
        return Attachment{}, err
    }
    return a, nil
}

func (r *AttachmentRepo) GetFile(attachmentID string) (string, []byte, error) {
    const q = `
      SELECT filename, file_data
        FROM ticket_attachments
       WHERE attachment_id = $1
    `
    row := r.db.QueryRowx(q, attachmentID)
    var fn string
    var data []byte
    if err := row.Scan(&fn, &data); err != nil {
        return "", nil, err
    }
    return fn, data, nil
}
