package game_test

import (
	"testing"

	"github.com/WilliamDann/AdaEngine/ada-chess/internal/board"
	"github.com/WilliamDann/AdaEngine/ada-chess/internal/fen"
)

func sq(rank, file int) board.Square {
	return board.NewSquare(rank, file)
}

func TestAttackers_Knight(t *testing.T) {
	pos, _ := fen.Parse("8/8/8/8/4N3/8/8/4K2k w - - 0 1")
	if pos.Attackers(sq(5, 3), board.White).Empty() {
		t.Error("d6 should be attacked by white knight on e4")
	}
	if pos.Attackers(sq(5, 5), board.White).Empty() {
		t.Error("f6 should be attacked by white knight on e4")
	}
	if !pos.Attackers(sq(3, 3), board.White).Empty() {
		t.Error("d4 should not be attacked by white knight on e4")
	}
}

func TestAttackers_Bishop(t *testing.T) {
	pos, _ := fen.Parse("8/8/8/8/3B4/8/8/4K2k w - - 0 1")
	if pos.Attackers(sq(4, 4), board.White).Empty() {
		t.Error("e5 should be attacked by white bishop on d4")
	}
	if pos.Attackers(sq(0, 0), board.White).Empty() {
		t.Error("a1 should be attacked by white bishop on d4")
	}
}

func TestAttackers_Rook(t *testing.T) {
	pos, _ := fen.Parse("8/8/8/8/3R4/8/8/4K2k w - - 0 1")
	if pos.Attackers(sq(3, 7), board.White).Empty() {
		t.Error("h4 should be attacked by white rook on d4")
	}
	if pos.Attackers(sq(7, 3), board.White).Empty() {
		t.Error("d8 should be attacked by white rook on d4")
	}
	if !pos.Attackers(sq(4, 4), board.White).Empty() {
		t.Error("e5 should not be attacked by rook on d4")
	}
}

func TestAttackers_RookBlocked(t *testing.T) {
	pos, _ := fen.Parse("8/8/8/8/8/P7/8/R3K2k w - - 0 1")
	if pos.Attackers(sq(1, 0), board.White).Empty() {
		t.Error("a2 should be attacked by rook on a1")
	}
	if !pos.Attackers(sq(3, 0), board.White).Empty() {
		t.Error("a4 should not be attacked by rook blocked by pawn on a3")
	}
}

func TestAttackers_Queen(t *testing.T) {
	pos, _ := fen.Parse("8/8/8/8/3Q4/8/8/4K2k w - - 0 1")
	if pos.Attackers(sq(3, 7), board.White).Empty() {
		t.Error("h4 should be attacked by queen on d4 (rook-like)")
	}
	if pos.Attackers(sq(6, 6), board.White).Empty() {
		t.Error("g7 should be attacked by queen on d4 (bishop-like)")
	}
}

func TestAttackers_WhitePawn(t *testing.T) {
	pos, _ := fen.Parse("8/8/8/8/4P3/8/8/4K2k w - - 0 1")
	if pos.Attackers(sq(4, 3), board.White).Empty() {
		t.Error("d5 should be attacked by white pawn on e4")
	}
	if pos.Attackers(sq(4, 5), board.White).Empty() {
		t.Error("f5 should be attacked by white pawn on e4")
	}
	if !pos.Attackers(sq(4, 4), board.White).Empty() {
		t.Error("e5 should not be attacked by pawn (forward is not an attack)")
	}
}

func TestAttackers_BlackPawn(t *testing.T) {
	pos, _ := fen.Parse("4k3/8/8/4p3/8/8/8/4K3 w - - 0 1")
	if pos.Attackers(sq(3, 3), board.Black).Empty() {
		t.Error("d4 should be attacked by black pawn on e5")
	}
	if pos.Attackers(sq(3, 5), board.Black).Empty() {
		t.Error("f4 should be attacked by black pawn on e5")
	}
}

func TestAttackers_King(t *testing.T) {
	pos, _ := fen.Parse("8/8/8/8/8/8/8/4K2k w - - 0 1")
	if pos.Attackers(sq(1, 4), board.White).Empty() {
		t.Error("e2 should be attacked by white king on e1")
	}
	if pos.Attackers(sq(1, 3), board.White).Empty() {
		t.Error("d2 should be attacked by white king on e1")
	}
}

func TestAttackers_Multiple(t *testing.T) {
	// Knights on c4 and g4 both attack e5
	pos, _ := fen.Parse("8/8/8/8/2N3N1/8/8/4K2k w - - 0 1")
	attackers := pos.Attackers(sq(4, 4), board.White)
	if attackers.Count() != 2 {
		t.Errorf("e5 should be attacked by 2 knights, got %d", attackers.Count())
	}
}

func TestIsAttacked(t *testing.T) {
	pos, _ := fen.Parse("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	if !pos.IsAttacked(sq(2, 4), board.White) {
		t.Error("e3 should be attacked by white in starting position")
	}
	if pos.IsAttacked(sq(3, 4), board.White) {
		t.Error("e4 should not be attacked by white in starting position")
	}
}

func TestInCheck_No(t *testing.T) {
	pos, _ := fen.Parse("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	if pos.InCheck() {
		t.Error("starting position should not be in check")
	}
}

func TestInCheck_RookCheck(t *testing.T) {
	pos, _ := fen.Parse("4r3/8/8/8/8/8/8/4K2k w - - 0 1")
	if !pos.InCheck() {
		t.Error("white king should be in check from black rook on e8")
	}
}

func TestInCheck_Blocked(t *testing.T) {
	pos, _ := fen.Parse("4r3/8/8/8/8/8/4P3/4K2k w - - 0 1")
	if pos.InCheck() {
		t.Error("white king should not be in check (pawn blocks rook)")
	}
}
