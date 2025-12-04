package auth

import (
	"time"

	"github.com/go-chi/jwtauth"
)

type Claims struct {
	UserID string `json:"sub"`
	Role   string `json:"role"`
	Exp    int64  `json:"exp"`
	Iat    int64  `json:"iat"`
}

type TokenGenerator struct {
	tokenAuth *jwtauth.JWTAuth
	expiresIn int
}

func NewTokenGenerator(tokenAuth *jwtauth.JWTAuth, expiresIn int) *TokenGenerator {
	return &TokenGenerator{
		tokenAuth: tokenAuth,
		expiresIn: expiresIn,
	}
}

func (tg *TokenGenerator) GenerateAccessToken(userID, role string) (string, time.Duration, error) {
	ttl := time.Duration(tg.expiresIn) * time.Second
	now := time.Now()

	claims := Claims{
		UserID: userID,
		Role:   role,
		Exp:    now.Add(ttl).Unix(),
		Iat:    now.Unix(),
	}

	claimsMap := map[string]interface{}{
		"sub":  claims.UserID,
		"role": claims.Role,
		"exp":  claims.Exp,
		"iat":  claims.Iat,
	}

	_, tokenString, err := tg.tokenAuth.Encode(claimsMap)
	if err != nil {
		return "", 0, err
	}

	return tokenString, ttl, nil
}
