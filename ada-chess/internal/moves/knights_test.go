package moves

import (
	"testing"

	"github.com/WilliamDann/AdaEngine/ada-chess/internal/board"
)

func TestKnightMoves(t *testing.T) {
	tests := []struct {
		name   string
		square board.Square
		expect board.Bitboard
	}{
		// corner — only 2 moves
		{"a1", sq(0, 0), bb(
			sq(2, 1), // NNE
			sq(1, 2), // ENE
		)},
		// corner — only 2 moves
		{"h8", sq(7, 7), bb(
			sq(5, 6), // SSW (well, from h8 perspective)
			sq(6, 5),
		)},
		// edge — 4 moves
		{"a4", sq(3, 0), bb(
			sq(5, 1), // NNE
			sq(4, 2), // ENE
			sq(2, 2), // ESE
			sq(1, 1), // SSE
		)},
		// center — all 8 moves
		{"d4", sq(3, 3), bb(
			sq(5, 4), // NNE
			sq(4, 5), // ENE
			sq(2, 5), // ESE
			sq(1, 4), // SSE
			sq(1, 2), // SSW
			sq(2, 1), // WSW
			sq(4, 1), // WNW
			sq(5, 2), // NNW
		)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := KnightMoves(tc.square)
			if got != tc.expect {
				t.Errorf("KnightMoves(%s)\ngot:\n%s\nwant:\n%s", tc.name, got.String(), tc.expect.String())
			}
		})
	}
}

func BenchmarkKnightMoves(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for s := board.Square(0); s < 64; s++ {
			KnightMoves(s)
		}
	}
}
