CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    level INTEGER NOT NULL DEFAULT 0
);

INSERT INTO roles (name, description, level) VALUES 
(
  'user',
  'A regular user can create posts and comments',
  1
),
(
  'moderator',
  'A moderator can update other users posts and comments',
  2
),
(
  'admin',
  'An admin can update and delete other users posts and comments',
  3
);


