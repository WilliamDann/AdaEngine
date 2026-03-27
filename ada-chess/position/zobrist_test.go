package position

import (
	"testing"

	"github.com/WilliamDann/AdaEngine/ada-chess/core"
)

func TestHashDeterministic(t *testing.T) {
	pos1 := setupPosition(map[core.Square]core.Piece{
		core.NewSquare(0, 4): core.NewPiece(core.King, core.White),
		core.NewSquare(7, 4): core.NewPiece(core.King, core.Black),
		core.NewSquare(1, 4): core.NewPiece(core.Pawn, core.White),
	}, core.White, NoCastling, core.InvalidSquare)
	pos1.Hash = ComputeHash(pos1)

	pos2 := setupPosition(map[core.Square]core.Piece{
		core.NewSquare(0, 4): core.NewPiece(core.King, core.White),
		core.NewSquare(7, 4): core.NewPiece(core.King, core.Black),
		core.NewSquare(1, 4): core.NewPiece(core.Pawn, core.White),
	}, core.White, NoCastling, core.InvalidSquare)
	pos2.Hash = ComputeHash(pos2)

	if pos1.Hash != pos2.Hash {
		t.Errorf("same position should produce same hash: %x != %x", pos1.Hash, pos2.Hash)
	}
}

func TestHashDiffersForDifferentPositions(t *testing.T) {
	pos1 := setupPosition(map[core.Square]core.Piece{
		core.NewSquare(0, 4): core.NewPiece(core.King, core.White),
		core.NewSquare(7, 4): core.NewPiece(core.King, core.Black),
		core.NewSquare(1, 4): core.NewPiece(core.Pawn, core.White),
	}, core.White, NoCastling, core.InvalidSquare)
	pos1.Hash = ComputeHash(pos1)

	pos2 := setupPosition(map[core.Square]core.Piece{
		core.NewSquare(0, 4): core.NewPiece(core.King, core.White),
		core.NewSquare(7, 4): core.NewPiece(core.King, core.Black),
		core.NewSquare(2, 4): core.NewPiece(core.Pawn, core.White), // pawn on e3 instead of e2
	}, core.White, NoCastling, core.InvalidSquare)
	pos2.Hash = ComputeHash(pos2)

	if pos1.Hash == pos2.Hash {
		t.Error("different positions should produce different hashes")
	}
}

func TestHashDiffersForSideToMove(t *testing.T) {
	pieces := map[core.Square]core.Piece{
		core.NewSquare(0, 4): core.NewPiece(core.King, core.White),
		core.NewSquare(7, 4): core.NewPiece(core.King, core.Black),
	}
	pos1 := setupPosition(pieces, core.White, NoCastling, core.InvalidSquare)
	pos1.Hash = ComputeHash(pos1)

	pos2 := setupPosition(pieces, core.Black, NoCastling, core.InvalidSquare)
	pos2.Hash = ComputeHash(pos2)

	if pos1.Hash == pos2.Hash {
		t.Error("hash should differ when side to move differs")
	}
}

func TestHashDiffersForCastlingRights(t *testing.T) {
	pieces := map[core.Square]core.Piece{
		core.NewSquare(0, 4): core.NewPiece(core.King, core.White),
		core.NewSquare(7, 4): core.NewPiece(core.King, core.Black),
	}
	pos1 := setupPosition(pieces, core.White, AllCastling, core.InvalidSquare)
	pos1.Hash = ComputeHash(pos1)

	pos2 := setupPosition(pieces, core.White, NoCastling, core.InvalidSquare)
	pos2.Hash = ComputeHash(pos2)

	if pos1.Hash == pos2.Hash {
		t.Error("hash should differ when castling rights differ")
	}
}

func TestHashDiffersForEnPassant(t *testing.T) {
	pieces := map[core.Square]core.Piece{
		core.NewSquare(0, 4): core.NewPiece(core.King, core.White),
		core.NewSquare(7, 4): core.NewPiece(core.King, core.Black),
	}
	pos1 := setupPosition(pieces, core.White, NoCastling, core.InvalidSquare)
	pos1.Hash = ComputeHash(pos1)

	pos2 := setupPosition(pieces, core.White, NoCastling, core.NewSquare(2, 4))
	pos2.Hash = ComputeHash(pos2)

	if pos1.Hash == pos2.Hash {
		t.Error("hash should differ when en passant square differs")
	}
}

