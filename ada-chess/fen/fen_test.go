package fen
import (
	"testing"
	"github.com/WilliamDann/AdaEngine/ada-chess/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/position"
)

var starting string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
var italian  string = "r1bqkb1r/pppp1ppp/2n2n2/4p3/2B1P3/3P1N2/PPP2PPP/RNBQK2R b KQkq - 0 4"
var ep       string = "rnbqkbnr/1pp1pppp/p7/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3"
var bcastle  string = "r1bqkb1r/ppp2ppp/2np1n2/4p3/2B1P3/3P1N2/PPP1KPPP/RNBQ3R b kq - 1 5"
var wcastle  string = "r1bq1b1r/ppppkppp/2n2n2/4p3/2B1P3/3P1N2/PPP2PPP/RNBQK2R w KQ - 1 5"
var nocastle string = "r1bq1b1r/ppppkppp/2n2n2/4p3/2B1P3/3P1N2/PPP1KPPP/RNBQ3R b - - 2 5"

func TestItalian(t *testing.T) {
	pos, err := Parse(italian)
	if err != nil {
		t.Fatal(err)
	}

	checks := []struct {
		sq core.Square
		piece core.Piece
	} {
		{ core.NewSquare(3, 2), core.NewPiece(core.Bishop, core.White) },
		{ core.NewSquare(2, 3), core.NewPiece(core.Pawn, core.White) },
		{ core.NewSquare(3, 4), core.NewPiece(core.Pawn, core.White) },
		{ core.NewSquare(2, 5), core.NewPiece(core.Knight, core.White) },

		{ core.NewSquare(5, 2), core.NewPiece(core.Knight, core.Black) },
		{ core.NewSquare(5, 5), core.NewPiece(core.Knight, core.Black) },
		{ core.NewSquare(4, 4), core.NewPiece(core.Pawn, core.Black) },
	}


	for _, tt := range checks {
		got := pos.Board.Check(tt.sq)
		if got != tt.piece {
			t.Errorf("at %s: got %s, expected %s", tt.sq.String(), got.String(), tt.piece.String())
		}
	}
}

func TestStartingPosition(t *testing.T) {
	pos, err := Parse(starting)
	if err != nil {
		t.Fatal(err)
	}

	// rank 0: RNBQKBNR
	rank0 := []core.PieceType{core.Rook, core.Knight, core.Bishop, core.Queen, core.King, core.Bishop, core.Knight, core.Rook}
	for file, pt := range rank0 {
		sq := core.NewSquare(0, file)
		got := pos.Board.Check(sq)
		expect := core.NewPiece(pt, core.White)
		if got != expect {
			t.Errorf("at %s: got %s, expected %s", sq.String(), got.String(), expect.String())
		}
	}

	// rank 1: white pawns
	for file := 0; file < 8; file++ {
		sq := core.NewSquare(1, file)
		got := pos.Board.Check(sq)
		expect := core.NewPiece(core.Pawn, core.White)
		if got != expect {
			t.Errorf("at %s: got %s, expected %s", sq.String(), got.String(), expect.String())
		}
	}

	// ranks 2-5: empty
	for rank := 2; rank <= 5; rank++ {
		for file := 0; file < 8; file++ {
			sq := core.NewSquare(rank, file)
			got := pos.Board.Check(sq)
			if got != core.None {
				t.Errorf("at %s: got %s, expected empty", sq.String(), got.String())
			}
		}
	}

	// rank 6: black pawns
	for file := 0; file < 8; file++ {
		sq := core.NewSquare(6, file)
		got := pos.Board.Check(sq)
		expect := core.NewPiece(core.Pawn, core.Black)
		if got != expect {
			t.Errorf("at %s: got %s, expected %s", sq.String(), got.String(), expect.String())
		}
	}

	// rank 7: rnbqkbnr
	for file, pt := range rank0 {
		sq := core.NewSquare(7, file)
		got := pos.Board.Check(sq)
		expect := core.NewPiece(pt, core.Black)
		if got != expect {
			t.Errorf("at %s: got %s, expected %s", sq.String(), got.String(), expect.String())
		}
	}
}

func TestCastling(t *testing.T) {
	tests := []struct {
		name     string
		fen      string
		castling position.CastlingRights
	}{
		{"starting", starting, position.AllCastling},
		{"italian", italian, position.AllCastling},
		{"ep", ep, position.AllCastling},
		{"bcastle", bcastle, position.BlackKingside | position.BlackQueenside},
		{"wcastle", wcastle, position.WhiteKingside | position.WhiteQueenside},
		{"nocastle", nocastle, position.NoCastling},
	}

	for _, tt := range tests {
		pos, err := Parse(tt.fen)
		if err != nil {
			t.Fatalf("%s: %v", tt.name, err)
		}
		if pos.Castling != tt.castling {
			t.Errorf("%s: got %s, expected %s", tt.name, pos.Castling.String(), tt.castling.String())
		}
	}
}

func TestActiveColor(t *testing.T) {
	tests := []struct {
		name  string
		fen   string
		color core.Color
	}{
		{"starting", starting, core.White},
		{"italian", italian, core.Black},
		{"ep", ep, core.White},
	}

	for _, tt := range tests {
		pos, err := Parse(tt.fen)
		if err != nil {
			t.Fatalf("%s: %v", tt.name, err)
		}
		if pos.ActiveColor != tt.color {
			t.Errorf("%s: got %d, expected %d", tt.name, pos.ActiveColor, tt.color)
		}
	}
}

func TestMoveClocks(t *testing.T) {
	tests := []struct {
		name      string
		fen       string
		halfmoves int
		fullmoves int
	}{
		{"starting", starting, 0, 1},
		{"italian", italian, 0, 4},
		{"ep", ep, 0, 3},
		{"bcastle", bcastle, 1, 5},
		{"wcastle", wcastle, 1, 5},
		{"nocastle", nocastle, 2, 5},
	}

	for _, tt := range tests {
		pos, err := Parse(tt.fen)
		if err != nil {
			t.Fatalf("%s: %v", tt.name, err)
		}
		if pos.Halfmoves != tt.halfmoves {
			t.Errorf("%s halfmoves: got %d, expected %d", tt.name, pos.Halfmoves, tt.halfmoves)
		}
		if pos.Fullmoves != tt.fullmoves {
			t.Errorf("%s fullmoves: got %d, expected %d", tt.name, pos.Fullmoves, tt.fullmoves)
		}
	}
}

func TestEnPassant(t *testing.T) {
	pos, err := Parse(ep)
	if err != nil {
		t.Fatal(err)
	}

	expect := core.NewSquare(5, 3) // d6
	if pos.EnPassant != expect {
		t.Errorf("en passant: got %s, expected %s", pos.EnPassant.String(), expect.String())
	}
}
