package moves

import (
	"testing"

	"github.com/WilliamDann/AdaEngine/ada-chess/internal/board"
)

func TestKingMoves(t *testing.T) {
	tests := []struct {
		name   string
		square board.Square
		expect board.Bitboard
	}{
		// corner — 3 moves
		{"a1", sq(0, 0), bb(
			sq(1, 0), // N
			sq(1, 1), // NE
			sq(0, 1), // E
		)},
		// corner — 3 moves
		{"h8", sq(7, 7), bb(
			sq(6, 7), // S
			sq(6, 6), // SW
			sq(7, 6), // W
		)},
		// edge — 5 moves
		{"a4", sq(3, 0), bb(
			sq(4, 0), // N
			sq(4, 1), // NE
			sq(3, 1), // E
			sq(2, 1), // SE
			sq(2, 0), // S
		)},
		// center — all 8 moves
		{"d4", sq(3, 3), bb(
			sq(4, 3), // N
			sq(4, 4), // NE
			sq(3, 4), // E
			sq(2, 4), // SE
			sq(2, 3), // S
			sq(2, 2), // SW
			sq(3, 2), // W
			sq(4, 2), // NW
		)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := KingMoves(tc.square)
			if got != tc.expect {
				t.Errorf("KingMoves(%s)\ngot:\n%s\nwant:\n%s", tc.name, got.String(), tc.expect.String())
			}
		})
	}
}

func BenchmarkKingMoves(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for s := board.Square(0); s < 64; s++ {
			KingMoves(s)
		}
	}
}
