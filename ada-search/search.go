package search

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/WilliamDann/AdaEngine/ada-chess/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/movegen"
	"github.com/WilliamDann/AdaEngine/ada-chess/position"
)

const (
	Mate   = 1_000_000
	Inf    = Mate + 1
	maxPly = 128
)

// Result holds the outcome of a search.
type Result struct {
	Move  core.Move
	Score int
	Depth int
	Nodes uint64
}

// searchState holds per-goroutine search state.
type searchState struct {
	tt      *TT
	stop    *atomic.Bool
	nodes   uint64
	killers [maxPly][2]core.Move
	history [2][64][64]int
}

func newSearchState(tt *TT, stop *atomic.Bool) *searchState {
	return &searchState{tt: tt, stop: stop}
}

// Search runs single-threaded iterative deepening alpha-beta to the given depth.
// The optional onDepth callback is called after each iteration completes.
func Search(pos *position.Position, depth int, onDepth ...func(Result)) Result {
	tt := NewTT(1 << 22)
	return searchMain(pos, depth, tt, nil, time.Time{}, onDepth...)
}

// SearchTimed runs single-threaded iterative deepening, stopping after the
// given duration. The search is aborted mid-depth if time runs out.
func SearchTimed(pos *position.Position, limit time.Duration, onDepth ...func(Result)) Result {
	tt := NewTT(1 << 22)
	stop := &atomic.Bool{}
	timer := time.AfterFunc(limit, func() { stop.Store(true) })
	defer timer.Stop()
	return searchMain(pos, maxPly, tt, stop, time.Time{}, onDepth...)
}

// SearchParallel runs Lazy SMP search with the given number of threads.
// If threads <= 0, uses runtime.NumCPU().
func SearchParallel(pos *position.Position, depth int, threads int, onDepth ...func(Result)) Result {
	if threads <= 0 {
		threads = runtime.NumCPU()
	}
	if threads == 1 {
		return Search(pos, depth, onDepth...)
	}

	tt := NewTT(1 << 22)
	stop := &atomic.Bool{}

	// Launch helper threads that fill the shared TT.
	// Each helper gets a different ID to create search diversity.
	var wg sync.WaitGroup
	helperNodes := make([]uint64, threads-1)
	for i := 0; i < threads-1; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			searchWorker(pos, depth, tt, stop, &helperNodes[id], id)
		}(i)
	}

	// Main thread runs the authoritative search.
	result := searchMain(pos, depth, tt, stop, time.Time{}, onDepth...)

	// Signal helpers to stop and wait for them.
	stop.Store(true)
	wg.Wait()

	for _, n := range helperNodes {
		result.Nodes += n
	}
	return result
}

// SearchTimedParallel runs Lazy SMP with a time limit instead of a depth limit.
// A timer sets the stop signal when time expires, aborting all threads mid-search.
func SearchTimedParallel(pos *position.Position, limit time.Duration, threads int, onDepth ...func(Result)) Result {
	if threads <= 0 {
		threads = runtime.NumCPU()
	}

	tt := NewTT(1 << 22)
	stop := &atomic.Bool{}
	timer := time.AfterFunc(limit, func() { stop.Store(true) })
	defer timer.Stop()

	var wg sync.WaitGroup
	helperNodes := make([]uint64, threads-1)
	for i := 0; i < threads-1; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			searchWorker(pos, maxPly, tt, stop, &helperNodes[id], id)
		}(i)
	}

	result := searchMain(pos, maxPly, tt, stop, time.Time{}, onDepth...)

	stop.Store(true)
	wg.Wait()

	for _, n := range helperNodes {
		result.Nodes += n
	}
	return result
}

