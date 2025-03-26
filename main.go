package main

import (
	"fmt"

	"github.com/WilliamDann/adachess/game"
)

func main() {
	pos := game.NewPosition(game.ItalianPosition)

	fmt.Println(pos.GetBoard())
	fmt.Println(pos.LegalMoves())
	fmt.Println(pos.Checkmate())
}
