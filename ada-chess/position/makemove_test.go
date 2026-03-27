package position

import (
	"testing"

	"github.com/WilliamDann/AdaEngine/ada-chess/core"
)

// helper to build a position with pieces on specific squares
func setupPosition(pieces map[core.Square]core.Piece, active core.Color, castling CastlingRights, ep core.Square) *Position {
	pos := &Position{
		Board:       core.NewChessboard(),
		ActiveColor: active,
		Castling:    castling,
		EnPassant:   ep,
		Halfmoves:   0,
		Fullmoves:   1,
	}
	for sq, p := range pieces {
		pos.Board.Set(sq, p)
	}
	return pos
}

func TestNormalMove(t *testing.T) {
	e2 := core.NewSquare(1, 4)
	e4 := core.NewSquare(3, 4)
	wPawn := core.NewPiece(core.Pawn, core.White)

	pos := setupPosition(map[core.Square]core.Piece{e2: wPawn}, core.White, NoCastling, core.InvalidSquare)
	m := core.NewMove(e2, e4)
	next := MakeMove(pos, m)

	if next.Board.Check(e2) != core.None {
		t.Error("source square should be empty after move")
	}
	if next.Board.Check(e4) != wPawn {
		t.Errorf("destination: got %s, want %s", next.Board.Check(e4).String(), wPawn.String())
	}
}

func TestOriginalPositionUnchanged(t *testing.T) {
	e2 := core.NewSquare(1, 4)
	e4 := core.NewSquare(3, 4)
	wPawn := core.NewPiece(core.Pawn, core.White)

	pos := setupPosition(map[core.Square]core.Piece{e2: wPawn}, core.White, NoCastling, core.InvalidSquare)
	m := core.NewMove(e2, e4)
	MakeMove(pos, m)

	if pos.Board.Check(e2) != wPawn {
		t.Error("original position should not be modified")
	}
	if pos.ActiveColor != core.White {
		t.Error("original active color should not change")
	}
}

func TestActiveColorFlips(t *testing.T) {
	e2 := core.NewSquare(1, 4)
	e3 := core.NewSquare(2, 4)
	wPawn := core.NewPiece(core.Pawn, core.White)

	pos := setupPosition(map[core.Square]core.Piece{e2: wPawn}, core.White, NoCastling, core.InvalidSquare)
	next := MakeMove(pos, core.NewMove(e2, e3))
	if next.ActiveColor != core.Black {
		t.Error("active color should flip to Black after White's move")
	}

	d7 := core.NewSquare(6, 3)
	d6 := core.NewSquare(5, 3)
	bPawn := core.NewPiece(core.Pawn, core.Black)
	pos2 := setupPosition(map[core.Square]core.Piece{d7: bPawn}, core.Black, NoCastling, core.InvalidSquare)
	next2 := MakeMove(pos2, core.NewMove(d7, d6))
	if next2.ActiveColor != core.White {
		t.Error("active color should flip to White after Black's move")
	}
}

func TestFullmoveIncrementsAfterBlack(t *testing.T) {
	e2 := core.NewSquare(1, 4)
	e3 := core.NewSquare(2, 4)
	wPawn := core.NewPiece(core.Pawn, core.White)

	pos := setupPosition(map[core.Square]core.Piece{e2: wPawn}, core.White, NoCastling, core.InvalidSquare)
	pos.Fullmoves = 5
	next := MakeMove(pos, core.NewMove(e2, e3))
	if next.Fullmoves != 5 {
		t.Errorf("fullmoves should not increment after White's move: got %d, want 5", next.Fullmoves)
	}

	d7 := core.NewSquare(6, 3)
	d6 := core.NewSquare(5, 3)
	bPawn := core.NewPiece(core.Pawn, core.Black)
	pos2 := setupPosition(map[core.Square]core.Piece{d7: bPawn}, core.Black, NoCastling, core.InvalidSquare)
	pos2.Fullmoves = 5
	next2 := MakeMove(pos2, core.NewMove(d7, d6))
	if next2.Fullmoves != 6 {
		t.Errorf("fullmoves should increment after Black's move: got %d, want 6", next2.Fullmoves)
	}
}

