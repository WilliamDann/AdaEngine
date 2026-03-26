package movegen

import (
	"github.com/WilliamDann/AdaEngine/ada-chess/internal/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/internal/position"
)

// LegalMoves generates all legal moves for the active color.
func LegalMoves(pos *position.Position) core.MoveList {
	var ml core.MoveList

	color := pos.ActiveColor
	enemy := color.Flip()
	friendly := pos.Board.ColorPieces(color)
	occupied := pos.Board.Occupied()
	kingSq := kingSquare(pos, color)

	checkers := Attackers(pos, kingSq, enemy)
	numCheckers := checkers.Count()
	pins := ComputePins(pos)

	// King moves are always candidates
	genKingMoves(pos, &ml, kingSq, enemy, friendly, occupied)

	// Double check: only king moves are legal
	if numCheckers > 1 {
		return ml
	}

	// Check mask: squares where non-king moves can go to resolve check.
	// If not in check, every square is valid.
	checkMask := core.Bitboard(0xFFFFFFFFFFFFFFFF)
	if numCheckers == 1 {
		var checkerSq core.Square
		for sq := range checkers.Squares() {
			checkerSq = core.Square(sq)
			break
		}
		// Must capture the checker or block the ray
		checkMask = core.Bitboard(0).Set(checkerSq)
		piece := pos.Board.Check(checkerSq).Type()
		if piece == core.Bishop || piece == core.Rook || piece == core.Queen {
			checkMask = checkMask.Union(between(kingSq, checkerSq))
		}
	}

	// Castling only when not in check
	if numCheckers == 0 {
		genCastling(pos, &ml, kingSq, enemy, occupied)
	}

	// Non-king piece moves
	genPawnMoves(pos, &ml, checkMask, pins, kingSq, color, enemy, friendly, occupied)
	genPieceMoves(pos, &ml, core.Knight, checkMask, pins, color, friendly, occupied)
	genPieceMoves(pos, &ml, core.Bishop, checkMask, pins, color, friendly, occupied)
	genPieceMoves(pos, &ml, core.Rook, checkMask, pins, color, friendly, occupied)
	genPieceMoves(pos, &ml, core.Queen, checkMask, pins, color, friendly, occupied)

	return ml
}

func kingSquare(pos *position.Position, color core.Color) core.Square {
	bb := pos.Board.Pieces(core.NewPiece(core.King, color))
	for sq := range bb.Squares() {
		return core.Square(sq)
	}
	return core.InvalidSquare
}

// isAttackedBy checks if sq is attacked by the given color, using a custom
// occupancy for sliding piece lookups. This is needed for king moves where
// the king is removed from occupancy so sliders can X-ray through.
func isAttackedBy(pos *position.Position, sq core.Square, by core.Color, occupied core.Bitboard) bool {
	if !KnightMoves(sq).Intersection(pos.Board.Pieces(core.NewPiece(core.Knight, by))).Empty() {
		return true
	}
	if !PawnAttacks(sq, by.Flip()).Intersection(pos.Board.Pieces(core.NewPiece(core.Pawn, by))).Empty() {
		return true
	}
	bq := pos.Board.Pieces(core.NewPiece(core.Bishop, by)).Union(pos.Board.Pieces(core.NewPiece(core.Queen, by)))
	if !BishopMoves(sq, occupied).Intersection(bq).Empty() {
		return true
	}
	rq := pos.Board.Pieces(core.NewPiece(core.Rook, by)).Union(pos.Board.Pieces(core.NewPiece(core.Queen, by)))
	if !RookMoves(sq, occupied).Intersection(rq).Empty() {
		return true
	}
	if !KingMoves(sq).Intersection(pos.Board.Pieces(core.NewPiece(core.King, by))).Empty() {
		return true
	}
	return false
}

// genKingMoves adds legal king moves (excluding castling).
func genKingMoves(pos *position.Position, ml *core.MoveList, kingSq core.Square, enemy core.Color, friendly, occupied core.Bitboard) {
	targets := KingMoves(kingSq).Subtract(friendly)
	occ := occupied.Clear(kingSq) // remove king so sliders see through
	for sq := range targets.Squares() {
		to := core.Square(sq)
		if !isAttackedBy(pos, to, enemy, occ) {
			ml.Add(core.NewMove(kingSq, to))
		}
	}
}

// genCastling adds legal castling moves. Only called when not in check.
func genCastling(pos *position.Position, ml *core.MoveList, kingSq core.Square, enemy core.Color, occupied core.Bitboard) {
	occ := occupied.Clear(kingSq)

	if pos.ActiveColor == core.White {
		if pos.Castling&position.WhiteKingside != 0 {
			f1, g1 := core.Square(5), core.Square(6)
			if !occupied.Check(f1) && !occupied.Check(g1) &&
				!isAttackedBy(pos, f1, enemy, occ) && !isAttackedBy(pos, g1, enemy, occ) {
				ml.Add(core.NewCastling(kingSq, g1))
			}
		}
		if pos.Castling&position.WhiteQueenside != 0 {
			b1, c1, d1 := core.Square(1), core.Square(2), core.Square(3)
			if !occupied.Check(b1) && !occupied.Check(c1) && !occupied.Check(d1) &&
				!isAttackedBy(pos, c1, enemy, occ) && !isAttackedBy(pos, d1, enemy, occ) {
				ml.Add(core.NewCastling(kingSq, c1))
			}
		}
	} else {
		if pos.Castling&position.BlackKingside != 0 {
			f8, g8 := core.Square(61), core.Square(62)
			if !occupied.Check(f8) && !occupied.Check(g8) &&
				!isAttackedBy(pos, f8, enemy, occ) && !isAttackedBy(pos, g8, enemy, occ) {
				ml.Add(core.NewCastling(kingSq, g8))
			}
		}
		if pos.Castling&position.BlackQueenside != 0 {
			b8, c8, d8 := core.Square(57), core.Square(58), core.Square(59)
			if !occupied.Check(b8) && !occupied.Check(c8) && !occupied.Check(d8) &&
				!isAttackedBy(pos, c8, enemy, occ) && !isAttackedBy(pos, d8, enemy, occ) {
				ml.Add(core.NewCastling(kingSq, c8))
			}
		}
	}
}

