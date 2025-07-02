-- migrations/0009_entities_seed.up.sql

-- 1) Роли
INSERT INTO roles(name) VALUES
  ('admin'),
  ('user'),
  ('manager'),
  ('support'),
  ('developer')
ON CONFLICT (name) DO NOTHING;

-- 2) Права
INSERT INTO permissions(name) VALUES
  ('create_ticket'),
  ('view_ticket'),
  ('update_ticket'),
  ('delete_ticket'),
  ('assign_ticket')
ON CONFLICT (name) DO NOTHING;

-- 3) Связь ролей и прав
INSERT INTO role_permissions(role_id, permission_id)
SELECT r.role_id, p.permission_id
  FROM roles r
  JOIN permissions p ON true
 WHERE r.name = 'admin'
ON CONFLICT DO NOTHING;

-- 4) Отделы
INSERT INTO departments(name) VALUES
  ('Support'),
  ('Development'),
  ('Sales'),
  ('HR'),
  ('Marketing')
ON CONFLICT (name) DO NOTHING;

-- 5) Статусы тикетов
INSERT INTO ticket_statuses(code, label) VALUES
  ('open',        'Open'),
  ('in_progress', 'In Progress'),
  ('resolved',    'Resolved'),
  ('closed',      'Closed'),
  ('reopened',    'Reopened')
ON CONFLICT (code) DO NOTHING;

-- 6) Приоритеты тикетов
INSERT INTO ticket_priorities(code, label, level) VALUES
  ('low',      'Low',      1),
  ('medium',   'Medium',   2),
  ('high',     'High',     3),
  ('critical', 'Critical', 4),
  ('urgent',   'Urgent',   5)
ON CONFLICT (code) DO NOTHING;
