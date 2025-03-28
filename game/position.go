package game

import (
	"math"
	"strings"
)

type Position struct {
	board *Board
	fen   Fen

	// position history for undos
	history []Fen
}

func (p *Position) GetBoard() *Board {
	return p.board
}

func (p *Position) GetFen() Fen {
	return p.fen
}

// if the king can be captured on the opponent's turn
func (p *Position) Check() bool {
	return len(CapturesPiece(ApplyRuleSet(p, CaptureOnlyRules), King)) != 0
}

// if a given player is in check
func (p *Position) InCheck(color Color) bool {
	temp := p.fen.ActiveColor
	p.fen.ActiveColor = !color

	val := p.Check()
	p.fen.ActiveColor = temp
	return val
}

// todo repetition
func (p *Position) IsDraw() bool {
	return p.fen.HalfmoveClock >= 50
}

// if there is check and the king has no escape
func (p *Position) Checkmate() bool {
	return p.InCheck(p.fen.ActiveColor) && len(p.LegalMoves()) == 0
}

// if there are no legal moves and there is no check
func (p *Position) Stalemate() bool {
	return !p.InCheck(p.fen.ActiveColor) && len(p.LegalMoves()) == 0
}

func (p *Position) Pass() {
	p.fen.ActiveColor = !p.fen.ActiveColor
}

// make a move
func (p *Position) Move(move Move) bool {
	// add current position to the history
	p.history = append(p.history, p.GetFen())

	// move the piece to it's new square
	p.board.Clear(move.From)
	p.board.Clear(move.To)

	p.board.Set(move.Piece, move.To)

	// move rook when castling
	if move.Castle {
		fromFile := "a"
		toFile := "c"
		if move.Side == &Kingside {
			fromFile = "h"
			toFile = "f"
		}

		rank := "1"
		if move.Piece.Color == Black {
			rank = "8"
		}

		// set rook for castling
		from := NewCoordSan(fromFile + rank)
		to := NewCoordSan(toFile + rank)

		p.board.Clear(*from)
		p.board.Clear(*to)

		p.board.Set(*NewPiece(move.Piece.Color, Rook), *to)
	}

	// set promote value when promoting
	if move.Promote {
		p.board.Clear(move.To)
		p.board.Set(*move.PromoteTarget, move.To)
	}

	// update fen striTng
	// 	is this a costly operation?
	//  this will be called many times per second
	//  removing it will probably be a good idea
	p.fen.PieceData = p.board.FenPieceData()

	// update fullmove clock
	if p.fen.ActiveColor {
		p.fen.FullmoveClock++
	}

	// change side to move
	p.fen.ActiveColor = !p.fen.ActiveColor

	// update castling rights
	if move.Piece.Type == King {
		if move.Piece.Color {
			p.fen.CastlingRights = strings.Replace(p.fen.CastlingRights, "K", "", 1)
			p.fen.CastlingRights = strings.Replace(p.fen.CastlingRights, "Q", "", 1)
		} else {
			p.fen.CastlingRights = strings.Replace(p.fen.CastlingRights, "k", "", 1)
			p.fen.CastlingRights = strings.Replace(p.fen.CastlingRights, "q", "", 1)
		}
	}

	if move.Piece.Type == Rook && move.From.X == 0 {
		if move.Piece.Color {
			p.fen.CastlingRights = strings.Replace(p.fen.CastlingRights, "Q", "", 1)
		} else {
			p.fen.CastlingRights = strings.Replace(p.fen.CastlingRights, "q", "", 1)
		}
	}
	if move.Piece.Type == Rook && move.From.X == 1 {
		if move.Piece.Color {
			p.fen.CastlingRights = strings.Replace(p.fen.CastlingRights, "K", "", 1)
		} else {
			p.fen.CastlingRights = strings.Replace(p.fen.CastlingRights, "k", "", 1)
		}
	}

	// en passant square
	if move.Piece.Type == Pawn && math.Abs(float64(move.From.Y)-float64(move.To.Y)) > 1 {
		if move.Piece.Color {
			p.fen.EnPassantSquare = NewCoord(move.From.X, move.From.Y-1)
		} else {
			p.fen.EnPassantSquare = NewCoord(move.From.X, move.From.Y+1)
		}
	} else {
		p.fen.EnPassantSquare = nil
	}

	// 50 move rule
	if move.Capture || move.Piece.Type == Pawn {
		p.fen.HalfmoveClock = 0
	} else {
		p.fen.HalfmoveClock++
	}

	// check if it's legal
	if p.Check() {
		// if it's not legal undo the history and the move
		p.Unmove()
		return false
	}

	return true
}

// undo a move
func (p *Position) Unmove() bool {
	if len(p.history) == 0 {
		return false
	}

	// get the last made move
	undoFen := p.history[len(p.history)-1]

	// remove that move from the history
	p.history = p.history[:len(p.history)-1]

	// reload the position
	p.fen = undoFen
	p.board = undoFen.GetBoard()

	return true
}

func (p *Position) Complete() bool {
	return p.Checkmate() || p.Stalemate() || p.IsDraw()
}

func (p *Position) PsuedolegalMoves() []Move {
	return ApplyRuleSet(p, StandardRules)
}

func (p *Position) LegalMoves() []Move {
	var moves []Move

	for _, move := range p.PsuedolegalMoves() {
		if p.Move(move) {
			p.Unmove()
			moves = append(moves, move)
		}
	}

	return moves
}

func (p *Position) String() string {
	return p.fen.String()
}

func NewEmptyPosition() *Position {
	return NewPosition(EmptyPosition)
}

func NewStartingPosition() *Position {
	return NewPosition(StartingPosition)
}

func NewPosition(fen string) *Position {
	var pos Position

	pos.fen = *NewFen(fen)
	pos.board = pos.fen.GetBoard()

	return &pos
}
