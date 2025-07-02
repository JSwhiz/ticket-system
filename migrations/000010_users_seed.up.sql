DO $$
DECLARE
  pw_hash TEXT := '$2a$10$7d9L6EoVgT8uEwJ1pQhDzeK2h3XjF4aY6BzoW2EKjUV5RmZ8YtW3O';
BEGIN
  INSERT INTO users (
    user_id,
    username,
    password_hash,
    email,
    role_id,
    department_id
  ) VALUES
    (
      uuid_generate_v4(),
      'admin1',
      pw_hash,
      'admin1@example.com',
      (SELECT role_id FROM roles WHERE name='admin'),
      (SELECT department_id FROM departments WHERE name='Support')
    ),
    (
      uuid_generate_v4(),
      'user1',
      pw_hash,
      'user1@example.com',
      (SELECT role_id FROM roles WHERE name='user'),
      (SELECT department_id FROM departments WHERE name='Sales')
    ),
    (
      uuid_generate_v4(),
      'manager1',
      pw_hash,
      'manager1@example.com',
      (SELECT role_id FROM roles WHERE name='manager'),
      (SELECT department_id FROM departments WHERE name='Development')
    ),
    (
      uuid_generate_v4(),
      'support1',
      pw_hash,
      'support1@example.com',
      (SELECT role_id FROM roles WHERE name='support'),
      (SELECT department_id FROM departments WHERE name='Support')
    ),
    (
      uuid_generate_v4(),
      'developer1',
      pw_hash,
      'developer1@example.com',
      (SELECT role_id FROM roles WHERE name='developer'),
      (SELECT department_id FROM departments WHERE name='Development')
    )
  ON CONFLICT (username) DO NOTHING;
END
$$;
