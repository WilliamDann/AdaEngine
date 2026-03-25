package board
import "testing"

var setTests = []struct {
	name   string
	square Square
	expect Bitboard
}{
	{"min-square", NewSquare(0, 0), 0b0000000000000000000000000000000000000000000000000000000000000001},
	{"mid-square", NewSquare(3, 3), 0b0000000000000000000000000000000000001000000000000000000000000000},
	{"max-square", NewSquare(7, 7), 0b1000000000000000000000000000000000000000000000000000000000000000},
}

func TestSet(t *testing.T) {
	for _, tt := range setTests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBitboard().Set(tt.square)
			if got != tt.expect {
				t.Errorf("Set(%v) = %064b, want %064b", tt.square, got, tt.expect)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	for _, tt := range setTests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBitboard().Set(tt.square).Check(tt.square)
			if got != true {
				t.Errorf("Check(%v) = %t, want %t", tt.square, got, true)
			}
		})
	}
}

func TestClear(t *testing.T) {
	for _, tt := range setTests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBitboard().Set(tt.square).Clear(tt.square)
			if got != 0 {
				t.Errorf("Set(%v) = %064b, want %064b", tt.square, got, uint64(0))
			}
		})
	}
}

func TestCount(t *testing.T) {
	count := 0
	bb := NewBitboard()

	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			bb = bb.Set(NewSquare(rank, file))
			count += 1

			if bb.Count() != count {
				t.Errorf("Count() = %d got %d", count, bb.Count())
			}
		}
	}
}

func TestSetOperations(t *testing.T) {
	a := NewBitboard().Set(NewSquare(0, 0)).Set(NewSquare(0, 1)) // squares 0,1
	b := NewBitboard().Set(NewSquare(0, 1)).Set(NewSquare(0, 2)) // squares 1,2

	tests := []struct {
		name   string
		got    Bitboard
		expect Bitboard
	}{
		// Union: 0,1 | 1,2 = 0,1,2
		{"union", a.Union(b),
			NewBitboard().Set(NewSquare(0, 0)).Set(NewSquare(0, 1)).Set(NewSquare(0, 2))},
		// Intersection: 0,1 & 1,2 = 1
		{"intersection", a.Intersection(b),
			NewBitboard().Set(NewSquare(0, 1))},
		// Difference: 0,1 ^ 1,2 = 0,2
		{"difference", a.Difference(b),
			NewBitboard().Set(NewSquare(0, 0)).Set(NewSquare(0, 2))},
		// Subtract: 0,1 &^ 1,2 = 0
		{"subtract", a.Subtract(b),
			NewBitboard().Set(NewSquare(0, 0))},
		// Subtract reverse: 1,2 &^ 0,1 = 2
		{"subtract-reverse", b.Subtract(a),
			NewBitboard().Set(NewSquare(0, 2))},
		// Invert empty = all bits set
		{"invert-empty", NewBitboard().Invert(),
			Bitboard(0xFFFFFFFFFFFFFFFF)},
		// Invert full = empty
		{"invert-full", Bitboard(0xFFFFFFFFFFFFFFFF).Invert(),
			NewBitboard()},
		// Union with self = self
		{"union-self", a.Union(a), a},
		// Intersection with self = self
		{"intersection-self", a.Intersection(a), a},
		// Difference with self = empty
		{"difference-self", a.Difference(a), NewBitboard()},
		// Subtract from self = empty
		{"subtract-self", a.Subtract(a), NewBitboard()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expect {
				t.Errorf("got %064b, want %064b", tt.got, tt.expect)
			}
		})
	}
}

func TestSquares(t *testing.T) {
	bb := NewBitboard()
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			bb = bb.Set(NewSquare(rank, file))
		}
	}

	count := 0
	for range bb.Squares() {
		count += 1
	}

	if count != 64 {
		t.Errorf("missing squares, n=%d expect 64", count)
	}
}
