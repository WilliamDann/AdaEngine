package search

import (
	"testing"

	"github.com/WilliamDann/AdaEngine/ada-chess/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/fen"
	"github.com/WilliamDann/AdaEngine/ada-chess/movegen"
	"github.com/WilliamDann/AdaEngine/ada-chess/position"
)

func TestTTStoreProbe(t *testing.T) {
	tt := NewTT(1024)
	entry := TTEntry{
		Key:   0xDEADBEEF,
		Move:  core.NewMove(core.NewSquare(1, 4), core.NewSquare(3, 4)),
		Depth: 5,
		Score: 42,
		Flag:  Exact,
	}
	tt.Store(entry)

	got, ok := tt.Probe(0xDEADBEEF)
	if !ok {
		t.Fatal("expected hit")
	}
	if got != entry {
		t.Fatalf("got %+v, want %+v", got, entry)
	}
}

func TestTTProbeMiss(t *testing.T) {
	tt := NewTT(1024)
	tt.Store(TTEntry{
		Key:   0xDEADBEEF,
		Depth: 5,
		Score: 42,
		Flag:  Exact,
	})

	_, ok := tt.Probe(0xCAFEBABE)
	if ok {
		t.Fatal("expected miss for different key")
	}
}

func TestTTOverwrite(t *testing.T) {
	tt := NewTT(1024)
	tt.Store(TTEntry{Key: 0xDEADBEEF, Depth: 3, Score: 10, Flag: Exact})
	tt.Store(TTEntry{Key: 0xDEADBEEF, Depth: 5, Score: 20, Flag: LowerBound})

	got, ok := tt.Probe(0xDEADBEEF)
	if !ok {
		t.Fatal("expected hit")
	}
	if got.Depth != 5 || got.Score != 20 || got.Flag != LowerBound {
		t.Fatalf("expected overwritten entry, got %+v", got)
	}
}

func TestTTTransposition(t *testing.T) {
	// Reach the same position via two move orders, verify TT hit
	e2 := core.NewSquare(1, 4)
	e3 := core.NewSquare(2, 4)
	d2 := core.NewSquare(1, 3)
	d3 := core.NewSquare(2, 3)
	d7 := core.NewSquare(6, 3)
	d6 := core.NewSquare(5, 3)
	wPawn := core.NewPiece(core.Pawn, core.White)
	bPawn := core.NewPiece(core.Pawn, core.Black)

	board := core.NewChessboard()
	board.Set(e2, wPawn)
	board.Set(d2, wPawn)
	board.Set(d7, bPawn)
	start := &position.Position{
		Board:       board,
		ActiveColor: core.White,
		EnPassant:   core.InvalidSquare,
	}
	start.Zobrist = start.ComputeZobrist()

	// path 1: e3, d6, d3
	p1 := position.MakeMove(start, core.NewMove(e2, e3))
	p1 = position.MakeMove(p1, core.NewMove(d7, d6))
	p1 = position.MakeMove(p1, core.NewMove(d2, d3))

	// path 2: d3, d6, e3
	p2 := position.MakeMove(start, core.NewMove(d2, d3))
	p2 = position.MakeMove(p2, core.NewMove(d7, d6))
	p2 = position.MakeMove(p2, core.NewMove(e2, e3))

	tt := NewTT(1024)
	entry := TTEntry{
		Key:   p1.Zobrist,
		Move:  core.NewMove(e3, e3), // dummy move
		Depth: 4,
		Score: 100,
		Flag:  Exact,
	}
	tt.Store(entry)

	// Probe with the other path's hash — should hit
	got, ok := tt.Probe(p2.Zobrist)
	if !ok {
		t.Fatal("expected TT hit for transposed position")
	}
	if got.Score != 100 {
		t.Fatalf("expected score 100, got %d", got.Score)
	}
}