func TestIncrementalHashNormalMove(t *testing.T) {
	e2 := core.NewSquare(1, 4)
	e4 := core.NewSquare(3, 4)
	pos := setupPosition(map[core.Square]core.Piece{
		e2:                  core.NewPiece(core.Pawn, core.White),
		core.NewSquare(0, 4): core.NewPiece(core.King, core.White),
		core.NewSquare(7, 4): core.NewPiece(core.King, core.Black),
	}, core.White, NoCastling, core.InvalidSquare)
	pos.Hash = ComputeHash(pos)

	next := MakeMove(pos, core.NewMove(e2, e4))
	expected := ComputeHash(next)
	if next.Hash != expected {
		t.Errorf("incremental hash %x != from-scratch %x", next.Hash, expected)
	}
}

func TestIncrementalHashCapture(t *testing.T) {
	d4 := core.NewSquare(3, 3)
	e5 := core.NewSquare(4, 4)
	pos := setupPosition(map[core.Square]core.Piece{
		d4:                  core.NewPiece(core.Knight, core.White),
		e5:                  core.NewPiece(core.Pawn, core.Black),
		core.NewSquare(0, 4): core.NewPiece(core.King, core.White),
		core.NewSquare(7, 4): core.NewPiece(core.King, core.Black),
	}, core.White, NoCastling, core.InvalidSquare)
	pos.Hash = ComputeHash(pos)

	next := MakeMove(pos, core.NewMove(d4, e5))
	expected := ComputeHash(next)
	if next.Hash != expected {
		t.Errorf("incremental hash %x != from-scratch %x after capture", next.Hash, expected)
	}
}

func TestIncrementalHashPromotion(t *testing.T) {
	e7 := core.NewSquare(6, 4)
	e8 := core.NewSquare(7, 4)
	pos := setupPosition(map[core.Square]core.Piece{
		e7:                  core.NewPiece(core.Pawn, core.White),
		core.NewSquare(0, 4): core.NewPiece(core.King, core.White),
		core.NewSquare(7, 0): core.NewPiece(core.King, core.Black),
	}, core.White, NoCastling, core.InvalidSquare)
	pos.Hash = ComputeHash(pos)

	next := MakeMove(pos, core.NewPromotion(e7, e8, core.Queen))
	expected := ComputeHash(next)
	if next.Hash != expected {
		t.Errorf("incremental hash %x != from-scratch %x after promotion", next.Hash, expected)
	}
}

func TestIncrementalHashEnPassant(t *testing.T) {
	e5 := core.NewSquare(4, 4)
	d5 := core.NewSquare(4, 3)
	d6 := core.NewSquare(5, 3)
	pos := setupPosition(map[core.Square]core.Piece{
		e5:                  core.NewPiece(core.Pawn, core.White),
		d5:                  core.NewPiece(core.Pawn, core.Black),
		core.NewSquare(0, 4): core.NewPiece(core.King, core.White),
		core.NewSquare(7, 4): core.NewPiece(core.King, core.Black),
	}, core.White, NoCastling, d6)
	pos.Hash = ComputeHash(pos)

	next := MakeMove(pos, core.NewEnPassant(e5, d6))
	expected := ComputeHash(next)
	if next.Hash != expected {
		t.Errorf("incremental hash %x != from-scratch %x after en passant", next.Hash, expected)
	}
}

func TestIncrementalHashCastling(t *testing.T) {
	tests := []struct {
		name    string
		king    core.Square
		rook    core.Square
		kingTo  core.Square
		color   core.Color
	}{
		{"white kingside", core.NewSquare(0, 4), core.NewSquare(0, 7), core.NewSquare(0, 6), core.White},
		{"white queenside", core.NewSquare(0, 4), core.NewSquare(0, 0), core.NewSquare(0, 2), core.White},
		{"black kingside", core.NewSquare(7, 4), core.NewSquare(7, 7), core.NewSquare(7, 6), core.Black},
		{"black queenside", core.NewSquare(7, 4), core.NewSquare(7, 0), core.NewSquare(7, 2), core.Black},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			otherKing := core.NewSquare(7, 4)
			otherColor := core.Black
			if tt.color == core.Black {
				otherKing = core.NewSquare(0, 4)
				otherColor = core.White
			}
			pos := setupPosition(map[core.Square]core.Piece{
				tt.king:   core.NewPiece(core.King, tt.color),
				tt.rook:   core.NewPiece(core.Rook, tt.color),
				otherKing: core.NewPiece(core.King, otherColor),
			}, tt.color, AllCastling, core.InvalidSquare)
			pos.Hash = ComputeHash(pos)

			next := MakeMove(pos, core.NewCastling(tt.king, tt.kingTo))
			expected := ComputeHash(next)
			if next.Hash != expected {
				t.Errorf("incremental hash %x != from-scratch %x", next.Hash, expected)
			}
		})
	}
}

