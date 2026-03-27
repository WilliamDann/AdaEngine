package core

// Move is a chess move packed into 16 bits.
//
//	bits  0-5: from square
//	bits  6-11: to square
//	bits 12-13: move type
//	bits 14-15: promotion piece
type Move uint16

const (
	MoveNormal    uint16 = 0
	MovePromotion uint16 = 1
	MoveEnPassant uint16 = 2
	MoveCastling  uint16 = 3
)

const (
	fromMask   uint16 = 0x3F
	toShift           = 6
	typeShift         = 12
	promoShift        = 14
)

const NoMove Move = 0

func NewMove(from, to Square) Move {
	return Move(uint16(from) | uint16(to)<<toShift)
}

func NewPromotion(from, to Square, piece PieceType) Move {
	promo := uint16(piece - Knight)
	return Move(uint16(from) | uint16(to)<<toShift | MovePromotion<<typeShift | promo<<promoShift)
}

func NewEnPassant(from, to Square) Move {
	return Move(uint16(from) | uint16(to)<<toShift | MoveEnPassant<<typeShift)
}

func NewCastling(from, to Square) Move {
	return Move(uint16(from) | uint16(to)<<toShift | MoveCastling<<typeShift)
}

func (m Move) From() Square {
	return Square(uint16(m) & fromMask)
}

func (m Move) To() Square {
	return Square((uint16(m) >> toShift) & 0x3F)
}

func (m Move) MoveType() uint16 {
	return (uint16(m) >> typeShift) & 0x03
}

func (m Move) PromoPiece() PieceType {
	return PieceType((uint16(m)>>promoShift)&0x03) + Knight
}

func (m Move) String() string {
	s := m.From().String() + m.To().String()
	if m.MoveType() == MovePromotion {
		s += string("nbrq"[(uint16(m)>>promoShift)&0x03])
	}
	return s
}

// MoveList is a pre-allocated fixed-size move list.
// Use instead of a slice to avoid heap allocations during generation.
type MoveList struct {
	moves [256]Move
	count int
}

func (ml *MoveList) Add(m Move) {
	ml.moves[ml.count] = m
	ml.count++
}

func (ml *MoveList) Count() int {
	return ml.count
}

func (ml *MoveList) Get(i int) Move {
	return ml.moves[i]
}

func (ml *MoveList) Clear() {
	ml.count = 0
}
