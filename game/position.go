package game

type Position struct {
	board *Board
	fen   Fen
}

func (p Position) GetBoard() *Board {
	return p.board
}

func (p Position) GetFen() Fen {
	return p.fen
}

func (p Position) LegalMoves() []Move {
	mg := NewMoveGenerator(p, StandardRuleSet)
	return mg.Generate()
}

func (p Position) String() string {
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
