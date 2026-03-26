// stores board information using a uint64

package core
import (
	"fmt"
	"iter"
	"math/bits"
	"strings"
)

// Bitboard maps each square to a bit in a uint64:
// example: white pawns on starting rank
//
//   8 | .  .  .  .  .  .  .  .
//   7 | .  .  .  .  .  .  .  .
//   6 | .  .  .  .  .  .  .  .
//   5 | .  .  .  .  .  .  .  .
//   4 | .  .  .  .  .  .  .  .
//   3 | .  .  .  .  .  .  .  .
//   2 | 1  1  1  1  1  1  1  1   = 0x000000000000FF00
//   1 | .  .  .  .  .  .  .  .
//       a  b  c  d  e  f  g  h
type Bitboard uint64
func NewBitboard() Bitboard {
	return Bitboard(0)
}


func (board Bitboard) Set(square Square) Bitboard {
	return board | 1 << square
}
func (board Bitboard) Clear(square Square) Bitboard {
	return board &^ (1 << square)
}
func (board Bitboard) Check(square Square) bool {
	return board&(1<<square) != 0
}

func (board Bitboard) Count() int {
	return bits.OnesCount64(uint64(board))
}


func (board Bitboard) Union(other Bitboard) Bitboard {
	return board | other
}
func (board Bitboard) Intersection(other Bitboard) Bitboard {
	return board & other
}
func (board Bitboard) Difference(other Bitboard) Bitboard {
	return board ^ other
}
func (board Bitboard) Subtract(other Bitboard) Bitboard {
	return board &^ other
}


func (board Bitboard) Invert() Bitboard {
	return board ^ 0xFFFFFFFFFFFFFFFF
}
func (board Bitboard) Empty() bool {
	return board == 0
}


// iterate over set squares using Least Signifigant Bit
func (board Bitboard) Squares() iter.Seq[int]  {
	return func(yield func(int) bool) {
		for board != 0 {
			sq := bits.TrailingZeros64(uint64(board))
			board &= board - 1
			if !yield(sq) {
				return
			}
		}
	}
}

func (board Bitboard) String() string {
	var sb strings.Builder
	for rank := 7; rank >= 0; rank-- {
		sb.WriteString(fmt.Sprintf("  %d |", rank+1))
		for file := 0; file < 8; file++ {
			sq := rank*8 + file
			if board.Check(Square(sq)) {
				sb.WriteString(" 1")
			} else {
				sb.WriteString(" .")
			}
		}
		sb.WriteString("\n")
	}
	sb.WriteString("      a b c d e f g h\n")
	return sb.String()
}
