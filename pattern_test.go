package pattern

import (
	"fmt"
	"math"
	"strconv"
	"testing"
)

var id string

func BenchmarkString(b *testing.B) {
	gen := New(
		Repeat(2, 5,
			Repeat(5, 5, OneOfByte([]byte("1234567890"))),
			Literal("-"),
			Repeat(3, 6, OneOfRune([]rune("あいうえおかきくけこさしすせそたちつてとなにぬねのはひふへほまみむめもやゆよらりるれろわをん"))),
			Literal("-"),
		),
		OneOfString([]string{"asd", "fgh", "jkl"}),
	)

	for i := 0; i < b.N; i++ {
		id = gen.String()
	}
}

func BenchmarkRepeat(b *testing.B) {

	benchs := []struct {
		name string
		gen  *gen
	}{
		{"Group(5)", New(Group(Literal("-"), Literal("-"), Literal("-"), Literal("-"), Literal("-")))},
		{"Repeat(5,5)", New(Repeat(5, 5, Literal("-")))},
		{"Repeat(50,100)", New(Repeat(50, 100, Literal("-")))},
		{"Optional()", New(Potentially(0.5, Literal("-")))},
		{"Repeat(0, 1)", New(Repeat(0, 1, Literal("-")))},
	}

	for _, bb := range benchs {
		b.Run(bb.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				id = bb.gen.String()
			}
		})
	}
}

