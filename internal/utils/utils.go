package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
)

func GenerateCUID() string {
	b := make([]byte, 12)
	rand.Read(b)
	return fmt.Sprintf("c%s", hex.EncodeToString(b))
}

func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Strict check for production, but flexible for test
		if os.Getenv("GIN_MODE") == "release" {
			panic("JWT_SECRET environment variable is not set in release mode")
		}
		return []byte("default_secret_key")
	}
	return []byte(secret)
}
