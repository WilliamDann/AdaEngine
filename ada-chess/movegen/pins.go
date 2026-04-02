package movegen

import (
	"github.com/WilliamDann/AdaEngine/ada-chess/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/position"
)

// PinState holds precomputed pin information for the active color.
type PinState struct {
	Pinned core.Bitboard
	Rays   [64]core.Bitboard
}

// ComputePins finds all pinned pieces for the active color and their allowed
// movement rays. A piece is pinned when it is the only friendly piece between
// the king and an enemy sliding attacker on that line.
func ComputePins(pos *position.Position) PinState {
	var state PinState

	kingBB := pos.Board.Pieces(core.NewPiece(core.King, pos.ActiveColor))
	var kingSq core.Square
	for sq := range kingBB.Squares() {
		kingSq = sq
		break
	}

	enemy := pos.ActiveColor.Flip()
	occupied := pos.Board.Occupied()
	friendly := pos.Board.ColorPieces(pos.ActiveColor)

	// Enemy sliding pieces that could pin along rank/file or diagonal
	enemyRQ := pos.Board.Pieces(core.NewPiece(core.Rook, enemy)).Union(
		pos.Board.Pieces(core.NewPiece(core.Queen, enemy)))
	enemyBQ := pos.Board.Pieces(core.NewPiece(core.Bishop, enemy)).Union(
		pos.Board.Pieces(core.NewPiece(core.Queen, enemy)))

	// Remove friendly pieces from occupancy so we can see through them to
	// find potential pinners — enemy pieces that would attack the king if
	// the friendly blockers weren't there.
	transparent := occupied.Subtract(friendly)

	// Rank/file pins
	pinners := RookMoves(kingSq, transparent).Intersection(enemyRQ)
	for pSq := range pinners.Squares() {
		pinnerSq := pSq
		between := rookBetween(kingSq, pinnerSq)
		pinned := between.Intersection(friendly)
		if pinned.Count() == 1 {
			for s := range pinned.Squares() {
				state.Pinned = state.Pinned.Set(s)
				state.Rays[s] = between.Set(pinnerSq)
			}
		}
	}

	// Diagonal pins
	pinners = BishopMoves(kingSq, transparent).Intersection(enemyBQ)
	for pSq := range pinners.Squares() {
		pinnerSq := pSq
		between := bishopBetween(kingSq, pinnerSq)
		pinned := between.Intersection(friendly)
		if pinned.Count() == 1 {
			for s := range pinned.Squares() {
				state.Pinned = state.Pinned.Set(s)
				state.Rays[s] = between.Set(pinnerSq)
			}
		}
	}

	return state
}

// rookBetween returns the squares strictly between a and b on a rank or file.
func rookBetween(a, b core.Square) core.Bitboard {
	occA := core.Bitboard(0).Set(a)
	occB := core.Bitboard(0).Set(b)
	return RookMoves(a, occB).Intersection(RookMoves(b, occA))
}

// bishopBetween returns the squares strictly between a and b on a diagonal.
func bishopBetween(a, b core.Square) core.Bitboard {
	occA := core.Bitboard(0).Set(a)
	occB := core.Bitboard(0).Set(b)
	return BishopMoves(a, occB).Intersection(BishopMoves(b, occA))
}
