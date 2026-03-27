package search

import (
	"testing"
	"time"

	"github.com/WilliamDann/AdaEngine/ada-chess/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/fen"
	"github.com/WilliamDann/AdaEngine/ada-chess/movegen"
	"github.com/WilliamDann/AdaEngine/ada-chess/position"
)

func TestTTStoreAndProbe(t *testing.T) {
	tt := NewTT(1024)
	hash := uint64(0xDEADBEEF)
	move := core.NewMove(core.NewSquare(1, 4), core.NewSquare(3, 4))

	tt.Store(hash, move, 150, 5, FlagExact)

	entry, hit := tt.Probe(hash)
	if !hit {
		t.Fatal("expected TT hit")
	}
	if entry.Move != move {
		t.Errorf("move: got %s, want %s", entry.Move, move)
	}
	if entry.Score != 150 {
		t.Errorf("score: got %d, want 150", entry.Score)
	}
	if entry.Depth != 5 {
		t.Errorf("depth: got %d, want 5", entry.Depth)
	}
	if entry.Flag != FlagExact {
		t.Errorf("flag: got %d, want FlagExact", entry.Flag)
	}
}

func TestTTProbeMiss(t *testing.T) {
	tt := NewTT(1024)
	tt.Store(0xAAAA, core.NoMove, 0, 1, FlagExact)

	_, hit := tt.Probe(0xBBBB)
	if hit {
		t.Error("expected TT miss for different hash")
	}
}

func TestTTOverwrite(t *testing.T) {
	tt := NewTT(1024)
	hash := uint64(0x1234)
	move1 := core.NewMove(core.NewSquare(1, 4), core.NewSquare(3, 4))
	move2 := core.NewMove(core.NewSquare(1, 3), core.NewSquare(3, 3))

	tt.Store(hash, move1, 100, 3, FlagAlpha)
	tt.Store(hash, move2, 200, 5, FlagExact)

	entry, hit := tt.Probe(hash)
	if !hit {
		t.Fatal("expected TT hit")
	}
	if entry.Move != move2 {
		t.Errorf("should have overwritten: got %s, want %s", entry.Move, move2)
	}
	if entry.Score != 200 {
		t.Errorf("score: got %d, want 200", entry.Score)
	}
}

func TestTTSizeRoundsToPowerOfTwo(t *testing.T) {
	tt := NewTT(1000)
	// 1000 rounds up to 1024
	if len(tt.slots) != 1024 {
		t.Errorf("size: got %d, want 1024", len(tt.slots))
	}
	if tt.mask != 1023 {
		t.Errorf("mask: got %d, want 1023", tt.mask)
	}
}

func TestTTFlags(t *testing.T) {
	tt := NewTT(1024)
	hash := uint64(0xCAFE)
	move := core.NewMove(core.NewSquare(0, 1), core.NewSquare(2, 2))

	flags := []TTFlag{FlagExact, FlagAlpha, FlagBeta}
	for _, flag := range flags {
		tt.Store(hash, move, 50, 3, flag)
		entry, hit := tt.Probe(hash)
		if !hit {
			t.Fatalf("expected hit for flag %d", flag)
		}
		if entry.Flag != flag {
			t.Errorf("flag: got %d, want %d", entry.Flag, flag)
		}
	}
}

