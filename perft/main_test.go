package perft

// performance test, move path enumeration
// https://www.chessprogramming.org/Perft

import (
	"fmt"
	"testing"

	"github.com/WilliamDann/adachess/game"
)

// sanity checks
func TestStartingPosition(t *testing.T) {
	pos := game.NewStartingPosition()
	value := Perft(pos, 1)

	if value.Nodes != 20 {
		t.Errorf("Starting position expected 20, for %d", value)
	}
}

// test if castling is generated
func TestCastlingPosition(t *testing.T) {
	pos := game.NewPosition("r1bqk2r/pppp1ppp/2n2n2/2b1p3/2B1P3/3P1N2/PPP2PPP/RNBQK2R w KQkq - 1 5")
	moves := pos.LegalMoves()

	found := false
	for _, move := range moves {
		if move.Castle && move.Side == &game.Kingside {
			found = true
		}
	}

	if !found {
		t.Errorf("Castling was not generated in position")
	}
}

// test if en passant is generated
func TestEnpassant(t *testing.T) {
	pos := game.NewPosition("7k/8/8/3pP3/8/8/8/7K w - d6 0 2")
	moves := pos.LegalMoves()

	found := false
	for _, move := range moves {
		if move.Capture && move.Piece.Type == game.Pawn && move.To.Equ(*game.NewCoordSan("d6")) {
			found = true
		}
	}

	if !found {
		t.Errorf("Castling was not generated in position. generated moves:")
		fmt.Println(moves)
	}
}

// perft positions
var perfts map[string]PerftFile = map[string]PerftFile{
	"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1": ReadPerftFile("perft_1.csv"),
}

func TestPerftDepth1(t *testing.T) {
	for fen, file := range perfts {
		pos := game.NewPosition(fen)
		expects := NewPerftResultsFromFile(file)

		depth := 0
		got := Perft(pos, depth+1)
		if !expects[depth].Equ(got) {
			t.Errorf("Perft info mismatch depth=%d position=%s expected=%d got=%d", depth, pos, expects[depth], got)
		}
	}
}

func TestPerftDepth2(t *testing.T) {
	for fen, file := range perfts {
		pos := game.NewPosition(fen)
		expects := NewPerftResultsFromFile(file)

		depth := 1
		got := Perft(pos, depth+1)
		if !expects[depth].Equ(got) {
			t.Errorf("Perft info mismatch depth=%d position=%s expected=%d got=%d", depth, pos, expects[depth], got)
		}
	}
}

func TestPerftDepth3(t *testing.T) {
	for fen, file := range perfts {
		pos := game.NewPosition(fen)
		expects := NewPerftResultsFromFile(file)

		depth := 2
		got := Perft(pos, depth+1)
		if !expects[depth].Equ(got) {
			t.Errorf("Perft info mismatch depth=%d position=%s expected=%d got=%d", depth, pos, expects[depth], got)
		}
	}
}

func TestPerftDepth4(t *testing.T) {
	for fen, file := range perfts {
		pos := game.NewPosition(fen)
		expects := NewPerftResultsFromFile(file)

		depth := 3
		got := Perft(pos, depth+1)
		if !expects[depth].Equ(got) {
			t.Errorf("Perft info mismatch depth=%d position=%s expected=%d got=%d", depth, pos, expects[depth], got)
		}
	}
}

func TestPerftDepth(t *testing.T) {
	for fen, file := range perfts {
		pos := game.NewPosition(fen)
		expects := NewPerftResultsFromFile(file)

		for depth, _ := range file {
			got := Perft(pos, depth+1)
			if !expects[depth].Equ(got) {
				t.Errorf("Perft info mismatch depth=%d position=%s expected=%d got=%d", depth, pos, expects[depth], got)
			}
		}
	}
}
