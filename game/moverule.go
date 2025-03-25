package game

import (
	"strings"
	"unicode"
)

// moverule for a slider in a given direction
func sliderRule(direction Coord) MoveRule {
	blocked := false
	return func(piece Piece, position Position, cursor Coord) *Coord {
		// if we've been blocked, this diection is complete
		if blocked {
			blocked = false
			return nil
		}

		// take a step in the direction
		cursor = cursor.Add(direction)

		// if we're off board, direction is complete
		if !cursor.Valid() {
			blocked = false
			return nil
		}

		// if there is a piece in our way
		if !position.board.IsEmpty(cursor) {
			blocker := position.board.Get(cursor)
			blocked = true

			// if it's our piece we've reached the end of this direction
			if blocker.Color == piece.Color {
				blocked = false
				return nil
			}
		}

		return &cursor
	}
}

// moverule for a slider for a single step
func stepRule(direction Coord) MoveRule {
	fired := false
	return func(piece Piece, position Position, cursor Coord) *Coord {
		// if this diection has already been examined, we're done
		if fired {
			fired = false
			return nil
		}
		fired = true

		// check for a capture
		cursor = cursor.Add(direction)

		// if the square is invalid, do not return the mvoe
		if !cursor.Valid() {
			fired = false
			return nil
		}

		if !position.board.IsEmpty(cursor) {
			blocker := position.board.Get(cursor)

			// if it's our piece we've reached the end of this direction
			if blocker.Color == piece.Color {
				fired = false
				return nil
			}
		}

		// return the move
		return &cursor
	}
}

// moverule for single steps that MUST be a capture
func pawnCaptureStep(direction Coord) MoveRule {
	fired := false
	return func(piece Piece, position Position, cursor Coord) *Coord {
		// if this diection has already been examined, we're done
		if fired {
			fired = false
			return nil
		}
		fired = true

		// determine the direction the pawn moves in
		if !piece.Color {
			direction.Y *= -1
		}

		// check for a capture
		cursor = cursor.Add(direction)
		if !position.board.IsEmpty(cursor) || (position.fen.EnPassantSquare != nil && position.fen.EnPassantSquare == &cursor) {
			blocker := position.board.Get(cursor)

			// the only case were this is a valid move is capturing an opponent piece
			if blocker.Color != piece.Color {
				return &cursor
			}
		}

		// not a capture
		fired = false
		return nil
	}
}

// generator for pawn moving a single step
//
//	this is different from the base single step as pawns cannot capture forward
func pawnStep() MoveRule {
	fired := false
	return func(piece Piece, position Position, cursor Coord) *Coord {
		// if this diection has already been examined, we're done
		// or if the pawn is not in the 2nd rand, the move is invalid
		if fired {
			fired = false
			return nil
		}
		fired = true

		// determine the direction the pawn moves in
		direction := North
		if !piece.Color {
			direction = South
		}

		// check for a capture
		cursor = cursor.Add(direction)
		if !position.board.IsEmpty(cursor) {
			fired = false
			return nil
		}

		// not a capture
		return &cursor
	}
}

// generator for pawn moving up two steps
func pawnTwoStep() MoveRule {
	fired := false
	return func(piece Piece, position Position, cursor Coord) *Coord {
		origin := 1
		if piece.Color == Black {
			origin = 6
		}

		// if this diection has already been examined, we're done
		// or if the pawn is not in the 2nd rand, the move is invalid
		if fired || cursor.Y != origin {
			fired = false
			return nil
		}
		fired = true

		// determine the direction the pawn moves in
		direction := North
		if !piece.Color {
			direction = South
		}

		// check for a capture
		cursor = cursor.Add(direction).Add(direction)
		if !position.board.IsEmpty(cursor) {
			fired = false
			return nil
		}

		// not a capture
		return &cursor
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

// gen move for castling
func castle(side Side) MoveRule {
	fired := false
	return func(piece Piece, position Position, cursor Coord) *Coord {
		// fire only once
		if fired {
			fired = false
			return nil
		}
		fired = true

		color := position.fen.ActiveColor

		// find correct castling direction
		direction := East
		if side == Queenside {
			direction = West
		}

		// if the king does not have the right to castle, the king has been moved or the rook has been moved
		// in that case the move is not legal
		// this also handles if the rook is in place or not
		if !checkCastlingRights(position, color, side) {
			fired = false
			return nil
		}

		origin := NewCoordSan("e1")
		if color == Black {
			origin = NewCoordSan("e8")
		}
		sq1 := origin.Add(direction)
		sq2 := sq1.Add(direction)
		sq3 := sq2.Add(direction)

		// if the path is blocked, no castle
		if !position.board.IsEmpty(sq1) || !position.board.IsEmpty(sq2) {
			fired = false
			return nil
		}

		// if we're castling queenside there is an extra square
		if side == Queenside && !position.board.IsEmpty(sq3) {
			fired = false
			return nil
		}

		// castle!
		return &sq2
	}
}
