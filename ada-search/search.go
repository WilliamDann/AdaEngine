package search

import (
	"github.com/WilliamDann/AdaEngine/ada-chess/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/movegen"
	"github.com/WilliamDann/AdaEngine/ada-chess/position"
)

const (
	Mate = 1_000_000
	Inf  = Mate + 1
)

// Result holds the outcome of a search.
type Result struct {
	Move  core.Move
	Score int
	Nodes uint64
}

// Search runs a minimax search to the given depth and returns the best move.
func Search(pos *position.Position, depth int) Result {
	var res Result
	res.Score = -Inf

	moves := movegen.LegalMoves(pos)
	for i := 0; i < moves.Count(); i++ {
		m := moves.Get(i)
		child := position.MakeMove(pos, m)
		res.Nodes++
		score := -minimax(child, depth-1, &res.Nodes)
		if score > res.Score {
			res.Score = score
			res.Move = m
		}
	}
	return res
}

func minimax(pos *position.Position, depth int, nodes *uint64) int {
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

	best := -Inf
	for i := 0; i < moves.Count(); i++ {
		child := position.MakeMove(pos, moves.Get(i))
		*nodes++
		score := -minimax(child, depth-1, nodes)
		if score > best {
			best = score
		}
	}
	return best
}