func TestTTReducesNodeCount(t *testing.T) {
	pos, _ := fen.Parse("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	// Search WITH TT (the normal path)
	resTT := Search(pos, 5)

	// Search WITHOUT TT: bare alpha-beta
	moves := movegen.LegalMoves(pos)
	n := moves.Count()
	ordered := make([]core.Move, n)
	scores := make([]int, n)
	for i := 0; i < n; i++ {
		ordered[i] = moves.Get(i)
	}
	var nodesNoTT uint64
	for d := 1; d <= 5; d++ {
		alpha := -Inf
		beta := Inf
		for i := 0; i < n; i++ {
			child := position.MakeMove(pos, ordered[i])
			nodesNoTT++
			score := -alphabetaNoTT(child, d-1, -beta, -alpha, &nodesNoTT)
			scores[i] = score
			if score > alpha {
				alpha = score
			}
		}
	}

	t.Logf("with TT: %d nodes, without TT: %d nodes", resTT.Nodes, nodesNoTT)
	if resTT.Nodes >= nodesNoTT {
		t.Errorf("TT should reduce node count: with=%d >= without=%d", resTT.Nodes, nodesNoTT)
	}
}

func TestTTNodesPerSecond(t *testing.T) {
	pos, _ := fen.Parse("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	depth := 5

	// Without TT
	moves := movegen.LegalMoves(pos)
	n := moves.Count()
	ordered := make([]core.Move, n)
	scores := make([]int, n)
	for i := 0; i < n; i++ {
		ordered[i] = moves.Get(i)
	}
	var nodesNoTT uint64
	startNoTT := time.Now()
	for d := 1; d <= depth; d++ {
		alpha := -Inf
		beta := Inf
		for i := 0; i < n; i++ {
			child := position.MakeMove(pos, ordered[i])
			nodesNoTT++
			score := -alphabetaNoTT(child, d-1, -beta, -alpha, &nodesNoTT)
			scores[i] = score
			if score > alpha {
				alpha = score
			}
		}
	}
	elapsedNoTT := time.Since(startNoTT)

	// With TT
	startTT := time.Now()
	resTT := Search(pos, depth)
	elapsedTT := time.Since(startTT)

	npsNoTT := float64(nodesNoTT) / elapsedNoTT.Seconds()
	npsTT := float64(resTT.Nodes) / elapsedTT.Seconds()

	t.Logf("without TT: %d nodes in %v (%.0f nps)", nodesNoTT, elapsedNoTT, npsNoTT)
	t.Logf("with    TT: %d nodes in %v (%.0f nps)", resTT.Nodes, elapsedTT, npsTT)
	t.Logf("node reduction: %.1f%%", 100*(1-float64(resTT.Nodes)/float64(nodesNoTT)))
	t.Logf("time reduction: %.1f%%", 100*(1-elapsedTT.Seconds()/elapsedNoTT.Seconds()))
}

func TestSearchParallel(t *testing.T) {
	pos, _ := fen.Parse("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	depth := 5

	// Single-threaded
	startST := time.Now()
	resST := Search(pos, depth)
	elapsedST := time.Since(startST)

	// Multi-threaded (4 threads)
	startMT := time.Now()
	resMT := SearchParallel(pos, depth, 4)
	elapsedMT := time.Since(startMT)

	t.Logf("single-thread: move=%s score=%d nodes=%d time=%v", resST.Move, resST.Score, resST.Nodes, elapsedST)
	t.Logf("4 threads:     move=%s score=%d nodes=%d time=%v", resMT.Move, resMT.Score, resMT.Nodes, elapsedMT)
	t.Logf("wall time reduction: %.1f%%", 100*(1-elapsedMT.Seconds()/elapsedST.Seconds()))

	if resMT.Move == core.NoMove {
		t.Error("parallel search should return a move")
	}
}

// bare alpha-beta for comparison — no TT
func alphabetaNoTT(pos *position.Position, depth int, alpha, beta int, nodes *uint64) int {
	moves := movegen.LegalMoves(pos)
	if moves.Count() == 0 {
		if movegen.InCheck(pos) {
			return -Mate
		}
		return 0
	}
	if depth == 0 {
		return Evaluate(pos)
	}
	for i := 0; i < moves.Count(); i++ {
		child := position.MakeMove(pos, moves.Get(i))
		*nodes++
		score := -alphabetaNoTT(child, depth-1, -beta, -alpha, nodes)
		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha = score
		}
	}
	return alpha
}
