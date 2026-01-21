package chess

import "strings"

// CastlingRights stores castling availability for both sides
type CastlingRights struct {
	whiteKingside  bool
	whiteQueenside bool
	blackKingside  bool
	blackQueenside bool
}

// NewCastlingRights creates castling rights from FEN notation
// FEN format: "KQkq" where K=white kingside, Q=white queenside, 
// k=black kingside, q=black queenside, "-" means no castling available
func NewCastlingRights(fen string) CastlingRights {
	if fen == "-" {
		return CastlingRights{}
	}

	return CastlingRights{
		whiteKingside:  strings.Contains(fen, "K"),
		whiteQueenside: strings.Contains(fen, "Q"),
		blackKingside:  strings.Contains(fen, "k"),
		blackQueenside: strings.Contains(fen, "q"),
	}
}

// CanCastleKingside checks if the given color can castle kingside
func (c CastlingRights) CanCastleKingside(color Color) bool {
	if color == White {
		return c.whiteKingside
	}
	return c.blackKingside
}

// CanCastleQueenside checks if the given color can castle queenside
func (c CastlingRights) CanCastleQueenside(color Color) bool {
	if color == White {
		return c.whiteQueenside
	}
	return c.blackQueenside
}

// RemoveKingside removes kingside castling rights for the given color
func (c *CastlingRights) RemoveKingside(color Color) {
	if color == White {
		c.whiteKingside = false
	} else {
		c.blackKingside = false
	}
}

// RemoveQueenside removes queenside castling rights for the given color
func (c *CastlingRights) RemoveQueenside(color Color) {
	if color == White {
		c.whiteQueenside = false
	} else {
		c.blackQueenside = false
	}
}

// RemoveAll removes all castling rights for the given color
func (c *CastlingRights) RemoveAll(color Color) {
	c.RemoveKingside(color)
	c.RemoveQueenside(color)
}

// ToFEN converts castling rights back to FEN notation
func (c CastlingRights) ToFEN() string {
	var result strings.Builder

	if c.whiteKingside {
		result.WriteString("K")
	}
	if c.whiteQueenside {
		result.WriteString("Q")
	}
	if c.blackKingside {
		result.WriteString("k")
	}
	if c.blackQueenside {
		result.WriteString("q")
	}

	if result.Len() == 0 {
		return "-"
	}
	return result.String()
}
