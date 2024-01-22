package types

import "github.com/golang-jwt/jwt/v5"

type ContextKey string

const UserKey ContextKey = "UserID"

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	UserID string `json:"userID"`
	jwt.RegisteredClaims
}
