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
		results.Add(Perft(position, depth-1))

		position.Unmove()

	}
	return &results
}