func TestHalfmoveResetOnPawnMove(t *testing.T) {
	e2 := core.NewSquare(1, 4)
	e3 := core.NewSquare(2, 4)
	wPawn := core.NewPiece(core.Pawn, core.White)

	pos := setupPosition(map[core.Square]core.Piece{e2: wPawn}, core.White, NoCastling, core.InvalidSquare)
	pos.Halfmoves = 10
	next := MakeMove(pos, core.NewMove(e2, e3))
	if next.Halfmoves != 0 {
		t.Errorf("halfmoves should reset on pawn move: got %d", next.Halfmoves)
	}
}

func TestHalfmoveResetOnCapture(t *testing.T) {
	d4 := core.NewSquare(3, 3)
	e5 := core.NewSquare(4, 4)
	wKnight := core.NewPiece(core.Knight, core.White)
	bPawn := core.NewPiece(core.Pawn, core.Black)

	pos := setupPosition(map[core.Square]core.Piece{d4: wKnight, e5: bPawn}, core.White, NoCastling, core.InvalidSquare)
	pos.Halfmoves = 7
	next := MakeMove(pos, core.NewMove(d4, e5))
	if next.Halfmoves != 0 {
		t.Errorf("halfmoves should reset on capture: got %d", next.Halfmoves)
	}
}

func TestHalfmoveIncrementsOnQuietMove(t *testing.T) {
	b1 := core.NewSquare(0, 1)
	c3 := core.NewSquare(2, 2)
	wKnight := core.NewPiece(core.Knight, core.White)

	pos := setupPosition(map[core.Square]core.Piece{b1: wKnight}, core.White, NoCastling, core.InvalidSquare)
	pos.Halfmoves = 3
	next := MakeMove(pos, core.NewMove(b1, c3))
	if next.Halfmoves != 4 {
		t.Errorf("halfmoves should increment on quiet non-pawn move: got %d, want 4", next.Halfmoves)
	}
}

func TestDoublePawnPushSetsEnPassant(t *testing.T) {
	tests := []struct {
		name   string
		from   core.Square
		to     core.Square
		epSq   core.Square
		color  core.Color
	}{
		{"white e2-e4", core.NewSquare(1, 4), core.NewSquare(3, 4), core.NewSquare(2, 4), core.White},
		{"white d2-d4", core.NewSquare(1, 3), core.NewSquare(3, 3), core.NewSquare(2, 3), core.White},
		{"black e7-e5", core.NewSquare(6, 4), core.NewSquare(4, 4), core.NewSquare(5, 4), core.Black},
		{"black d7-d5", core.NewSquare(6, 3), core.NewSquare(4, 3), core.NewSquare(5, 3), core.Black},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pawn := core.NewPiece(core.Pawn, tt.color)
			pos := setupPosition(map[core.Square]core.Piece{tt.from: pawn}, tt.color, NoCastling, core.InvalidSquare)
			next := MakeMove(pos, core.NewMove(tt.from, tt.to))
			if next.EnPassant != tt.epSq {
				t.Errorf("en passant: got %s, want %s", next.EnPassant.String(), tt.epSq.String())
			}
		})
	}
}

func TestSinglePawnPushNoEnPassant(t *testing.T) {
	e2 := core.NewSquare(1, 4)
	e3 := core.NewSquare(2, 4)
	wPawn := core.NewPiece(core.Pawn, core.White)

	pos := setupPosition(map[core.Square]core.Piece{e2: wPawn}, core.White, NoCastling, core.InvalidSquare)
	next := MakeMove(pos, core.NewMove(e2, e3))
	if next.EnPassant != core.InvalidSquare {
		t.Errorf("single pawn push should not set en passant: got %s", next.EnPassant.String())
	}
}

func TestEnPassantClearedAfterMove(t *testing.T) {
	b1 := core.NewSquare(0, 1)
	c3 := core.NewSquare(2, 2)
	wKnight := core.NewPiece(core.Knight, core.White)

	pos := setupPosition(map[core.Square]core.Piece{b1: wKnight}, core.White, NoCastling, core.NewSquare(5, 3))
	next := MakeMove(pos, core.NewMove(b1, c3))
	if next.EnPassant != core.InvalidSquare {
		t.Errorf("en passant should be cleared after non-double-push move: got %s", next.EnPassant.String())
	}
}

