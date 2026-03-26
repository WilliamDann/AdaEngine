package movegen

import (
	"github.com/WilliamDann/AdaEngine/ada-chess/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/position"
)

func PseudoKnightMoves(pos *position.Position, sq core.Square) core.Bitboard {
	friendly := pos.Board.ColorPieces(pos.ActiveColor)
	return KnightMoves(sq).Subtract(friendly)
}

func PseudoRookMoves(pos *position.Position, sq core.Square) core.Bitboard {
	friendly := pos.Board.ColorPieces(pos.ActiveColor)
	occupied := pos.Board.Occupied()
	return RookMoves(sq, occupied).Subtract(friendly)
}

func PseudoBishopMoves(pos *position.Position, sq core.Square) core.Bitboard {
	friendly := pos.Board.ColorPieces(pos.ActiveColor)
	occupied := pos.Board.Occupied()
	return BishopMoves(sq, occupied).Subtract(friendly)
}

func PseudoQueenMoves(pos *position.Position, sq core.Square) core.Bitboard {
	friendly := pos.Board.ColorPieces(pos.ActiveColor)
	occupied := pos.Board.Occupied()
	return QueenMoves(sq, occupied).Subtract(friendly)
}

func PseudoKingMoves(pos *position.Position, sq core.Square) core.Bitboard {
	friendly := pos.Board.ColorPieces(pos.ActiveColor)
	return KingMoves(sq).Subtract(friendly)
}

func PseudoPawnMoves(pos *position.Position, sq core.Square) core.Bitboard {
	if pos.ActiveColor == core.White {
		return pseudoWhitePawnMoves(pos, sq)
	}
	return pseudoBlackPawnMoves(pos, sq)
}

func pseudoWhitePawnMoves(pos *position.Position, sq core.Square) core.Bitboard {
	occupied := pos.Board.Occupied()
	enemy := pos.Board.ColorPieces(core.Black)
	var result core.Bitboard

	// single push
	single := core.Bitboard(1 << (sq + 8)).Subtract(occupied)
	result = result.Union(single)

	// double push from rank 2
	if sq.Rank() == 1 && !single.Empty() {
		double := core.Bitboard(1 << (sq + 16)).Subtract(occupied)
		result = result.Union(double)
	}

	// captures
	if sq.File() < 7 {
		result = result.Union(core.Bitboard(1 << (sq + 9)).Intersection(enemy))
	}
	if sq.File() > 0 {
		result = result.Union(core.Bitboard(1 << (sq + 7)).Intersection(enemy))
	}

	// en passant
	if pos.EnPassant.Valid() {
		if sq.File() < 7 && core.Square(sq+9) == pos.EnPassant {
			result = result.Set(pos.EnPassant)
		}
		if sq.File() > 0 && core.Square(sq+7) == pos.EnPassant {
			result = result.Set(pos.EnPassant)
		}
	}

	return result
}

func pseudoBlackPawnMoves(pos *position.Position, sq core.Square) core.Bitboard {
	occupied := pos.Board.Occupied()
	enemy := pos.Board.ColorPieces(core.White)
	var result core.Bitboard

	// single push
	single := core.Bitboard(1 << (sq - 8)).Subtract(occupied)
	result = result.Union(single)

	// double push from rank 7
	if sq.Rank() == 6 && !single.Empty() {
		double := core.Bitboard(1 << (sq - 16)).Subtract(occupied)
		result = result.Union(double)
	}

	// captures
	if sq.File() > 0 {
		result = result.Union(core.Bitboard(1 << (sq - 9)).Intersection(enemy))
	}
	if sq.File() < 7 {
		result = result.Union(core.Bitboard(1 << (sq - 7)).Intersection(enemy))
	}

	// en passant
	if pos.EnPassant.Valid() {
		if sq.File() > 0 && core.Square(sq-9) == pos.EnPassant {
			result = result.Set(pos.EnPassant)
		}
		if sq.File() < 7 && core.Square(sq-7) == pos.EnPassant {
			result = result.Set(pos.EnPassant)
		}
	}

	return result
}
