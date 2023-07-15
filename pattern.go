package pattern

import (
	"sync/atomic"

	"github.com/sollniss/pattern/internal"
)

// Part is a part of a pattern.
type Part interface {
	// Append appends the Part to the output pattern.
	Append([]byte) []byte
}

type gen struct {
	parts []Part
}

// New returns a new pattern generator.
// The generator implements the Part interface, which means it can be used as a Part of another pattern.
func New(p ...Part) *gen {

	parts := make([]Part, 0, len(p))

	for i := 0; i < len(p); i++ {
		switch v := p[i].(type) {
		case nullpart:
			// Skip nullparts.
			continue
		case group:
			// Unwrap Group.
			parts = append(parts, v...)
		default:
			parts = append(parts, v)
		}

	}

	return &gen{
		parts: parts,
	}
}

// String returns a random pattern based on the Parts used to initialize the generator.
func (g gen) String() string {
	b := make([]byte, 0, 100)
	for _, p := range g.parts {
		b = p.Append(b)
	}
	return string(b)
}

// Append appends the generated pattern to b.
//
// Implements the Part interface.
func (g gen) Append(b []byte) []byte {
	for _, p := range g.parts {
		b = p.Append(b)
	}
	return b
}

type nullpart struct{}

func (p nullpart) Append(b []byte) []byte {
	return b
}

// Group returns a Part that wraps p into a single Part.
func Group(p ...Part) Part {
	// A Group of one is just the Part.
	if len(p) == 1 {
		return p[0]
	}

	return group(p)
}

type group []Part

func (p group) Append(b []byte) []byte {
	for _, p := range p {
		b = p.Append(b)
	}
	return b
}

// Repeat returns a Part that repeats p between min and max times randomly.
// If min == max, the Part will be repeated exactly max times in each iteration.
func Repeat(min uint32, max uint32, p ...Part) Part {
	if max == 0 {
		panic("max must be > 0")
	}

	if max < min {
		panic("max must be >= min")
	}

	// A constant repeat is a Group.
	if min > 0 && min == max {
		g := make(group, 0, len(p)*int(max))
		for i := uint32(0); i < max; i++ {
			g = append(g, p...)
		}
		return g
	}

	// A repeat with min == 0 and max == 1 is an Optional.
	if min == 0 && max == 1 {
		return potentially50{
			part: Group(p...),
		}
	}

	return repeat{
		parts: p,
		min:   min,
		maxr:  (max - min) + 1,
	}
}

type repeat struct {
	parts []Part
	min   uint32
	// maxr is the value needed to generate [min, max] with the RNG.
	maxr uint32
}

func (p repeat) Append(b []byte) []byte {
	n := internal.RandN(p.maxr) + p.min
	for i := uint32(0); i < n; i++ {
		for _, p := range p.parts {
			b = p.Append(b)
		}
	}

	return b
}

// Potentially returns a Part that will include p with probability c.
//
// Panics if c is < 0.
func Potentially(c float64, p Part) Part {
	if c < 0 {
		panic("chance must be > 0")
	}

	if c == 0 {
		return nullpart{}
	}

	// An Potentially with c >= 1 can never not be included.
	if c >= 1 {
		return p
	}

	// Fast path for c == 0.5, since we can check the last bit of the random number.
	if c == 0.5 {
		return potentially50{
			part: p,
		}
	}

	return potentiallyP{
		part:    p,
		percent: c,
	}
}

type potentially50 struct {
	part Part
}

func (p potentially50) Append(b []byte) []byte {
	if internal.Fastrand()&1 == 1 {
		b = p.part.Append(b)
	}
	return b
}

type potentiallyP struct {
	part    Part
	percent float64
}

func (p potentiallyP) Append(b []byte) []byte {
	if internal.RandFloat64() <= p.percent {
		b = p.part.Append(b)
	}
	return b
}

type literal []byte

// Literal returns a Part that will always output s.
func Literal(s string) Part {
	return literal(s)
}

func (p literal) Append(b []byte) []byte {
	return append(b, p...)
}

