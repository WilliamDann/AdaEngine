package game

import (
	"strconv"
	"strings"
	"unicode"
)

const (
	EmptyPosition    = "8/8/8/8/8/8/8/8 w - - 0 1"
	StartingPosition = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	ItalianPosition  = "r1bqkbnr/pppp1ppp/2n5/4p3/2B1P3/5N2/PPPP1PPP/RNBQK2R b KQkq - 3 3"
)

type Fen struct {
	PieceData       string
	ActiveColor     Color
	CastlingRights  string
	EnPassantSquare *Coord
	HalfmoveClock   int
	FullmoveClock   int
}

func (f Fen) GetBoard() *Board {
	board := NewBoard()
	rows := strings.Split(f.PieceData, "/")

	x := 7
	y := 7
	for _, row := range rows {
		for _, piece := range row {
			if unicode.IsDigit(piece) {
				x -= int(piece - '0')
			} else {
				piece := *NewPieceFromChar(piece)
				coord := *NewCoord(x, y)

				board.Set(piece, coord)
				x--
			}
		}
		x = 7
		y--
	}

	return board
}

func NewFen(fen string) *Fen {
	var obj Fen

	segments := strings.Split(fen, " ")

	// 1. Piece data
	obj.PieceData = segments[0]

	// 2. active color
	if segments[1] == "w" {
		obj.ActiveColor = White
	} else {
		obj.ActiveColor = Black
	}

	// 3. Castling rights
	obj.CastlingRights = segments[2]

	// 4. ep target square
	if segments[3] != "-" {
		obj.EnPassantSquare = NewCoordSan(segments[3])
	}

	// 5. halfmove clock
	obj.HalfmoveClock, _ = strconv.Atoi(segments[4])

	// 6. fullmove clock
	obj.FullmoveClock, _ = strconv.Atoi(segments[5])

	return &obj
}

func (f Fen) String() string {
	var sb strings.Builder

	// 1. Piece Placement Data
	sb.WriteString(f.PieceData)
	sb.WriteRune(' ')

	// 2. Active color
	if f.ActiveColor {
		sb.WriteString("w ")
	} else {
		sb.WriteString("b ")
	}

	// 3. Castling Rights
	sb.WriteString(f.CastlingRights)
	sb.WriteRune(' ')

	// 4. En passant target square
	if f.EnPassantSquare == nil {
		sb.WriteString("- ")
	} else {
		sb.WriteString(f.EnPassantSquare.String())
		sb.WriteRune(' ')
	}

	// 5. Halfmove clock (50 move rule)
	sb.WriteString(strconv.Itoa(f.HalfmoveClock))
	sb.WriteRune(' ')

	// 6. Fullmove clock
	sb.WriteString(strconv.Itoa(f.FullmoveClock))
	sb.WriteRune(' ')

	return sb.String()
}
