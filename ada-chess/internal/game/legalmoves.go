package game

import (
	"github.com/WilliamDann/AdaEngine/ada-chess/internal/board"
	"github.com/WilliamDann/AdaEngine/ada-chess/internal/moves"
)

// LegalMoves generates all legal moves for the active color.
func (pos *Position) LegalMoves() board.MoveList {
	var ml board.MoveList

	color := pos.ActiveColor
	enemy := color.Flip()
	friendly := pos.Board.ColorPieces(color)
	occupied := pos.Board.Occupied()
	kingSq := pos.kingSquare(color)

	checkers := pos.Attackers(kingSq, enemy)
	numCheckers := checkers.Count()
	pins := pos.ComputePins()

	// King moves are always candidates
	pos.genKingMoves(&ml, kingSq, enemy, friendly, occupied)

	// Double check: only king moves are legal
	if numCheckers > 1 {
		return ml
	}

	// Check mask: squares where non-king moves can go to resolve check.
	// If not in check, every square is valid.
	checkMask := board.Bitboard(0xFFFFFFFFFFFFFFFF)
	if numCheckers == 1 {
		var checkerSq board.Square
		for sq := range checkers.Squares() {
			checkerSq = board.Square(sq)
			break
		}
		// Must capture the checker or block the ray
		checkMask = board.Bitboard(0).Set(checkerSq)
		piece := pos.Board.Check(checkerSq).Type()
		if piece == board.Bishop || piece == board.Rook || piece == board.Queen {
			checkMask = checkMask.Union(between(kingSq, checkerSq))
		}
	}

	// Castling only when not in check
	if numCheckers == 0 {
		pos.genCastling(&ml, kingSq, enemy, occupied)
	}

	// Non-king piece moves
	pos.genPawnMoves(&ml, checkMask, pins, kingSq, color, enemy, friendly, occupied)
	pos.genPieceMoves(&ml, board.Knight, checkMask, pins, color, friendly, occupied)
	pos.genPieceMoves(&ml, board.Bishop, checkMask, pins, color, friendly, occupied)
	pos.genPieceMoves(&ml, board.Rook, checkMask, pins, color, friendly, occupied)
	pos.genPieceMoves(&ml, board.Queen, checkMask, pins, color, friendly, occupied)

	return ml
}

func (pos *Position) kingSquare(color board.Color) board.Square {
	bb := pos.Board.Pieces(board.NewPiece(board.King, color))
	for sq := range bb.Squares() {
		return board.Square(sq)
	}
	return board.InvalidSquare
}

// isAttackedBy checks if sq is attacked by the given color, using a custom
// occupancy for sliding piece lookups. This is needed for king moves where
// the king is removed from occupancy so sliders can X-ray through.
func (pos *Position) isAttackedBy(sq board.Square, by board.Color, occupied board.Bitboard) bool {
	if !moves.KnightMoves(sq).Intersection(pos.Board.Pieces(board.NewPiece(board.Knight, by))).Empty() {
		return true
	}
	if !moves.PawnAttacks(sq, by.Flip()).Intersection(pos.Board.Pieces(board.NewPiece(board.Pawn, by))).Empty() {
		return true
	}
	bq := pos.Board.Pieces(board.NewPiece(board.Bishop, by)).Union(pos.Board.Pieces(board.NewPiece(board.Queen, by)))
	if !moves.BishopMoves(sq, occupied).Intersection(bq).Empty() {
		return true
	}
	rq := pos.Board.Pieces(board.NewPiece(board.Rook, by)).Union(pos.Board.Pieces(board.NewPiece(board.Queen, by)))
	if !moves.RookMoves(sq, occupied).Intersection(rq).Empty() {
		return true
	}
	if !moves.KingMoves(sq).Intersection(pos.Board.Pieces(board.NewPiece(board.King, by))).Empty() {
		return true
	}
	return false
}

