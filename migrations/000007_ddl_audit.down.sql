-- ================================================
-- 0006_ticket_history_triggers.down.sql
-- Откат row-trigger’а истории
-- ================================================
DROP TRIGGER IF EXISTS trg_ticket_history ON tickets;
DROP FUNCTION IF EXISTS fn_ticket_history();
