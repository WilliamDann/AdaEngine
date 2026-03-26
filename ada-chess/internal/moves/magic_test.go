package moves

import (
	"testing"

	"github.com/WilliamDann/AdaEngine/ada-chess/internal/board"
)

func TestRookMovesOpenBoard(t *testing.T) {
	blockers := board.NewBitboard()

	tests := []struct {
		name   string
		square board.Square
		expect board.Bitboard
	}{
		{"a1", sq(0, 0), bb(
			sq(1, 0), sq(2, 0), sq(3, 0), sq(4, 0), sq(5, 0), sq(6, 0), sq(7, 0), // N
			sq(0, 1), sq(0, 2), sq(0, 3), sq(0, 4), sq(0, 5), sq(0, 6), sq(0, 7), // E
		)},
		{"d4", sq(3, 3), bb(
			sq(4, 3), sq(5, 3), sq(6, 3), sq(7, 3), // N
			sq(2, 3), sq(1, 3), sq(0, 3), // S
			sq(3, 4), sq(3, 5), sq(3, 6), sq(3, 7), // E
			sq(3, 2), sq(3, 1), sq(3, 0), // W
		)},
		{"h8", sq(7, 7), bb(
			sq(6, 7), sq(5, 7), sq(4, 7), sq(3, 7), sq(2, 7), sq(1, 7), sq(0, 7), // S
			sq(7, 6), sq(7, 5), sq(7, 4), sq(7, 3), sq(7, 2), sq(7, 1), sq(7, 0), // W
		)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := RookMoves(tc.square, blockers)
			if got != tc.expect {
				t.Errorf("RookMoves(%s, empty)\ngot:\n%s\nwant:\n%s", tc.name, got.String(), tc.expect.String())
			}
		})
	}
}

func TestRookMovesWithBlockers(t *testing.T) {
	tests := []struct {
		name     string
		square   board.Square
		blockers board.Bitboard
		expect   board.Bitboard
	}{
		// rook on d4, blocked on d6 and f4
		{"d4 blocked d6 f4", sq(3, 3),
			bb(sq(5, 3), sq(3, 5)),
			bb(
				sq(4, 3), sq(5, 3), // N (stops at d6, inclusive)
				sq(2, 3), sq(1, 3), sq(0, 3), // S
				sq(3, 4), sq(3, 5), // E (stops at f4, inclusive)
				sq(3, 2), sq(3, 1), sq(3, 0), // W
			),
		},
		// rook on a1, blocked on a3 and d1
		{"a1 blocked a3 d1", sq(0, 0),
			bb(sq(2, 0), sq(0, 3)),
			bb(
				sq(1, 0), sq(2, 0), // N (stops at a3)
				sq(0, 1), sq(0, 2), sq(0, 3), // E (stops at d1)
			),
		},
		// rook on e4, blocked on all sides one square away
		{"e4 surrounded", sq(3, 4),
			bb(sq(4, 4), sq(2, 4), sq(3, 5), sq(3, 3)),
			bb(
				sq(4, 4), // N
				sq(2, 4), // S
				sq(3, 5), // E
				sq(3, 3), // W
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := RookMoves(tc.square, tc.blockers)
			if got != tc.expect {
				t.Errorf("RookMoves\ngot:\n%s\nwant:\n%s", got.String(), tc.expect.String())
			}
		})
	}
}

func TestBishopMovesOpenBoard(t *testing.T) {
	blockers := board.NewBitboard()

	tests := []struct {
		name   string
		square board.Square
		expect board.Bitboard
	}{
		{"a1", sq(0, 0), bb(
			sq(1, 1), sq(2, 2), sq(3, 3), sq(4, 4), sq(5, 5), sq(6, 6), sq(7, 7), // NE
		)},
		{"d4", sq(3, 3), bb(
			sq(4, 4), sq(5, 5), sq(6, 6), sq(7, 7), // NE
			sq(4, 2), sq(5, 1), sq(6, 0), // NW
			sq(2, 4), sq(1, 5), sq(0, 6), // SE
			sq(2, 2), sq(1, 1), sq(0, 0), // SW
		)},
		{"h8", sq(7, 7), bb(
			sq(6, 6), sq(5, 5), sq(4, 4), sq(3, 3), sq(2, 2), sq(1, 1), sq(0, 0), // SW
		)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := BishopMoves(tc.square, blockers)
			if got != tc.expect {
				t.Errorf("BishopMoves(%s, empty)\ngot:\n%s\nwant:\n%s", tc.name, got.String(), tc.expect.String())
			}
		})
	}
}

