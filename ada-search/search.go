package search

import (
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

// Search runs iterative deepening alpha-beta to the given depth.
// The optional onDepth callback is called after each iteration completes.
func Search(pos *position.Position, depth int, onDepth ...func(Result)) Result {
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
		alpha := -Inf
		beta := Inf
		for i := 0; i < n; i++ {
			child := position.MakeMove(pos, ordered[i])
			best.Nodes++
			score := -alphabeta(child, d-1, -beta, -alpha, &best.Nodes)
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

		if len(onDepth) > 0 && onDepth[0] != nil {
			onDepth[0](best)
		}

		// Sort moves descending by score for next iteration
		sortMoves(ordered, scores, n)
	}
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

func alphabeta(pos *position.Position, depth int, alpha, beta int, nodes *uint64) int {
	moves := movegen.LegalMoves(pos)

	// Terminal: no legal moves
	if moves.Count() == 0 {
		if movegen.InCheck(pos) {
			return -Mate // Checkmated
		}
		return 0 // Stalemate
	}

	if depth == 0 {
		return Evaluate(pos)
	}

	for i := 0; i < moves.Count(); i++ {
		child := position.MakeMove(pos, moves.Get(i))
		*nodes++
		score := -alphabeta(child, depth-1, -beta, -alpha, nodes)
		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha = score
		}
	}
	return alpha
}
