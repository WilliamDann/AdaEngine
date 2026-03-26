package movegen

import "github.com/WilliamDann/AdaEngine/ada-chess/internal/core"

var (
	kingOffsets = [8]Direction{North, NE, East, SE, South, SW, West, NW}
	kingTable   = [64]core.Bitboard{}
)

func init() {
	for sq := core.Square(0); sq < 64; sq++ {
		rank := sq.Rank()
		file := sq.File()
		var bb core.Bitboard
		for _, d := range kingOffsets {
			r := rank + d.Rank
			f := file + d.File
			if r >= 0 && r < 8 && f >= 0 && f < 8 {
				bb = bb.Set(core.NewSquare(r, f))
			}
		}
		kingTable[sq] = bb
	}
}

func KingMoves(square core.Square) core.Bitboard {
	return kingTable[square]
}
