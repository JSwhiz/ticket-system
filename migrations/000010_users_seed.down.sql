DELETE FROM users
    WHERE username IN (
        'admin1',
        'user1',
        'manager1',
        'support1',
        'developer1'
    );