// searchMain runs iterative deepening on the calling goroutine.
// When used with a timed search, the stop signal is set by an external
// timer. If stop fires mid-depth, the partial depth is discarded and
// the result from the last fully completed depth is returned.
func searchMain(pos *position.Position, depth int, tt *TT, stop *atomic.Bool, deadline time.Time, onDepth ...func(Result)) Result {
	state := newSearchState(tt, stop)
	_ = deadline // kept for API compatibility; timing now uses stop signal

	var best Result
	best.Score = -Inf

	moves := movegen.LegalMoves(pos)
	n := moves.Count()
	ordered := make([]core.Move, n)
	scores := make([]int, n)
	for i := 0; i < n; i++ {
		ordered[i] = moves.Get(i)
	}

	for d := 1; d <= depth; d++ {
		if stop != nil && stop.Load() {
			break
		}
		alpha := -Inf
		beta := Inf
		aborted := false
		for i := 0; i < n; i++ {
			child := position.MakeMove(pos, ordered[i])
			state.nodes++
			score := -state.alphabeta(child, d-1, 1, -beta, -alpha, true)
			if stop != nil && stop.Load() {
				aborted = true
				break
			}
			scores[i] = score
			if score > alpha {
				alpha = score
			}
		}
		if aborted {
			break // discard partial depth, keep last complete result
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
		best.Nodes = state.nodes

		if len(onDepth) > 0 && onDepth[0] != nil {
			onDepth[0](best)
		}

		// Sort moves descending by score for next iteration
		sortMoves(ordered, scores, n)
	}
	best.Nodes = state.nodes
	return best
}

// searchWorker runs iterative deepening as a helper, filling the shared TT
// until stop is signalled. Each helper gets a unique id used to rotate the
// initial move order, so different threads explore different branches first
// and produce useful TT entries for the main thread.
func searchWorker(pos *position.Position, depth int, tt *TT, stop *atomic.Bool, nodes *uint64, id int) {
	state := newSearchState(tt, stop)

	moves := movegen.LegalMoves(pos)
	n := moves.Count()
	ordered := make([]core.Move, n)
	for i := 0; i < n; i++ {
		ordered[i] = moves.Get(i)
	}

	// Rotate move list by thread ID so each helper starts with a
	// different root move, creating natural search diversity.
	offset := (id + 1) % n
	rotated := make([]core.Move, n)
	for i := 0; i < n; i++ {
		rotated[i] = ordered[(i+offset)%n]
	}
	copy(ordered, rotated)

	scores := make([]int, n)

	for d := 1; d <= depth; d++ {
		if stop.Load() {
			return
		}
		alpha := -Inf
		beta := Inf
		for i := 0; i < n; i++ {
			if stop.Load() {
				return
			}
			child := position.MakeMove(pos, ordered[i])
			state.nodes++
			score := -state.alphabeta(child, d-1, 1, -beta, -alpha, true)
			scores[i] = score
			if score > alpha {
				alpha = score
			}
		}
		sortMoves(ordered, scores, n)
	}
	*nodes = state.nodes
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

// scoreToTT adjusts a mate score for TT storage by converting from
// root-relative to position-relative.
func scoreToTT(score, ply int) int {
	if score > Mate-maxPly {
		return score + ply
	}
	if score < -Mate+maxPly {
		return score - ply
	}
	return score
}

// scoreFromTT reverses the adjustment when reading a mate score from the TT.
func scoreFromTT(score, ply int) int {
	if score > Mate-maxPly {
		return score - ply
	}
	if score < -Mate+maxPly {
		return score + ply
	}
	return score
}

func (s *searchState) alphabeta(pos *position.Position, depth, ply int, alpha, beta int, allowNull bool) int {
	if s.stop != nil && s.stop.Load() {
		return 0
	}

	moves := movegen.LegalMoves(pos)

	// Terminal: no legal moves
	if moves.Count() == 0 {
		if movegen.InCheck(pos) {
			return -(Mate - ply) // Checkmated — prefer shorter mates
		}
		return 0 // Stalemate
	}

	if depth == 0 {
		return Evaluate(pos)
	}

	// TT probe
	entry, hit := s.tt.Probe(pos.Hash)
	if hit && int(entry.Depth) >= depth {
		ttScore := scoreFromTT(int(entry.Score), ply)
		switch entry.Flag {
		case FlagExact:
			return ttScore
		case FlagAlpha:
			if ttScore <= alpha {
				return alpha
			}
		case FlagBeta:
			if ttScore >= beta {
				return ttScore
			}
		}
	}

	inCheck := movegen.InCheck(pos)

	// Null-move pruning: if giving the opponent an extra turn still
	// results in a beta cutoff, this position is likely very strong
	// and we can prune the subtree.
	if allowNull && !inCheck && depth >= 3 {
		null := position.NullMove(pos)
		s.nodes++
		score := -s.alphabeta(null, depth-3, ply+1, -beta, -beta+1, false)
		if score >= beta {
			return beta
		}
	}

	// --- Move ordering ---
	orderStart := 0

	// 1. TT best move first
	if hit && entry.Move != core.NoMove {
		for j := 0; j < moves.Count(); j++ {
			if moves.Get(j) == entry.Move {
				moves.Swap(orderStart, j)
				orderStart++
				break
			}
		}
	}

	// 2. Killer moves next
	if ply < maxPly {
		for k := 0; k < 2; k++ {
			killer := s.killers[ply][k]
			if killer == core.NoMove {
				continue
			}
			for j := orderStart; j < moves.Count(); j++ {
				if moves.Get(j) == killer {
					moves.Swap(orderStart, j)
					orderStart++
					break
				}
			}
		}
	}

	origAlpha := alpha
	var bestMove core.Move
	searched := 0

	ci := 0
	if pos.ActiveColor == core.Black {
		ci = 1
	}

	for i := 0; i < moves.Count(); i++ {
		// 3. History-based ordering for remaining moves (selection sort)
		if i >= orderStart {
			bestJ := i
			bestH := s.history[ci][moves.Get(i).From()][moves.Get(i).To()]
			for j := i + 1; j < moves.Count(); j++ {
				mj := moves.Get(j)
				h := s.history[ci][mj.From()][mj.To()]
				if h > bestH {
					bestH = h
					bestJ = j
				}
			}
			if bestJ != i {
				moves.Swap(i, bestJ)
			}
		}

		mv := moves.Get(i)
		child := position.MakeMove(pos, mv)
		s.nodes++

		isCapture := pos.Board.HasPiece(mv.To()) || mv.MoveType() == core.MoveEnPassant

		// Late move reductions: search later quiet moves at reduced
		// depth first; if they look good, re-search at full depth.
		var score int
		doFull := true
		if searched >= 3 && depth >= 3 && !isCapture && mv.MoveType() != core.MovePromotion && !inCheck {
			score = -s.alphabeta(child, depth-2, ply+1, -(alpha + 1), -alpha, true)
			doFull = score > alpha
		}

		if doFull {
			score = -s.alphabeta(child, depth-1, ply+1, -beta, -alpha, true)
		}
		searched++

		if score >= beta {
			s.tt.Store(pos.Hash, mv, scoreToTT(score, ply), depth, FlagBeta)
			// Update killers and history for quiet cutoff moves
			if !isCapture {
				if ply < maxPly && s.killers[ply][0] != mv {
					s.killers[ply][1] = s.killers[ply][0]
					s.killers[ply][0] = mv
				}
				s.history[ci][mv.From()][mv.To()] += depth * depth
			}
			return beta
		}
		if score > alpha {
			alpha = score
			bestMove = mv
		}
	}

	// Store result
	flag := FlagAlpha
	if alpha > origAlpha {
		flag = FlagExact
	}
	s.tt.Store(pos.Hash, bestMove, scoreToTT(alpha, ply), depth, flag)

	return alpha
}
