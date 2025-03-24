package game

import "unicode"

type PieceType = rune
type Color = bool

const (
	Pawn   PieceType = 'p'
	Knight PieceType = 'n'
	Bishop PieceType = 'b'
	Rook   PieceType = 'r'
	Queen  PieceType = 'q'
	King   PieceType = 'k'
	None   PieceType = 0

	White Color = true
	Black Color = false
)

type Piece struct {
	Type  PieceType
	Color Color
}

func (p Piece) Is(other Piece) bool {
	return p.Type == other.Type && p.Color == other.Color
}

func (p Piece) IsNone() bool {
	return p.Type == None
}

func (p Piece) String() string {
	// if square is empty, we don't want to return the '0' character
	if p.Type == 0 {
		return "_"
	}

	if p.Color {
		return string(unicode.ToUpper(p.Type))
	}
	return string(p.Type)
}

func NewPieceFromChar(char rune) *Piece {
	color := true
	if !unicode.IsUpper(char) {
		color = false
	}
	char = unicode.ToLower(char)
	return &Piece{char, color}
}

func NewPiece(color Color, pieceType PieceType) *Piece {
	return &Piece{pieceType, color}
}

func NoPiece() *Piece {
	return &Piece{None, Black}
}
