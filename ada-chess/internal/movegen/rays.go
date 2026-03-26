package movegen

import (
	"github.com/WilliamDann/AdaEngine/ada-chess/internal/core"
)

// direction for a ray on the board
type Direction struct {
	Rank, File int
}

// direction definitions
var (
	North = Direction{1, 0}
	South = Direction{-1, 0}
	East  = Direction{0, 1}
	West  = Direction{0, -1}
	NE    = Direction{1, 1}
	NW    = Direction{1, -1}
	SE    = Direction{-1, 1}
	SW    = Direction{-1, -1}

	rank1 = boardRay(core.NewSquare(0, 0), East).Set(core.NewSquare(0, 0))
	rank8 = boardRay(core.NewSquare(7, 0), East).Set(core.NewSquare(7, 0))
	fileA = boardRay(core.NewSquare(0, 0), North).Set(core.NewSquare(0, 0))
	fileH = boardRay(core.NewSquare(0, 7), North).Set(core.NewSquare(0, 7))
)

// calculate a directional ray on the chessboard
func boardRay(start core.Square, direction Direction) core.Bitboard {
	rank := start.Rank() + direction.Rank
	file := start.File() + direction.File
	bits := core.NewBitboard()

	for rank >= 0 && rank <= 7 && file >= 0 && file <= 7 {
		bits = bits.Set(core.NewSquare(rank, file))
		rank += direction.Rank
		file += direction.File
	}

	return bits
}

// boardRay that stops when blocked
//   this is the slow move generation that magic bitboards replaces
//   exists for generation of magic numbers
func boardRayWithBlockers(start core.Square, direction Direction, blockers core.Bitboard) core.Bitboard {
	rank := start.Rank() + direction.Rank
	file := start.File() + direction.File
	bits := core.NewBitboard()

	for rank >= 0 && rank <= 7 && file >= 0 && file <= 7 {
		square := core.NewSquare(rank, file)
		bits = bits.Set(square)
		if blockers.Check(square) {
			break
		}
		rank += direction.Rank
		file += direction.File
	}

	return bits
}

// gets a mask of relevant squares for a piece on a given square
func rookMask(start core.Square) core.Bitboard {
	return boardRay(start, North).Subtract(rank8).
	Union(boardRay(start, South).Subtract(rank1)).
	Union(boardRay(start, East).Subtract(fileH)).
	Union(boardRay(start, West).Subtract(fileA))
}
func bishopMask(start core.Square) core.Bitboard {
	return boardRay(start, NE).Subtract(rank8.Union(fileH)).
	Union(boardRay(start, NW).Subtract(rank8.Union(fileA))).
	Union(boardRay(start, SE).Subtract(rank1.Union(fileH))).
	Union(boardRay(start, SW).Subtract(rank1.Union(fileA)))
}

// finds all the squares a sliding piece can see given a starting square and a set of blockers
func rookAttacks(start core.Square, blockers core.Bitboard) core.Bitboard {
	return boardRayWithBlockers(start, North, blockers).
	Union(boardRayWithBlockers(start, South, blockers)).
	Union(boardRayWithBlockers(start, East, blockers)).
	Union(boardRayWithBlockers(start, West, blockers))
}
func bishopAttacks(start core.Square, blockers core.Bitboard) core.Bitboard {
	return boardRayWithBlockers(start, NE, blockers).
	Union(boardRayWithBlockers(start, NW, blockers)).
	Union(boardRayWithBlockers(start, SE, blockers)).
	Union(boardRayWithBlockers(start, SW, blockers))
}
