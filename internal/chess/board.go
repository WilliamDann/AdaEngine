package chess

import (
	"strings"
	"fmt"
)

// function to iterate over pieces in a board
type PieceIterator func(piece Piece, square int) bool

// stores board state
type Board struct {
	// piece data
	blockers		Bitboard				// map of all pieces

	white				Bitboard				// map of white pieces
	black				Bitboard				// map of black pieces

	pieces 			map[int]Piece 	// square -> Piece code


	// fen data
	activeColor Color
	castling    CastlingRights	// if a player can castle
	enPassant   int 						// en passant square (-1 if none)
	halfmove    int
	fullmove		int
}

// create an empty board
func NewBoard() *Board {
	b := Board{}
	b.Reset()
	return &b
}

// Iterate over all pieces on the board
func (b *Board) ForEachPiece(fn PieceIterator) {
	for square := 0; square < 64; square++ {
		if !b.IsClear(square) {
			if !fn(b.pieces[square], square) {
				return
			}
		}
	}
}

// Iterate over pieces of a specific color
func (b *Board) ForEachColorPiece(color Color, fn PieceIterator) {
    // Choose the appropriate bitboard
    var colorBoard Bitboard
    if color == White {
        colorBoard = b.white
    } else {
        colorBoard = b.black
    }

    // Iterate using PopLSB
    for colorBoard != 0 {
        var square int
        colorBoard, square = colorBoard.PopLSB()
        if !fn(b.pieces[square], square) {
            return
        }
    }
}

// reset the board to starting state
func (b *Board) Reset() {
	b.blockers = 0
	b.white    = 0
	b.black    = 0
	b.pieces   = map[int]Piece{}
}

// get the peice at a given square
func (b *Board) Read(square int) Piece {
	if (b.IsClear(square)) {
		return NewPiece(None, White)
	}
	return b.pieces[square]
}

// place a piece at a given square
func (b *Board) Place(piece Piece, square int) bool {
	// cannot place a piece where one already exists
	if (!b.IsClear(square)) {
		return false
	}

	// add piece to list
	b.pieces[square] = piece

	// update tracking Bitboards
	if (piece.Color() == White) {
		b.white = b.white.Set(square)
	} else {
		b.black = b.black.Set(square)
	}
	b.blockers = b.blockers.Set(square)
	return true
}

// move a piece from one location to another
//  (note! this does not check if the move is legal)
func (b *Board) Move(from, to int) bool {
	if (b.IsClear(from) || !b.IsClear(to)) {
		return false
	}

	b.Place(b.Read(from), to)
	b.Clear(from)

	return true
}

// helpers for piece lookup & move generation
func (b *Board) IsClear(square int) bool {
	return !b.blockers.Check(square)
}
func (b *Board) HasWhite(square int) bool {
	return b.white.Check(square)
}
func (b *Board) HasBlack(squrae int) bool {
	return b.black.Check(squrae)
}

// clear a given square
func (b *Board) Clear(square int) {
	b.blockers = b.blockers.Clear(square)
	b.white = b.white.Clear(square)
	b.black = b.black.Clear(square)
	b.pieces[square] = NewPiece(None, White)
}

// IsInCheck returns true if the given color's king is under attack
func (b *Board) IsInCheck(color Color) bool {
    // Find the king
    kingSquare := -1
    b.ForEachColorPiece(color, func(piece Piece, square int) bool {
        if piece.Type() == King {
            kingSquare = square
            return false // stop iterating
        }
        return true
    })
    
    if kingSquare == -1 {
        return false // No king found (shouldn't happen in valid position)
    }
    
    // Check if any enemy piece can attack the king square
    enemyColor := Black
    if color == Black {
        enemyColor = White
    }
    
    return b.IsSquareAttacked(kingSquare, enemyColor)
}

// IsSquareAttacked checks if a square is attacked by the given color
func (b *Board) IsSquareAttacked(square int, byColor Color) bool {
    // Check each enemy piece to see if it attacks this square
    attacked := false
    
    b.ForEachColorPiece(byColor, func(piece Piece, from int) bool {
        var attacks Bitboard
        
        switch piece.Type() {
        case Pawn:
            attacks = PawnAttacks(from, b, byColor)
        case Knight:
            attacks = KnightMoves(from, b)
        case Bishop:
            attacks = BishopMoves(from, b)
        case Rook:
            attacks = RookMoves(from, b)
        case Queen:
            attacks = QueenMoves(from, b)
        case King:
            attacks = KingMoves(from, b)
        }
        
        if attacks.Check(square) {
            attacked = true
            return false // stop iterating
        }
        return true
    })
    
    return attacked
}

// PawnAttacks returns squares a pawn attacks (not moves - pawns attack diagonally)
func PawnAttacks(square int, board *Board, color Color) Bitboard {
    var result Bitboard = 0
    file, rank := Coords(square)
    
    var forward int
    if color == White {
        forward = 1
    } else {
        forward = -1
    }
    
    // Left diagonal attack
    if file > 0 && rank+forward >= 0 && rank+forward < 8 {
        result = result.Set(Square(file-1, rank+forward))
    }
    
    // Right diagonal attack
    if file < 7 && rank+forward >= 0 && rank+forward < 8 {
        result = result.Set(Square(file+1, rank+forward))
    }
    
    return result
}

// tostring (generated by Claude)
func (b *Board) String() string {
	var sb strings.Builder
	
	// Iterate through ranks from 8 to 1 (top to bottom)
	for rank := 7; rank >= 0; rank-- {
		// Rank number
		sb.WriteString(fmt.Sprintf("%d ", rank+1))
		
		// Iterate through files a-h (left to right)
		for file := 0; file < 8; file++ {
			square := rank*8 + file
			piece := b.Read(square)
			
			// Get piece symbol
			var symbol string
			if piece.Type() == None {
				symbol = "."
			} else {
				symbol = pieceSymbol(piece)
			}
			
			sb.WriteString(symbol + " ")
		}
		
		sb.WriteString("\n")
	}
	
	// File labels
	sb.WriteString("  a b c d e f g h\n")
	
	// Add game info
	sb.WriteString("\n")
	if b.activeColor == White {
		sb.WriteString("Side to move: White\n")
	} else {
		sb.WriteString("Side to move: Black\n")
	}
	sb.WriteString(fmt.Sprintf("Castling: %s\n", b.castling.ToFEN()))
	sb.WriteString(fmt.Sprintf("Move: %d\n", b.fullmove))
	
	return sb.String()
}

// Helper function to get piece symbol
func pieceSymbol(piece Piece) string {
	var symbol string
	
	switch piece.Type() {
	case Pawn:
		symbol = "P"
	case Knight:
		symbol = "N"
	case Bishop:
		symbol = "B"
	case Rook:
		symbol = "R"
	case Queen:
		symbol = "Q"
	case King:
		symbol = "K"
	default:
		symbol = "?"
	}
	
	// Use lowercase for black pieces
	if piece.Color() == Black {
		symbol = strings.ToLower(symbol)
	}
	
	return symbol
}
