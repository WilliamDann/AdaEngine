package movegen

import (
	"testing"

	"github.com/WilliamDann/AdaEngine/ada-chess/internal/core"
)

func bb(squares ...core.Square) core.Bitboard {
	b := core.NewBitboard()
	for _, sq := range squares {
		b = b.Set(sq)
	}
	return b
}

func sq(rank, file int) core.Square {
	return core.NewSquare(rank, file)
}

func TestBoardRays(t *testing.T) {
	tests := []struct {
		name   string
		start  core.Square
		dir    Direction
		expect core.Bitboard
	}{
		// North: rank+1
		{"N from a1", sq(0, 0), North, bb(sq(1, 0), sq(2, 0), sq(3, 0), sq(4, 0), sq(5, 0), sq(6, 0), sq(7, 0))},
		{"N from a8", sq(7, 0), North, bb()},
		{"N from h8", sq(7, 7), North, bb()},
		{"N from d4", sq(3, 3), North, bb(sq(4, 3), sq(5, 3), sq(6, 3), sq(7, 3))},
		{"N from d1", sq(0, 3), North, bb(sq(1, 3), sq(2, 3), sq(3, 3), sq(4, 3), sq(5, 3), sq(6, 3), sq(7, 3))},

		// South: rank-1
		{"S from a8", sq(7, 0), South, bb(sq(6, 0), sq(5, 0), sq(4, 0), sq(3, 0), sq(2, 0), sq(1, 0), sq(0, 0))},
		{"S from a1", sq(0, 0), South, bb()},
		{"S from h1", sq(0, 7), South, bb()},
		{"S from d4", sq(3, 3), South, bb(sq(2, 3), sq(1, 3), sq(0, 3))},
		{"S from h8", sq(7, 7), South, bb(sq(6, 7), sq(5, 7), sq(4, 7), sq(3, 7), sq(2, 7), sq(1, 7), sq(0, 7))},

		// East: file+1
		{"E from a1", sq(0, 0), East, bb(sq(0, 1), sq(0, 2), sq(0, 3), sq(0, 4), sq(0, 5), sq(0, 6), sq(0, 7))},
		{"E from h1", sq(0, 7), East, bb()},
		{"E from h8", sq(7, 7), East, bb()},
		{"E from d4", sq(3, 3), East, bb(sq(3, 4), sq(3, 5), sq(3, 6), sq(3, 7))},
		{"E from a4", sq(3, 0), East, bb(sq(3, 1), sq(3, 2), sq(3, 3), sq(3, 4), sq(3, 5), sq(3, 6), sq(3, 7))},

		// West: file-1
		{"W from h1", sq(0, 7), West, bb(sq(0, 6), sq(0, 5), sq(0, 4), sq(0, 3), sq(0, 2), sq(0, 1), sq(0, 0))},
		{"W from a1", sq(0, 0), West, bb()},
		{"W from a8", sq(7, 0), West, bb()},
		{"W from d4", sq(3, 3), West, bb(sq(3, 2), sq(3, 1), sq(3, 0))},
		{"W from h8", sq(7, 7), West, bb(sq(7, 6), sq(7, 5), sq(7, 4), sq(7, 3), sq(7, 2), sq(7, 1), sq(7, 0))},

		// NE: rank+1, file+1
		{"NE from a1", sq(0, 0), NE, bb(sq(1, 1), sq(2, 2), sq(3, 3), sq(4, 4), sq(5, 5), sq(6, 6), sq(7, 7))},
		{"NE from h8", sq(7, 7), NE, bb()},
		{"NE from h1", sq(0, 7), NE, bb()},
		{"NE from a8", sq(7, 0), NE, bb()},
		{"NE from d4", sq(3, 3), NE, bb(sq(4, 4), sq(5, 5), sq(6, 6), sq(7, 7))},
		{"NE from a4", sq(3, 0), NE, bb(sq(4, 1), sq(5, 2), sq(6, 3), sq(7, 4))},

		// NW: rank+1, file-1
		{"NW from h1", sq(0, 7), NW, bb(sq(1, 6), sq(2, 5), sq(3, 4), sq(4, 3), sq(5, 2), sq(6, 1), sq(7, 0))},
		{"NW from a1", sq(0, 0), NW, bb()},
		{"NW from a8", sq(7, 0), NW, bb()},
		{"NW from h8", sq(7, 7), NW, bb()},
		{"NW from d4", sq(3, 3), NW, bb(sq(4, 2), sq(5, 1), sq(6, 0))},
		{"NW from h4", sq(3, 7), NW, bb(sq(4, 6), sq(5, 5), sq(6, 4), sq(7, 3))},

		// SE: rank-1, file+1
		{"SE from a8", sq(7, 0), SE, bb(sq(6, 1), sq(5, 2), sq(4, 3), sq(3, 4), sq(2, 5), sq(1, 6), sq(0, 7))},
		{"SE from h1", sq(0, 7), SE, bb()},
		{"SE from a1", sq(0, 0), SE, bb()},
		{"SE from h8", sq(7, 7), SE, bb()},
		{"SE from d4", sq(3, 3), SE, bb(sq(2, 4), sq(1, 5), sq(0, 6))},
		{"SE from a4", sq(3, 0), SE, bb(sq(2, 1), sq(1, 2), sq(0, 3))},

		// SW: rank-1, file-1
		{"SW from h8", sq(7, 7), SW, bb(sq(6, 6), sq(5, 5), sq(4, 4), sq(3, 3), sq(2, 2), sq(1, 1), sq(0, 0))},
		{"SW from a1", sq(0, 0), SW, bb()},
		{"SW from h1", sq(0, 7), SW, bb()},
		{"SW from a8", sq(7, 0), SW, bb()},
		{"SW from d4", sq(3, 3), SW, bb(sq(2, 2), sq(1, 1), sq(0, 0))},
		{"SW from h4", sq(3, 7), SW, bb(sq(2, 6), sq(1, 5), sq(0, 4))},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := boardRay(tc.start, tc.dir)
			// t.Logf("boardRay(%s, %s)\n%s", tc.start, tc.name, got.String())
			if got != tc.expect {
				t.Errorf("expected:\n%s", tc.expect.String())
			}
		})
	}
}

