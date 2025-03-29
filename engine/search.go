package engine

import (
	"math"

	"github.com/WilliamDann/adachess/game"
)

func Search(position *game.Position, depth int) (*game.Move, float64) {
	if position.GetFen().ActiveColor {
		return alphaBetaMax(position, math.Inf(-1), math.Inf(1), depth)
	}
	return alphaBetaMin(position, math.Inf(-1), math.Inf(1), depth)
}

func alphaBetaMax(position *game.Position, alpha float64, beta float64, depth int) (*game.Move, float64) {
	// if depth == 0 we're at a leaf node
	if depth == 0 {
		return nil, Eval(position)
	}

	var bestMove game.Move
	bestValue := math.Inf(-1)
	for _, move := range position.LegalMoves() {

		position.Move(move)
		_, score := alphaBetaMin(position, alpha, beta, depth-1)
		position.Unmove()

		if score > alpha {
			bestMove = move
			bestValue = alpha
		}
		if score >= beta {
			return &move, score
		}
	}

	return &bestMove, bestValue
}

func alphaBetaMin(position *game.Position, alpha float64, beta float64, depth int) (*game.Move, float64) {
	if depth == 0 {
		return nil, -Eval(position)
	}

	var bestMove game.Move
	bestValue := math.Inf(1)
	for _, move := range position.LegalMoves() {

		position.Move(move)
		_, score := alphaBetaMax(position, alpha, beta, depth-1)
		position.Unmove()

		if score < bestValue {
			beta = score
		}
		if score <= alpha {
			return &move, score
		}
	}

	return &bestMove, bestValue
}
