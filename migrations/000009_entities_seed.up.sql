INSERT INTO roles(name) VALUES
    ('admin'),
    ('user'),
    ('manager'),
    ('support'),
    ('developer');

INSERT INTO permissions(name) VALUES
    ('create_ticket'),
    ('view_ticket'),
    ('update_ticket'),
    ('delete_ticket'),
    ('assign_ticket');

INSERT INTO role_permissions(role_id, permission_id)
    SELECT r.role_id, p.permission_id
        FROM roles r
        JOIN permissions p ON true
    WHERE r.name = 'admin';

INSERT INTO departments(name) VALUES
    ('Support'),
    ('Development'),
    ('Sales'),
    ('HR'),
    ('Marketing');

INSERT INTO ticket_statuses(code, label) VALUES
    ('open',         'Open'),
    ('in_progress',  'In Progress'),
    ('resolved',     'Resolved'),
    ('closed',       'Closed'),
    ('reopened',     'Reopened');

INSERT INTO ticket_priorities(code, label, level) VALUES
    ('low',      'Low',      1),
    ('medium',   'Medium',   2),
    ('high',     'High',     3),
    ('critical', 'Critical', 4),
    ('urgent',   'Urgent',   5);
