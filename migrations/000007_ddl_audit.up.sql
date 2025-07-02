CREATE TABLE ddl_audit (
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
  FOR rec IN SELECT * FROM pg_event_trigger_ddl_commands() LOOP
    INSERT INTO ddl_audit(
      user_name, command_tag, schema_name, object_identity, command_text
    ) VALUES (
      session_user,
      rec.command_tag,
      rec.schema_name,
      rec.object_identity,
      rec.command
    );
  END LOOP;
END;
$$ LANGUAGE plpgsql;

CREATE EVENT TRIGGER trg_ddl_audit
  ON ddl_command_end
  EXECUTE FUNCTION fn_ddl_audit();