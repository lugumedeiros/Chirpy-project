package auth

import (
	"github.com/alexedwards/argon2id"
)

func HashPassword(psw string) (string, error) {
	return argon2id.CreateHash(psw, argon2id.DefaultParams)
}

func CheckPasswordHash(psw string, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(psw, hash)
}
