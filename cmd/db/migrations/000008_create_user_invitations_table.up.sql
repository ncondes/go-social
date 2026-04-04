CREATE TABLE IF NOT EXISTS user_invitations (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  method TEXT NOT NULL CHECK (method IN ('email', 'sms')),
  token_hash BYTEA NOT NULL,
  status TEXT NOT NULL CHECK (status IN ('pending', 'consumed', 'revoked')),
  attempt_count INT NOT NULL DEFAULT 0 CHECK (attempt_count >= 0),
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMP(0) WITH TIME ZONE NOT NULL,
  consumed_at TIMESTAMP(0) WITH TIME ZONE,

  CONSTRAINT fk_user_invitations_user_id
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- At most one *pending* invitation per user per delivery method.
CREATE UNIQUE INDEX idx_user_invitations_pending_user_method
  ON user_invitations (user_id, method)
  WHERE status = 'pending';

CREATE INDEX idx_user_invitations_user_id ON user_invitations (user_id);
CREATE INDEX idx_user_invitations_expires_at ON user_invitations (expires_at)
  WHERE status = 'pending';