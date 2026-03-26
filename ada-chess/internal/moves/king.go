package moves

import "github.com/WilliamDann/AdaEngine/ada-chess/internal/board"

var (
	kingOffsets = [8]Direction{North, NE, East, SE, South, SW, West, NW}
	kingTable   = [64]board.Bitboard{}
)

func init() {
	for sq := board.Square(0); sq < 64; sq++ {
		rank := sq.Rank()
		file := sq.File()
		var bb board.Bitboard
		for _, d := range kingOffsets {
			r := rank + d.Rank
			f := file + d.File
			if r >= 0 && r < 8 && f >= 0 && f < 8 {
				bb = bb.Set(board.NewSquare(r, f))
			}
		}
		kingTable[sq] = bb
	}
}

func KingMoves(square board.Square) board.Bitboard {
	return kingTable[square]
}
