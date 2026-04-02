package position

import "github.com/WilliamDann/AdaEngine/ada-chess/core"

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
		Zobrist:     pos.Zobrist,
	}
	if pos.ActiveColor == core.Black {
		next.Fullmoves++
	}

	from := m.From()
	to := m.To()
	piece := pos.Board.Check(from)
	captured := next.Board.Clear(to)

	// Reset halfmove clock on pawn move or capture
	if piece.Type() == core.Pawn || captured != core.None {
		next.Halfmoves = 0
	}

	switch m.MoveType() {
	case core.MoveNormal:
		next.Board.Clear(from)
		next.Board.Set(to, piece)

		next.Zobrist ^= pieceSquareKeys[piece][from]
		next.Zobrist ^= pieceSquareKeys[piece][to]
		if captured != core.None {
			next.Zobrist ^= pieceSquareKeys[captured][to]
		}

		// Double pawn push sets en passant square
		if piece.Type() == core.Pawn {
			diff := int(to) - int(from)
			if diff == 16 || diff == -16 {
				next.EnPassant = core.Square((int(from) + int(to)) / 2)
			}
		}

	case core.MovePromotion:
		newPiece := core.NewPiece(m.PromoPiece(), pos.ActiveColor)

		next.Board.Clear(from)
		next.Board.Set(to, newPiece)

		next.Zobrist ^= pieceSquareKeys[piece][from]
		next.Zobrist ^= pieceSquareKeys[newPiece][to]
		if captured != core.None {
			next.Zobrist ^= pieceSquareKeys[captured][to]
		}

	case core.MoveEnPassant:
		next.Board.Clear(from)
		next.Board.Set(to, piece)
		// Remove the captured pawn
		dir := 8
		if pos.ActiveColor == core.Black {
			dir = -8
		}
		next.Board.Clear(core.Square(int(to) - dir))
		next.Halfmoves = 0

		next.Zobrist ^= pieceSquareKeys[piece][from]
		next.Zobrist ^= pieceSquareKeys[piece][to]
		next.Zobrist ^= pieceSquareKeys[core.NewPiece(core.Pawn, pos.ActiveColor.Flip())][core.Square(int(to) - dir)]

	case core.MoveCastling:
		next.Board.Clear(from)
		next.Board.Set(to, piece)
		// Move the rook
		switch to {
		case core.Square(6): // white kingside
			rook := next.Board.Clear(core.Square(7))
			next.Board.Set(core.Square(5), rook)
			next.Zobrist ^= pieceSquareKeys[rook][core.Square(7)]
			next.Zobrist ^= pieceSquareKeys[rook][core.Square(5)]

		case core.Square(2): // white queenside
			rook := next.Board.Clear(core.Square(0))
			next.Board.Set(core.Square(3), rook)
			next.Zobrist ^= pieceSquareKeys[rook][core.Square(0)]
			next.Zobrist ^= pieceSquareKeys[rook][core.Square(3)]
		
		case core.Square(62): // black kingside
			rook := next.Board.Clear(core.Square(63))
			next.Board.Set(core.Square(61), rook)
			next.Zobrist ^= pieceSquareKeys[rook][core.Square(63)]
			next.Zobrist ^= pieceSquareKeys[rook][core.Square(61)]

		case core.Square(58): // black queenside
			rook := next.Board.Clear(core.Square(56))
			next.Board.Set(core.Square(59), rook)
			next.Zobrist ^= pieceSquareKeys[rook][core.Square(56)]
			next.Zobrist ^= pieceSquareKeys[rook][core.Square(59)]
		}

		next.Zobrist ^= pieceSquareKeys[piece][from]
		next.Zobrist ^= pieceSquareKeys[piece][to]
	}

	// Update castling rights
	updateCastling(next, from, to)

	// update position hash

	// always toggle side to move
	next.Zobrist ^= sideToMoveKey

	next.Zobrist ^= enPassantKey(pos)
	next.Zobrist ^= enPassantKey(next)

	next.Zobrist ^= castlingKey(pos)
	next.Zobrist ^= castlingKey(next)

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
