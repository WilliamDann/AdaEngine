package chess

import ("unicode")

type PieceType  uint8
type Color      uint8
type Piece      uint8

// piece definitions
const (
	None    PieceType = 0b0000
	Pawn 		PieceType = 0b0001
	Knight	PieceType = 0b0010
	Bishop  PieceType = 0b0011
	Rook    PieceType = 0b0100
	Queen   PieceType = 0b0101
	King    PieceType = 0b0110
	Invalid PieceType = 0b0111

	White   Color     = 0b1000
	Black   Color     = 0b0000
)

// maps PieceType to a character
var pieceTypeToRune = map[PieceType]rune {
	Pawn : 'p',
	Knight : 'n',
	Bishop: 'b',
	Rook: 'r',
	Queen: 'q',
	King: 'k',
	None: '_',
	Invalid: 'x',
}

// peice consttructor
func NewPiece(tp PieceType, color Color) Piece {
	return Piece(uint8(color) | uint8(tp))
}

// get the PieceType of a piece code
func (code Piece) Type() PieceType {
	return PieceType(0b0111 & uint8(code))
}

// get the Color of a peiceCode
func (code Piece) Color() Color {
	return Color(0b1000 & uint8(code))
}

// get the piece as a string
func (p Piece) String() string {
	rune := pieceTypeToRune[p.Type()]
	if p.Color() == White {
		rune = unicode.ToUpper(rune)
	}
	return string(rune)
}
