package engine

import (
	"github.com/WilliamDann/adachess/game"
)

var MaterialValue map[game.PieceType]float64 = map[game.PieceType]float64{
	game.Pawn:   1,
	game.Bishop: 3,
	game.Knight: 3,
	game.Rook:   5,
	game.Queen:  9,
	game.King:   0,
}

func Mobility(position *game.Position) float64 {
	restore := position.GetFen().ActiveColor
	var score float64

	if !restore {
		position.Pass()
	}

	score += float64(len(game.ApplyRuleSet(position, game.StandardRules)))
	position.Pass()
	score -= float64(len(game.ApplyRuleSet(position, game.StandardRules)))

	if position.GetFen().ActiveColor != restore {
		position.Pass()
	}

	return score
}

func Material(position *game.Position) float64 {
	var value float64 = 0

	for piece, coords := range position.GetBoard().Pieces() {
		var mod float64 = float64(len(coords))
		if !piece.Color {
			mod *= -1
		}
		value += MaterialValue[piece.Type] * mod
	}

	return value
}

func Eval(position *game.Position) float64 {
	return Material(position) + 0.05*Mobility(position)
}
