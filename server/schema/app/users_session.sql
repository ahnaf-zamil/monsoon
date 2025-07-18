CREATE TABLE IF NOT EXISTS users_session (
  session_id BIGINT PRIMARY KEY,
  user_id BIGINT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
  refresh_token TEXT NOT NULL,
  created_at BIGINT NOT NULL
  -- TODO: Will add more stuff for auth, such as 2fa phone
)
