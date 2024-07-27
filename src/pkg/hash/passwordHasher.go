package hash

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=password.go -destination=../../internal/tests/unitTests/serviceTests/mocks/mockPasswordHasher.go --package=mocks

// IPasswordHasher provides hashing logic to securely store passwords.
type IPasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}

type PasswordHasher struct {
	salt string
}

func NewPasswordHasher(salt string) *PasswordHasher {
	return &PasswordHasher{salt: salt}
}

// Hash creates bcrypt hash of the given password.
func (ph *PasswordHasher) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+ph.salt), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("[!] ERROR! Error hashing password: %v", err)
	}

	return string(hashedPassword), nil
}

// Compare compares hashed password with the plain password.
func (ph *PasswordHasher) Compare(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+ph.salt))
	if err != nil {
		return fmt.Errorf("[!] ERROR! Wrong password")
	}

	return nil
}
