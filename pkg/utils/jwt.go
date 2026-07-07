package utils

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = getJWTKey()

// getJWTKey returns the JWT key from environment.
// If not set, return nil so callers can detect misconfiguration.
func getJWTKey() []byte {
	key := os.Getenv("JWT_KEY")
	if key == "" {
		log.Println("WARNING: JWT_KEY environment variable tidak diset. Set JWT_KEY sebelum deploy ke production.")
		return nil
	}
	return []byte(key)
}

type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Nama   string `json:"nama"`
	NIP    string `json:"nip,omitempty"`
	jwt.RegisteredClaims
}

func GenerateToken(claims jwt.Claims) (string, error) {
	if len(jwtKey) == 0 {
		return "", fmt.Errorf("server misconfigured: JWT_KEY environment variable not set")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("gagal membuat token: %v", err)
	}
	return tokenString, nil
}

// ParseToken memvalidasi dan mem-parse JWT token.
// Memvalidasi algoritma signing untuk mencegah serangan alg:none.
func ParseToken(tokenStr string) (*TokenClaims, error) {
	tokenStr = strings.Replace(tokenStr, "Bearer ", "", 1)
	if len(jwtKey) == 0 {
		return nil, errors.New("server misconfigured: JWT_KEY environment variable not set")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validasi algoritma signing — hanya izinkan HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode signing tidak valid: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("token tidak valid")
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, errors.New("token claims error")
	}

	return claims, nil
}

// RandString generates a random string of given length (not for cryptographic use)
func RandString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
