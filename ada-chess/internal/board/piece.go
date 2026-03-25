package board

// PieceCode a union of PieceType and Color
type Piece 			uint8
type PieceType 	uint8
type Color 			uint8

// code definitions
const (
	None Piece     = 0b0000
	Pawn PieceType = iota     // 001
	Knight										// 010
	Bishop										// 011
	Rook											// 100
	Queen											// 101
	King											// 110

	White Color = 0b0000
	Black Color = 0b1000
)

// 0b1101 where
//    the first bit is color (black, in this case)
//    the remaining bits are piece type (Queen in this case)
func NewPiece(pieceType PieceType, color Color) Piece {
	return Piece(pieceType) | Piece(color)
}


func (piece Piece) Type() PieceType {
	return PieceType(0b0111 & piece)
}
func (piece Piece) Color() Color {
	return Color(0b1000 & piece)
}

func (color Color) Flip() Color {
	return color ^ Black
}

func (piece Piece) String() string {
	white := [...]string{".", "P", "N", "B", "R", "Q", "K"}
	black := [...]string{".", "p", "n", "b", "r", "q", "k"}
	if piece.Color() == Black {
		return black[piece.Type()]
	}
	return white[piece.Type()]
}
