package lib

import "github.com/matthewhartstonge/argon2"

type IPasswordHasher interface {
	Hash(password string) ([]byte, error)
	Verify(password string, hashed_password string) (bool, error)
}

type PasswordHasher struct {
	argon argon2.Config
}

func (p *PasswordHasher) Hash(password string) ([]byte, error) {
	pw_hash, err := p.argon.HashEncoded([]byte(password))
	return pw_hash, err
}

func (p *PasswordHasher) Verify(password string, hashed_password string) (bool, error) {
	ok, err := argon2.VerifyEncoded([]byte(password), []byte(hashed_password))
	return ok, err
}

func GetPasswordHasher() IPasswordHasher {
	return &PasswordHasher{argon: argon2.DefaultConfig()}
}
