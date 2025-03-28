package perft

import (
	"github.com/WilliamDann/adachess/game"
)

// Perft function
// https://www.chessprogramming.org/Perft
func Perft(position *game.Position, depth int) *PerftResults {
	var results PerftResults

	if depth == 0 {
		return NewPerftResultsNodes(1)
	}

	moves := position.LegalMoves()

	for _, move := range moves {
		position.Move(move)

		if position.InCheck(position.GetFen().ActiveColor) {
			results.Checks += 1
		}
		if position.Checkmate() {
			results.Checkmates += 1
		}
		if move.Castle {
			results.Castles += 1
		}
		if move.Capture {
			results.Captures += 1
		}
		if move.Enpassant {
			results.Enpassant += 1
		}
		if move.Promote {
			results.Promos += 1
		}
		results.Add(Perft(position, depth-1))

		position.Unmove()

	}
	return &results
}
