//go:build go1.19 && !unsafe

package internal

import (
	"hash/maphash"
)

func Fastrand() uint64 {
	var h maphash.Hash
	h.SetSeed(maphash.MakeSeed())
	return h.Sum64()
}