func TestRookRay(t *testing.T) {
	tests := []struct {
		name   string
		start  core.Square
		expect core.Bitboard
	}{
		// center
		{"d4", sq(3, 3), bb(
			sq(4, 3), sq(5, 3), sq(6, 3), // N (d5-d7)
			sq(2, 3), sq(1, 3), // S (d3-d2)
			sq(3, 4), sq(3, 5), sq(3, 6), // E (e4-g4)
			sq(3, 2), sq(3, 1), // W (c4-b4)
		)},
		// corners
		{"a1", sq(0, 0), bb(
			sq(1, 0), sq(2, 0), sq(3, 0), sq(4, 0), sq(5, 0), sq(6, 0), // N (a2-a7)
			sq(0, 1), sq(0, 2), sq(0, 3), sq(0, 4), sq(0, 5), sq(0, 6), // E (b1-g1)
		)},
		{"h8", sq(7, 7), bb(
			sq(6, 7), sq(5, 7), sq(4, 7), sq(3, 7), sq(2, 7), sq(1, 7), // S (h7-h2)
			sq(7, 6), sq(7, 5), sq(7, 4), sq(7, 3), sq(7, 2), sq(7, 1), // W (g8-b8)
		)},
		{"a8", sq(7, 0), bb(
			sq(6, 0), sq(5, 0), sq(4, 0), sq(3, 0), sq(2, 0), sq(1, 0), // S (a7-a2)
			sq(7, 1), sq(7, 2), sq(7, 3), sq(7, 4), sq(7, 5), sq(7, 6), // E (b8-g8)
		)},
		{"h1", sq(0, 7), bb(
			sq(1, 7), sq(2, 7), sq(3, 7), sq(4, 7), sq(5, 7), sq(6, 7), // N (h2-h7)
			sq(0, 6), sq(0, 5), sq(0, 4), sq(0, 3), sq(0, 2), sq(0, 1), // W (g1-b1)
		)},
		// edges
		{"a4", sq(3, 0), bb(
			sq(4, 0), sq(5, 0), sq(6, 0), // N (a5-a7)
			sq(2, 0), sq(1, 0), // S (a3-a2)
			sq(3, 1), sq(3, 2), sq(3, 3), sq(3, 4), sq(3, 5), sq(3, 6), // E (b4-g4)
		)},
		{"d1", sq(0, 3), bb(
			sq(1, 3), sq(2, 3), sq(3, 3), sq(4, 3), sq(5, 3), sq(6, 3), // N (d2-d7)
			sq(0, 4), sq(0, 5), sq(0, 6), // E (e1-g1)
			sq(0, 2), sq(0, 1), // W (c1-b1)
		)},
		{"h4", sq(3, 7), bb(
			sq(4, 7), sq(5, 7), sq(6, 7), // N (h5-h7)
			sq(2, 7), sq(1, 7), // S (h3-h2)
			sq(3, 6), sq(3, 5), sq(3, 4), sq(3, 3), sq(3, 2), sq(3, 1), // W (g4-b4)
		)},
		{"d8", sq(7, 3), bb(
			sq(6, 3), sq(5, 3), sq(4, 3), sq(3, 3), sq(2, 3), sq(1, 3), // S (d7-d2)
			sq(7, 4), sq(7, 5), sq(7, 6), // E (e8-g8)
			sq(7, 2), sq(7, 1), // W (c8-b8)
		)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := rookMask(tc.start)
			if got != tc.expect {
				t.Errorf("rookMask(%s)\ngot:\n%s\nwant:\n%s", tc.start, got.String(), tc.expect.String())
			}
		})
	}
}

