// File: helper/jwt.go
package helper

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = getJWTKey()

// getJWTKey returns the JWT key from environment or uses a fallback
func getJWTKey() []byte {
	key := os.Getenv("JWT_KEY")
	if key == "" {
		key = "secret-key-boleh-diubah" // Fallback key
	}
	return []byte(key)
}

type TokenClaims struct {
	UserID     string `json:"user_id"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Kode_prodi string `json:"kode_prodi"`
	Nama       string `json:"nama"`
	NIP        string `json:"nip"`
	NIM        string `json:"nim"`
	jwt.RegisteredClaims
}

func GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("gagal membuat token: %v", err)
	}
	return tokenString, nil
}

func ParseToken(tokenStr string) (*TokenClaims, error) {
	tokenStr = strings.Replace(tokenStr, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
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

func RandString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func ExtractToken(tokenString string) map[string]interface{} {
	// Bersihkan token dari prefix "Bearer " jika ada
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validasi metode enkripsi
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode token tidak valid: %v", token.Method.Alg())
		}
		return jwtKey, nil
	})

	if err != nil {
		fmt.Printf("Error parsing token: %v\n", err)
		return nil
	}

	if !token.Valid {
		fmt.Println("Token tidak valid")
		return nil
	}

	// Ambil klaim
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Validasi claims yang diperlukan
		requiredClaims := []string{"user_id", "role", "email"}
		for _, claim := range requiredClaims {
			if _, exists := claims[claim]; !exists {
				fmt.Printf("Claim %s tidak ditemukan\n", claim)
				return nil
			}
		}

		// Validasi role
		role, ok := claims["role"].(string)
		if !ok {
			fmt.Println("Role tidak valid")
			return nil
		}

		validRoles := map[string]bool{
			"admin_akademik": true,
			"admin_prodi":    true,
			"dosen":          true,
			"mahasiswa":      true,
		}
		if !validRoles[role] {
			fmt.Printf("Role %s tidak valid\n", role)
			return nil
		}

		// Validasi expired time
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				fmt.Println("Token sudah expired")
				return nil
			}
		}

		return claims
	}

	fmt.Println("Gagal mengkonversi claims")
	return nil
}
