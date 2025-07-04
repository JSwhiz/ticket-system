-- ================================================
-- 0006_ticket_history_triggers.up.sql
-- row-trigger для аудита изменений в tickets
-- ================================================
CREATE OR REPLACE FUNCTION fn_ticket_history() RETURNS trigger AS $$
DECLARE
  changed_by UUID := current_setting('app.current_user_id', true)::UUID;
BEGIN
  IF OLD.title IS DISTINCT FROM NEW.title THEN
    INSERT INTO ticket_history(
      ticket_id, changed_by, field_name, old_value, new_value
    ) VALUES (
      OLD.ticket_id, changed_by, 'title',
      OLD.title::text, NEW.title::text
    );
  END IF;

  IF OLD.description IS DISTINCT FROM NEW.description THEN
    INSERT INTO ticket_history(
      ticket_id, changed_by, field_name, old_value, new_value
    ) VALUES (
      OLD.ticket_id, changed_by, 'description',
      OLD.description, NEW.description
    );
  END IF;

  IF OLD.status_id IS DISTINCT FROM NEW.status_id THEN
    INSERT INTO ticket_history(
      ticket_id, changed_by, field_name, old_value, new_value
    ) VALUES (
      OLD.ticket_id, changed_by, 'status_id',
      OLD.status_id::text, NEW.status_id::text
    );
  END IF;

  IF OLD.priority_id IS DISTINCT FROM NEW.priority_id THEN
    INSERT INTO ticket_history(
      ticket_id, changed_by, field_name, old_value, new_value
    ) VALUES (
      OLD.ticket_id, changed_by, 'priority_id',
      OLD.priority_id::text, NEW.priority_id::text
    );
  END IF;

  IF OLD.assignee_id IS DISTINCT FROM NEW.assignee_id THEN
    INSERT INTO ticket_history(
      ticket_id, changed_by, field_name, old_value, new_value
    ) VALUES (
      OLD.ticket_id, changed_by, 'assignee_id',
      OLD.assignee_id::text, NEW.assignee_id::text
    );
  END IF;

  IF OLD.department_id IS DISTINCT FROM NEW.department_id THEN
    INSERT INTO ticket_history(
      ticket_id, changed_by, field_name, old_value, new_value
    ) VALUES (
      OLD.ticket_id, changed_by, 'department_id',
      OLD.department_id::text, NEW.department_id::text
    );
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_ticket_history ON tickets;
CREATE TRIGGER trg_ticket_history
  AFTER UPDATE ON tickets
  FOR EACH ROW
  EXECUTE FUNCTION fn_ticket_history();