func TestPromotion(t *testing.T) {
	e7 := core.NewSquare(6, 4)
	e8 := core.NewSquare(7, 4)
	wPawn := core.NewPiece(core.Pawn, core.White)

	promos := []struct {
		piece core.PieceType
		name  string
	}{
		{core.Queen, "queen"},
		{core.Rook, "rook"},
		{core.Bishop, "bishop"},
		{core.Knight, "knight"},
	}

	for _, tt := range promos {
		t.Run(tt.name, func(t *testing.T) {
			pos := setupPosition(map[core.Square]core.Piece{e7: wPawn}, core.White, NoCastling, core.InvalidSquare)
			m := core.NewPromotion(e7, e8, tt.piece)
			next := MakeMove(pos, m)

			if next.Board.Check(e7) != core.None {
				t.Error("source square should be empty after promotion")
			}
			expect := core.NewPiece(tt.piece, core.White)
			got := next.Board.Check(e8)
			if got != expect {
				t.Errorf("promoted piece: got %s, want %s", got.String(), expect.String())
			}
		})
	}
}

func TestBlackPromotion(t *testing.T) {
	d2 := core.NewSquare(1, 3)
	d1 := core.NewSquare(0, 3)
	bPawn := core.NewPiece(core.Pawn, core.Black)

	pos := setupPosition(map[core.Square]core.Piece{d2: bPawn}, core.Black, NoCastling, core.InvalidSquare)
	m := core.NewPromotion(d2, d1, core.Queen)
	next := MakeMove(pos, m)

	expect := core.NewPiece(core.Queen, core.Black)
	got := next.Board.Check(d1)
	if got != expect {
		t.Errorf("black promotion: got %s, want %s", got.String(), expect.String())
	}
}

func TestEnPassantCapture(t *testing.T) {
	// White pawn on e5 captures en passant on d6, removing black pawn on d5
	e5 := core.NewSquare(4, 4)
	d6 := core.NewSquare(5, 3)
	d5 := core.NewSquare(4, 3)
	wPawn := core.NewPiece(core.Pawn, core.White)
	bPawn := core.NewPiece(core.Pawn, core.Black)

	pos := setupPosition(map[core.Square]core.Piece{e5: wPawn, d5: bPawn}, core.White, NoCastling, d6)
	m := core.NewEnPassant(e5, d6)
	next := MakeMove(pos, m)

	if next.Board.Check(e5) != core.None {
		t.Error("source square should be empty")
	}
	if next.Board.Check(d6) != wPawn {
		t.Error("pawn should be on en passant target square")
	}
	if next.Board.Check(d5) != core.None {
		t.Error("captured pawn should be removed")
	}
	if next.Halfmoves != 0 {
		t.Error("halfmoves should reset on en passant capture")
	}
}

func TestBlackEnPassantCapture(t *testing.T) {
	// Black pawn on d4 captures en passant on e3, removing white pawn on e4
	d4 := core.NewSquare(3, 3)
	e3 := core.NewSquare(2, 4)
	e4 := core.NewSquare(3, 4)
	bPawn := core.NewPiece(core.Pawn, core.Black)
	wPawn := core.NewPiece(core.Pawn, core.White)

	pos := setupPosition(map[core.Square]core.Piece{d4: bPawn, e4: wPawn}, core.Black, NoCastling, e3)
	m := core.NewEnPassant(d4, e3)
	next := MakeMove(pos, m)

	if next.Board.Check(d4) != core.None {
		t.Error("source square should be empty")
	}
	if next.Board.Check(e3) != bPawn {
		t.Error("pawn should be on en passant target square")
	}
	if next.Board.Check(e4) != core.None {
		t.Error("captured pawn should be removed")
	}
}

