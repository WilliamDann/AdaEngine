package moves

import "github.com/WilliamDann/AdaEngine/ada-chess/internal/board"

var (
	NNE = Direction{2, 1}
	ENE = Direction{1, 2}
	ESE = Direction{-1, 2}
	SSE = Direction{-2, 1}
	SSW = Direction{-2, -1}
	WSW = Direction{-1, -2}
	WNW = Direction{1, -2}
	NNW = Direction{2, -1}

	knightOffsets = [8]Direction{NNE, ENE, ESE, SSE, SSW, WSW, WNW, NNW}
	knightTable   = [64]board.Bitboard{}
)

func init() {
	for sq := board.Square(0); sq < 64; sq++ {
		rank := sq.Rank()
		file := sq.File()
		var bb board.Bitboard
		for _, d := range knightOffsets {
			r := rank + d.Rank
			f := file + d.File
			if r >= 0 && r < 8 && f >= 0 && f < 8 {
				bb = bb.Set(board.NewSquare(r, f))
			}
		}
		knightTable[sq] = bb
	}
}

func KnightMoves(start board.Square) board.Bitboard {
	return knightTable[start]
}