func TestTranspositionHashMatch(t *testing.T) {
	// Reach the same position via two different move orders using knights:
	// Path 1: 1.Nc3 Nc6 2.Nf3 Nf6
	// Path 2: 1.Nf3 Nf6 2.Nc3 Nc6
	wKing := core.NewPiece(core.King, core.White)
	bKing := core.NewPiece(core.King, core.Black)

	b1 := core.NewSquare(0, 1) // white queenside knight
	g1 := core.NewSquare(0, 6) // white kingside knight
	b8 := core.NewSquare(7, 1) // black queenside knight
	g8 := core.NewSquare(7, 6) // black kingside knight
	c3 := core.NewSquare(2, 2)
	f3 := core.NewSquare(2, 5)
	c6 := core.NewSquare(5, 2)
	f6 := core.NewSquare(5, 5)
	e1 := core.NewSquare(0, 4)
	e8 := core.NewSquare(7, 4)

	wN := core.NewPiece(core.Knight, core.White)
	bN := core.NewPiece(core.Knight, core.Black)

	pieces := map[core.Square]core.Piece{
		b1: wN, g1: wN, b8: bN, g8: bN, e1: wKing, e8: bKing,
	}

	// Path 1: 1.Nc3 Nc6 2.Nf3 Nf6
	pos1 := setupPosition(pieces, core.White, NoCastling, core.InvalidSquare)
	pos1.Hash = ComputeHash(pos1)
	pos1 = MakeMove(pos1, core.NewMove(b1, c3))
	pos1 = MakeMove(pos1, core.NewMove(b8, c6))
	pos1 = MakeMove(pos1, core.NewMove(g1, f3))
	pos1 = MakeMove(pos1, core.NewMove(g8, f6))

	// Path 2: 1.Nf3 Nf6 2.Nc3 Nc6
	pos2 := setupPosition(pieces, core.White, NoCastling, core.InvalidSquare)
	pos2.Hash = ComputeHash(pos2)
	pos2 = MakeMove(pos2, core.NewMove(g1, f3))
	pos2 = MakeMove(pos2, core.NewMove(g8, f6))
	pos2 = MakeMove(pos2, core.NewMove(b1, c3))
	pos2 = MakeMove(pos2, core.NewMove(b8, c6))

	if pos1.Hash != pos2.Hash {
		t.Errorf("same position via different move orders should have same hash: %x != %x", pos1.Hash, pos2.Hash)
	}
}

func TestIncrementalHashMultipleMovesDeep(t *testing.T) {
	// Play several moves and verify the hash stays consistent at each step.
	wKing := core.NewPiece(core.King, core.White)
	bKing := core.NewPiece(core.King, core.Black)
	pos := setupPosition(map[core.Square]core.Piece{
		core.NewSquare(0, 4): wKing,
		core.NewSquare(7, 4): bKing,
		core.NewSquare(1, 4): core.NewPiece(core.Pawn, core.White),
		core.NewSquare(0, 1): core.NewPiece(core.Knight, core.White),
		core.NewSquare(6, 4): core.NewPiece(core.Pawn, core.Black),
		core.NewSquare(7, 1): core.NewPiece(core.Knight, core.Black),
	}, core.White, NoCastling, core.InvalidSquare)
	pos.Hash = ComputeHash(pos)

	moves := []core.Move{
		core.NewMove(core.NewSquare(1, 4), core.NewSquare(3, 4)), // e4
		core.NewMove(core.NewSquare(6, 4), core.NewSquare(4, 4)), // e5
		core.NewMove(core.NewSquare(0, 1), core.NewSquare(2, 2)), // Nc3
		core.NewMove(core.NewSquare(7, 1), core.NewSquare(5, 2)), // Nc6
	}

	for i, m := range moves {
		pos = MakeMove(pos, m)
		expected := ComputeHash(pos)
		if pos.Hash != expected {
			t.Errorf("move %d: incremental hash %x != from-scratch %x", i, pos.Hash, expected)
		}
	}
}
