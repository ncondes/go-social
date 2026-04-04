ALTER TABLE users RENAME CONSTRAINT users_email_key TO uq_users_email;
ALTER TABLE users RENAME CONSTRAINT users_username_key TO uq_users_username;