func TestBishopRay(t *testing.T) {
	tests := []struct {
		name   string
		start  core.Square
		expect core.Bitboard
	}{
		// center
		{"d4", sq(3, 3), bb(
			sq(4, 4), sq(5, 5), sq(6, 6), // NE (e5-g7)
			sq(4, 2), sq(5, 1), // NW (c5-b6)
			sq(2, 4), sq(1, 5), // SE (e3-f2)
			sq(2, 2), sq(1, 1), // SW (c3-b2)
		)},
		// corners
		{"a1", sq(0, 0), bb(
			sq(1, 1), sq(2, 2), sq(3, 3), sq(4, 4), sq(5, 5), sq(6, 6), // NE (b2-g7)
		)},
		{"h8", sq(7, 7), bb(
			sq(6, 6), sq(5, 5), sq(4, 4), sq(3, 3), sq(2, 2), sq(1, 1), // SW (g7-b2)
		)},
		{"a8", sq(7, 0), bb(
			sq(6, 1), sq(5, 2), sq(4, 3), sq(3, 4), sq(2, 5), sq(1, 6), // SE (b7-g2)
		)},
		{"h1", sq(0, 7), bb(
			sq(1, 6), sq(2, 5), sq(3, 4), sq(4, 3), sq(5, 2), sq(6, 1), // NW (g2-b7)
		)},
		// edges
		{"a4", sq(3, 0), bb(
			sq(4, 1), sq(5, 2), sq(6, 3), // NE (b5-d7)
			sq(2, 1), sq(1, 2), // SE (b3-c2)
		)},
		{"d1", sq(0, 3), bb(
			sq(1, 4), sq(2, 5), sq(3, 6), // NE (e2-g4)
			sq(1, 2), sq(2, 1), // NW (c2-b3)
		)},
		{"h4", sq(3, 7), bb(
			sq(4, 6), sq(5, 5), sq(6, 4), // NW (g5-e7)
			sq(2, 6), sq(1, 5), // SW (g3-f2)
		)},
		{"d8", sq(7, 3), bb(
			sq(6, 4), sq(5, 5), sq(4, 6), // SE (e7-g5)
			sq(6, 2), sq(5, 1), // SW (c7-b6)
		)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := bishopMask(tc.start)
			if got != tc.expect {
				t.Errorf("bishopRay(%s)\ngot:\n%s\nwant:\n%s", tc.start, got.String(), tc.expect.String())
			}
		})
	}
}
