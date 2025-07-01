DELETE FROM ticket_priorities
    WHERE code IN ('low','medium','high','critical','urgent');

DELETE FROM ticket_statuses
    WHERE code IN ('open','in_progress','resolved','closed','reopened');

DELETE FROM departments
    WHERE name IN ('Support','Development','Sales','HR','Marketing');

DELETE FROM role_permissions rp
    USING roles r, permissions p
    WHERE rp.role_id = r.role_id
        AND rp.permission_id = p.permission_id
        AND r.name = 'admin'
        AND p.name IN ('create_ticket','view_ticket','update_ticket','delete_ticket','assign_ticket');

DELETE FROM permissions
    WHERE name IN ('create_ticket','view_ticket','update_ticket','delete_ticket','assign_ticket');

DELETE FROM roles
    WHERE name IN ('admin','user','manager','support','developer');
