ALTER TABLE users RENAME CONSTRAINT uq_users_email TO users_email_key;
ALTER TABLE users RENAME CONSTRAINT uq_users_username TO users_username_key;