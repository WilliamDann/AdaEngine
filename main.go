package main

import (
	"github.com/WilliamDann/adachess/game"
	"github.com/WilliamDann/adachess/perft"
)

func main() {
	pos := game.NewStartingPosition()
	perft.Perft(pos, 2)
}
