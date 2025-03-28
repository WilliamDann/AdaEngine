package game

import "strings"

type Move struct {
	Piece Piece

	From Coord
	To   Coord

	Capture       bool
	CaptureTarget *Piece
	Enpassant     bool

	Promote       bool
	PromoteTarget *Piece

	Castle bool
	Side   *Side
}

func (m Move) String() string {
	var sb strings.Builder

	if m.Castle {
		if *m.Side == Kingside {
			return "0-0"
		}
		return "0-0-0"
	}

	if m.Piece.Type != Pawn {
		sb.WriteString(strings.ToUpper(m.Piece.String()))
	}

	if m.Capture {
		sb.WriteString("x")
	}

	sb.WriteString(m.To.String())

	if m.Promote {
		sb.WriteString("=")
		sb.WriteString(m.PromoteTarget.String())
	}

	return sb.String()
}

type moveBuilder struct {
	move Move
}

func MoveBuilder() *moveBuilder {
	var move Move
	return &moveBuilder{move}
}

func (m *moveBuilder) Piece(value Piece) *moveBuilder {
	m.move.Piece = value
	return m
}

func (m *moveBuilder) Enpassant(value bool) *moveBuilder {
	m.move.Enpassant = value
	return m
}

func (m *moveBuilder) From(coord Coord) *moveBuilder {
	m.move.From = coord
	return m
}

func (m *moveBuilder) To(coord Coord) *moveBuilder {
	m.move.To = coord
	return m
}

func (m *moveBuilder) Capture(value bool, target *Piece) *moveBuilder {
	m.move.Capture = value
	m.move.CaptureTarget = target
	return m
}

func (m *moveBuilder) Promote(value bool, piece *Piece) *moveBuilder {
	m.move.Promote = value
	m.move.PromoteTarget = piece
	return m
}

func (m *moveBuilder) Castle(value bool, side *Side) *moveBuilder {
	m.move.Castle = value
	m.move.Side = side
	return m
}

func (m moveBuilder) Move() Move {
	return m.move
}
