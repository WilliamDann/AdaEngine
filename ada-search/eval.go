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
	for pt := core.PieceType(1); pt <= 5; pt++ {
		val := pieceValue[pt]
		white := pos.Board.Pieces(core.NewPiece(pt, core.White)).Count()
		black := pos.Board.Pieces(core.NewPiece(pt, core.Black)).Count()
		if pos.ActiveColor == core.White {
			score += val * (white - black)
		} else {
			score += val * (black - white)
		}
	}
	return score
}
