package auth

import (
	"time"
	"fmt"
	"net/http"
	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	// "github.com/google/uuid"
	"strings"
	"crypto/rand"
	"encoding/hex"
)

func HashPassword(psw string) (string, error) {
	return argon2id.CreateHash(psw, argon2id.DefaultParams)
}

func CheckPasswordHash(psw string, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(psw, hash)
}

func MakeJWT(userID string, tokenSecret string, expiresIn time.Duration)(string, error){
	timeNow := time.Now()
	timeThen := timeNow.Add(expiresIn)
	claim := jwt.RegisteredClaims{
		Issuer: "chirpy-access",
		Subject: userID,
		IssuedAt: jwt.NewNumericDate(timeNow),
		ExpiresAt: jwt.NewNumericDate(timeThen),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (string, error){
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func (token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("invalid signing method")
			}
			return []byte(tokenSecret), nil
		},)
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	userID := claims.Subject
	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	const prefix = "Bearer "
	authValRaw := headers.Get("Authorization")
	val := strings.Replace(authValRaw, prefix, "", 1)
	if val == "" {
		return "", fmt.Errorf("Invalid auth key")
	} else {
		return val, nil
	}
}

func MakeRefreshToken() string{
	val := make([]byte, 32)
	rand.Read(val)
	return hex.EncodeToString(val)
}

func GetApiKey(headers http.Header) (string, error){
	const prefix = "ApiKey "
	authValRaw := headers.Get("Authorization")
	val := strings.Replace(authValRaw, prefix, "", 1)
	if val == "" {
		return "", fmt.Errorf("Invalid auth key")
	} else {
		return val, nil
	}
}