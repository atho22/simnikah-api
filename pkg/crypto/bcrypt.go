package crypto

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword melakukan hashing password dengan bcrypt
func HashPassword(password string) (string, error) {
	// Validasi input
	if password == "" {
		return "", errors.New("password tidak boleh kosong")
	}

	// Bcrypt dengan cost 12 (aman dan cukup cepat)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// VerifyPassword memverifikasi password dengan hash bcrypt
func VerifyPassword(password, hashedPassword string) error {
	// Validasi input
	if password == "" || hashedPassword == "" {
		return errors.New("password atau hash tidak boleh kosong")
	}

	// Verifikasi dengan bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return errors.New("password tidak valid")
	}

	return nil
}