// OneOf returns a Part that selects one of p randomly in each iteration.
func OneOf(p ...Part) Part {
	// OneOf with one Part is just the Part.
	if len(p) == 1 {
		return p[0]
	}

	return anyOf{
		parts: p,
		len:   uint32(len(p)),
	}
}

type anyOf struct {
	parts []Part
	len   uint32
}

func (p anyOf) Append(b []byte) []byte {
	n := internal.RandN(p.len)
	return p.parts[n].Append(b)
}

// OneOfString returns a Part that will output one of s randomly in each iteration.
func OneOfString(s []string) Part {
	return anyOfString{
		alphabet: s,
		len:      uint32(len(s)),
	}
}

type anyOfString struct {
	alphabet []string
	len      uint32
}

func (p anyOfString) Append(b []byte) []byte {
	n := internal.RandN(p.len)
	return append(b, p.alphabet[n]...)
}

// OneOfByte returns a Part that will select one of b randomly in each iteration.
func OneOfByte(b []byte) Part {
	return anyOfByte{
		alphabet: b,
		len:      uint32(len(b)),
	}
}

type anyOfByte struct {
	alphabet []byte
	len      uint32
}

func (p anyOfByte) Append(b []byte) []byte {
	n := internal.RandN(p.len)
	return append(b, p.alphabet[n])
}

// OneOfRune returns a Part that will select one of r randomly in each iteration.
// The length of the alphabet must be less than 2^32.
func OneOfRune(r []rune) Part {
	return anyOfRune{
		alphabet: r,
		len:      uint32(len(r)),
	}
}

type anyOfRune struct {
	alphabet []rune
	len      uint32
}

func (p anyOfRune) Append(b []byte) []byte {
	n := internal.RandN(p.len)
	return append(b, string(p.alphabet[n])...)
}

// Shuffle returns a Part that randomly rearranges p in each iteration.
// Uses the Fisher-Yates shuffle to generate permutations.
//
// https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
func Shuffle(p ...Part) Part {
	return shuffle{
		parts: p,
		len:   uint32(len(p)),
	}
}

type shuffle struct {
	parts []Part
	len   uint32
}

func (p shuffle) Append(b []byte) []byte {

	// Fisher-Yates shuffle: https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
	for i := p.len - 1; i > 0; i-- {
		j := internal.RandN(i + 1)
		p.parts[i], p.parts[j] = p.parts[j], p.parts[i]
	}

	for i := uint32(0); i < p.len; i++ {
		b = p.parts[i].Append(b)
	}

	return b
}

// Sequence returns a Part that will on each iteration increment a number from start to max.
// The number will be zero-padded to width.
// The output number will reset to start when max is reached.
func Sequence(start uint64, max uint64, width int) Part {
	if max < start {
		panic("max must be >= min")
	}

	var curr uint64 = start - 1
	return sequence{
		start: start,
		max:   max,
		width: width,
		curr:  &curr,
	}
}

type sequence struct {
	start uint64
	max   uint64
	width int
	curr  *uint64
}

func (p sequence) Append(b []byte) []byte {
	for {
		last := atomic.LoadUint64(p.curr)
		curr := last + 1
		if curr > p.max || curr < p.start {
			curr = p.start
		}

		if atomic.CompareAndSwapUint64(p.curr, last, curr) {
			return appendInt(b, curr, p.width)
		}
	}
}

func appendInt(b []byte, u uint64, width int) []byte {
	// Compute the number of decimal digits.
	var n int
	if u == 0 {
		n = 1
	}
	for u2 := u; u2 > 0; u2 /= 10 {
		n++
	}

	// Add 0-padding.
	for pad := width - n; pad > 0; pad-- {
		b = append(b, '0')
	}

	// Ensure capacity.
	if len(b)+n <= cap(b) {
		b = b[:len(b)+n]
	} else {
		b = append(b, make([]byte, n)...)
	}

	// Assemble decimal in reverse order.
	i := len(b) - 1
	for u >= 10 && i > 0 {
		q := u / 10
		b[i] = itob(u - q*10)
		u = q
		i--
	}
	b[i] = itob(u)
	return b
}

func itob(u uint64) byte {
	return '0' + byte(u)
}
