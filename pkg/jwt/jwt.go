package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your_secret_key")

func ValidateJWT(tokenStr string) (jwt.MapClaims, bool) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	return claims, ok
}

func GenerateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID, // стандартное поле для ID
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(), // неделя
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func RefreshToken(refreshToken string) (string, jwt.MapClaims, bool) {
	claims, ok := ValidateJWT(refreshToken)
	if !ok {
		return "", nil, false
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", nil, false
	}

	newAccess, err := GenerateAccessToken(sub)
	if err != nil {
		return "", nil, false
	}

	return newAccess, claims, true
}