// genKingMoves adds legal king moves (excluding castling).
func (pos *Position) genKingMoves(ml *board.MoveList, kingSq board.Square, enemy board.Color, friendly, occupied board.Bitboard) {
	targets := moves.KingMoves(kingSq).Subtract(friendly)
	occ := occupied.Clear(kingSq) // remove king so sliders see through
	for sq := range targets.Squares() {
		to := board.Square(sq)
		if !pos.isAttackedBy(to, enemy, occ) {
			ml.Add(board.NewMove(kingSq, to))
		}
	}
}

// genCastling adds legal castling moves. Only called when not in check.
func (pos *Position) genCastling(ml *board.MoveList, kingSq board.Square, enemy board.Color, occupied board.Bitboard) {
	occ := occupied.Clear(kingSq)

	if pos.ActiveColor == board.White {
		if pos.Castling&WhiteKingside != 0 {
			f1, g1 := board.Square(5), board.Square(6)
			if !occupied.Check(f1) && !occupied.Check(g1) &&
				!pos.isAttackedBy(f1, enemy, occ) && !pos.isAttackedBy(g1, enemy, occ) {
				ml.Add(board.NewCastling(kingSq, g1))
			}
		}
		if pos.Castling&WhiteQueenside != 0 {
			b1, c1, d1 := board.Square(1), board.Square(2), board.Square(3)
			if !occupied.Check(b1) && !occupied.Check(c1) && !occupied.Check(d1) &&
				!pos.isAttackedBy(c1, enemy, occ) && !pos.isAttackedBy(d1, enemy, occ) {
				ml.Add(board.NewCastling(kingSq, c1))
			}
		}
	} else {
		if pos.Castling&BlackKingside != 0 {
			f8, g8 := board.Square(61), board.Square(62)
			if !occupied.Check(f8) && !occupied.Check(g8) &&
				!pos.isAttackedBy(f8, enemy, occ) && !pos.isAttackedBy(g8, enemy, occ) {
				ml.Add(board.NewCastling(kingSq, g8))
			}
		}
		if pos.Castling&BlackQueenside != 0 {
			b8, c8, d8 := board.Square(57), board.Square(58), board.Square(59)
			if !occupied.Check(b8) && !occupied.Check(c8) && !occupied.Check(d8) &&
				!pos.isAttackedBy(c8, enemy, occ) && !pos.isAttackedBy(d8, enemy, occ) {
				ml.Add(board.NewCastling(kingSq, c8))
			}
		}
	}
}

// genPieceMoves adds legal moves for knights, bishops, rooks, and queens.
func (pos *Position) genPieceMoves(ml *board.MoveList, pieceType board.PieceType, checkMask board.Bitboard, pins PinState, color board.Color, friendly, occupied board.Bitboard) {
	pieces := pos.Board.Pieces(board.NewPiece(pieceType, color))
	for sq := range pieces.Squares() {
		from := board.Square(sq)

		var targets board.Bitboard
		switch pieceType {
		case board.Knight:
			// A pinned knight can never move (no knight move stays on a pin ray)
			if pins.Pinned.Check(from) {
				continue
			}
			targets = moves.KnightMoves(from)
		case board.Bishop:
			targets = moves.BishopMoves(from, occupied)
		case board.Rook:
			targets = moves.RookMoves(from, occupied)
		case board.Queen:
			targets = moves.QueenMoves(from, occupied)
		}

		targets = targets.Subtract(friendly).Intersection(checkMask)
		if pins.Pinned.Check(from) {
			targets = targets.Intersection(pins.Rays[from])
		}

		for to := range targets.Squares() {
			ml.Add(board.NewMove(from, board.Square(to)))
		}
	}
}