func TestBishopMovesWithBlockers(t *testing.T) {
	tests := []struct {
		name     string
		square   board.Square
		blockers board.Bitboard
		expect   board.Bitboard
	}{
		// bishop on d4, blocked on f6 and b2
		{"d4 blocked f6 b2", sq(3, 3),
			bb(sq(5, 5), sq(1, 1)),
			bb(
				sq(4, 4), sq(5, 5), // NE (stops at f6)
				sq(4, 2), sq(5, 1), sq(6, 0), // NW
				sq(2, 4), sq(1, 5), sq(0, 6), // SE
				sq(2, 2), sq(1, 1), // SW (stops at b2)
			),
		},
		// bishop on a1, blocked on c3
		{"a1 blocked c3", sq(0, 0),
			bb(sq(2, 2)),
			bb(
				sq(1, 1), sq(2, 2), // NE (stops at c3)
			),
		},
		// bishop on e4, blocked on all diagonals one square away
		{"e4 surrounded", sq(3, 4),
			bb(sq(4, 5), sq(4, 3), sq(2, 5), sq(2, 3)),
			bb(
				sq(4, 5), // NE
				sq(4, 3), // NW
				sq(2, 5), // SE
				sq(2, 3), // SW
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := BishopMoves(tc.square, tc.blockers)
			if got != tc.expect {
				t.Errorf("BishopMoves\ngot:\n%s\nwant:\n%s", got.String(), tc.expect.String())
			}
		})
	}
}

func BenchmarkRookMovesMagic(b *testing.B) {
	blockers := bb(sq(3, 5), sq(5, 3), sq(1, 3))
	for i := 0; i < b.N; i++ {
		for s := board.Square(0); s < 64; s++ {
			RookMoves(s, blockers)
		}
	}
}

func BenchmarkRookMovesSlow(b *testing.B) {
	blockers := bb(sq(3, 5), sq(5, 3), sq(1, 3))
	for i := 0; i < b.N; i++ {
		for s := board.Square(0); s < 64; s++ {
			rookAttacks(s, blockers)
		}
	}
}

func BenchmarkBishopMovesMagic(b *testing.B) {
	blockers := bb(sq(5, 5), sq(1, 1), sq(5, 1))
	for i := 0; i < b.N; i++ {
		for s := board.Square(0); s < 64; s++ {
			BishopMoves(s, blockers)
		}
	}
}

func BenchmarkBishopMovesSlow(b *testing.B) {
	blockers := bb(sq(5, 5), sq(1, 1), sq(5, 1))
	for i := 0; i < b.N; i++ {
		for s := board.Square(0); s < 64; s++ {
			bishopAttacks(s, blockers)
		}
	}
}

// verify magic lookups match the slow ray-walking method for all squares
func TestMagicMatchesSlowRook(t *testing.T) {
	for s := board.Square(0); s < 64; s++ {
		mask := rookMask(s)
		blockers := board.Bitboard(0)
		for {
			fast := RookMoves(s, blockers)
			slow := rookAttacks(s, blockers)
			if fast != slow {
				t.Errorf("rook sq=%d blockers=0x%x\nfast:\n%s\nslow:\n%s",
					s, uint64(blockers), fast.String(), slow.String())
			}
			blockers = board.Bitboard((uint64(blockers) - uint64(mask)) & uint64(mask))
			if blockers == 0 {
				break
			}
		}
	}
}

func TestMagicMatchesSlowBishop(t *testing.T) {
	for s := board.Square(0); s < 64; s++ {
		mask := bishopMask(s)
		blockers := board.Bitboard(0)
		for {
			fast := BishopMoves(s, blockers)
			slow := bishopAttacks(s, blockers)
			if fast != slow {
				t.Errorf("bishop sq=%d blockers=0x%x\nfast:\n%s\nslow:\n%s",
					s, uint64(blockers), fast.String(), slow.String())
			}
			blockers = board.Bitboard((uint64(blockers) - uint64(mask)) & uint64(mask))
			if blockers == 0 {
				break
			}
		}
	}
}
