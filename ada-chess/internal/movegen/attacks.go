package movegen

import (
	"github.com/WilliamDann/AdaEngine/ada-chess/internal/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/internal/position"
)

// Attackers returns a bitboard of all pieces of the given color that attack sq.
// Uses reverse lookup: look up attacks FROM sq for each piece type, then
// intersect with enemy pieces of that type.
func Attackers(pos *position.Position, sq core.Square, by core.Color) core.Bitboard {
	occupied := pos.Board.Occupied()

	knights := pos.Board.Pieces(core.NewPiece(core.Knight, by))
	bishops := pos.Board.Pieces(core.NewPiece(core.Bishop, by))
	rooks := pos.Board.Pieces(core.NewPiece(core.Rook, by))
	queens := pos.Board.Pieces(core.NewPiece(core.Queen, by))
	king := pos.Board.Pieces(core.NewPiece(core.King, by))
	pawns := pos.Board.Pieces(core.NewPiece(core.Pawn, by))

	attackers := KnightMoves(sq).Intersection(knights)
	attackers = attackers.Union(BishopMoves(sq, occupied).Intersection(bishops.Union(queens)))
	attackers = attackers.Union(RookMoves(sq, occupied).Intersection(rooks.Union(queens)))
	attackers = attackers.Union(KingMoves(sq).Intersection(king))
	attackers = attackers.Union(PawnAttacks(sq, by.Flip()).Intersection(pawns))

	return attackers
}

// IsAttacked returns true if sq is attacked by any piece of the given color.
func IsAttacked(pos *position.Position, sq core.Square, by core.Color) bool {
	return !Attackers(pos, sq, by).Empty()
}

// InCheck returns true if the active color's king is in check.
func InCheck(pos *position.Position) bool {
	kingBB := pos.Board.Pieces(core.NewPiece(core.King, pos.ActiveColor))
	for sq := range kingBB.Squares() {
		return IsAttacked(pos, core.Square(sq), pos.ActiveColor.Flip())
	}
	return false
}