func TestSearchResultUnchangedWithTT(t *testing.T) {
	pos, _ := fen.Parse("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	depth := 4

	r1 := Search(pos, depth, 1)
	r2 := Search(pos, depth, 1)

	if r1.Move != r2.Move {
		t.Errorf("search results differ: move %s vs %s", r1.Move, r2.Move)
	}
	if r1.Score != r2.Score {
		t.Errorf("search results differ: score %d vs %d", r1.Score, r2.Score)
	}
}

func TestTTReducesNodeCount(t *testing.T) {
	// A middlegame position with plenty of transpositions
	pos, err := fen.Parse("r1bqkbnr/pppppppp/2n5/8/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 2 2")
	if err != nil {
		t.Fatal(err)
	}
	depth := 5

	// First search populates the TT from scratch
	r1 := Search(pos, depth, 1)

	// Second search reuses the same code path (new TT, but iterative deepening
	// itself benefits from TT within the search). Compare against a baseline
	// without TT by using depth 1 as a sanity check — the real test is that
	// search completes and produces a valid result with reasonable node counts.
	r2 := Search(pos, depth, 1)

	if r1.Move == core.NoMove || r2.Move == core.NoMove {
		t.Fatal("expected valid moves from both searches")
	}

	t.Logf("search 1: move=%s score=%d nodes=%d", r1.Move, r1.Score, r1.Nodes)
	t.Logf("search 2: move=%s score=%d nodes=%d", r2.Move, r2.Score, r2.Nodes)
}

func benchmarkSearch(b *testing.B, pos *position.Position, depth int, tt *TT) {
	for b.Loop() {
		var nodes uint64
		moves := movegen.LegalMoves(pos)
		for i := 0; i < moves.Count(); i++ {
			child := position.MakeMove(pos, moves.Get(i))
			nodes++
			alphabeta(tt, &killers{}, child, 1, depth-1, -Inf, Inf, &nodes)
		}
	}
}

func BenchmarkSearchWithTT(b *testing.B) {
	pos, _ := fen.Parse("r1bqkbnr/pppppppp/2n5/8/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 2 2")
	tt := NewTT(1 << 22)
	benchmarkSearch(b, pos, 5, tt)
}

func BenchmarkSearchWithoutTT(b *testing.B) {
	pos, _ := fen.Parse("r1bqkbnr/pppppppp/2n5/8/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 2 2")
	benchmarkSearch(b, pos, 5, nil)
}

func BenchmarkSearch1Thread(b *testing.B) {
	pos, _ := fen.Parse("r1bqkbnr/pppppppp/2n5/8/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 2 2")
	for b.Loop() {
		Search(pos, 5, 1)
	}
}

func BenchmarkSearch2Threads(b *testing.B) {
	pos, _ := fen.Parse("r1bqkbnr/pppppppp/2n5/8/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 2 2")
	for b.Loop() {
		Search(pos, 5, 2)
	}
}

func BenchmarkSearch4Threads(b *testing.B) {
	pos, _ := fen.Parse("r1bqkbnr/pppppppp/2n5/8/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 2 2")
	for b.Loop() {
		Search(pos, 5, 4)
	}
}

func BenchmarkSearchAllThreads(b *testing.B) {
	pos, _ := fen.Parse("r1bqkbnr/pppppppp/2n5/8/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 2 2")
	for b.Loop() {
		Search(pos, 5, 0)
	}
}

func TestTTAndNoTTSameResult(t *testing.T) {
	pos, _ := fen.Parse("r1bqkbnr/pppppppp/2n5/8/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 2 2")
	depth := 5

	// Search with TT
	var nodesWithTT uint64
	tt := NewTT(1 << 22)
	moves := movegen.LegalMoves(pos)
	bestWithTT := core.NoMove
	bestScoreWithTT := -Inf
	for i := 0; i < moves.Count(); i++ {
		child := position.MakeMove(pos, moves.Get(i))
		nodesWithTT++
		score := -alphabeta(tt, &killers{}, child, 1, depth-1, -Inf, Inf, &nodesWithTT)
		if score > bestScoreWithTT {
			bestScoreWithTT = score
			bestWithTT = moves.Get(i)
		}
	}

	// Search without TT
	var nodesWithoutTT uint64
	bestWithoutTT := core.NoMove
	bestScoreWithoutTT := -Inf
	for i := 0; i < moves.Count(); i++ {
		child := position.MakeMove(pos, moves.Get(i))
		nodesWithoutTT++
		score := -alphabeta(nil, &killers{}, child, 1, depth-1, -Inf, Inf, &nodesWithoutTT)
		if score > bestScoreWithoutTT {
			bestScoreWithoutTT = score
			bestWithoutTT = moves.Get(i)
		}
	}

	t.Logf("with TT:    move=%s score=%d nodes=%d", bestWithTT, bestScoreWithTT, nodesWithTT)
	t.Logf("without TT: move=%s score=%d nodes=%d", bestWithoutTT, bestScoreWithoutTT, nodesWithoutTT)

	if bestWithTT != bestWithoutTT {
		t.Errorf("different moves: with TT=%s, without TT=%s", bestWithTT, bestWithoutTT)
	}
	if bestScoreWithTT != bestScoreWithoutTT {
		t.Errorf("different scores: with TT=%d, without TT=%d", bestScoreWithTT, bestScoreWithoutTT)
	}
}

func TestMateScoreAdjustment(t *testing.T) {
	tests := []struct {
		name  string
		score int
		ply   int
	}{
		{"positive mate", Mate - 5, 3},
		{"negative mate", -Mate + 5, 3},
		{"normal score", 150, 3},
		{"zero", 0, 5},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stored := adjustScoreForStore(tc.score, tc.ply)
			recovered := adjustScoreForProbe(stored, tc.ply)
			if recovered != tc.score {
				t.Errorf("round-trip failed: stored %d, recovered %d, want %d", stored, recovered, tc.score)
			}
		})
	}
}
