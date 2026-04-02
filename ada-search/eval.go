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

// Advancement bonus per rank in centipawns. Rank 0 (home rank) gets nothing,
// rank 7 (promotion rank) gets the most.
var advanceBonus = [7]int{
	0,  // None
	5,  // Pawn
	3,  // Knight
	2,  // Bishop
	1,  // Rook
	1,  // Queen
	0,  // King
}

// Proximity bonus per piece type, awarded per unit of closeness to the
// opponent king (14 - manhattan distance). Higher values make the piece
// gravitate toward the enemy king.
var proximityBonus = [7]int{
	0,  // None
	0,  // Pawn
	2,  // Knight
	1,  // Bishop
	1,  // Rook
	3,  // Queen
	0,  // King
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func findKing(pos *position.Position, color core.Color) core.Square {
	bb := pos.Board.Pieces(core.NewPiece(core.King, color))
	for sq := range bb.Squares() {
		return sq
	}
	return core.Square(0)
}

// Evaluate returns a score in centipawns from the active color's perspective.
// Positive means the active color is better.
func Evaluate(pos *position.Position) int {
	wksq := findKing(pos, core.White)
	bksq := findKing(pos, core.Black)

	score := 0
	for pt := core.PieceType(1); pt <= 5; pt++ {
		val := pieceValue[pt]
		adv := advanceBonus[pt]
		prox := proximityBonus[pt]

		white := pos.Board.Pieces(core.NewPiece(pt, core.White))
		black := pos.Board.Pieces(core.NewPiece(pt, core.Black))

		ws := 0
		for sq := range white.Squares() {
			r, f := sq.Rank(), sq.File()
			ws += val + adv*r
			if prox > 0 {
				ws += prox * (14 - abs(r-bksq.Rank()) - abs(f-bksq.File()))
			}
		}

		bs := 0
		for sq := range black.Squares() {
			r, f := sq.Rank(), sq.File()
			bs += val + adv*(7-r)
			if prox > 0 {
				bs += prox * (14 - abs(r-wksq.Rank()) - abs(f-wksq.File()))
			}
		}

		if pos.ActiveColor == core.White {
			score += ws - bs
		} else {
			score += bs - ws
		}
	}
	return score
}
