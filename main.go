package main

import (
	"fmt"

	"github.com/WilliamDann/adachess/game"
)

func main() {
	pos := game.NewPosition(game.ItalianPosition)
	// pos.GetBoard().Set(*game.NewPiece(game.Black, game.Pawn), *game.NewCoord(2, 3))

	mg := game.NewMoveGenerator(*pos)

	for i := 0; i < 1_000_000; i++ {
		mg.Generate()
	}

	// for _, move := range moves {
	// 	pos.GetBoard().Set(game.Piece{'x', false}, move.To)
	// }

	fmt.Println(mg.Generate())
	fmt.Println(pos.GetBoard())
}