// genPieceMoves adds legal moves for knights, bishops, rooks, and queens.
func genPieceMoves(pos *position.Position, ml *core.MoveList, pieceType core.PieceType, checkMask core.Bitboard, pins PinState, color core.Color, friendly, occupied core.Bitboard) {
	pieces := pos.Board.Pieces(core.NewPiece(pieceType, color))
	for sq := range pieces.Squares() {
		from := core.Square(sq)

		var targets core.Bitboard
		switch pieceType {
		case core.Knight:
			// A pinned knight can never move (no knight move stays on a pin ray)
			if pins.Pinned.Check(from) {
				continue
			}
			targets = KnightMoves(from)
		case core.Bishop:
			targets = BishopMoves(from, occupied)
		case core.Rook:
			targets = RookMoves(from, occupied)
		case core.Queen:
			targets = QueenMoves(from, occupied)
		}

		targets = targets.Subtract(friendly).Intersection(checkMask)
		if pins.Pinned.Check(from) {
			targets = targets.Intersection(pins.Rays[from])
		}

		for to := range targets.Squares() {
			ml.Add(core.NewMove(from, core.Square(to)))
		}
	}
}

// genPawnMoves adds legal pawn moves including double pushes, captures,
// en passant, and promotions.
func genPawnMoves(pos *position.Position, ml *core.MoveList, checkMask core.Bitboard, pins PinState, kingSq core.Square, color core.Color, enemy core.Color, friendly, occupied core.Bitboard) {
	pawns := pos.Board.Pieces(core.NewPiece(core.Pawn, color))
	enemies := pos.Board.ColorPieces(enemy)

	promoRank := 7
	startRank := 1
	pushDir := 8
	if color == core.Black {
		promoRank = 0
		startRank = 6
		pushDir = -8
	}

	for sq := range pawns.Squares() {
		from := core.Square(sq)

		restriction := checkMask
		if pins.Pinned.Check(from) {
			restriction = restriction.Intersection(pins.Rays[from])
		}

		// Single push
		to := core.Square(int(from) + pushDir)
		if to.Valid() && !occupied.Check(to) {
			if restriction.Check(to) {
				if to.Rank() == promoRank {
					addPromotions(ml, from, to)
				} else {
					ml.Add(core.NewMove(from, to))
				}
			}

			// Double push from starting rank
			if from.Rank() == startRank {
				to2 := core.Square(int(from) + 2*pushDir)
				if !occupied.Check(to2) && restriction.Check(to2) {
					ml.Add(core.NewMove(from, to2))
				}
			}
		}

		// Captures (including promotions)
		attacks := PawnAttacks(from, color).Intersection(enemies).Intersection(restriction)
		for capSq := range attacks.Squares() {
			capTo := core.Square(capSq)
			if capTo.Rank() == promoRank {
				addPromotions(ml, from, capTo)
			} else {
				ml.Add(core.NewMove(from, capTo))
			}
		}

		// En passant
		if pos.EnPassant.Valid() && PawnAttacks(from, color).Check(pos.EnPassant) {
			// The captured pawn sits on the same rank as our pawn
			capturedSq := core.Square(int(pos.EnPassant) - pushDir)

			// En passant resolves check if we capture the checker or land on a blocking square
			epValid := restriction.Check(pos.EnPassant) || checkMask.Check(capturedSq)

			// Pin check: if pinned, the ep square must be on the pin ray
			if pins.Pinned.Check(from) && !pins.Rays[from].Check(pos.EnPassant) {
				epValid = false
			}

			// Horizontal discovery: removing both pawns from the rank may expose the king
			if epValid {
				epValid = isEPSafe(pos, from, pos.EnPassant, capturedSq, kingSq, enemy, occupied)
			}

			if epValid {
				ml.Add(core.NewEnPassant(from, pos.EnPassant))
			}
		}
	}
}

// isEPSafe checks that an en passant capture doesn't reveal a horizontal
// attack on the king. This is the one case pin detection doesn't catch:
// both the capturing and captured pawns leave the same rank.
func isEPSafe(pos *position.Position, from, to, capturedSq, kingSq core.Square, enemy core.Color, occupied core.Bitboard) bool {
	if kingSq.Rank() != from.Rank() {
		return true
	}
	occ := occupied.Clear(from).Clear(capturedSq).Set(to)
	rq := pos.Board.Pieces(core.NewPiece(core.Rook, enemy)).Union(
		pos.Board.Pieces(core.NewPiece(core.Queen, enemy)))
	return RookMoves(kingSq, occ).Intersection(rq).Empty()
}

func addPromotions(ml *core.MoveList, from, to core.Square) {
	ml.Add(core.NewPromotion(from, to, core.Queen))
	ml.Add(core.NewPromotion(from, to, core.Rook))
	ml.Add(core.NewPromotion(from, to, core.Bishop))
	ml.Add(core.NewPromotion(from, to, core.Knight))
}

// between returns the squares strictly between two aligned squares.
func between(a, b core.Square) core.Bitboard {
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
