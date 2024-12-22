package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashToken(token string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(token))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
