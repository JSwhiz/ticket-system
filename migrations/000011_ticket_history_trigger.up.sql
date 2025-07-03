CREATE OR REPLACE FUNCTION fn_ticket_history() RETURNS trigger AS $$
DECLARE
  v_changed_by UUID;
BEGIN
  -- Проверяем, установлена ли настройка app.current_user_id
  BEGIN
    v_changed_by := current_setting('app.current_user_id', true)::UUID;
  EXCEPTION WHEN OTHERS THEN
    v_changed_by := NULL; -- Если настройка отсутствует или некорректна, устанавливаем NULL
    RAISE NOTICE 'Failed to get app.current_user_id: %', SQLERRM;
  END;

  -- Логирование для отладки
  RAISE NOTICE 'Processing ticket history for ticket_id: %', OLD.ticket_id;

  IF OLD.title IS DISTINCT FROM NEW.title THEN
    INSERT INTO ticket_history(history_id, ticket_id, changed_by, field_name, old_value, new_value, changed_at)
    VALUES (uuid_generate_v4(), OLD.ticket_id, v_changed_by, 'title', OLD.title, NEW.title, now());
  END IF;

  IF OLD.description IS DISTINCT FROM NEW.description THEN
    INSERT INTO ticket_history(history_id, ticket_id, changed_by, field_name, old_value, new_value, changed_at)
    VALUES (uuid_generate_v4(), OLD.ticket_id, v_changed_by, 'description', OLD.description, NEW.description, now());
  END IF;

  IF OLD.status_id IS DISTINCT FROM NEW.status_id THEN
    INSERT INTO ticket_history(history_id, ticket_id, changed_by, field_name, old_value, new_value, changed_at)
    VALUES (uuid_generate_v4(), OLD.ticket_id, v_changed_by, 'status_id', OLD.status_id::text, NEW.status_id::text, now());
  END IF;

  IF OLD.priority_id IS DISTINCT FROM NEW.priority_id THEN
    INSERT INTO ticket_history(history_id, ticket_id, changed_by, field_name, old_value, new_value, changed_at)
    VALUES (uuid_generate_v4(), OLD.ticket_id, v_changed_by, 'priority_id', OLD.priority_id::text, NEW.priority_id::text, now());
  END IF;

  IF OLD.assignee_id IS DISTINCT FROM NEW.assignee_id THEN
    INSERT INTO ticket_history(history_id, ticket_id, changed_by, field_name, old_value, new_value, changed_at)
    VALUES (uuid_generate_v4(), OLD.ticket_id, v_changed_by, 'assignee_id', COALESCE(OLD.assignee_id::text, ''), COALESCE(NEW.assignee_id::text, ''), now());
  END IF;

  IF OLD.department_id IS DISTINCT FROM NEW.department_id THEN
    INSERT INTO ticket_history(history_id, ticket_id, changed_by, field_name, old_value, new_value, changed_at)
    VALUES (uuid_generate_v4(), OLD.ticket_id, v_changed_by, 'department_id', COALESCE(OLD.department_id::text, ''), COALESCE(NEW.department_id::text, ''), now());
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trg_ticket_history
  AFTER UPDATE ON tickets
  FOR EACH ROW
  EXECUTE FUNCTION fn_ticket_history();