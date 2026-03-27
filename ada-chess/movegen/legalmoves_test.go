package movegen_test

import (
	"testing"

	"github.com/WilliamDann/AdaEngine/ada-chess/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/fen"
	"github.com/WilliamDann/AdaEngine/ada-chess/movegen"
	"github.com/WilliamDann/AdaEngine/ada-chess/position"
)

func TestLegalMoves_StartingPosition(t *testing.T) {
	pos, _ := fen.Parse("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	ml := movegen.LegalMoves(pos)
	// 16 pawn moves + 4 knight moves = 20
	if ml.Count() != 20 {
		t.Errorf("starting position: got %d moves, want 20", ml.Count())
		printMoves(t, ml)
	}
}

func TestLegalMoves_StartingPositionBlack(t *testing.T) {
	pos, _ := fen.Parse("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1")
	ml := movegen.LegalMoves(pos)
	if ml.Count() != 20 {
		t.Errorf("black starting response: got %d moves, want 20", ml.Count())
		printMoves(t, ml)
	}
}

func TestLegalMoves_KingInCheck_MustResolve(t *testing.T) {
	// White king e1, black rook e8, open file. Only legal moves resolve the check.
	pos, _ := fen.Parse("4r3/8/8/8/8/8/8/4K2k w - - 0 1")
	ml := movegen.LegalMoves(pos)
	// King can move to d1, d2, f2, f1 (not e2 — still attacked by rook)
	for i := 0; i < ml.Count(); i++ {
		m := ml.Get(i)
		if m.To() == sq2(1, 4) { // e2
			t.Error("e2 should not be legal (rook still attacks it)")
		}
	}
	if ml.Count() == 0 {
		t.Error("should have legal moves to escape check")
	}
}

func TestLegalMoves_DoubleCheck_OnlyKingMoves(t *testing.T) {
	// Double check: king must move, no blocking
	pos, _ := fen.Parse("4r3/8/8/8/1b6/8/8/4K2k w - - 0 1")
	// b4 checks e1 diag? b4=(3,1), e1=(0,4), diff=(3,3) — yes.
	// e8 rook checks e1 on file — yes.
	// Two checkers! Only king moves.
	ml := movegen.LegalMoves(pos)
	for i := 0; i < ml.Count(); i++ {
		m := ml.Get(i)
		if m.From() != sq2(0, 4) { // all moves must be from e1 (king)
			t.Errorf("double check: non-king move found: %v", m)
		}
	}
}

func TestLegalMoves_PinnedPieceCantMoveFreely(t *testing.T) {
	// Knight on e4 pinned by rook on e8 to king on e1.
	// Knight has no legal moves (no knight move stays on the e-file).
	pos, _ := fen.Parse("4r3/8/8/8/4N3/8/8/4K2k w - - 0 1")
	ml := movegen.LegalMoves(pos)
	for i := 0; i < ml.Count(); i++ {
		m := ml.Get(i)
		if m.From() == sq2(3, 4) { // e4 = the knight
			t.Errorf("pinned knight should have no moves, but found: %v", m)
		}
	}
}

func TestLegalMoves_PinnedRookCanMoveAlongPin(t *testing.T) {
	// White king e1, white rook e4, black rook e8.
	// Rook is pinned but can move along the e-file.
	pos, _ := fen.Parse("4r3/8/8/8/4R3/8/8/4K2k w - - 0 1")
	ml := movegen.LegalMoves(pos)
	foundRookMove := false
	for i := 0; i < ml.Count(); i++ {
		m := ml.Get(i)
		if m.From() == sq2(3, 4) { // e4 rook
			// Must stay on e-file
			if m.To().File() != 4 {
				t.Errorf("pinned rook moved off pin ray: %v", m)
			}
			foundRookMove = true
		}
	}
	if !foundRookMove {
		t.Error("pinned rook should be able to move along pin ray")
	}
}

func TestLegalMoves_Castling(t *testing.T) {
	// Both sides can castle
	pos, _ := fen.Parse("r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq - 0 1")
	ml := movegen.LegalMoves(pos)
	foundKingside := false
	foundQueenside := false
	for i := 0; i < ml.Count(); i++ {
		m := ml.Get(i)
		if m.MoveType() == core.MoveCastling {
			if m.To() == sq2(0, 6) {
				foundKingside = true
			}
			if m.To() == sq2(0, 2) {
				foundQueenside = true
			}
		}
	}
	if !foundKingside {
		t.Error("white kingside castling should be legal")
	}
	if !foundQueenside {
		t.Error("white queenside castling should be legal")
	}
}

func TestLegalMoves_CastlingBlockedByAttack(t *testing.T) {
	// Black rook on f8 attacks f1 (f-pawn removed), preventing kingside castling
	pos, _ := fen.Parse("5r2/8/8/8/8/8/PPPPP1PP/R3K2R w KQ - 0 1")
	ml := movegen.LegalMoves(pos)
	for i := 0; i < ml.Count(); i++ {
		m := ml.Get(i)
		if m.MoveType() == core.MoveCastling && m.To() == sq2(0, 6) {
			t.Error("kingside castling should be blocked (f1 attacked)")
		}
	}
}

func TestLegalMoves_NoCastlingInCheck(t *testing.T) {
	// King in check from black rook on e8 (e-pawn removed)
	pos, _ := fen.Parse("4r3/8/8/8/8/8/PPPP1PPP/R3K2R w KQ - 0 1")
	ml := movegen.LegalMoves(pos)
	for i := 0; i < ml.Count(); i++ {
		m := ml.Get(i)
		if m.MoveType() == core.MoveCastling {
			t.Error("cannot castle while in check")
		}
	}
}

func TestLegalMoves_Promotion(t *testing.T) {
	// White pawn on e7, can push to e8 or capture d8/f8
	pos, _ := fen.Parse("3r1r2/4P3/8/8/8/8/8/4K2k w - - 0 1")
	ml := movegen.LegalMoves(pos)
	promoCount := 0
	for i := 0; i < ml.Count(); i++ {
		m := ml.Get(i)
		if m.MoveType() == core.MovePromotion {
			promoCount++
		}
	}
	// e8 push (4 promos) + d8 capture (4 promos) + f8 capture (4 promos) = 12
	if promoCount != 12 {
		t.Errorf("expected 12 promotion moves, got %d", promoCount)
		printMoves(t, ml)
	}
}

func TestLegalMoves_EnPassant(t *testing.T) {
	// White pawn e5, black pawn just double-pushed to d5, ep square d6
	pos, _ := fen.Parse("4k3/8/8/3pP3/8/8/8/4K3 w - d6 0 1")
	ml := movegen.LegalMoves(pos)
	foundEP := false
	for i := 0; i < ml.Count(); i++ {
		m := ml.Get(i)
		if m.MoveType() == core.MoveEnPassant {
			foundEP = true
			if m.From() != sq2(4, 4) || m.To() != sq2(5, 3) {
				t.Errorf("unexpected ep move: %v", m)
			}
		}
	}
	if !foundEP {
		t.Error("en passant should be legal")
	}
}

func TestLegalMoves_EnPassantHorizontalDiscovery(t *testing.T) {
	// King on a5, enemy rook on h5, pawns on e5 (white) and d5 (black, just pushed).
	// En passant removes both pawns from rank 5, exposing king to rook.
	pos, _ := fen.Parse("4k3/8/8/K2pP2r/8/8/8/8 w - d6 0 1")
	ml := movegen.LegalMoves(pos)
	for i := 0; i < ml.Count(); i++ {
		m := ml.Get(i)
		if m.MoveType() == core.MoveEnPassant {
			t.Error("en passant should be illegal (horizontal discovery)")
		}
	}
}

func TestLegalMoves_Checkmate(t *testing.T) {
	// Back rank mate: king h1, black rook on a1, pawns blocking escape
	pos, _ := fen.Parse("4k3/8/8/8/8/8/5PPP/r5K1 w - - 0 1")
	ml := movegen.LegalMoves(pos)
	if ml.Count() != 0 {
		t.Errorf("checkmate: expected 0 legal moves, got %d", ml.Count())
		printMoves(t, ml)
	}
}

func TestLegalMoves_Stalemate(t *testing.T) {
	// King a1, black queen c2 — king has no legal moves, not in check
	pos, _ := fen.Parse("4k3/8/8/8/8/8/2q5/K7 w - - 0 1")
	ml := movegen.LegalMoves(pos)
	if ml.Count() != 0 {
		t.Errorf("stalemate: expected 0 legal moves, got %d", ml.Count())
		printMoves(t, ml)
	}
}

// Perft: count leaf nodes at a given depth. The gold standard for move gen correctness.
func perft(pos *position.Position, depth int) uint64 {
	if depth == 0 {
		return 1
	}
	ml := movegen.LegalMoves(pos)
	if depth == 1 {
		return uint64(ml.Count())
	}
	var nodes uint64
	for i := 0; i < ml.Count(); i++ {
		child := makeMove(pos, ml.Get(i))
		nodes += perft(child, depth-1)
	}
	return nodes
}

// makeMove creates a new position after applying a move.
// Minimal implementation for perft — copies the board and applies the move.
func makeMove(pos *position.Position, m core.Move) *position.Position {
	newPos := &position.Position{
		Board:       copyBoard(pos.Board),
		ActiveColor: pos.ActiveColor.Flip(),
		Castling:    pos.Castling,
		EnPassant:   core.InvalidSquare,
		Halfmoves:   pos.Halfmoves + 1,
		Fullmoves:   pos.Fullmoves,
	}
	if pos.ActiveColor == core.Black {
		newPos.Fullmoves++
	}

	from := m.From()
	to := m.To()
	piece := pos.Board.Check(from)

	// Move the piece
	newPos.Board.Clear(from)
	newPos.Board.Clear(to)

	switch m.MoveType() {
	case core.MovePromotion:
		newPos.Board.Set(to, core.NewPiece(m.PromoPiece(), pos.ActiveColor))
	case core.MoveEnPassant:
		// Remove captured pawn
		dir := 8
		if pos.ActiveColor == core.Black {
			dir = -8
		}
		capturedSq := core.Square(int(to) - dir)
		newPos.Board.Clear(capturedSq)
		newPos.Board.Set(to, piece)
	case core.MoveCastling:
		newPos.Board.Set(to, piece)
		// Move the rook
		switch to {
		case core.Square(6): // g1
			rook := newPos.Board.Check(core.Square(7))
			newPos.Board.Clear(core.Square(7))
			newPos.Board.Set(core.Square(5), rook)
		case core.Square(2): // c1
			rook := newPos.Board.Check(core.Square(0))
			newPos.Board.Clear(core.Square(0))
			newPos.Board.Set(core.Square(3), rook)
		case core.Square(62): // g8
			rook := newPos.Board.Check(core.Square(63))
			newPos.Board.Clear(core.Square(63))
			newPos.Board.Set(core.Square(61), rook)
		case core.Square(58): // c8
			rook := newPos.Board.Check(core.Square(56))
			newPos.Board.Clear(core.Square(56))
			newPos.Board.Set(core.Square(59), rook)
		}
	default:
		newPos.Board.Set(to, piece)
	}

	// Update castling rights
	// King moves
	if piece.Type() == core.King {
		if pos.ActiveColor == core.White {
			newPos.Castling &^= position.WhiteKingside | position.WhiteQueenside
		} else {
			newPos.Castling &^= position.BlackKingside | position.BlackQueenside
		}
	}
	// Rook moves or captures
	switch from {
	case core.Square(0):
		newPos.Castling &^= position.WhiteQueenside
	case core.Square(7):
		newPos.Castling &^= position.WhiteKingside
	case core.Square(56):
		newPos.Castling &^= position.BlackQueenside
	case core.Square(63):
		newPos.Castling &^= position.BlackKingside
	}
	switch to {
	case core.Square(0):
		newPos.Castling &^= position.WhiteQueenside
	case core.Square(7):
		newPos.Castling &^= position.WhiteKingside
	case core.Square(56):
		newPos.Castling &^= position.BlackQueenside
	case core.Square(63):
		newPos.Castling &^= position.BlackKingside
	}

	// En passant: set if double pawn push
	if piece.Type() == core.Pawn {
		diff := int(to) - int(from)
		if diff == 16 || diff == -16 {
			newPos.EnPassant = core.Square((int(from) + int(to)) / 2)
		}
		newPos.Halfmoves = 0
	}
	// Reset halfmove on captures
	if pos.Board.Check(to) != core.None {
		newPos.Halfmoves = 0
	}

	return newPos
}

func copyBoard(b *core.Chessboard) *core.Chessboard {
	nb := core.NewChessboard()
	for sq := core.Square(0); sq < 64; sq++ {
		piece := b.Check(sq)
		if piece != core.None {
			nb.Set(sq, piece)
		}
	}
	return nb
}

func TestPerft_StartingPosition(t *testing.T) {
	pos, _ := fen.Parse("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	tests := []struct {
		depth int
		want  uint64
	}{
		{1, 20},
		{2, 400},
		{3, 8902},
		{4, 197281},
	}
	for _, tc := range tests {
		got := perft(pos, tc.depth)
		if got != tc.want {
			t.Errorf("perft(%d) = %d, want %d", tc.depth, got, tc.want)
		}
	}
}

func printMoves(t *testing.T, ml core.MoveList) {
	t.Helper()
	for i := 0; i < ml.Count(); i++ {
		t.Logf("  %s", ml.Get(i))
	}
}
