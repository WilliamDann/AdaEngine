package search

import (
	"sync"
	"runtime"
	"github.com/WilliamDann/AdaEngine/ada-chess/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/movegen"
	"github.com/WilliamDann/AdaEngine/ada-chess/position"
)

const (
	Mate = 30_000
	Inf  = Mate + 1
)

// Result holds the outcome of a search.
type Result struct {
	Move  core.Move
	Score int
	Depth int
	Nodes uint64
}

func quiesce(tt *TT, pos *position.Position, ply int, alpha, beta int, nodes *uint64) int {
	entry, found := tt.Probe(pos.Zobrist)
	startAlpha   := alpha
	if found {
		score := adjustScoreForProbe(entry.Score, ply)

		if entry.Flag == Exact {
			return score
		}
		if entry.Flag == LowerBound {
			alpha = max(alpha, score)
		}
		if entry.Flag == UpperBound {
			beta = min(beta, score)
		}
		if alpha >= beta {
			return score
		}
	}

	stand := Evaluate(pos)
	if stand >= beta {
		return beta
	}
	if stand > alpha {
		alpha = stand
	}

	captures := movegen.LegalCaptures(pos)
	for i := 0; i < captures.Count(); i++ {
		captured := pos.Board.Check(captures.Get(i).To()).Type()
		if stand + pieceValue[captured] + 200 < alpha {
			continue
		}

		child := position.MakeMove(pos, captures.Get(i))
		*nodes++

		score := -quiesce(tt, child, ply+1, -beta, -alpha, nodes)
		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha = score
		}
	}


	// store move in the transposition table
	var flagType SearchFlag = Exact
	if alpha <= startAlpha {
		flagType = UpperBound
	} else if alpha >= beta {
		flagType = LowerBound
	}

	tt.Store(TTEntry{
		Key: pos.Zobrist,
		Move: core.NoMove,
		Depth: 0,
		Score: adjustScoreForStore(alpha, ply),
		Flag: flagType,
	})

	return alpha
}

func adjustScoreForStore(score int, ply int) int16 {
	if score > Mate-100 {
		return int16(score + ply)
	}
	if score < -Mate+100 {
		return int16(score - ply)
	}
	return int16(score)
}

func adjustScoreForProbe(score int16, ply int) int {
	s := int(score)
	if s > Mate-100 {
		return s - ply
	}
	if s < -Mate+100 {
		return s + ply
	}
	return s
}

// score for diff in attacking vs attacked peices
func mvvlva(pos *position.Position, m core.Move) int {
	victim := pos.Board.Check(m.To()).Type()
	attacker := pos.Board.Check(m.From()).Type()
	if victim == 0 {
		return 0 // quiet move
	}
	return pieceValue[victim] - pieceValue[attacker]
}

// Search runs iterative deepening alpha-beta to the given depth.
// The optional onDepth callback is called after each iteration completes.
func Search(pos *position.Position, depth int, threads int, onDepth ...func(Result)) Result {
	tt := NewTT(1 << 22)

	numThreads := threads
	if numThreads <= 0 {
		numThreads = runtime.NumCPU()
	}
	
	results    := make([]Result, numThreads)
	var wg sync.WaitGroup

	for t := 0; t < numThreads; t++ {
		wg.Add(1)
		var cb func(Result)
		if t == 0 && len(onDepth) > 0 {
			cb = onDepth[0]
		}
		go func(thread int, callback func(Result)) {
			defer wg.Done()
			results[thread] = searchWorker(tt, pos, depth, thread, callback)
		}(t, cb)
	}
	wg.Wait()

	best := results[0]
	for _, r := range results[1:] {
		if r.Depth > best.Depth || (r.Depth == best.Depth && r.Score > best.Score) {
			best = r
		}
	}

	return best
}

