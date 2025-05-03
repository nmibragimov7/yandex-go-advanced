package logger

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInit(t *testing.T) {
	t.Run("Logger initialized", func(t *testing.T) {
		sgr := Init()
		assert.NotNil(t, sgr)
	})

	t.Run("Logger initialized with error", func(t *testing.T) {
		origNewProduction := newProduction
		defer func() { newProduction = origNewProduction }()

		newProduction = func(opts ...zap.Option) (*zap.Logger, error) {
			return nil, errors.New("mock error")
		}

		logger := Init()
		assert.Nil(t, logger)
	})
}
