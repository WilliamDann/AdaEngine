package board
import (
	"testing"
)

func TestSetBounds(t *testing.T) {
	board := NewChessboard()
	piece := NewPiece(Pawn, White)

	for x := 0; x < 8; x++ {
	for y := 0; y < 8; y++ {
		square := NewSquare(x, y)
		board.Set(square, piece)
	}
	}

	if board.pieces[piece].Count() != 64 {
		t.Errorf("expected 64 pawns, got %d", board.pieces[piece].Count())
	}
	if board.white != NewBitboard().Invert() {
		t.Errorf("all squares white expected got, \n%s", board.white.String())
	}
}

func TestClearBounds(t *testing.T) {
	board := NewChessboard()
	piece := NewPiece(Pawn, Black)

	// set all squares to contain pieces
	board.pieces[piece] = board.pieces[piece].Invert()
	board.black         = NewBitboard().Invert()

	for x := 0; x < 8; x++ {
	for y := 0; y < 8; y++ {
		square := NewSquare(x, y)
		board.Clear(square)
	}
	}

	if board.pieces[piece].Count() != 0 {
		t.Errorf("expected 0 pawns, got %d", board.pieces[piece].Count())
	}
	if board.black != 0 {
		t.Errorf("expected 0 black squares, got %s", board.black.String())
	}
}

func TestCheckBounds(t *testing.T) {
	board := NewChessboard()
	piece := NewPiece(Queen, White)

	// set all squares to contain pieces
	board.pieces[piece] = board.pieces[piece].Invert()
	board.white         = NewBitboard().Invert()
	
	for x := 0; x < 8; x++ {
	for y := 0; y < 8; y++ {
		square := NewSquare(x, y)
		if board.Check(square) != piece {
			t.Errorf("expected %s @ %s, got %s", piece.String(), square.String(), board.Check(square).String())
		}
	}
	}
}

func TestHasPieceBounds(t *testing.T) {
	board := NewChessboard()
	piece := NewPiece(Bishop, Black)

	// chek empty squares
	for x := 0; x < 8; x++ {
	for y := 0; y < 8; y++ {
		square := NewSquare(x, y)
		if board.HasPiece(square) {
			t.Errorf("unexpected piece at %s", square)
		}
	}
	}

	// set all squares to contain pieces
	board.pieces[piece] = board.pieces[piece].Invert()
	board.black         = NewBitboard().Invert()

	// check full squares
	for x := 0; x < 8; x++ {
	for y := 0; y < 8; y++ {
		square := NewSquare(x, y)
		if !board.HasPiece(square) {
			t.Errorf("expected %s @ %s, got %s", piece.String(), square.String(), board.Check(square).String())
		}
	}
	}
}

func TestHasColorPieceBounds(t *testing.T) {
	board := NewChessboard()
	piece := NewPiece(Knight, Black)

	// set all squares to contain pieces
	board.pieces[piece] = board.pieces[piece].Invert()
	board.black         = NewBitboard().Invert()

	for x := 0; x < 8; x++ {
	for y := 0; y < 8; y++ {
		square := NewSquare(x, y)
		if !board.HasColorPiece(square, Black) {
			t.Errorf("failed to find paice at %s", square.String())
		}
		if board.HasColorPiece(square, White) {
			t.Errorf("unexpected piece at %s", square.String())
		}
	}
	}
}