func BenchmarkOneOf(b *testing.B) {

	benchs := []struct {
		name string
		gen  *gen
	}{
		{"OneOfByte urlsafe", New(OneOfByte([]byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-")))},
		{"OneOfRune urlsafe", New(OneOfRune([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-")))},
		{"OneOfString urlsafe", New(OneOfString([]string{
			"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
			"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
			"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "_", "-"}))},
		{"OneOfString", New(OneOfString([]string{"asd", "fgh", "jkl", "123", "あいう", "ḇ̴̨̢̪͉̱̪̝͕̩̳̋ͅņ̵̱͍̥̝̹̹͚͖̣͓̮̎̊ͅm̶̛̰͉̘̯̗̬̪̻̫̩̘͓̀̎̏̅͑̔̒̚", "😀😆🙃"}))},

		{"OneOfRune zalgo", New(OneOfRune([]rune("t̵̡̛͈̪͙͎̙̘̙̥͇̉̂̈̒͊͋̑͒͘ë̸̢̟͇͓̲̝̣̗̳́̂͆̉̐͂͗͊̈́̓͝s̴̬̳͖̝͊͗̒͊̏t̸͉̬̼̳̯̞͖̯͚̦̎̂̊̓̆̚̚͝+/")))},
	}

	for _, bb := range benchs {
		b.Run(bb.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				id = bb.gen.String()
			}
		})
	}
}

func TestAppend(t *testing.T) {
	gen := New(Literal("o"))

	b := []byte("hello")
	b = gen.Append(b)
	b = gen.Append(b)

	if string(b) != "hellooo" {
		t.Errorf("Append returned invalid value: want \"hellooo\", got %s", strconv.Quote(string(b)))
	}
}

func TestGroup(t *testing.T) {
	gen := New(Group())
	p := gen.String()
	if p != "" {
		t.Errorf("Group returned invalid value: want \"\", got %s", strconv.Quote(p))
	}

	gen = New(Group(Literal("o")))
	p = gen.String()
	if p != "o" {
		t.Errorf("Group returned invalid value: want \"o\", got %s", strconv.Quote(p))
	}

	gen = New(Group(Literal("o"), Literal("o"), Literal("o")))
	p = gen.String()
	if p != "ooo" {
		t.Errorf("Group returned invalid value: want \"ooo\", got %s", strconv.Quote(p))
	}

	// Prevent the Group from being unwrapped.
	gen = New(Repeat(1, 1, (Group(Literal("o"), Literal("o"), Literal("o")))))
	p = gen.String()
	if p != "ooo" {
		t.Errorf("Group returned invalid value: want \"ooo\", got %s", strconv.Quote(p))
	}
}

func TestRepeat(t *testing.T) {
	tests := []struct {
		name string
		min  uint32
		max  uint32
	}{
		{"optional", 0, 1},
		{"const length", 5, 5},
		{"normal", 5, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := New(Repeat(tt.min, tt.max, Literal("o")))

			foundMin := uint32(math.MaxUint32)
			foundMax := uint32(0)

			for i := 0; i < 10000; i++ {
				len := uint32(len(gen.String()))
				if len < tt.min || len > tt.max {
					t.Errorf("Repeat has invalid length: want [%d,%d], got %d", tt.min, tt.max, len)
				}

				if len < foundMin {
					foundMin = len
				} else if len > foundMax {
					foundMax = len
				}
			}

			// This could statistically fail, but it's unlikely (0.1^10000).
			if foundMin != tt.min {
				t.Errorf("Repeat never returned min length: want %d, got %d", tt.min, foundMin)
			}

			if foundMax != tt.max {
				t.Errorf("Repeat never returned max length: want %d, got %d", tt.max, foundMax)
			}

		})
	}
}

func TestRepeatEmpty(t *testing.T) {
	gen := New(Repeat(10, 100))
	p := gen.String()
	if p != "" {
		t.Errorf("Repeat returned invalid value: want \"\", got %s", strconv.Quote(p))
	}
}

func TestConstRepeatEqualToGroup(t *testing.T) {
	g1 := New(Group(Literal("o"), Literal("o"), Literal("o"), Literal("o"), Literal("o")))
	g2 := New(Repeat(5, 5, Literal("o")))
	v1 := g1.String()
	v2 := g2.String()
	if v1 != v2 {
		t.Errorf("constant Repeat and equivalent group did not return the same ID: want %s, got %s", strconv.Quote(v1), strconv.Quote(v2))
	}
}

func TestRepeatConstLength(t *testing.T) {
	const length uint32 = 1000

	gen := New(Repeat(length, length, Literal("o")))
	len := uint32(len(gen.String()))
	if len != length {
		t.Errorf("constant Repeat has invalid length: want %d, got %d", length, len)
	}
}

func TestRepeatPanic(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Repeat with max == 0 did not panic")
			}
		}()

		New(Repeat(0, 0, Literal("o")))
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Repeat with max < min did not panic")
			}
		}()

		New(Repeat(2, 1, Literal("o")))
	}()
}

func TestPotentially(t *testing.T) {
	var chances []float64 = []float64{0.1, 0.5, 0.9, 0.99}

	for _, c := range chances {

		gen := New(Potentially(c, Literal("o")))

		res := make(map[string]bool, 2)

		for i := 0; i < 1000; i++ {
			res[gen.String()] = true
		}

		if len(res) != 2 {
			t.Errorf("Potentially(%f) does not have two states: want 2, got %d", c, len(res))
		}

		if !res[""] {
			t.Errorf("Potentially(%f) never returned empty string", c)
		}

		if !res["o"] {
			t.Errorf("Potentially(%f) never returned Part", c)
		}
	}
}

func TestPotentiallyZero(t *testing.T) {
	gen := New(Potentially(0, Literal("o")))

	res := make(map[string]bool, 1)

	for i := 0; i < 100; i++ {
		res[gen.String()] = true
	}

	if len(res) != 1 {
		t.Errorf("Potentially does not have one state: want 1, got %d", len(res))
	}

	if !res[""] {
		t.Errorf("Potentially never returned empty string, even though c was 0")
	}

	if res["o"] {
		t.Errorf("Potentially returned Part, even though c was 0")
	}
}

func TestPotentiallyOne(t *testing.T) {
	gen := New(Potentially(1, Literal("o")))

	res := make(map[string]bool, 1)

	for i := 0; i < 100; i++ {
		res[gen.String()] = true
	}

	if len(res) != 1 {
		t.Errorf("Potentially does not have one state: want 1, got %d", len(res))
	}

	if res[""] {
		t.Errorf("Potentially returned empty string eben though c was 1")
	}

	if !res["o"] {
		t.Errorf("Potentially never Part, even though c was 0")
	}
}

func TestPotentiallyPanic(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Potentially with negative c did not panic")
			}
		}()

		New(Potentially(-1, Literal("o")))
	}()
}

func TestOneOf(t *testing.T) {
	var alphabets [][]string = [][]string{
		{""},
		{"a"},
		{"a", "b", "c"},
	}

	for _, alphabet := range alphabets {
		parts := make([]Part, 0, len(alphabet))
		hitmap := make(map[string]bool, len(alphabet))
		for _, v := range alphabet {
			parts = append(parts, Literal(v))
			hitmap[v] = false
		}

		gen := New(OneOf(parts...))
		for i := 0; i < 100; i++ {
			v := gen.String()
			if _, ok := hitmap[v]; !ok {
				t.Errorf("OneOf returned invalid value: want one of %v, got %s", alphabet, strconv.Quote(v))
			}
			hitmap[v] = true
		}

		for s, found := range hitmap {
			if !found {
				t.Errorf("OneOf with Literals %v never returned %v", alphabet, strconv.Quote(s))
			}
		}
	}
}

func TestOneOfByte(t *testing.T) {
	var alphabet []byte = []byte("\naB%1 ")

	bytemap := make(map[byte]bool, len(alphabet))
	for _, v := range alphabet {
		bytemap[v] = false
	}

	gen := New(OneOfByte(alphabet))
	for i := 0; i < 100; i++ {
		v := []byte(gen.String())[0]
		if _, ok := bytemap[v]; !ok {
			t.Errorf("OneOfByte returned invalid byte: want one of %s, got %s", strconv.Quote(string(alphabet)), strconv.Quote(string(v)))
		}
		bytemap[v] = true
	}

	for b, found := range bytemap {
		if !found {
			t.Errorf("OneOfByte with alphabet %s never returned %v", strconv.Quote(string(alphabet)), strconv.Quote(string(b)))
		}
	}
}

