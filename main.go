package main

import (
	"fmt"

	"github.com/WilliamDann/adachess/game"
)

func main() {
	pos := game.NewPosition("r2k3r/p6p/8/8/8/8/P6P/R2K3R w KQkq - 0 1")

	fmt.Println(pos.LegalMoves())
	fmt.Println(pos.GetBoard())
}
