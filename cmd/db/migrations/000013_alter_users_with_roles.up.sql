ALTER TABLE users
ADD COLUMN role_id BIGINT REFERENCES roles(id) ON DELETE RESTRICT;

UPDATE users SET role_id = (SELECT id FROM roles WHERE name = 'user')
WHERE role_id IS NULL;

ALTER TABLE users
ALTER COLUMN role_id SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id);
