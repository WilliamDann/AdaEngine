package chess

import (
	"strings"
	"fmt"
	"unicode"
)

func (b *Board) LoadFEN(fen string) error {
	// Split FEN into components
	parts := strings.Fields(fen)
	if len(parts) < 1 {
		return fmt.Errorf("invalid FEN: empty string")
	}

	b.Reset()

	// Parse piece placement (first part of FEN)
	ranks := strings.Split(parts[0], "/")
	if len(ranks) != 8 {
		return fmt.Errorf("invalid FEN: expected 8 ranks, got %d", len(ranks))
	}

	// Iterate through ranks from 8 to 1
	for rankIdx, rankStr := range ranks {
		rank := 7 - rankIdx // Convert to 0-indexed rank (7 = rank 8, 0 = rank 1)
		file := 0

		for _, char := range rankStr {
			if file >= 8 {
				return fmt.Errorf("invalid FEN: too many files in rank %d", rank+1)
			}

			// Check if it's a digit (empty squares)
			if char >= '1' && char <= '8' {
				file += int(char - '0')
				continue
			}

			// Parse piece
			piece, err := parseFENPiece(char)
			if err != nil {
				return fmt.Errorf("invalid FEN: %v", err)
			}

			// Place piece on board
			square := rank*8 + file
			b.Place(piece, square)
			file++
		}

		if file != 8 {
			return fmt.Errorf("invalid FEN: rank %d has %d files, expected 8", rank+1, file)
		}
	}

	// Parse active color (second part - optional)
	if len(parts) >= 2 {
		switch parts[1] {
		case "w":
			b.activeColor = White
		case "b":
			b.activeColor = Black
		default:
			return fmt.Errorf("invalid FEN: invalid active color '%s'", parts[1])
		}
	}

	// Parse castling rights (third part - optional)
	if len(parts) >= 3 {
		b.castling = NewCastlingRights(parts[2])
	}

	// Parse en passant target square (fourth part - optional)
	if len(parts) >= 4 {
		if parts[3] != "-" {
			square, err := parseSquare(parts[3])
			if err != nil {
				return fmt.Errorf("invalid FEN: invalid en passant square '%s': %v", parts[3], err)
			}
			b.enPassant = square
		} else {
			b.enPassant = -1
		}
	}

	// Parse halfmove clock (fifth part - optional)
	if len(parts) >= 5 {
		var halfmove int
		_, err := fmt.Sscanf(parts[4], "%d", &halfmove)
		if err != nil {
			return fmt.Errorf("invalid FEN: invalid halfmove clock '%s'", parts[4])
		}
		b.halfmove = halfmove
	}

	// Parse fullmove number (sixth part - optional)
	if len(parts) >= 6 {
		var fullmove int
		_, err := fmt.Sscanf(parts[5], "%d", &fullmove)
		if err != nil {
			return fmt.Errorf("invalid FEN: invalid fullmove number '%s'", parts[5])
		}
		b.fullmove = fullmove
	}

	return nil
}

// Helper function to parse a FEN piece character
func parseFENPiece(char rune) (Piece, error) {
	var pieceType PieceType
	var color Color

	// Determine color
	if unicode.IsUpper(char) {
		color = White
		char = unicode.ToLower(char)
	} else {
		color = Black
	}

	// Determine piece type
	switch char {
	case 'p':
		pieceType = Pawn
	case 'n':
		pieceType = Knight
	case 'b':
		pieceType = Bishop
	case 'r':
		pieceType = Rook
	case 'q':
		pieceType = Queen
	case 'k':
		pieceType = King
	default:
		return 0, fmt.Errorf("unknown piece character: %c", char)
	}

	return NewPiece(pieceType, color), nil
}

// Helper function to parse algebraic notation square (e.g., "e4" -> 28)
func parseSquare(s string) (int, error) {
	if len(s) != 2 {
		return -1, fmt.Errorf("square must be 2 characters")
	}

	file := int(s[0] - 'a')
	rank := int(s[1] - '1')

	if file < 0 || file > 7 || rank < 0 || rank > 7 {
		return -1, fmt.Errorf("square out of bounds")
	}

	return rank*8 + file, nil
}
