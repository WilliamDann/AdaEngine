package position

import "github.com/WilliamDann/AdaEngine/ada-chess/core"

// NullMove returns a position with the side to move flipped and en passant
// cleared. No pieces move. The board is shared (not cloned) since nothing
// changes on it.
func NullMove(pos *Position) *Position {
	h := pos.Hash
	h ^= sideKey
	if pos.EnPassant.Valid() {
		h ^= epKeys[pos.EnPassant.File()]
		h ^= epKeys[8]
	}
	return &Position{
		Board:       pos.Board,
		ActiveColor: pos.ActiveColor.Flip(),
		Castling:    pos.Castling,
		EnPassant:   core.InvalidSquare,
		Halfmoves:   pos.Halfmoves,
		Fullmoves:   pos.Fullmoves,
		Hash:        h,
	}
}

// MakeMove applies a move to a position and returns the resulting position.
// The original position is not modified.
func MakeMove(pos *Position, m core.Move) *Position {
	next := &Position{
		Board:       pos.Board.Clone(),
		ActiveColor: pos.ActiveColor.Flip(),
		Castling:    pos.Castling,
		EnPassant:   core.InvalidSquare,
		Halfmoves:   pos.Halfmoves + 1,
		Fullmoves:   pos.Fullmoves,
	}
	if pos.ActiveColor == core.Black {
		next.Fullmoves++
	}

	from := m.From()
	to := m.To()
	piece := pos.Board.Check(from)
	captured := next.Board.Clear(to)

	// Begin incremental hash update
	h := pos.Hash
	h ^= sideKey
	h ^= castlingKeys[pos.Castling]
	if pos.EnPassant.Valid() {
		h ^= epKeys[pos.EnPassant.File()]
	} else {
		h ^= epKeys[8]
	}
	h ^= pieceKeys[piece][from]
	if captured != core.None {
		h ^= pieceKeys[captured][to]
	}

	// Reset halfmove clock on pawn move or capture
	if piece.Type() == core.Pawn || captured != core.None {
		next.Halfmoves = 0
	}

	switch m.MoveType() {
	case core.MoveNormal:
		next.Board.Clear(from)
		next.Board.Set(to, piece)
		h ^= pieceKeys[piece][to]

		// Double pawn push sets en passant square
		if piece.Type() == core.Pawn {
			diff := int(to) - int(from)
			if diff == 16 || diff == -16 {
				next.EnPassant = core.Square((int(from) + int(to)) / 2)
			}
		}

	case core.MovePromotion:
		next.Board.Clear(from)
		promoPiece := core.NewPiece(m.PromoPiece(), pos.ActiveColor)
		next.Board.Set(to, promoPiece)
		h ^= pieceKeys[promoPiece][to]

	case core.MoveEnPassant:
		next.Board.Clear(from)
		next.Board.Set(to, piece)
		h ^= pieceKeys[piece][to]
		// Remove the captured pawn
		dir := 8
		if pos.ActiveColor == core.Black {
			dir = -8
		}
		epCapSq := core.Square(int(to) - dir)
		epPawn := next.Board.Clear(epCapSq)
		h ^= pieceKeys[epPawn][epCapSq]
		next.Halfmoves = 0

	case core.MoveCastling:
		next.Board.Clear(from)
		next.Board.Set(to, piece)
		h ^= pieceKeys[piece][to]
		// Move the rook
		switch to {
		case core.Square(6): // white kingside
			rook := next.Board.Clear(core.Square(7))
			next.Board.Set(core.Square(5), rook)
			h ^= pieceKeys[rook][core.Square(7)]
			h ^= pieceKeys[rook][core.Square(5)]
		case core.Square(2): // white queenside
			rook := next.Board.Clear(core.Square(0))
			next.Board.Set(core.Square(3), rook)
			h ^= pieceKeys[rook][core.Square(0)]
			h ^= pieceKeys[rook][core.Square(3)]
		case core.Square(62): // black kingside
			rook := next.Board.Clear(core.Square(63))
			next.Board.Set(core.Square(61), rook)
			h ^= pieceKeys[rook][core.Square(63)]
			h ^= pieceKeys[rook][core.Square(61)]
		case core.Square(58): // black queenside
			rook := next.Board.Clear(core.Square(56))
			next.Board.Set(core.Square(59), rook)
			h ^= pieceKeys[rook][core.Square(56)]
			h ^= pieceKeys[rook][core.Square(59)]
		}
	}

	// Update castling rights
	updateCastling(next, from, to)

	// Finish hash: XOR in new castling and EP
	h ^= castlingKeys[next.Castling]
	if next.EnPassant.Valid() {
		h ^= epKeys[next.EnPassant.File()]
	} else {
		h ^= epKeys[8]
	}
	next.Hash = h

	return next
}

func updateCastling(pos *Position, from, to core.Square) {
	// King moves revoke both sides
	if from == core.Square(4) {
		pos.Castling &^= WhiteKingside | WhiteQueenside
	}
	if from == core.Square(60) {
		pos.Castling &^= BlackKingside | BlackQueenside
	}
	// Rook moves or captures revoke that side
	if from == core.Square(0) || to == core.Square(0) {
		pos.Castling &^= WhiteQueenside
	}
	if from == core.Square(7) || to == core.Square(7) {
		pos.Castling &^= WhiteKingside
	}
	if from == core.Square(56) || to == core.Square(56) {
		pos.Castling &^= BlackQueenside
	}
	if from == core.Square(63) || to == core.Square(63) {
		pos.Castling &^= BlackKingside
	}
}
