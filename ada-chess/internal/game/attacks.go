package game

import (
	"github.com/WilliamDann/AdaEngine/ada-chess/internal/board"
	"github.com/WilliamDann/AdaEngine/ada-chess/internal/moves"
)

// Attackers returns a bitboard of all pieces of the given color that attack sq.
// Uses reverse lookup: look up attacks FROM sq for each piece type, then
// intersect with enemy pieces of that type.
func (pos *Position) Attackers(sq board.Square, by board.Color) board.Bitboard {
	occupied := pos.Board.Occupied()

	knights := pos.Board.Pieces(board.NewPiece(board.Knight, by))
	bishops := pos.Board.Pieces(board.NewPiece(board.Bishop, by))
	rooks := pos.Board.Pieces(board.NewPiece(board.Rook, by))
	queens := pos.Board.Pieces(board.NewPiece(board.Queen, by))
	king := pos.Board.Pieces(board.NewPiece(board.King, by))
	pawns := pos.Board.Pieces(board.NewPiece(board.Pawn, by))

	attackers := moves.KnightMoves(sq).Intersection(knights)
	attackers = attackers.Union(moves.BishopMoves(sq, occupied).Intersection(bishops.Union(queens)))
	attackers = attackers.Union(moves.RookMoves(sq, occupied).Intersection(rooks.Union(queens)))
	attackers = attackers.Union(moves.KingMoves(sq).Intersection(king))
	attackers = attackers.Union(moves.PawnAttacks(sq, by.Flip()).Intersection(pawns))

	return attackers
}

// IsAttacked returns true if sq is attacked by any piece of the given color.
func (pos *Position) IsAttacked(sq board.Square, by board.Color) bool {
	return !pos.Attackers(sq, by).Empty()
}

// InCheck returns true if the active color's king is in check.
func (pos *Position) InCheck() bool {
	kingBB := pos.Board.Pieces(board.NewPiece(board.King, pos.ActiveColor))
	for sq := range kingBB.Squares() {
		return pos.IsAttacked(board.Square(sq), pos.ActiveColor.Flip())
	}
	return false
}
