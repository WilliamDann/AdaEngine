package movegen

import (
	"github.com/WilliamDann/AdaEngine/ada-chess/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/position"
)

// LegalCaptures generates all legal capture moves and promotions for the active color.
func LegalCaptures(pos *position.Position) core.MoveList {
	var ml core.MoveList

	color := pos.ActiveColor
	enemy := color.Flip()
	friendly := pos.Board.ColorPieces(color)
	enemies := pos.Board.ColorPieces(enemy)
	occupied := pos.Board.Occupied()
	kingSq := kingSquare(pos, color)

	checkers := Attackers(pos, kingSq, enemy)
	numCheckers := checkers.Count()
	pins := ComputePins(pos)

	// King captures are always candidates
	genKingCaptures(pos, &ml, kingSq, enemy, friendly, enemies, occupied)

	// Double check: only king moves are legal
	if numCheckers > 1 {
		return ml
	}

	// Check mask
	checkMask := core.Bitboard(0xFFFFFFFFFFFFFFFF)
	if numCheckers == 1 {
		var checkerSq core.Square
		for sq := range checkers.Squares() {
			checkerSq = sq
			break
		}
		checkMask = core.Bitboard(0).Set(checkerSq)
		piece := pos.Board.Check(checkerSq).Type()
		if piece == core.Bishop || piece == core.Rook || piece == core.Queen {
			checkMask = checkMask.Union(between(kingSq, checkerSq))
		}
	}

	// Non-king piece captures
	genPawnCaptures(pos, &ml, checkMask, pins, kingSq, color, enemy, friendly, enemies, occupied)
	genPieceCaptures(pos, &ml, core.Knight, checkMask, pins, color, friendly, enemies, occupied)
	genPieceCaptures(pos, &ml, core.Bishop, checkMask, pins, color, friendly, enemies, occupied)
	genPieceCaptures(pos, &ml, core.Rook, checkMask, pins, color, friendly, enemies, occupied)
	genPieceCaptures(pos, &ml, core.Queen, checkMask, pins, color, friendly, enemies, occupied)

	return ml
}

// genKingCaptures adds legal king captures only.
func genKingCaptures(pos *position.Position, ml *core.MoveList, kingSq core.Square, enemy core.Color, friendly, enemies, occupied core.Bitboard) {
	targets := KingMoves(kingSq).Intersection(enemies)
	occ := occupied.Clear(kingSq)
	for sq := range targets.Squares() {
		if !isAttackedBy(pos, sq, enemy, occ) {
			ml.Add(core.NewMove(kingSq, sq))
		}
	}
}

// genPieceCaptures adds legal captures for knights, bishops, rooks, and queens.
func genPieceCaptures(pos *position.Position, ml *core.MoveList, pieceType core.PieceType, checkMask core.Bitboard, pins PinState, color core.Color, friendly, enemies, occupied core.Bitboard) {
	pieces := pos.Board.Pieces(core.NewPiece(pieceType, color))
	for sq := range pieces.Squares() {
		var targets core.Bitboard
		switch pieceType {
		case core.Knight:
			if pins.Pinned.Check(sq) {
				continue
			}
			targets = KnightMoves(sq)
		case core.Bishop:
			targets = BishopMoves(sq, occupied)
		case core.Rook:
			targets = RookMoves(sq, occupied)
		case core.Queen:
			targets = QueenMoves(sq, occupied)
		}

		targets = targets.Subtract(friendly).Intersection(enemies).Intersection(checkMask)
		if pins.Pinned.Check(sq) {
			targets = targets.Intersection(pins.Rays[sq])
		}

		for to := range targets.Squares() {
			ml.Add(core.NewMove(sq, to))
		}
	}
}

// genPawnCaptures adds legal pawn captures, en passant, and promotions.
func genPawnCaptures(pos *position.Position, ml *core.MoveList, checkMask core.Bitboard, pins PinState, kingSq core.Square, color core.Color, enemy core.Color, friendly, enemies, occupied core.Bitboard) {
	pawns := pos.Board.Pieces(core.NewPiece(core.Pawn, color))

	promoRank := 7
	pushDir := 8
	if color == core.Black {
		promoRank = 0
		pushDir = -8
	}

	for sq := range pawns.Squares() {
		restriction := checkMask
		if pins.Pinned.Check(sq) {
			restriction = restriction.Intersection(pins.Rays[sq])
		}

		// Promotion pushes (non-capture but tactical)
		to := core.Square(int(sq) + pushDir)
		if to.Valid() && !occupied.Check(to) && to.Rank() == promoRank && restriction.Check(to) {
			addPromotions(ml, sq, to)
		}

		// Captures (including promotion captures)
		attacks := PawnAttacks(sq, color).Intersection(enemies).Intersection(restriction)
		for capSq := range attacks.Squares() {
			if capSq.Rank() == promoRank {
				addPromotions(ml, sq, capSq)
			} else {
				ml.Add(core.NewMove(sq, capSq))
			}
		}

		// En passant
		if pos.EnPassant.Valid() && PawnAttacks(sq, color).Check(pos.EnPassant) {
			capturedSq := core.Square(int(pos.EnPassant) - pushDir)

			epValid := restriction.Check(pos.EnPassant) || checkMask.Check(capturedSq)

			if pins.Pinned.Check(sq) && !pins.Rays[sq].Check(pos.EnPassant) {
				epValid = false
			}

			if epValid {
				epValid = isEPSafe(pos, sq, pos.EnPassant, capturedSq, kingSq, enemy, occupied)
			}

			if epValid {
				ml.Add(core.NewEnPassant(sq, pos.EnPassant))
			}
		}
	}
}
