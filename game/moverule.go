package game

import (
	"strings"
	"unicode"
)

// get only the captures from a set of moves
func Captures(moves []Move) []Move {
	arr := []Move{}
	for _, move := range moves {
		if move.Capture {
			arr = append(arr, move)
		}
	}
	return arr
}

func NoCaptures(moves []Move) []Move {
	arr := []Move{}
	for _, move := range moves {
		if !move.Capture {
			arr = append(arr, move)
		}
	}
	return arr
}

func CapturesPiece(moves []Move, target PieceType) []Move {
	arr := []Move{}
	for _, move := range moves {
		if move.Capture && move.CaptureTarget.Type == target {
			arr = append(arr, move)
		}
	}
	return arr
}

// moverule for a slider in a given direction
func slider(position Position, start Coord, direction Direction) []Move {
	var moves []Move

	piece := position.board.Get(start)
	cursor := start

	blocked := false
	for {
		cursor = cursor.Add(direction)

		// if we're blocked by a piece or the edge of the board
		if blocked || !cursor.Valid() {
			break
		}

		// if we're blocked by a piece
		if !position.board.IsEmpty(cursor) {
			blocked = true
			blocker := position.board.Get(cursor)

			// if we can capture the piece, we get a capture
			if blocker.Color != piece.Color {
				moves = append(moves, MoveBuilder().
					Piece(piece).
					From(start).
					To(cursor).
					Capture(true, &blocker).
					Move())
			}

			// don't generate a normal move to the blocked square
			continue
		}

		// we're unblocked, so generate a move
		moves = append(moves, MoveBuilder().
			Piece(piece).
			From(start).
			To(cursor).
			Move())
	}

	return moves
}

// move rule for a step in a given direction
func step(position Position, start Coord, direction Direction) []Move {
	piece := position.board.Get(start)
	to := start.Add(direction)

	// if we're off the board, fail
	if !to.Valid() {
		return []Move{}
	}

	// if we're blocked by a piece
	if !position.board.IsEmpty(to) {
		blocker := position.board.Get(to)

		// if we can capture, that's the move
		if blocker.Color != piece.Color {
			return []Move{MoveBuilder().
				Piece(piece).
				From(start).
				To(to).
				Capture(true, &blocker).
				Move()}
		}

		// cannot capture our own peices!
		return []Move{}
	}

	// if the square is clear we can move there
	return []Move{
		MoveBuilder().
			Piece(piece).
			From(start).
			To(to).
			Move(),
	}
}

// check if castling rights exist for a color in a direction
func checkCastlingRights(position Position, color Color, side Side) bool {
	look := 'k'
	if side == Queenside {
		look = 'q'
	}
	if color {
		look = unicode.ToUpper(look)
	}

	return strings.Contains(position.fen.CastlingRights, string(look))
}

// generate moves for castling
func castling(position Position, start Coord) []Move {
	var moves []Move

	piece := position.board.Get(start)

	// check if the king can castle kingside
	if checkCastlingRights(position, piece.Color, Kingside) {
		sq1 := start.Add(East)
		sq2 := sq1.Add(East)

		// if the two squares in the way are not clear, no castle!
		if position.board.IsEmpty(sq1) && position.board.IsEmpty(sq2) {
			moves = append(moves, MoveBuilder().
				Piece(piece).
				From(start).
				To(start.Add(East).Add(East)).
				Castle(true, &Kingside).
				Move())
		}
	}

	// check if the king can castle kingside
	if checkCastlingRights(position, piece.Color, Queenside) {
		sq1 := start.Add(West)
		sq2 := sq1.Add(West)
		sq3 := sq2.Add(West)

		// if the two squares in the way are not clear, no castle!
		if position.board.IsEmpty(sq1) && position.board.IsEmpty(sq2) && position.board.IsEmpty(sq3) {
			moves = append(moves, MoveBuilder().
				Piece(piece).
				From(start).
				To(start.Add(East).Add(East)).
				Castle(true, &Queenside).
				Move())
		}
	}

	return moves
}

// this is distinct from step because of en passant
func pawnCap(position Position, start Coord, direction Direction) []Move {
	piece := position.board.Get(start)
	to := start.Add(direction)

	// if we're off the board, fail
	if !to.Valid() {
		return []Move{}
	}

	// if we're blocked by a piece, or the en passant square is our target
	ep := (position.fen.EnPassantSquare != nil && position.fen.EnPassantSquare.Equ(to))
	if !position.board.IsEmpty(to) || ep {
		blocker := position.board.Get(to)

		// if we can capture, that's the move
		if blocker.Color != piece.Color || ep {
			return []Move{MoveBuilder().
				Piece(piece).
				From(start).
				To(to).
				Capture(true, &blocker).
				Enpassant(ep).
				Move()}
		}

		// cannot capture our own peices!
		return []Move{}
	}

	// if the square is clear we can move there
	return []Move{
		MoveBuilder().
			Piece(piece).
			From(start).
			To(to).
			Move(),
	}
}

// generate pawn moves
func pawn(position Position, start Coord) []Move {
	piece := position.board.Get(start)
	up := Coord{1, 1}
	origin := 1
	if piece.Color == Black {
		up.Y = -1
		origin = 6
	}

	var moves []Move

	// single step up
	moves = append(moves, NoCaptures(step(position, start, North.Mul(up)))...)

	// 2 steps up
	if start.Y == origin && len(moves) != 0 {
		moves = append(moves, NoCaptures(step(position, start, North.Add(North).Mul(up)))...)
	}

	// pawn captures
	moves = append(moves, Captures(pawnCap(position, start, Northwest.Mul(up)))...)
	moves = append(moves, Captures(pawnCap(position, start, Northwest.Mul(up)))...)

	return moves
}