func TestWhiteKingsideCastling(t *testing.T) {
	e1 := core.NewSquare(0, 4)
	g1 := core.NewSquare(0, 6)
	f1 := core.NewSquare(0, 5)
	h1 := core.NewSquare(0, 7)
	wKing := core.NewPiece(core.King, core.White)
	wRook := core.NewPiece(core.Rook, core.White)

	pos := setupPosition(map[core.Square]core.Piece{e1: wKing, h1: wRook}, core.White, AllCastling, core.InvalidSquare)
	m := core.NewCastling(e1, g1)
	next := MakeMove(pos, m)

	if next.Board.Check(e1) != core.None {
		t.Error("king source should be empty")
	}
	if next.Board.Check(h1) != core.None {
		t.Error("rook source should be empty")
	}
	if next.Board.Check(g1) != wKing {
		t.Error("king should be on g1")
	}
	if next.Board.Check(f1) != wRook {
		t.Error("rook should be on f1")
	}
}

func TestWhiteQueensideCastling(t *testing.T) {
	e1 := core.NewSquare(0, 4)
	c1 := core.NewSquare(0, 2)
	d1 := core.NewSquare(0, 3)
	a1 := core.NewSquare(0, 0)
	wKing := core.NewPiece(core.King, core.White)
	wRook := core.NewPiece(core.Rook, core.White)

	pos := setupPosition(map[core.Square]core.Piece{e1: wKing, a1: wRook}, core.White, AllCastling, core.InvalidSquare)
	m := core.NewCastling(e1, c1)
	next := MakeMove(pos, m)

	if next.Board.Check(e1) != core.None {
		t.Error("king source should be empty")
	}
	if next.Board.Check(a1) != core.None {
		t.Error("rook source should be empty")
	}
	if next.Board.Check(c1) != wKing {
		t.Error("king should be on c1")
	}
	if next.Board.Check(d1) != wRook {
		t.Error("rook should be on d1")
	}
}

func TestBlackKingsideCastling(t *testing.T) {
	e8 := core.NewSquare(7, 4)
	g8 := core.NewSquare(7, 6)
	f8 := core.NewSquare(7, 5)
	h8 := core.NewSquare(7, 7)
	bKing := core.NewPiece(core.King, core.Black)
	bRook := core.NewPiece(core.Rook, core.Black)

	pos := setupPosition(map[core.Square]core.Piece{e8: bKing, h8: bRook}, core.Black, AllCastling, core.InvalidSquare)
	m := core.NewCastling(e8, g8)
	next := MakeMove(pos, m)

	if next.Board.Check(g8) != bKing {
		t.Error("king should be on g8")
	}
	if next.Board.Check(f8) != bRook {
		t.Error("rook should be on f8")
	}
	if next.Board.Check(e8) != core.None {
		t.Error("king source should be empty")
	}
	if next.Board.Check(h8) != core.None {
		t.Error("rook source should be empty")
	}
}

func TestBlackQueensideCastling(t *testing.T) {
	e8 := core.NewSquare(7, 4)
	c8 := core.NewSquare(7, 2)
	d8 := core.NewSquare(7, 3)
	a8 := core.NewSquare(7, 0)
	bKing := core.NewPiece(core.King, core.Black)
	bRook := core.NewPiece(core.Rook, core.Black)

	pos := setupPosition(map[core.Square]core.Piece{e8: bKing, a8: bRook}, core.Black, AllCastling, core.InvalidSquare)
	m := core.NewCastling(e8, c8)
	next := MakeMove(pos, m)

	if next.Board.Check(c8) != bKing {
		t.Error("king should be on c8")
	}
	if next.Board.Check(d8) != bRook {
		t.Error("rook should be on d8")
	}
}

func TestCastlingRightsRevokedOnKingMove(t *testing.T) {
	e1 := core.NewSquare(0, 4)
	d1 := core.NewSquare(0, 3)
	wKing := core.NewPiece(core.King, core.White)

	pos := setupPosition(map[core.Square]core.Piece{e1: wKing}, core.White, AllCastling, core.InvalidSquare)
	next := MakeMove(pos, core.NewMove(e1, d1))

	if next.Castling&WhiteKingside != 0 || next.Castling&WhiteQueenside != 0 {
		t.Errorf("white castling should be revoked after king move: got %s", next.Castling.String())
	}
	if next.Castling&BlackKingside == 0 || next.Castling&BlackQueenside == 0 {
		t.Error("black castling should be preserved when white king moves")
	}
}

