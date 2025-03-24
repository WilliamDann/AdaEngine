package main

import (
	"fmt"

	"github.com/WilliamDann/adachess/game"
)

func main() {
	pos := game.NewPosition(game.ItalianPosition)
	mg := game.NewMoveGenerator(*pos)

	fmt.Println(mg.Generate())
	fmt.Println(pos.GetBoard())
}
