package core

import (
	"fmt"
	"strings"
)

type Chessboard struct {
	pieces [15]Bitboard
	white Bitboard
	black Bitboard
}

func NewChessboard() *Chessboard {
	return &Chessboard{}
}

func (b *Chessboard) Set(sq Square, piece Piece) {
	b.pieces[piece] = b.pieces[piece].Set(sq)
	if piece.Color() == White {
		b.white = b.white.Set(sq)
	} else {
		b.black = b.black.Set(sq)
	}
}
func (b *Chessboard) Clear(sq Square) Piece {
	for piece, bb := range b.pieces {
		b.pieces[piece] = b.pieces[piece].Clear(sq)
		b.white = b.white.Clear(sq)
		b.black = b.black.Clear(sq)

		if bb.Check(sq) {
			return Piece(piece)
		}
	}
	return None
}
func (b *Chessboard) Check(sq Square) Piece {
	if !b.HasPiece(sq) {
		return None
	}

	for piece, bb := range b.pieces {
		if bb.Check(sq) {
			return Piece(piece)
		}
	}
	return None
}


func (b *Chessboard) Pieces(piece Piece) Bitboard {
	return b.pieces[piece]
}
func (b *Chessboard) ColorPieces(color Color) Bitboard {
	if color == White {
		return b.white
	}
	return b.black
}
func (b *Chessboard) Occupied() Bitboard {
	return b.white.Union(b.black)
}


func (b *Chessboard) HasPiece(sq Square) bool {
	bb := b.white.Union(b.black)
	return bb.Check(sq)
}
func (b *Chessboard) HasColorPiece(sq Square, color Color) bool {
	if color == White {
		return b.white.Check(sq)
	}
	return b.black.Check(sq)
}


func (b *Chessboard) String() string {
	var sb strings.Builder
	for rank := 7; rank >= 0; rank-- {
		sb.WriteString(fmt.Sprintf("  %d |", rank+1))
		for file := 0; file < 8; file++ {
			sq := Square(rank*8 + file)
			sb.WriteString(" " + b.Check(sq).String())
		}
		sb.WriteString("\n")
	}
	sb.WriteString("      a b c d e f g h\n")
	return sb.String()
}
