CREATE TABLE IF NOT EXISTS users_auth (
  id BIGINT PRIMARY KEY REFERENCES users(id),
  email STRING UNIQUE NOT NULL,
  pw_hash STRING NOT NULL
  -- TODO: Will add more stuff for auth, such as 2fa phone
)
