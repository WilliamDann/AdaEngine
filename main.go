package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/WilliamDann/adachess/engine"
	"github.com/WilliamDann/adachess/game"
)

// find a given move in the legal move list
func findMove(position *game.Position, from game.Coord, to game.Coord) *game.Move {
	for _, move := range position.LegalMoves() {
		if move.From.Equ(from) && move.To.Equ(to) {
			return &move
		}
	}

	// not legal
	return nil
}

func main() {
	position := game.NewStartingPosition()
	depth := 2
	playing := game.Black

	if len(os.Args) > 1 {
		if os.Args[2] == "start" {
			position = game.NewStartingPosition()
		} else {
			position = game.NewPosition(os.Args[2])
		}
	}

	sidestr := "white"
	if !playing {
		sidestr = "black"
	}

	fmt.Println("Ada Engine depth " + strconv.Itoa(depth) + " playing " + sidestr)

	for {
		fmt.Println(position.GetBoard())
		fmt.Println(position.GetFen())

		if position.GetFen().ActiveColor == playing {
			move, eval := engine.Search(position, depth)
			// fmt.Println(position.GetBoard())
			fmt.Println(move)
			fmt.Println(eval)

			position.Move(*move)
		} else {
			var from string
			var to string

			fmt.Scan(&from, &to)

			fromSquare := game.NewCoordSan(from)
			toSquare := game.NewCoordSan(to)

			move := findMove(position, *fromSquare, *toSquare)
			if move == nil {
				fmt.Println("Invalid move " + fromSquare.String() + " " + toSquare.String())
				continue
			}

			position.Move(*move)
		}
	}
}
