package movegen

import (
	"testing"

	"github.com/WilliamDann/AdaEngine/ada-chess/internal/core"
)

func TestWhitePawnAttacks_Center(t *testing.T) {
	// e4: attacks d5, f5
	got := PawnAttacks(sq(3, 4), core.White)
	want := bb(sq(4, 3), sq(4, 5))
	if got != want {
		t.Errorf("white pawn e4\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestBlackPawnAttacks_Center(t *testing.T) {
	// e5: attacks d4, f4
	got := PawnAttacks(sq(4, 4), core.Black)
	want := bb(sq(3, 3), sq(3, 5))
	if got != want {
		t.Errorf("black pawn e5\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestWhitePawnAttacks_AFile(t *testing.T) {
	// a2: attacks b3 only (no left capture)
	got := PawnAttacks(sq(1, 0), core.White)
	want := bb(sq(2, 1))
	if got != want {
		t.Errorf("white pawn a2\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestWhitePawnAttacks_HFile(t *testing.T) {
	// h4: attacks g5 only (no right capture)
	got := PawnAttacks(sq(3, 7), core.White)
	want := bb(sq(4, 6))
	if got != want {
		t.Errorf("white pawn h4\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestBlackPawnAttacks_AFile(t *testing.T) {
	// a7: attacks b6 only
	got := PawnAttacks(sq(6, 0), core.Black)
	want := bb(sq(5, 1))
	if got != want {
		t.Errorf("black pawn a7\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestBlackPawnAttacks_HFile(t *testing.T) {
	// h5: attacks g4 only
	got := PawnAttacks(sq(4, 7), core.Black)
	want := bb(sq(3, 6))
	if got != want {
		t.Errorf("black pawn h5\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestWhitePawnAttacks_Rank8(t *testing.T) {
	// rank 8 (rank index 7): no attacks (pawn can't be here normally, but table should be empty)
	for file := 0; file < 8; file++ {
		got := PawnAttacks(sq(7, file), core.White)
		if !got.Empty() {
			t.Errorf("white pawn rank 8 file %d should have no attacks, got:\n%s", file, got)
		}
	}
}

func TestBlackPawnAttacks_Rank1(t *testing.T) {
	// rank 1 (rank index 0): no attacks
	for file := 0; file < 8; file++ {
		got := PawnAttacks(sq(0, file), core.Black)
		if !got.Empty() {
			t.Errorf("black pawn rank 1 file %d should have no attacks, got:\n%s", file, got)
		}
	}
}

func TestPawnAttacks_AllSquaresCount(t *testing.T) {
	for s := core.Square(0); s < 64; s++ {
		rank := int(s) / 8
		file := int(s) % 8

		// white
		wc := PawnAttacks(s, core.White).Count()
		wantW := 0
		if rank < 7 {
			if file > 0 {
				wantW++
			}
			if file < 7 {
				wantW++
			}
		}
		if wc != wantW {
			t.Errorf("white pawn sq=%d: got %d attacks, want %d", s, wc, wantW)
		}

		// black
		bc := PawnAttacks(s, core.Black).Count()
		wantB := 0
		if rank > 0 {
			if file > 0 {
				wantB++
			}
			if file < 7 {
				wantB++
			}
		}
		if bc != wantB {
			t.Errorf("black pawn sq=%d: got %d attacks, want %d", s, bc, wantB)
		}
	}
}

func BenchmarkPawnAttacks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for s := core.Square(0); s < 64; s++ {
			PawnAttacks(s, core.White)
			PawnAttacks(s, core.Black)
		}
	}
}
