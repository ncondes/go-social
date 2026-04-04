DROP INDEX IF EXISTS idx_user_invitations_pending_user_method;
DROP INDEX IF EXISTS idx_user_invitations_expires_at;

ALTER TABLE user_invitations DROP COLUMN IF EXISTS consumed_at;
ALTER TABLE user_invitations DROP COLUMN IF EXISTS attempt_count;
ALTER TABLE user_invitations DROP COLUMN IF EXISTS status;

ALTER TABLE user_invitations
  ADD CONSTRAINT uq_user_invitations_user_method UNIQUE (user_id, method);

CREATE INDEX idx_user_invitations_expires_at ON user_invitations (expires_at);