CREATE INDEX idx_tickets_status     ON tickets(status_id);
CREATE INDEX idx_tickets_priority   ON tickets(priority_id);
CREATE INDEX idx_tickets_assignee   ON tickets(assignee_id);
CREATE INDEX idx_tickets_department ON tickets(department_id);

CREATE INDEX idx_tickets_assignee_status_deleted
  ON tickets(assignee_id, status_id)
  WHERE deleted_at IS NULL;

CREATE INDEX idx_tickets_created_status
  ON tickets(created_at, status_id)
  WHERE deleted_at IS NULL;

ALTER TABLE tickets ADD COLUMN IF NOT EXISTS search_vector tsvector;
CREATE INDEX idx_tickets_search ON tickets USING GIN(search_vector);

CREATE TRIGGER trg_tickets_search
  BEFORE INSERT OR UPDATE ON tickets
  FOR EACH ROW EXECUTE FUNCTION
    tsvector_update_trigger('search_vector', 'pg_catalog.english', 'title', 'description');

CREATE FUNCTION fn_update_timestamp() RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_tickets
  BEFORE UPDATE ON tickets
  FOR EACH ROW EXECUTE FUNCTION fn_update_timestamp();