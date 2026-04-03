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

// Search runs iterative deepening alpha-beta to the given depth.
// The optional onDepth callback is called after each iteration completes.
func Search(pos *position.Position, depth int, onDepth ...func(Result)) Result {
	var best Result
	best.Score = -Inf

	tt := NewTT(1 << 22)

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
			score := -alphabeta(tt, child, 1, d-1, -beta, -alpha, &best.Nodes)
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

	if depth == 0 {
		return quiesce(tt, pos, ply, alpha, beta, nodes)
	}

	for i := 0; i < moves.Count(); i++ {
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
