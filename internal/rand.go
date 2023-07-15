package internal

import (
	"crypto/rand"
	"encoding/binary"
	"math/bits"
)

const (
	int53Mask = 1<<53 - 1
	f53Mul    = 0x1.0p-53
)

// RandN returns a random uint32 in [0, n).
func RandN(n uint32) uint32 {
	res, _ := bits.Mul64(uint64(n), Fastrand())
	return uint32(res)
}

// Float64 returns a random float64 in [0.0, 1.0).
func RandFloat64() float64 {
	return float64(Fastrand()&int53Mask) * f53Mul
}

func SecureRandomReader(b []byte) int {
	// Ignore error, we might not have gotten all bytes,
	// but can use what we got.
	len, _ := rand.Reader.Read(b)
	return len
}

func FastRandomReader(b []byte) int {
	n := 0
	for ; n+8 <= len(b); n += 8 {
		binary.LittleEndian.PutUint64(b[n:n+8], Fastrand())
	}
	if n < len(b) {
		val := Fastrand()
		for ; n < len(b); n++ {
			b[n] = byte(val)
			val >>= 8
		}
	}

	return n
}
