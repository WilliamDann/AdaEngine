package game

import (
	"github.com/WilliamDann/AdaEngine/ada-chess/internal/board"
	"github.com/WilliamDann/AdaEngine/ada-chess/internal/moves"
)

// PinState holds precomputed pin information for the active color.
type PinState struct {
	Pinned board.Bitboard
	Rays   [64]board.Bitboard
}

// ComputePins finds all pinned pieces for the active color and their allowed
// movement rays. A piece is pinned when it is the only friendly piece between
// the king and an enemy sliding attacker on that line.
func (pos *Position) ComputePins() PinState {
	var state PinState

	kingBB := pos.Board.Pieces(board.NewPiece(board.King, pos.ActiveColor))
	var kingSq board.Square
	for sq := range kingBB.Squares() {
		kingSq = board.Square(sq)
		break
	}

	enemy := pos.ActiveColor.Flip()
	occupied := pos.Board.Occupied()
	friendly := pos.Board.ColorPieces(pos.ActiveColor)

	// Enemy sliding pieces that could pin along rank/file or diagonal
	enemyRQ := pos.Board.Pieces(board.NewPiece(board.Rook, enemy)).Union(
		pos.Board.Pieces(board.NewPiece(board.Queen, enemy)))
	enemyBQ := pos.Board.Pieces(board.NewPiece(board.Bishop, enemy)).Union(
		pos.Board.Pieces(board.NewPiece(board.Queen, enemy)))

	// Remove friendly pieces from occupancy so we can see through them to
	// find potential pinners — enemy pieces that would attack the king if
	// the friendly blockers weren't there.
	transparent := occupied.Subtract(friendly)

	// Rank/file pins
	pinners := moves.RookMoves(kingSq, transparent).Intersection(enemyRQ)
	for pSq := range pinners.Squares() {
		pinnerSq := board.Square(pSq)
		between := rookBetween(kingSq, pinnerSq)
		pinned := between.Intersection(friendly)
		if pinned.Count() == 1 {
			for s := range pinned.Squares() {
				state.Pinned = state.Pinned.Set(board.Square(s))
				state.Rays[s] = between.Set(pinnerSq)
			}
		}
	}

	// Diagonal pins
	pinners = moves.BishopMoves(kingSq, transparent).Intersection(enemyBQ)
	for pSq := range pinners.Squares() {
		pinnerSq := board.Square(pSq)
		between := bishopBetween(kingSq, pinnerSq)
		pinned := between.Intersection(friendly)
		if pinned.Count() == 1 {
			for s := range pinned.Squares() {
				state.Pinned = state.Pinned.Set(board.Square(s))
				state.Rays[s] = between.Set(pinnerSq)
			}
		}
	}

	return state
}

// rookBetween returns the squares strictly between a and b on a rank or file.
func rookBetween(a, b board.Square) board.Bitboard {
	occA := board.Bitboard(0).Set(a)
	occB := board.Bitboard(0).Set(b)
	return moves.RookMoves(a, occB).Intersection(moves.RookMoves(b, occA))
}

// bishopBetween returns the squares strictly between a and b on a diagonal.
func bishopBetween(a, b board.Square) board.Bitboard {
	occA := board.Bitboard(0).Set(a)
	occB := board.Bitboard(0).Set(b)
	return moves.BishopMoves(a, occB).Intersection(moves.BishopMoves(b, occA))
}
