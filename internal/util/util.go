package util

import (
	"crypto/rand"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func VerifyPassword(password, hashed string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err
}

func RandomInt(min, max int) int {
	nBig, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	return int(nBig.Int64()) + min
}
