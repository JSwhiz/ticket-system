CREATE TABLE roles (
  role_id   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name      TEXT NOT NULL UNIQUE
);
CREATE TABLE permissions (
  permission_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name          TEXT NOT NULL UNIQUE
);
CREATE TABLE role_permissions (
  role_id       UUID NOT NULL REFERENCES roles(role_id) ON DELETE CASCADE,
  permission_id UUID NOT NULL REFERENCES permissions(permission_id) ON DELETE CASCADE,
  PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE departments (
  department_id SMALLSERIAL PRIMARY KEY,
  name          TEXT NOT NULL UNIQUE
);

CREATE TABLE ticket_statuses (
  status_id SMALLSERIAL PRIMARY KEY,
  code      TEXT NOT NULL UNIQUE,
  label     TEXT NOT NULL
);

CREATE TABLE ticket_priorities (
  priority_id SMALLSERIAL PRIMARY KEY,
  code        TEXT NOT NULL UNIQUE,
  label       TEXT NOT NULL,
  level       SMALLINT NOT NULL CHECK (level BETWEEN 1 AND 5)
);