package search

import (
	"github.com/WilliamDann/AdaEngine/ada-chess/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/position"
)

// Piece values in centipawns.
var pieceValue = [7]int{
	0,    // None
	100,  // Pawn
	320,  // Knight
	330,  // Bishop
	500,  // Rook
	900,  // Queen
	0,    // King (not counted)
}

// Evaluate returns a score in centipawns from the active color's perspective.
// Positive means the active color is better.
func Evaluate(pos *position.Position) int {
	score := 0
	for sq := core.Square(0); sq < 64; sq++ {
		piece := pos.Board.Check(sq)
		if piece == core.None {
			continue
		}
		val := pieceValue[piece.Type()]
		if piece.Color() == pos.ActiveColor {
			score += val
		} else {
			score -= val
		}
	}
	return score
}
