package game

import (
	"testing"
)

func TestPlacePiece(t *testing.T) {
	board := NewBoard()
	values := map[Coord]Piece{
		{0, 0}: *NewPiece(White, Pawn),
		{2, 2}: *NewPiece(White, Pawn),
		{7, 7}: *NewPiece(White, Pawn),
	}

	// place pieces
	for coord, piece := range values {
		board.Set(piece, coord)
	}

	// check
	for coord, piece := range values {
		got := board.Get(coord)
		if !piece.Is(got) {
			t.Errorf("Expected piece %s got %s", piece, got)
		}
	}
}

func TestClearPiece(t *testing.T) {
	board := NewBoard()
	values := map[Coord]Piece{
		{0, 0}: *NewPiece(White, Pawn),
		{2, 2}: *NewPiece(White, Pawn),
		{4, 4}: *NewPiece(White, Pawn),
		{7, 7}: *NewPiece(White, Pawn),
	}

	// place pieces
	for coord, piece := range values {
		board.Set(piece, coord)
	}

	// check
	for coord := range values {
		board.Clear(coord)
		got := board.Get(coord)
		if !got.IsNone() {
			t.Errorf("Expected piece %s got %s", NoPiece(), got)
		}
	}
}
