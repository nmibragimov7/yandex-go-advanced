package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	t.Run("Logger initialized", func(t *testing.T) {
		sgr := Init()
		assert.NotNil(t, sgr)
	})
}
