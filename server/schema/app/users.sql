CREATE TABLE IF NOT EXISTS users (
  id BIGINT PRIMARY KEY,
  username TEXT NOT NULL,
  display_name TEXT NOT NULL,
  created_at BIGINT NOT NULL,
  updated_at BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS users_auth (
  id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  email TEXT UNIQUE NOT NULL,
  pw_hash BYTEA NOT NULL,
  pw_salt BYTEA NOT NULL, -- Password salt
  enc_salt BYTEA NOT NULL, -- Key Encryption salt
  key_seed_cipher BYTEA NOT NULL, -- Encrypted Curve25519 keypair seed
  nonce BYTEA NOT NULL -- 12 byte nonce for AES-GCM
);

CREATE TABLE IF NOT EXISTS users_key (
  user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  pub_sig_key BYTEA NOT NULL,
  pub_enc_key BYTEA NOT NULL
);

CREATE TABLE IF NOT EXISTS users_session (
  session_id BIGINT PRIMARY KEY,
  user_id BIGINT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
  refresh_token TEXT NOT NULL,
  created_at BIGINT NOT NULL
  -- TODO: Will add more stuff for auth, such as 2fa phone
);
