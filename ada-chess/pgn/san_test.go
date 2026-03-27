package pgn

import (
	"testing"

	"github.com/WilliamDann/AdaEngine/ada-chess/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/fen"
	"github.com/WilliamDann/AdaEngine/ada-chess/movegen"
	"github.com/WilliamDann/AdaEngine/ada-chess/position"
)

func findMove(pos *position.Position, from, to string) core.Move {
	ml := movegen.LegalMoves(pos)
	for i := 0; i < ml.Count(); i++ {
		m := ml.Get(i)
		if m.From().String()+m.To().String() == from+to {
			return m
		}
	}
	return core.NoMove
}

func findMoveStr(pos *position.Position, uci string) core.Move {
	ml := movegen.LegalMoves(pos)
	for i := 0; i < ml.Count(); i++ {
		m := ml.Get(i)
		if m.String() == uci {
			return m
		}
	}
	return core.NoMove
}

func TestSANPawnMove(t *testing.T) {
	pos, _ := fen.Parse("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	m := findMove(pos, "e2", "e4")
	got := SAN(pos, m)
	if got != "e4" {
		t.Errorf("expected e4, got %s", got)
	}
}

func TestSANKnightMove(t *testing.T) {
	pos, _ := fen.Parse("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	m := findMove(pos, "g1", "f3")
	got := SAN(pos, m)
	if got != "Nf3" {
		t.Errorf("expected Nf3, got %s", got)
	}
}

func TestSANCapture(t *testing.T) {
	pos, _ := fen.Parse("rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 2")
	m := findMove(pos, "e4", "d5")
	got := SAN(pos, m)
	if got != "exd5" {
		t.Errorf("expected exd5, got %s", got)
	}
}

func TestSANCastling(t *testing.T) {
	pos, _ := fen.Parse("r1bqk2r/pppp1ppp/2n2n2/2b1p3/2B1P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4")
	m := findMoveStr(pos, "e1g1")
	got := SAN(pos, m)
	if got != "O-O" {
		t.Errorf("expected O-O, got %s", got)
	}
}

func TestSANPromotion(t *testing.T) {
	pos, _ := fen.Parse("8/P7/8/8/8/8/4k3/2K5 w - - 0 1")
	m := findMoveStr(pos, "a7a8q")
	got := SAN(pos, m)
	if got != "a8=Q" {
		t.Errorf("expected a8=Q, got %s", got)
	}
}

func TestSANDisambiguationFile(t *testing.T) {
	// Two rooks on a1 and h1, king on e2 — both rooks can reach d1
	pos, _ := fen.Parse("6k1/8/8/8/8/8/4K3/R6R w - - 0 1")
	m := findMove(pos, "a1", "d1")
	got := SAN(pos, m)
	if got != "Rad1" {
		t.Errorf("expected Rad1, got %s", got)
	}
}

func TestSANCheck(t *testing.T) {
	pos, _ := fen.Parse("4k3/8/8/8/8/8/8/4K2R w K - 0 1")
	m := findMove(pos, "h1", "h8")
	got := SAN(pos, m)
	if got != "Rh8+" {
		t.Errorf("expected Rh8+, got %s", got)
	}
}

func TestSANCheckmate(t *testing.T) {
	// Scholar's mate final move: Qxf7#
	pos, _ := fen.Parse("r1bqkb1r/pppp1ppp/2n2n2/4p2Q/2B1P3/8/PPPP1PPP/RNB1K1NR w KQkq - 4 4")
	m := findMove(pos, "h5", "f7")
	got := SAN(pos, m)
	if got != "Qxf7#" {
		t.Errorf("expected Qxf7#, got %s", got)
	}
}

func TestGamePGN(t *testing.T) {
	pos, _ := fen.Parse("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	g := NewGame(pos)
	g.White = "Human"
	g.Black = "AdaEngine"
	g.Date = "2026.03.26"

	// 1. e4 e5
	m1 := findMove(pos, "e2", "e4")
	g.AddMove(pos, m1)
	pos = position.MakeMove(pos, m1)

	m2 := findMove(pos, "e7", "e5")
	g.AddMove(pos, m2)
	pos = position.MakeMove(pos, m2)

	// 2. Nf3
	m3 := findMove(pos, "g1", "f3")
	g.AddMove(pos, m3)

	g.Result = "*"
	pgn := g.String()

	if g.MoveCount() != 3 {
		t.Errorf("expected 3 moves, got %d", g.MoveCount())
	}

	expected := `[Event "AdaEngine Game"]
[Site "?"]
[Date "2026.03.26"]
[Round "?"]
[White "Human"]
[Black "AdaEngine"]
[Result "*"]

1. e4 e5 2. Nf3 *
`
	if pgn != expected {
		t.Errorf("PGN mismatch.\nExpected:\n%s\nGot:\n%s", expected, pgn)
	}
}
