CREATE TABLE IF NOT EXISTS users_auth (
  id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE NOT NULL,
  email TEXT UNIQUE NOT NULL,
  pw_hash TEXT NOT NULL
  -- TODO: Will add more stuff for auth, such as 2fa phone
)
