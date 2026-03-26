package position

import "strings"

// stores ability to castle in different directions
type CastlingRights uint8
const (
	WhiteKingside CastlingRights = 1 << iota
	WhiteQueenside
	BlackKingside
	BlackQueenside

	NoCastling	CastlingRights = 0
	AllCastling CastlingRights = WhiteKingside | WhiteQueenside | BlackKingside | BlackQueenside
)

func (c CastlingRights) String() string {
	if c == NoCastling {
		return "-"
	}
	var sb strings.Builder
	if c&WhiteKingside != 0 {
		sb.WriteRune('K')
	}
	if c&WhiteQueenside != 0 {
		sb.WriteRune('Q')
	}
	if c&BlackKingside != 0 {
		sb.WriteRune('k')
	}
	if c&BlackQueenside != 0 {
		sb.WriteRune('q')
	}
	return sb.String()
}
