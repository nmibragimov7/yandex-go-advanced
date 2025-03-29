package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	t.Run("Main run no error", func(t *testing.T) {
		err := run()
		assert.NoError(t, err)
	})
}
