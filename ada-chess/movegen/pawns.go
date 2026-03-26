package movegen

import "github.com/WilliamDann/AdaEngine/ada-chess/core"

var whitePawnAttacks [64]core.Bitboard
var blackPawnAttacks [64]core.Bitboard

func init() {
	for sq := 0; sq < 64; sq++ {
		rank := sq / 8
		file := sq % 8

		var white core.Bitboard
		if rank < 7 && file > 0 {
			white = white.Set(core.Square(sq + 7))
		}
		if rank < 7 && file < 7 {
			white = white.Set(core.Square(sq + 9))
		}
		whitePawnAttacks[sq] = white

		var black core.Bitboard
		if rank > 0 && file > 0 {
			black = black.Set(core.Square(sq - 9))
		}
		if rank > 0 && file < 7 {
			black = black.Set(core.Square(sq - 7))
		}
		blackPawnAttacks[sq] = black
	}
}

func PawnAttacks(sq core.Square, color core.Color) core.Bitboard {
	if color == core.White {
		return whitePawnAttacks[sq]
	}
	return blackPawnAttacks[sq]
}
