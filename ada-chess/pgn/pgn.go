package pgn

import (
	"fmt"
	"strings"
	"time"

	"github.com/WilliamDann/AdaEngine/ada-chess/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/position"
)

// Game holds the metadata and move history needed to produce a PGN string.
type Game struct {
	Event  string
	Site   string
	Date   string
	White  string
	Black  string
	Result string

	// startPos is the position before any moves were made.
	startPos *position.Position
	// moves recorded during play, paired with the position before each move.
	moves    []core.Move
	positions []*position.Position
}

// NewGame creates a game starting from the given position.
func NewGame(start *position.Position) *Game {
	return &Game{
		Event:    "AdaEngine Game",
		Site:     "?",
		Date:     time.Now().Format("2006.01.02"),
		White:    "?",
		Black:    "?",
		Result:   "*",
		startPos: start,
	}
}

// AddMove records a move. pos must be the position before the move is made.
func (g *Game) AddMove(pos *position.Position, m core.Move) {
	g.positions = append(g.positions, pos)
	g.moves = append(g.moves, m)
}

// MoveCount returns the number of moves recorded.
func (g *Game) MoveCount() int {
	return len(g.moves)
}

// String returns the complete PGN text for the game.
func (g *Game) String() string {
	var sb strings.Builder

	// Seven Tag Roster
	writeTag(&sb, "Event", g.Event)
	writeTag(&sb, "Site", g.Site)
	writeTag(&sb, "Date", g.Date)
	writeTag(&sb, "Round", "?")
	writeTag(&sb, "White", g.White)
	writeTag(&sb, "Black", g.Black)
	writeTag(&sb, "Result", g.Result)
	sb.WriteString("\n")

	// Move text
	line := 0
	for i, m := range g.moves {
		pos := g.positions[i]
		san := SAN(pos, m)

		if pos.ActiveColor == core.White {
			token := fmt.Sprintf("%d. %s", pos.Fullmoves, san)
			if line > 0 && line+1+len(token) > 80 {
				sb.WriteString("\n")
				line = 0
			} else if line > 0 {
				sb.WriteString(" ")
				line++
			}
			sb.WriteString(token)
			line += len(token)
		} else {
			// Black move — if this is the first recorded move, add move number with "..."
			if i == 0 {
				token := fmt.Sprintf("%d... %s", pos.Fullmoves, san)
				sb.WriteString(token)
				line += len(token)
			} else {
				if line > 0 && line+1+len(san) > 80 {
					sb.WriteString("\n")
					line = 0
				} else if line > 0 {
					sb.WriteString(" ")
					line++
				}
				sb.WriteString(san)
				line += len(san)
			}
		}
	}

	// Result
	if line > 0 {
		sb.WriteString(" ")
	}
	sb.WriteString(g.Result)
	sb.WriteString("\n")

	return sb.String()
}

func writeTag(sb *strings.Builder, name, value string) {
	sb.WriteString(fmt.Sprintf("[%s \"%s\"]\n", name, value))
}
