//go:build unsafe

package internal

import (
	_ "unsafe"
)

//go:linkname Fastrand runtime.fastrand64
func Fastrand() uint64
