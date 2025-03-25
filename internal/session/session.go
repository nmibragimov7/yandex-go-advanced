package session

import (
	"fmt"
	"time"

	"yandex-go-advanced/internal/config"

	"github.com/gin-gonic/gin"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// SessionProvider - struct that contains the necessary session settings
type SessionProvider struct {
	*config.Config
}

const (
	cookieName = "user_token"
)

// Claims - struct that contains jwt settings
type Claims struct {
	jwtv5.RegisteredClaims
	UserID int64
}

var hashKey = []byte("my-secret-hash-key")

// GenerateToken - func for generate token
func (p *SessionProvider) GenerateToken(userID int64) (string, error) {
	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, Claims{
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
		UserID: userID,
	})

	signed, err := token.SignedString(hashKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signed, nil
}

// ParseCookie - func for parse cookie
func (p *SessionProvider) ParseCookie(c *gin.Context) (int64, error) {
	cookie, err := c.Cookie(cookieName)
	if err != nil || cookie == "" {
		return 0, fmt.Errorf("failed to parse cookie: %w", err)
	}

	claims := &Claims{}
	token, err := jwtv5.ParseWithClaims(cookie, claims,
		func(t *jwtv5.Token) (interface{}, error) {
			return hashKey, nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return 0, jwtv5.ErrTokenNotValidYet
	}

	if claims.UserID == 0 {
		return 0, jwtv5.ErrInvalidKey
	}

	return claims.UserID, nil
}

// CheckCookie - func for check cookie
func (p *SessionProvider) CheckCookie(cookie string) error {
	claims := &Claims{}
	token, err := jwtv5.ParseWithClaims(cookie, claims,
		func(t *jwtv5.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwtv5.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return hashKey, nil
		},
	)
	if err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return jwtv5.ErrTokenNotValidYet
	}

	if claims.UserID == 0 {
		return jwtv5.ErrInvalidKey
	}

	return nil
}
