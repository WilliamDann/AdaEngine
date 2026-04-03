package movegen_test

import (
	"testing"

	"github.com/WilliamDann/AdaEngine/ada-chess/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/fen"
	"github.com/WilliamDann/AdaEngine/ada-chess/movegen"
)

// isCapture checks if a move is a capture or promotion in the given position.
func isCapture(pos *core.Chessboard, m core.Move, enemy core.Color) bool {
	if m.MoveType() == core.MovePromotion {
		return true
	}
	if m.MoveType() == core.MoveEnPassant {
		return true
	}
	to := m.To()
	piece := pos.Check(to)
	return piece != core.None && piece.Color() == enemy
}

func TestLegalCaptures_StartingPosition(t *testing.T) {
	pos, _ := fen.Parse("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	ml := movegen.LegalCaptures(pos)
	if ml.Count() != 0 {
		t.Errorf("starting position: got %d captures, want 0", ml.Count())
	}
}

func TestLegalCaptures_SubsetOfLegalMoves(t *testing.T) {
	positions := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"r1bqkbnr/pppppppp/2n5/8/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 2 2",
		"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
		"rnbqkb1r/pp1p1pPp/8/2p1pP2/1P1P4/3P3P/P1P1P3/RNBQKBNR w KQkq e6 0 1",
		"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
	}

	for _, f := range positions {
		pos, err := fen.Parse(f)
		if err != nil {
			t.Fatalf("bad FEN %q: %v", f, err)
		}

		captures := movegen.LegalCaptures(pos)
		allMoves := movegen.LegalMoves(pos)

		// Build set of all legal moves
		legalSet := make(map[core.Move]bool)
		for i := 0; i < allMoves.Count(); i++ {
			legalSet[allMoves.Get(i)] = true
		}

		// Every capture must be in legal moves
		for i := 0; i < captures.Count(); i++ {
			m := captures.Get(i)
			if !legalSet[m] {
				t.Errorf("FEN %q: capture %s not in legal moves", f, m)
			}
		}
	}
}

func TestLegalCaptures_OnlyTacticalMoves(t *testing.T) {
	positions := []string{
		"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
		"rnbqkb1r/pp1p1pPp/8/2p1pP2/1P1P4/3P3P/P1P1P3/RNBQKBNR w KQkq e6 0 1",
		"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
	}

	for _, f := range positions {
		pos, err := fen.Parse(f)
		if err != nil {
			t.Fatalf("bad FEN %q: %v", f, err)
		}

		enemy := pos.ActiveColor.Flip()
		captures := movegen.LegalCaptures(pos)

		for i := 0; i < captures.Count(); i++ {
			m := captures.Get(i)
			if !isCapture(pos.Board, m, enemy) {
				t.Errorf("FEN %q: move %s is not a capture or promotion", f, m)
			}
		}
	}
}

func TestLegalCaptures_AllCapturesFound(t *testing.T) {
	positions := []string{
		"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
		"rnbqkb1r/pp1p1pPp/8/2p1pP2/1P1P4/3P3P/P1P1P3/RNBQKBNR w KQkq e6 0 1",
		"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
		"r1bqkbnr/pppppppp/2n5/8/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 2 2",
	}

	for _, f := range positions {
		pos, err := fen.Parse(f)
		if err != nil {
			t.Fatalf("bad FEN %q: %v", f, err)
		}

		enemy := pos.ActiveColor.Flip()
		captures := movegen.LegalCaptures(pos)
		allMoves := movegen.LegalMoves(pos)

		// Count tactical moves in full legal move list
		captureSet := make(map[core.Move]bool)
		for i := 0; i < captures.Count(); i++ {
			captureSet[captures.Get(i)] = true
		}

		for i := 0; i < allMoves.Count(); i++ {
			m := allMoves.Get(i)
			if isCapture(pos.Board, m, enemy) && !captureSet[m] {
				t.Errorf("FEN %q: capture %s missing from LegalCaptures", f, m)
			}
		}
	}
}

func TestLegalCaptures_EnPassant(t *testing.T) {
	// White pawn on e5, black just played d7-d5
	pos, _ := fen.Parse("rnbqkbnr/ppp1pppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3")
	captures := movegen.LegalCaptures(pos)

	found := false
	for i := 0; i < captures.Count(); i++ {
		if captures.Get(i).MoveType() == core.MoveEnPassant {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected en passant in captures")
	}
}

func TestLegalCaptures_InCheck(t *testing.T) {
	// White king in check, only legal moves are to block or capture the checker
	pos, _ := fen.Parse("rnb1kbnr/pppp1ppp/8/4p3/6Pq/5P2/PPPPP2P/RNBQKBNR w KQkq - 1 3")
	captures := movegen.LegalCaptures(pos)
	allMoves := movegen.LegalMoves(pos)

	// All captures should be legal
	legalSet := make(map[core.Move]bool)
	for i := 0; i < allMoves.Count(); i++ {
		legalSet[allMoves.Get(i)] = true
	}
	for i := 0; i < captures.Count(); i++ {
		if !legalSet[captures.Get(i)] {
			t.Errorf("capture %s is not legal while in check", captures.Get(i))
		}
	}
}