func searchWorker(tt *TT, pos *position.Position, depth int, thread int, onDepth func(Result)) Result {
	var best Result
	best.Score = -Inf
	var nodes uint64

	moves := movegen.LegalMoves(pos)
	n := moves.Count()
	ordered := make([]core.Move, n)
	scores := make([]int, n)
	for i := 0; i < n; i++ {
		ordered[i] = moves.Get(i)
	}

	for d := 1; d <= depth; d++ {
		alpha := -Inf
		beta := Inf
		for i := 0; i < n; i++ {
			child := position.MakeMove(pos, ordered[i])
			nodes++
			score := -alphabeta(tt, child, 1, d-1, -beta, -alpha, &nodes)
			scores[i] = score
			if score > alpha {
				alpha = score
			}
		}
		// Find best move for this iteration
		bestIdx := 0
		for i := 1; i < n; i++ {
			if scores[i] > scores[bestIdx] {
				bestIdx = i
			}
		}
		best.Move = ordered[bestIdx]
		best.Score = scores[bestIdx]
		best.Depth = d

		if onDepth != nil {
			onDepth(best)
		}

		// Sort moves descending by score for next iteration
		sortMoves(ordered, scores, n)
	}

	best.Nodes = nodes
	return best
}

// sortMoves does an insertion sort of moves by descending score.
func sortMoves(moves []core.Move, scores []int, n int) {
	for i := 1; i < n; i++ {
		m, s := moves[i], scores[i]
		j := i
		for j > 0 && scores[j-1] < s {
			moves[j], scores[j] = moves[j-1], scores[j-1]
			j--
		}
		moves[j], scores[j] = m, s
	}
}

func alphabeta(tt *TT, pos *position.Position, ply int, depth int, alpha, beta int, nodes *uint64) int {
	// look up in transposition table
	entry, found := tt.Probe(pos.Zobrist)
	startAlpha   := alpha
	bestMove     := core.NoMove
	if found {
		if entry.Depth >= int8(depth) {
			score := adjustScoreForProbe(entry.Score, ply)

			if entry.Flag == Exact {
				return score
			}
			if entry.Flag == LowerBound {
				alpha = max(alpha, score)
			}
			if entry.Flag == UpperBound {
				beta = min(beta, score)
			}
			if alpha >= beta {
				return score
			}
		}
	}

	moves := movegen.LegalMoves(pos)

	// Terminal: no legal moves
	if moves.Count() == 0 {
		if movegen.InCheck(pos) {
			return -Mate // Checkmated
		}
		return 0 // Stalemate
	}

	// extend line when in check
	if movegen.InCheck(pos) {
		depth++
	}

	if depth == 0 {
		return quiesce(tt, pos, ply, alpha, beta, nodes)
	}

	// null move pruning (if we can skip a move and be winning just prune)
	if depth >= 3 && !movegen.InCheck(pos) {
		nullChild := position.MakeNullMove(pos)
		nullScore := -alphabeta(tt, nullChild, ply+1, depth-3, -beta, -beta+1, nodes)
		if nullScore >= beta {
			return beta
		}
	}

	// search best tt move
	if found && entry.Move != core.NoMove {
		child := position.MakeMove(pos, entry.Move)
		*nodes++
		score := -alphabeta(tt, child, ply+1, depth-1, -beta, -alpha, nodes)
		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha = score
			bestMove = entry.Move
		}
	}

	for i := 0; i < moves.Count(); i++ {
		// skip TT move
		if found && moves.Get(i) == entry.Move {
			continue
		}

		bestIdx   := i
		bestScore := mvvlva(pos, moves.Get(i))
		for j := i + 1; j < moves.Count(); j++ {
			if found && moves.Get(j) == entry.Move {
				continue
			}
			s := mvvlva(pos, moves.Get(j))
			if s > bestScore {
				bestScore = s
				bestIdx   = j
			}
		}
		moves.Swap(i, bestIdx)

		child := position.MakeMove(pos, moves.Get(i))
		*nodes++
		score := -alphabeta(tt, child, ply+1, depth-1, -beta, -alpha, nodes)
		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha    = score
			bestMove = moves.Get(i)
		}
	}

	// store move in the transposition table
	var flagType SearchFlag = Exact
	if alpha <= startAlpha {
		flagType = UpperBound
	} else if alpha >= beta {
		flagType = LowerBound
	}

	tt.Store(TTEntry{
		Key: pos.Zobrist,
		Move: bestMove,
		Depth: int8(depth),
		Score: adjustScoreForStore(alpha, ply),
		Flag: flagType,
	})

	return alpha
}
