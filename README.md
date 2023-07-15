# pattern
Go library to generate string patterns.

# Example

```go
import (
	"fmt"

	"github.com/sollniss/pattern"
)

func main() {
	gen := pattern.New(
		pattern.Repeat(2, 4,
			pattern.Repeat(5, 5, pattern.OneOfByte([]byte("1234567890"))),
			pattern.Literal("-"),
			pattern.Repeat(2, 4, pattern.OneOfRune([]rune("xyz"))),
			pattern.Literal("-"),
		),
		pattern.Potentially(0.3, pattern.OneOfString([]string{"AAAA", "BBBB", "CCCC", "DDDD"})),
		pattern.Literal("-"),
		pattern.Shuffle(pattern.Literal("a"), pattern.Literal("b"), pattern.Literal("c")),
		pattern.Sequence(1, 999, 4),
	)
	for i := 0; i < 10; i++ {
		fmt.Println(gen.String())
	}
}
```
Outputs:
```
58525-zzz-30662-yxx-71228-xz-21148-xyz--bca0001
22813-yzyx-74510-yx-40801-yzyz--cab0002
61736-yx-76951-yxxz--cab0003
79908-zxx-96321-zxz-71442-yz-DDDD-acb0004
97571-xz-87210-xy--cab0005
11168-zzz-34216-xy-92184-zz-13246-zxy--abc0006
76919-yx-13975-xx-34534-yyz--cab0007
99940-yzz-31787-yyx--acb0008
13386-zz-60143-yzy-47725-zzxx-CCCC-cab0009
53765-zzyx-87695-zyxz-74415-yxyz--cab0010
```

# Usage

The package exposes an interface `Part` which all generator fuctions implement. The generator returned by `New` also implements this interface, which allows combination of multiple generators.
```go
type Part interface {
	// Append appends the Part to the output pattern.
	Append([]byte) []byte
}
```

## Functions

```go
Group(p ...Part) Part
```
Group returns a `Part` that wraps `p` into a single `Part`.

```go
Repeat(min uint32, max uint32, p ...Part) Part
```
Repeat returns a `Part` that repeats `p` between `min` and `max` times randomly.

```go
Potentially(c float64, p Part) Part
```
Potentially returns a `Part` that will include `p` with probability `c`.

```go
Literal(s string) Part
```
Literal returns a `Part` that will always output `s`.

```go
OneOf(p ...Part) Part
```
OneOf returns a `Part` that selects one of `p` randomly in each iteration.
The package also provides the convenience functions `OneOfString`, `OneOfByte` and `OneOfRune`.

```go
Shuffle(p ...Part) Part
```
Shuffle returns a `Part` that randomly rearranges `p` in each iteration.

```go
Sequence(start uint64, max uint64, width int) Part
```
Sequence returns a `Part` that will on each iteration increment a number from `start` to `max` and display the number zero-padded to `width`.
Sequence is thread safe.