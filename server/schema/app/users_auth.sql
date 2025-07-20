CREATE TABLE IF NOT EXISTS users_auth (
  id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE NOT NULL,
  email TEXT UNIQUE NOT NULL,
  pw_hash BYTEA NOT NULL,
  pw_salt BYTEA NOT NULL, -- Password salt
  enc_salt BYTEA NOT NULL, -- Key Encryption salt
  key_seed_cipher BYTEA NOT NULL, -- Encrypted Curve25519 keypair seed
  nonce BYTEA NOT NULL -- 12 byte nonce for AES-GCM
)
