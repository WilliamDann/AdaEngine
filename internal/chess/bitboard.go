package chess

import (
	"fmt"
	"math/bits"
	"strings"
)

// storage of binary board info
type Bitboard uint64

// manipulation operations
func (b Bitboard) Set(square int) Bitboard {
	return b | (1 << square)
}
func (b Bitboard) Clear(square int) Bitboard {
	return b &^ (1 << square)
}
func (b Bitboard) Flip(square int) Bitboard {
	return b ^ (1 << square)
}
func (b Bitboard) Check(square int) bool {
	return (b & (1 << square)) != 0
}

func (b Bitboard) Union(other Bitboard) Bitboard {
	return b | other
}
func (b Bitboard) Difference(other Bitboard) Bitboard {
	return b &^ other
}

// for iteration - you pop LSB for getting the next peice.
//   while b != { b, square := b.PopLSB() } // for example
func (b Bitboard) LSB() int {
	if b == 0 {
		return -1
	}
	return bits.TrailingZeros64(uint64(b))
}
func (b Bitboard) PopLSB() (Bitboard, int) {
	sq := b.LSB()
	return b & (b - 1), sq
}

func (b Bitboard) String() string {
	var result strings.Builder
	
	// Iterate from rank 8 down to rank 1 (top to bottom visually)
	for rank := 7; rank >= 0; rank-- {
		result.WriteString(fmt.Sprintf("%d ", rank+1))
		
		for file := 0; file < 8; file++ {
			square := rank*8 + file
			if b.Check(square) {
				result.WriteString("1 ")
			} else {
				result.WriteString(". ")
			}
		}
		
		result.WriteString("\n")
	}
	
	// Board footer
	result.WriteString("  a b c d e f g h\n")
	
	return result.String()
}
