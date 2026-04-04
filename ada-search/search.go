package search

import (
	"math"
	"sync"
	"runtime"
	"github.com/WilliamDann/AdaEngine/ada-chess/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/movegen"
	"github.com/WilliamDann/AdaEngine/ada-chess/position"
)

const maxDepth = 64
const maxMoves = 256

// Precomputed LMR reduction table: lmrTable[depth][moveIndex]
var lmrTable [maxDepth][maxMoves]int

func init() {
	for d := 1; d < maxDepth; d++ {
		for m := 1; m < maxMoves; m++ {
			lmrTable[d][m] = int(0.5 + math.Log(float64(d))*math.Log(float64(m))/2.0)
		}
	}
}

const (
	Mate   = 30_000
	Inf    = Mate + 1
	maxPly = 128
)

// killers stores two killer moves per ply — quiet moves that caused
// beta cutoffs in sibling nodes at the same depth.
type killers [maxPly][2]core.Move

func (k *killers) store(ply int, m core.Move) {
	if ply >= maxPly {
		return
	}
	if k[ply][0] != m {
		k[ply][1] = k[ply][0]
		k[ply][0] = m
	}
}

func (k *killers) isKiller(ply int, m core.Move) bool {
	if ply >= maxPly {
		return false
	}
	return k[ply][0] == m || k[ply][1] == m
}

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
	var k killers

	moves := movegen.LegalMoves(pos)
	n := moves.Count()
	ordered := make([]core.Move, n)
	scores := make([]int, n)
	for i := 0; i < n; i++ {
		ordered[i] = moves.Get(i)
	}

	const aspirationWindow = 50

	for d := 1; d <= depth; d++ {
		alpha := -Inf
		beta := Inf

		// Aspiration window: use previous score to narrow the search
		if d >= 4 && best.Score > -Mate+100 && best.Score < Mate-100 {
			alpha = best.Score - aspirationWindow
			beta = best.Score + aspirationWindow
		}

	research:
		for i := 0; i < n; i++ {
			child := position.MakeMove(pos, ordered[i])
			nodes++
			score := -alphabeta(tt, &k, child, 1, d-1, -beta, -alpha, &nodes)
			scores[i] = score
			if score > alpha {
				alpha = score
			}
		}

		// Check if aspiration window failed — re-search with full window
		bestIdx := 0
		for i := 1; i < n; i++ {
			if scores[i] > scores[bestIdx] {
				bestIdx = i
			}
		}
		if d >= 4 && (scores[bestIdx] <= best.Score-aspirationWindow || scores[bestIdx] >= best.Score+aspirationWindow) {
			alpha = -Inf
			beta = Inf
			goto research
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

func alphabeta(tt *TT, k *killers, pos *position.Position, ply int, depth int, alpha, beta int, nodes *uint64) int {
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

	if depth == 0 {
		return quiesce(tt, pos, ply, alpha, beta, nodes)
	}

	inCheck := movegen.InCheck(pos)

	// null move pruning (if we can skip a move and be winning just prune)
	if depth >= 3 && !inCheck {
		nullChild := position.MakeNullMove(pos)
		nullScore := -alphabeta(tt, k, nullChild, ply+1, depth-3, -beta, -beta+1, nodes)
		if nullScore >= beta {
			return beta
		}
	}

	// futility pruning: at shallow depths, skip quiet moves that can't beat alpha
	futile := false
	if depth <= 2 && !inCheck {
		if Evaluate(pos)+depth*150 <= alpha {
			futile = true
		}
	}

	// search best tt move
	if found && entry.Move != core.NoMove {
		child := position.MakeMove(pos, entry.Move)
		*nodes++
		score := -alphabeta(tt, k, child, ply+1, depth-1, -beta, -alpha, nodes)
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
		if k.isKiller(ply, moves.Get(i)) {
			bestScore += 50
		}
		for j := i + 1; j < moves.Count(); j++ {
			if found && moves.Get(j) == entry.Move {
				continue
			}
			s := mvvlva(pos, moves.Get(j))
			if k.isKiller(ply, moves.Get(j)) {
				s += 50
			}
			if s > bestScore {
				bestScore = s
				bestIdx   = j
			}
		}
		moves.Swap(i, bestIdx)
		mv := moves.Get(i)
		child := position.MakeMove(pos, mv)
		*nodes++

		isCapture := pos.Board.Check(mv.To()).Type() != 0
		givesCheck := movegen.InCheck(child)

		if futile && !isCapture && !givesCheck {
			continue
		}

		// LMR: reduced search for late quiet moves
		var score int
		doFull := true
		if i >= 3 && depth >= 3 && !isCapture && !givesCheck {
			R := lmrTable[min(depth, maxDepth-1)][min(i, maxMoves-1)]
			if R < 1 {
				R = 1
			}
			if R >= depth {
				R = depth - 1
			}
			score = -alphabeta(tt, k, child, ply+1, depth-1-R, -(alpha+1), -alpha, nodes)
			doFull = score > alpha
		}
		if doFull {
			score = -alphabeta(tt, k, child, ply+1, depth-1, -beta, -alpha, nodes)
		}
		if score >= beta {
			if !isCapture {
				k.store(ply, mv)
			}
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
