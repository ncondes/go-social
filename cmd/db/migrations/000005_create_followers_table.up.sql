CREATE TABLE IF NOT EXISTS followers (
  user_id BIGINT NOT NULL,
  follower_id BIGINT NOT NULL,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),

  CONSTRAINT pk_followers PRIMARY KEY (user_id, follower_id), -- composite primary key
  CONSTRAINT fk_followers_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_followers_follower_id FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
  CHECK (user_id != follower_id)
);

CREATE INDEX idx_user_id ON followers(user_id);
CREATE INDEX idx_follower_id ON followers(follower_id);