func TestCastlingRightsRevokedOnBlackKingMove(t *testing.T) {
	e8 := core.NewSquare(7, 4)
	d8 := core.NewSquare(7, 3)
	bKing := core.NewPiece(core.King, core.Black)

	pos := setupPosition(map[core.Square]core.Piece{e8: bKing}, core.Black, AllCastling, core.InvalidSquare)
	next := MakeMove(pos, core.NewMove(e8, d8))

	if next.Castling&BlackKingside != 0 || next.Castling&BlackQueenside != 0 {
		t.Errorf("black castling should be revoked after king move: got %s", next.Castling.String())
	}
	if next.Castling&WhiteKingside == 0 || next.Castling&WhiteQueenside == 0 {
		t.Error("white castling should be preserved when black king moves")
	}
}

func TestCastlingRightsRevokedOnRookMove(t *testing.T) {
	tests := []struct {
		name    string
		rookSq  core.Square
		moveTo  core.Square
		color   core.Color
		revoked CastlingRights
		kept    CastlingRights
	}{
		{"white a-rook", core.NewSquare(0, 0), core.NewSquare(1, 0), core.White, WhiteQueenside, WhiteKingside | BlackKingside | BlackQueenside},
		{"white h-rook", core.NewSquare(0, 7), core.NewSquare(1, 7), core.White, WhiteKingside, WhiteQueenside | BlackKingside | BlackQueenside},
		{"black a-rook", core.NewSquare(7, 0), core.NewSquare(6, 0), core.Black, BlackQueenside, WhiteKingside | WhiteQueenside | BlackKingside},
		{"black h-rook", core.NewSquare(7, 7), core.NewSquare(6, 7), core.Black, BlackKingside, WhiteKingside | WhiteQueenside | BlackQueenside},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rook := core.NewPiece(core.Rook, tt.color)
			pos := setupPosition(map[core.Square]core.Piece{tt.rookSq: rook}, tt.color, AllCastling, core.InvalidSquare)
			next := MakeMove(pos, core.NewMove(tt.rookSq, tt.moveTo))

			if next.Castling&tt.revoked != 0 {
				t.Errorf("castling right should be revoked: got %s", next.Castling.String())
			}
			if next.Castling&tt.kept != tt.kept {
				t.Errorf("other castling rights should be kept: got %s", next.Castling.String())
			}
		})
	}
}

func TestCastlingRightsRevokedOnRookCapture(t *testing.T) {
	// White bishop captures black rook on a8 — should revoke black queenside
	a8 := core.NewSquare(7, 0)
	b7 := core.NewSquare(6, 1)
	wBishop := core.NewPiece(core.Bishop, core.White)
	bRook := core.NewPiece(core.Rook, core.Black)

	pos := setupPosition(map[core.Square]core.Piece{b7: wBishop, a8: bRook}, core.White, AllCastling, core.InvalidSquare)
	next := MakeMove(pos, core.NewMove(b7, a8))

	if next.Castling&BlackQueenside != 0 {
		t.Errorf("black queenside castling should be revoked when a8 rook captured: got %s", next.Castling.String())
	}
}

func TestCastlingRightsPreservedOnIrrelevantMove(t *testing.T) {
	b1 := core.NewSquare(0, 1)
	c3 := core.NewSquare(2, 2)
	wKnight := core.NewPiece(core.Knight, core.White)

	pos := setupPosition(map[core.Square]core.Piece{b1: wKnight}, core.White, AllCastling, core.InvalidSquare)
	next := MakeMove(pos, core.NewMove(b1, c3))

	if next.Castling != AllCastling {
		t.Errorf("castling should be preserved on irrelevant move: got %s, want %s", next.Castling.String(), AllCastling.String())
	}
}

func TestCapture(t *testing.T) {
	d4 := core.NewSquare(3, 3)
	e5 := core.NewSquare(4, 4)
	wKnight := core.NewPiece(core.Knight, core.White)
	bPawn := core.NewPiece(core.Pawn, core.Black)

	pos := setupPosition(map[core.Square]core.Piece{d4: wKnight, e5: bPawn}, core.White, NoCastling, core.InvalidSquare)
	next := MakeMove(pos, core.NewMove(d4, e5))

	if next.Board.Check(d4) != core.None {
		t.Error("source should be empty")
	}
	if next.Board.Check(e5) != wKnight {
		t.Error("capturing piece should occupy target square")
	}
}
