-- ================================================
-- 0007_ddl_audit.up.sql
-- event-trigger для аудита любых DDL-команд
-- ================================================
CREATE TABLE IF NOT EXISTS ddl_audit (
  event_id        BIGSERIAL PRIMARY KEY,
  event_time      TIMESTAMPTZ NOT NULL DEFAULT now(),
  user_name       TEXT,
  command_tag     TEXT,
  schema_name     TEXT,
  object_identity TEXT,
  command_text    TEXT
);

CREATE OR REPLACE FUNCTION fn_ddl_audit() RETURNS event_trigger AS $$
DECLARE
  rec record;
BEGIN
  FOR rec IN
    SELECT
      tg.event_tag                     AS command_tag,
      n.nspname                        AS schema_name,
      (tg.object_id::regclass)::text   AS object_identity,
      pg_get_ddl_command(
        tg.object_id, tg.event_tag
      )                                AS command_text
    FROM pg_event_trigger_ddl_commands() tg
    LEFT JOIN pg_class cls
      ON cls.oid = tg.object_id
    LEFT JOIN pg_namespace n
      ON cls.relnamespace = n.oid
  LOOP
    INSERT INTO ddl_audit(
      user_name, command_tag,
      schema_name, object_identity, command_text
    ) VALUES (
      session_user,
      rec.command_tag,
      rec.schema_name,
      rec.object_identity,
      rec.command_text
    );
  END LOOP;
END;
$$ LANGUAGE plpgsql;

DROP EVENT TRIGGER IF EXISTS trg_ddl_audit;
CREATE EVENT TRIGGER trg_ddl_audit
  ON ddl_command_end
  EXECUTE FUNCTION fn_ddl_audit();
