CREATE TABLE ticket_history (
  history_id  UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
  ticket_id   UUID        NOT NULL REFERENCES tickets(ticket_id) ON DELETE CASCADE,
  changed_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  changed_by  UUID        REFERENCES users(user_id),
  field_name  TEXT        NOT NULL,
  old_value   TEXT,
  new_value   TEXT,
  change_data JSONB
);

CREATE TABLE ticket_comments (
  comment_id  UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
  ticket_id   UUID        NOT NULL REFERENCES tickets(ticket_id) ON DELETE CASCADE,
  author_id   UUID        NOT NULL REFERENCES users(user_id),
  content     TEXT        NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE ticket_attachments (
  attachment_id UUID      PRIMARY KEY DEFAULT uuid_generate_v4(),
  ticket_id     UUID      NOT NULL REFERENCES tickets(ticket_id) ON DELETE CASCADE,
  filename      TEXT      NOT NULL,
  metadata      JSONB,
  file_data     BYTEA,
  file_path     TEXT,
  uploaded_by   UUID      NOT NULL REFERENCES users(user_id),
  uploaded_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  CHECK (
    (file_data IS NOT NULL)::int + (file_path IS NOT NULL)::int = 1
  )
);