package internal

import "testing"

var v uint32

func BenchmarkRandN(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v = RandN(100)
	}
}
