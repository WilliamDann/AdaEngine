package movegen

import (
	"github.com/WilliamDann/AdaEngine/ada-chess/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/position"
)

// LegalCaptures generates legal captures, en passant, and promotions.
// When in check it returns all legal moves (must escape check).
func LegalCaptures(pos *position.Position) core.MoveList {
	if InCheck(pos) {
		return LegalMoves(pos)
	}

	var ml core.MoveList

	color := pos.ActiveColor
	enemy := color.Flip()
	enemies := pos.Board.ColorPieces(enemy)
	occupied := pos.Board.Occupied()
	kingSq := kingSquare(pos, color)
	pins := ComputePins(pos)
	checkMask := core.Bitboard(0xFFFFFFFFFFFFFFFF)

	// King captures
	targets := KingMoves(kingSq).Intersection(enemies)
	occ := occupied.Clear(kingSq)
	for sq := range targets.Squares() {
		to := core.Square(sq)
		if !isAttackedBy(pos, to, enemy, occ) {
			ml.Add(core.NewMove(kingSq, to))
		}
	}

	// Piece captures (knight, bishop, rook, queen)
	for _, pt := range []core.PieceType{core.Knight, core.Bishop, core.Rook, core.Queen} {
		pieces := pos.Board.Pieces(core.NewPiece(pt, color))
		for sq := range pieces.Squares() {
			from := core.Square(sq)

			var t core.Bitboard
			switch pt {
			case core.Knight:
				if pins.Pinned.Check(from) {
					continue
				}
				t = KnightMoves(from)
			case core.Bishop:
				t = BishopMoves(from, occupied)
			case core.Rook:
				t = RookMoves(from, occupied)
			case core.Queen:
				t = QueenMoves(from, occupied)
			}

			t = t.Intersection(enemies).Intersection(checkMask)
			if pins.Pinned.Check(from) {
				t = t.Intersection(pins.Rays[from])
			}

			for to := range t.Squares() {
				ml.Add(core.NewMove(from, core.Square(to)))
			}
		}
	}

	// Pawn captures, en passant, and promotions (including push-promotions)
	pawns := pos.Board.Pieces(core.NewPiece(core.Pawn, color))

	promoRank := 7
	pushDir := 8
	if color == core.Black {
		promoRank = 0
		pushDir = -8
	}

	for sq := range pawns.Squares() {
		from := core.Square(sq)

		restriction := checkMask
		if pins.Pinned.Check(from) {
			restriction = restriction.Intersection(pins.Rays[from])
		}

		// Promotion push (quiet but creates material)
		to := core.Square(int(from) + pushDir)
		if to.Valid() && !occupied.Check(to) && to.Rank() == promoRank {
			if restriction.Check(to) {
				addPromotions(&ml, from, to)
			}
		}

		// Captures
		attacks := PawnAttacks(from, color).Intersection(enemies).Intersection(restriction)
		for capSq := range attacks.Squares() {
			capTo := core.Square(capSq)
			if capTo.Rank() == promoRank {
				addPromotions(&ml, from, capTo)
			} else {
				ml.Add(core.NewMove(from, capTo))
			}
		}

		// En passant
		if pos.EnPassant.Valid() && PawnAttacks(from, color).Check(pos.EnPassant) {
			capturedSq := core.Square(int(pos.EnPassant) - pushDir)
			epValid := restriction.Check(pos.EnPassant) || checkMask.Check(capturedSq)
			if pins.Pinned.Check(from) && !pins.Rays[from].Check(pos.EnPassant) {
				epValid = false
			}
			if epValid {
				epValid = isEPSafe(pos, from, pos.EnPassant, capturedSq, kingSq, enemy, occupied)
			}
			if epValid {
				ml.Add(core.NewEnPassant(from, pos.EnPassant))
			}
		}
	}

	return ml
}
