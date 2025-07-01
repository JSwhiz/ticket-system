CREATE MATERIALIZED VIEW mv_departments AS
  SELECT department_id AS id, name FROM departments;

CREATE MATERIALIZED VIEW mv_ticket_statuses AS
  SELECT status_id AS id, label FROM ticket_statuses;

CREATE MATERIALIZED VIEW mv_ticket_priorities AS
  SELECT priority_id AS id, label FROM ticket_priorities;