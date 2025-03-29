package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkGetKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetKey()
	}
}

func TestGetKey(t *testing.T) {
	t.Run("GetKey test", func(t *testing.T) {
		key := GetKey()
		assert.Equal(t, len(key), 8)
	})
}
