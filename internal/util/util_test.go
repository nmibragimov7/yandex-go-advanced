package util

import "testing"

func BenchmarkGetKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetKey()
	}
}
