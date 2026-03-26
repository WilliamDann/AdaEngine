package game

import (
	"github.com/WilliamDann/AdaEngine/ada-chess/internal/board"
	"github.com/WilliamDann/AdaEngine/ada-chess/internal/moves"
)

func (pos *Position) KnightMoves(sq board.Square) board.Bitboard {
	friendly := pos.Board.ColorPieces(pos.ActiveColor)
	return moves.KnightMoves(sq).Subtract(friendly)
}

func (pos *Position) RookMoves(sq board.Square) board.Bitboard {
	friendly := pos.Board.ColorPieces(pos.ActiveColor)
	occupied := pos.Board.Occupied()
	return moves.RookMoves(sq, occupied).Subtract(friendly)
}

func (pos *Position) BishopMoves(sq board.Square) board.Bitboard {
	friendly := pos.Board.ColorPieces(pos.ActiveColor)
	occupied := pos.Board.Occupied()
	return moves.BishopMoves(sq, occupied).Subtract(friendly)
}

func (pos *Position) QueenMoves(sq board.Square) board.Bitboard {
	friendly := pos.Board.ColorPieces(pos.ActiveColor)
	occupied := pos.Board.Occupied()
	return moves.QueenMoves(sq, occupied).Subtract(friendly)
}

func (pos *Position) KingMoves(sq board.Square) board.Bitboard {
	friendly := pos.Board.ColorPieces(pos.ActiveColor)
	return moves.KingMoves(sq).Subtract(friendly)
}

func (pos *Position) PawnMoves(sq board.Square) board.Bitboard {
	if pos.ActiveColor == board.White {
		return pos.whitePawnMoves(sq)
	}
	return pos.blackPawnMoves(sq)
}

func (pos *Position) whitePawnMoves(sq board.Square) board.Bitboard {
	occupied := pos.Board.Occupied()
	enemy := pos.Board.ColorPieces(board.Black)
	var result board.Bitboard

	// single push
	single := board.Bitboard(1 << (sq + 8)).Subtract(occupied)
	result = result.Union(single)

	// double push from rank 2
	if sq.Rank() == 1 && !single.Empty() {
		double := board.Bitboard(1 << (sq + 16)).Subtract(occupied)
		result = result.Union(double)
	}

	// captures
	if sq.File() < 7 {
		result = result.Union(board.Bitboard(1 << (sq + 9)).Intersection(enemy))
	}
	if sq.File() > 0 {
		result = result.Union(board.Bitboard(1 << (sq + 7)).Intersection(enemy))
	}

	// en passant
	if pos.EnPassant.Valid() {
		if sq.File() < 7 && board.Square(sq+9) == pos.EnPassant {
			result = result.Set(pos.EnPassant)
		}
		if sq.File() > 0 && board.Square(sq+7) == pos.EnPassant {
			result = result.Set(pos.EnPassant)
		}
	}

	return result
}

func (pos *Position) blackPawnMoves(sq board.Square) board.Bitboard {
	occupied := pos.Board.Occupied()
	enemy := pos.Board.ColorPieces(board.White)
	var result board.Bitboard

	// single push
	single := board.Bitboard(1 << (sq - 8)).Subtract(occupied)
	result = result.Union(single)

	// double push from rank 7
	if sq.Rank() == 6 && !single.Empty() {
		double := board.Bitboard(1 << (sq - 16)).Subtract(occupied)
		result = result.Union(double)
	}

	// captures
	if sq.File() > 0 {
		result = result.Union(board.Bitboard(1 << (sq - 9)).Intersection(enemy))
	}
	if sq.File() < 7 {
		result = result.Union(board.Bitboard(1 << (sq - 7)).Intersection(enemy))
	}

	// en passant
	if pos.EnPassant.Valid() {
		if sq.File() > 0 && board.Square(sq-9) == pos.EnPassant {
			result = result.Set(pos.EnPassant)
		}
		if sq.File() < 7 && board.Square(sq-7) == pos.EnPassant {
			result = result.Set(pos.EnPassant)
		}
	}

	return result
}
