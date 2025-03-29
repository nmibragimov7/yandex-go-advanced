package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	t.Run("Generate token", func(t *testing.T) {
		ssp := SessionProvider{
			Config: nil,
		}

		token, err := ssp.GenerateToken(int64(1))
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestParseCookie(t *testing.T) {
	t.Run("Parse cookie", func(t *testing.T) {
		ssp := SessionProvider{
			Config: nil,
		}

		token, err := ssp.GenerateToken(int64(1))
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		userID, err := ssp.ParseCookie(token)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), userID)
	})

	t.Run("Parse cookie with error", func(t *testing.T) {
		ssp := SessionProvider{
			Config: nil,
		}

		token, err := ssp.GenerateToken(int64(1))
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		userID, err := ssp.ParseCookie(token)
		assert.NoError(t, err)
		assert.NotEqual(t, int64(0), userID)
	})
}

func TestCheckCookie(t *testing.T) {
	t.Run("Check cookie", func(t *testing.T) {
		ssp := SessionProvider{
			Config: nil,
		}

		token, err := ssp.GenerateToken(int64(1))
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		err = ssp.CheckCookie(token)
		assert.NoError(t, err)
	})

	t.Run("Check cookie with error", func(t *testing.T) {
		ssp := SessionProvider{
			Config: nil,
		}

		token, err := ssp.GenerateToken(int64(0))
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		err = ssp.CheckCookie(token)
		assert.Error(t, err)
	})
}