// genPawnMoves adds legal pawn moves including double pushes, captures,
// en passant, and promotions.
func (pos *Position) genPawnMoves(ml *board.MoveList, checkMask board.Bitboard, pins PinState, kingSq board.Square, color board.Color, enemy board.Color, friendly, occupied board.Bitboard) {
	pawns := pos.Board.Pieces(board.NewPiece(board.Pawn, color))
	enemies := pos.Board.ColorPieces(enemy)

	promoRank := 7
	startRank := 1
	pushDir := 8
	if color == board.Black {
		promoRank = 0
		startRank = 6
		pushDir = -8
	}

	for sq := range pawns.Squares() {
		from := board.Square(sq)

		restriction := checkMask
		if pins.Pinned.Check(from) {
			restriction = restriction.Intersection(pins.Rays[from])
		}

		// Single push
		to := board.Square(int(from) + pushDir)
		if to.Valid() && !occupied.Check(to) {
			if restriction.Check(to) {
				if to.Rank() == promoRank {
					addPromotions(ml, from, to)
				} else {
					ml.Add(board.NewMove(from, to))
				}
			}

			// Double push from starting rank
			if from.Rank() == startRank {
				to2 := board.Square(int(from) + 2*pushDir)
				if !occupied.Check(to2) && restriction.Check(to2) {
					ml.Add(board.NewMove(from, to2))
				}
			}
		}

		// Captures (including promotions)
		attacks := moves.PawnAttacks(from, color).Intersection(enemies).Intersection(restriction)
		for capSq := range attacks.Squares() {
			capTo := board.Square(capSq)
			if capTo.Rank() == promoRank {
				addPromotions(ml, from, capTo)
			} else {
				ml.Add(board.NewMove(from, capTo))
			}
		}

		// En passant
		if pos.EnPassant.Valid() && moves.PawnAttacks(from, color).Check(pos.EnPassant) {
			// The captured pawn sits on the same rank as our pawn
			capturedSq := board.Square(int(pos.EnPassant) - pushDir)

			// En passant resolves check if we capture the checker or land on a blocking square
			epValid := restriction.Check(pos.EnPassant) || checkMask.Check(capturedSq)

			// Pin check: if pinned, the ep square must be on the pin ray
			if pins.Pinned.Check(from) && !pins.Rays[from].Check(pos.EnPassant) {
				epValid = false
			}

			// Horizontal discovery: removing both pawns from the rank may expose the king
			if epValid {
				epValid = pos.isEPSafe(from, pos.EnPassant, capturedSq, kingSq, enemy, occupied)
			}

			if epValid {
				ml.Add(board.NewEnPassant(from, pos.EnPassant))
			}
		}
	}
}

// isEPSafe checks that an en passant capture doesn't reveal a horizontal
// attack on the king. This is the one case pin detection doesn't catch:
// both the capturing and captured pawns leave the same rank.
func (pos *Position) isEPSafe(from, to, capturedSq, kingSq board.Square, enemy board.Color, occupied board.Bitboard) bool {
	if kingSq.Rank() != from.Rank() {
		return true
	}
	occ := occupied.Clear(from).Clear(capturedSq).Set(to)
	rq := pos.Board.Pieces(board.NewPiece(board.Rook, enemy)).Union(
		pos.Board.Pieces(board.NewPiece(board.Queen, enemy)))
	return moves.RookMoves(kingSq, occ).Intersection(rq).Empty()
}

func addPromotions(ml *board.MoveList, from, to board.Square) {
	ml.Add(board.NewPromotion(from, to, board.Queen))
	ml.Add(board.NewPromotion(from, to, board.Rook))
	ml.Add(board.NewPromotion(from, to, board.Bishop))
	ml.Add(board.NewPromotion(from, to, board.Knight))
}

// between returns the squares strictly between two aligned squares.
func between(a, b board.Square) board.Bitboard {
	if a.Rank() == b.Rank() || a.File() == b.File() {
		return rookBetween(a, b)
	}
	ar, af := a.Rank(), a.File()
	br, bf := b.Rank(), b.File()
	dr := ar - br
	df := af - bf
	if dr < 0 {
		dr = -dr
	}
	if df < 0 {
		df = -df
	}
	if dr == df {
		return bishopBetween(a, b)
	}
	return 0
}
