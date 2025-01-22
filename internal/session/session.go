package session

import (
	"fmt"
	"time"
	"yandex-go-advanced/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type SessionProvider struct {
	*config.Config
}

const (
	cookieName = "user_token"
)

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
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signed, nil
}

func (p *SessionProvider) ParseToken(c *gin.Context) (int64, error) {
	cookie, err := c.Cookie(cookieName)
	if err != nil {
		return 0, fmt.Errorf("failed to parse cookie: %w", err)
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(cookie, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(*p.SercretKey), nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return 0, jwt.ErrTokenNotValidYet
	}
	if claims.UserID == 0 {
		return 0, jwt.ErrInvalidKey
	}

	return claims.UserID, nil
}
