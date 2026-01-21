package main

import (
	"github.com/WilliamDann/AdaEngine/internal/chess"
	"fmt"
)

func main() {
	board := chess.NewBoard()
	board.LoadFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	println(board.String())


	for move := range chess.LegalMoves(board) {
    fmt.Printf("Move: %s -> %s\n", chess.SAN(move.From), chess.SAN(move.To))
}
}
