package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	t.Run("Config initialized", func(t *testing.T) {
		cnf := Init()
		assert.NotNil(t, cnf)
	})
}
