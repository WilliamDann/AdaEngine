package main

import (
	"fmt"

	"github.com/WilliamDann/adachess/game"
)

func main() {
	pos := game.NewPosition("2R1R2K/8/8/7R/3k4/8/8/8 b - - 0 1")

	fmt.Println(pos.GetBoard())
	fmt.Println(pos.LegalMoves())
}
