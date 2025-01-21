package session

import (
	"fmt"
	"time"
	"yandex-go-advanced/internal/config"

	"github.com/golang-jwt/jwt/v4"
)

type SessionProvider struct {
	*config.Config
}

type Claims struct {
	jwt.RegisteredClaims
	UserID int64
}

func (p *SessionProvider) GenerateToken(userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
		UserID: userID,
	})

	signed, err := token.SignedString([]byte(*p.SercretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %s", err.Error())
	}

	return signed, nil
}

func (p *SessionProvider) ValidateToken(tokenString string) (int64, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(*p.SercretKey), nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %s", err.Error())
	}

	if !token.Valid {
		return 0, jwt.ErrTokenNotValidYet
	}
	if claims.UserID == 0 {
		return 0, jwt.ErrInvalidKey
	}

	return claims.UserID, nil
}
