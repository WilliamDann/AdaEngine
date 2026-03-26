package position

import (
	"fmt"
	"strings"

	"github.com/WilliamDann/AdaEngine/ada-chess/core"
)

// state for an active chesss game
type Position struct {
	Board       *core.Chessboard
	ActiveColor core.Color
	Castling    CastlingRights
	EnPassant   core.Square
	Halfmoves   int
	Fullmoves   int
}


func NewPosition() *Position {
	return &Position{}
}

func (pos *Position) String() string {
	var sb strings.Builder

	// board
	if pos.Board != nil {
		sb.WriteString(pos.Board.String())
	}

	// game state
	color := "w"
	if pos.ActiveColor == core.Black {
		color = "b"
	}

	ep := "-"
	if pos.EnPassant.Valid() {
		ep = pos.EnPassant.String()
	}

	sb.WriteString(fmt.Sprintf("\n  Turn: %s  Castling: %s  En Passant: %s\n", color, pos.Castling, ep))
	sb.WriteString(fmt.Sprintf("  Halfmoves: %d  Fullmoves: %d\n", pos.Halfmoves, pos.Fullmoves))

	return sb.String()
}
