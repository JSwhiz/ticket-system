INSERT INTO roles(name) VALUES ('admin'), ('user');
INSERT INTO permissions(name) VALUES ('create_ticket'), ('close_ticket'), ('view_reports');
INSERT INTO role_permissions(role_id, permission_id)
  SELECT r.role_id, p.permission_id
  FROM roles r CROSS JOIN permissions p WHERE r.name='admin';
INSERT INTO departments(name) VALUES ('Support'), ('Development'), ('Sales');
INSERT INTO ticket_statuses(code, label)
  VALUES ('open','Open'),('in_progress','In Progress'),('resolved','Resolved'),('closed','Closed');
INSERT INTO ticket_priorities(code, label, level)
  VALUES ('low','Low',1),('medium','Medium',2),('high','High',3),('critical','Critical',5);