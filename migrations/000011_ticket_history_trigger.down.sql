-- ================================================
-- 0007_ddl_audit.down.sql
-- Откат DDL-аудита
-- ================================================
DROP EVENT TRIGGER IF EXISTS trg_ddl_audit;
DROP FUNCTION IF EXISTS fn_ddl_audit();
DROP TABLE IF EXISTS ddl_audit;
