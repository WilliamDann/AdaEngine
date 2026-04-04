package search

import (
	"testing"

	"github.com/WilliamDann/AdaEngine/ada-chess/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/fen"
)

func TestSearchStartingPosition(t *testing.T) {
	pos, _ := fen.Parse("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	res := Search(pos, 1, 1, 0)
	if res.Move == core.NoMove {
		t.Fatal("expected a move from the starting position")
	}
	t.Logf("depth 1: move=%s score=%d nodes=%d", res.Move, res.Score, res.Nodes)
}

func TestSearchDepth3(t *testing.T) {
	pos, _ := fen.Parse("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	res := Search(pos, 3, 1, 0)
	if res.Move == core.NoMove {
		t.Fatal("expected a move")
	}
	t.Logf("depth 3: move=%s score=%d nodes=%d", res.Move, res.Score, res.Nodes)
}

func TestSearchMateIn1(t *testing.T) {
	// White to move, Qh7# is mate in 1
	pos, err := fen.Parse("r1bqkb1r/pppp1ppp/2n2n2/4p2Q/2B1P3/8/PPPP1PPP/RNB1K1NR w KQkq - 4 4")
	if err != nil {
		t.Fatal(err)
	}
	res := Search(pos, 1, 1, 0)
	// Qxf7# — the queen on h5 captures f7
	if res.Move.To() != core.NewSquare(6, 5) {
		t.Errorf("expected mate move Qxf7#, got %s", res.Move)
	}
	if res.Score != Mate {
		t.Errorf("expected mate score %d, got %d", Mate, res.Score)
	}
	t.Logf("mate in 1: move=%s score=%d nodes=%d", res.Move, res.Score, res.Nodes)
}
