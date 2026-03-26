package moves

import "github.com/WilliamDann/AdaEngine/ada-chess/internal/board"

var whitePawnAttacks [64]board.Bitboard
var blackPawnAttacks [64]board.Bitboard

func init() {
	for sq := 0; sq < 64; sq++ {
		rank := sq / 8
		file := sq % 8

		var white board.Bitboard
		if rank < 7 && file > 0 {
			white = white.Set(board.Square(sq + 7))
		}
		if rank < 7 && file < 7 {
			white = white.Set(board.Square(sq + 9))
		}
		whitePawnAttacks[sq] = white

		var black board.Bitboard
		if rank > 0 && file > 0 {
			black = black.Set(board.Square(sq - 9))
		}
		if rank > 0 && file < 7 {
			black = black.Set(board.Square(sq - 7))
		}
		blackPawnAttacks[sq] = black
	}
}

func PawnAttacks(sq board.Square, color board.Color) board.Bitboard {
	if color == board.White {
		return whitePawnAttacks[sq]
	}
	return blackPawnAttacks[sq]
}