func TestOneOfRune(t *testing.T) {
	var alphabet []rune = []rune("\naB%1 ó̸̡̮̠̯͉̩͉͈͔̳̯̠̪͕͙̀䯂☺😀")

	bytemap := make(map[rune]bool, len(alphabet))
	for _, v := range alphabet {
		bytemap[v] = false
	}

	gen := New(OneOfRune(alphabet))
	for i := 0; i < 1000; i++ {
		v := []rune(gen.String())[0]
		if _, ok := bytemap[v]; !ok {
			t.Errorf("OneOfRune returned invalid rune: want one of %s, got %s", strconv.Quote(string(alphabet)), strconv.Quote(string(v)))
		}
		bytemap[v] = true
	}

	for r, found := range bytemap {
		if !found {
			t.Errorf("OneOfRune with alphabet %s never returned %v", strconv.Quote(string(alphabet)), strconv.Quote(string(r)))
		}
	}
}

func TestOneOfString(t *testing.T) {
	var alphabet []string = []string{"aaa", "bbb", "ccc"}

	hitmap := make(map[string]bool, len(alphabet))
	for _, v := range alphabet {
		hitmap[v] = false
	}

	gen := New(OneOfString(alphabet))
	for i := 0; i < 100; i++ {
		v := gen.String()
		if _, ok := hitmap[v]; !ok {
			t.Errorf("OneOfString returned invalid rune: want one of %v, got %s", alphabet, strconv.Quote(v))
		}
		hitmap[v] = true
	}

	for s, found := range hitmap {
		if !found {
			t.Errorf("OneOfString with alphabet %v never returned %v", alphabet, strconv.Quote(s))
		}
	}
}

func TestSequence(t *testing.T) {
	start := 0
	max := 100
	width := 4

	gen := New(Sequence(uint64(start), uint64(max), width))
	for i := start; i <= max; i++ {
		v := gen.String()
		if len(v) != width {
			t.Errorf("Sequence has invalid width: want %d, got %d", width, len(v))
		}
		if v != fmt.Sprintf(fmt.Sprintf("%%0%dd", width), i) {
			t.Errorf("Sequence returned invalid value: want %d, got %s", i, v)
		}
	}

	// Test overflow behaviour.
	v := gen.String()
	if v != fmt.Sprintf(fmt.Sprintf("%%0%dd", width), start) {
		t.Errorf("Sequence returned invalid ID: want %d, got %s", max, v)
	}
}

func TestSequencePanic(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Sequence with max < start did not panic")
			}
		}()

		New(Sequence(2, 1, 0))
	}()
}

func TestSequenceUint64Overflow(t *testing.T) {
	val := uint64(math.MaxUint64)
	gen := New(sequence{
		start: 1,
		max:   uint64(math.MaxUint64),
		width: 0,
		curr:  &val,
	})

	v := gen.String()
	if v != "1" {
		t.Errorf("Sequence returned invalid ID: want 1, got %s", v)
	}
}

func TestShuffle(t *testing.T) {
	gen := New(Shuffle(Literal("a"), Literal("b"), Literal("c")))

	hitmap := map[string]bool{
		"abc": false,
		"acb": false,
		"bac": false,
		"bca": false,
		"cab": false,
		"cba": false,
	}

	for i := 0; i < 100; i++ {
		v := gen.String()
		if _, ok := hitmap[v]; !ok {
			t.Errorf("Shuffle returned invalid permutation with Literals \"a\", \"b\", \"c\": got %s", strconv.Quote(string(v)))
		}
		hitmap[v] = true
	}

	for v, found := range hitmap {
		if !found {
			t.Errorf("Shuffle with Literals \"a\", \"b\", \"c\" never returned %s", strconv.Quote(v))
		}
	}
}

// Not a real test, just a way to preview generated strings.
func TestPreviewID(t *testing.T) {
	t.Skip()

	gen := New(
		Repeat(2, 4,
			Repeat(5, 5, OneOfByte([]byte("1234567890"))),
			Literal("-"),
			Repeat(2, 4, OneOfRune([]rune("あいうえお"))),
			Literal("-"),
		),
		Potentially(0.3, OneOfString([]string{"aaaa", "bbbb", "cccc", "dddd"})),
		Literal("-"),
		Shuffle(Literal("a"), Literal("b"), Literal("c")),
		Sequence(1, 999, 2),
	)
	for i := 0; i < 10; i++ {
		t.Error(gen.String())
	}
}
