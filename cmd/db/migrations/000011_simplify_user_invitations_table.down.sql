DROP INDEX IF EXISTS idx_user_invitations_expires_at;

ALTER TABLE user_invitations DROP CONSTRAINT IF EXISTS uq_user_invitations_user_method;

ALTER TABLE user_invitations
  ADD COLUMN status TEXT NOT NULL DEFAULT 'pending'
    CHECK (status IN ('pending', 'consumed', 'revoked'));

ALTER TABLE user_invitations
  ADD COLUMN attempt_count INT NOT NULL DEFAULT 0 CHECK (attempt_count >= 0);

ALTER TABLE user_invitations
  ADD COLUMN consumed_at TIMESTAMP(0) WITH TIME ZONE;

CREATE UNIQUE INDEX idx_user_invitations_pending_user_method
  ON user_invitations (user_id, method)
  WHERE status = 'pending';

CREATE INDEX idx_user_invitations_expires_at ON user_invitations (expires_at)
  WHERE status = 'pending';