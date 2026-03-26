package game_test

import (
	"testing"

	"github.com/WilliamDann/AdaEngine/ada-chess/internal/board"
	"github.com/WilliamDann/AdaEngine/ada-chess/internal/fen"
)

func TestPins_RookPinOnFile(t *testing.T) {
	// White king e1, white knight e4, black rook e8
	pos, _ := fen.Parse("4r3/8/8/8/4N3/8/8/4K2k w - - 0 1")
	pins := pos.ComputePins()

	e4 := sq(3, 4)
	if !pins.Pinned.Check(e4) {
		t.Fatal("knight on e4 should be pinned")
	}

	// Pin ray should include squares between king and rook, plus the rook
	ray := pins.Rays[e4]
	e8 := sq(7, 4)
	if !ray.Check(e8) {
		t.Error("pin ray should include the pinner (e8)")
	}
	// Squares between: e2, e3, e5, e6, e7
	for _, s := range []board.Square{sq(1, 4), sq(2, 4), sq(4, 4), sq(5, 4), sq(6, 4)} {
		if !ray.Check(s) {
			t.Errorf("pin ray should include %v", s)
		}
	}
	// Off-ray squares should not be included
	if ray.Check(sq(3, 3)) {
		t.Error("d4 should not be in pin ray")
	}
}

func TestPins_RookPinOnRank(t *testing.T) {
	// White king e1, white bishop c1, black rook a1
	pos, _ := fen.Parse("7k/8/8/8/8/8/8/r1B1K3 w - - 0 1")
	pins := pos.ComputePins()

	c1 := sq(0, 2)
	if !pins.Pinned.Check(c1) {
		t.Fatal("bishop on c1 should be pinned by rook on a1")
	}

	ray := pins.Rays[c1]
	a1 := sq(0, 0)
	if !ray.Check(a1) {
		t.Error("pin ray should include the pinner (a1)")
	}
	// b1 and d1 are between
	if !ray.Check(sq(0, 1)) {
		t.Error("pin ray should include b1")
	}
	if !ray.Check(sq(0, 3)) {
		t.Error("pin ray should include d1")
	}
}

func TestPins_BishopPinOnDiagonal(t *testing.T) {
	// White king e1, white knight d2, black bishop a5
	pos, _ := fen.Parse("7k/8/8/b7/8/8/3N4/4K3 w - - 0 1")
	pins := pos.ComputePins()

	d2 := sq(1, 3)
	if !pins.Pinned.Check(d2) {
		t.Fatal("knight on d2 should be pinned by bishop on a5")
	}

	ray := pins.Rays[d2]
	a5 := sq(4, 0)
	if !ray.Check(a5) {
		t.Error("pin ray should include the pinner (a5)")
	}
	// c3 and b4 are between d2 and a5
	if !ray.Check(sq(2, 2)) {
		t.Error("pin ray should include c3")
	}
	if !ray.Check(sq(3, 1)) {
		t.Error("pin ray should include b4")
	}
}

func TestPins_QueenPinsLikeRook(t *testing.T) {
	// White king e1, white rook e5, black queen e8
	pos, _ := fen.Parse("4q3/8/8/4R3/8/8/8/4K2k w - - 0 1")
	pins := pos.ComputePins()

	e5 := sq(4, 4)
	if !pins.Pinned.Check(e5) {
		t.Fatal("rook on e5 should be pinned by queen on e8")
	}
}

func TestPins_QueenPinsLikeBishop(t *testing.T) {
	// White king e1, white pawn f2, black queen h4
	pos, _ := fen.Parse("7k/8/8/8/7q/8/5P2/4K3 w - - 0 1")
	pins := pos.ComputePins()

	f2 := sq(1, 5)
	if !pins.Pinned.Check(f2) {
		t.Fatal("pawn on f2 should be pinned by queen on h4")
	}

	ray := pins.Rays[f2]
	h4 := sq(3, 7)
	if !ray.Check(h4) {
		t.Error("pin ray should include the pinner (h4)")
	}
	if !ray.Check(sq(2, 6)) {
		t.Error("pin ray should include g3")
	}
}

func TestPins_NoPinTwoPiecesBetween(t *testing.T) {
	// White king e1, white knight e3, white bishop e5, black rook e8
	pos, _ := fen.Parse("4r3/8/8/4B3/8/4N3/8/4K2k w - - 0 1")
	pins := pos.ComputePins()

	if pins.Pinned.Count() != 0 {
		t.Errorf("no pins expected with two pieces between, got %d", pins.Pinned.Count())
	}
}

func TestPins_NoPinEnemyBlocks(t *testing.T) {
	// White king e1, black pawn e3 blocks the line, white knight e5, black rook e8
	pos, _ := fen.Parse("4r3/8/8/4N3/8/4p3/8/4K2k w - - 0 1")
	pins := pos.ComputePins()

	if pins.Pinned.Count() != 0 {
		t.Errorf("no pins expected when enemy piece blocks, got %d", pins.Pinned.Count())
	}
}

func TestPins_NoPinStartingPosition(t *testing.T) {
	pos, _ := fen.Parse("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	pins := pos.ComputePins()

	if pins.Pinned.Count() != 0 {
		t.Errorf("no pins in starting position, got %d", pins.Pinned.Count())
	}
}

func TestPins_MultiplePins(t *testing.T) {
	// White king e1, pinned by rook on e8 through e4, and by bishop on a5 through d2
	pos, _ := fen.Parse("4r3/8/8/b7/4N3/8/3N4/4K2k w - - 0 1")
	pins := pos.ComputePins()

	if pins.Pinned.Count() != 2 {
		t.Fatalf("expected 2 pins, got %d", pins.Pinned.Count())
	}
	if !pins.Pinned.Check(sq(3, 4)) {
		t.Error("knight on e4 should be pinned")
	}
	if !pins.Pinned.Check(sq(1, 3)) {
		t.Error("knight on d2 should be pinned")
	}
}

func TestPins_BlackPins(t *testing.T) {
	// Black to move. Black king e8, black knight e5, white rook e1
	pos, _ := fen.Parse("4k3/8/8/4n3/8/8/8/4R2K b - - 0 1")
	pins := pos.ComputePins()

	e5 := sq(4, 4)
	if !pins.Pinned.Check(e5) {
		t.Fatal("black knight on e5 should be pinned by white rook on e1")
	}
}
