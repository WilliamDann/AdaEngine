package chess

import "iter"

// represents a move
type Move struct {
    From, To int
}

// defines a direction to move in on the board
type Direction = [2]int // [file, rank]

// direction definitions
var (
    // Slider directions (rook, bishop, queen)
    North     Direction = Direction{0, 1}
    Northeast Direction = Direction{1, 1}
    East      Direction = Direction{1, 0}
    Southeast Direction = Direction{1, -1}
    South     Direction = Direction{0, -1}
    Southwest Direction = Direction{-1, -1}
    West      Direction = Direction{-1, 0}
    Northwest Direction = Direction{-1, 1}
    // Knight moves
    Knight_NNE Direction = Direction{1, 2}
    Knight_ENE Direction = Direction{2, 1}
    Knight_ESE Direction = Direction{2, -1}
    Knight_SSE Direction = Direction{1, -2}
    Knight_SSW Direction = Direction{-1, -2}
    Knight_WSW Direction = Direction{-2, -1}
    Knight_WNW Direction = Direction{-2, 1}
    Knight_NNW Direction = Direction{-1, 2}
)

// MoveRule generates legal moves for a piece considering board state
type MoveRule = func(square int, board *Board) Bitboard
type MoveIter = func(from, to int) bool

// generation for pieces that slide
func slider(square int, step Direction, board *Board) Bitboard {
    var result Bitboard = 0
    file, rank := Coords(square)
    nextFile := file + step[0]
    nextRank := rank + step[1]
    
    for nextFile >= 0 && nextFile < 8 && nextRank >= 0 && nextRank < 8 {
        nextSquare := Square(nextFile, nextRank)
        result = result.Set(nextSquare)
 
        // Stop if we hit any piece
        if board.blockers.Check(nextSquare) {
            break
        }
 
        nextFile += step[0]
        nextRank += step[1]
    }
    
    // Remove friendly pieces
    piece := board.Read(square)
    if piece.Color() == White {
        return result.Difference(board.white)
    }
    return result.Difference(board.black)
}

// generation for pieces that step
func stepper(square int, step Direction, board *Board) Bitboard {
    var result Bitboard = 0
    file, rank := Coords(square)
    nextFile := file + step[0]
    nextRank := rank + step[1]

		if nextFile >= 0 && nextFile < 8 && nextRank >= 0 && nextRank < 8 {
        result = result.Set(Square(nextFile, nextRank))
    }

    // Remove friendly pieces based on the piece at the origin square
    piece := board.Read(square)
    if piece.Color() == White {
        return result.Difference(board.white)
    }
    return result.Difference(board.black)
}

// Legal moves for pieces
func BishopMoves(square int, board *Board) Bitboard {
    var result Bitboard = 0
    result = result.Union(slider(square, Northeast, board))
    result = result.Union(slider(square, Southeast, board))
    result = result.Union(slider(square, Southwest, board))
    result = result.Union(slider(square, Northwest, board))
    return result
}
func RookMoves(square int, board *Board) Bitboard {
    var result Bitboard = 0
    result = result.Union(slider(square, North, board))
    result = result.Union(slider(square, East, board))
    result = result.Union(slider(square, South, board))
    result = result.Union(slider(square, West, board))
    return result
}
func QueenMoves(square int, board *Board) Bitboard {
    return BishopMoves(square, board).Union(RookMoves(square, board))
}
func KnightMoves(square int, board *Board) Bitboard {
  	var result Bitboard = 0
    result = result.Union(stepper(square, Knight_NNE, board))
    result = result.Union(stepper(square, Knight_ENE, board))
    result = result.Union(stepper(square, Knight_ESE, board))
    result = result.Union(stepper(square, Knight_SSE, board))
    result = result.Union(stepper(square, Knight_SSW, board))
    result = result.Union(stepper(square, Knight_WSW, board))
    result = result.Union(stepper(square, Knight_WNW, board))
    result = result.Union(stepper(square, Knight_NNW, board))
    return result
}
func KingMoves(square int, board *Board) Bitboard {
    var result Bitboard = 0
    result = result.Union(stepper(square, North, board))
    result = result.Union(stepper(square, Northeast, board))
    result = result.Union(stepper(square, East, board))
    result = result.Union(stepper(square, Southeast, board))
    result = result.Union(stepper(square, South, board))
    result = result.Union(stepper(square, Southwest, board))
    result = result.Union(stepper(square, West, board))
    result = result.Union(stepper(square, Northwest, board))

    // Castling
    piece := board.pieces[square]
    color := piece.Color()

    // Check if king is on starting square
    file, rank := Coords(square)
    expectedRank := 0
    if color == Black {
        expectedRank = 7
    }

    if file == 4 && rank == expectedRank {
        // Kingside castling
        if board.castling.CanCastleKingside(color) {
            // Check squares between king and rook are empty
            if !board.blockers.Check(Square(5, rank)) && 
               !board.blockers.Check(Square(6, rank)) {
                // TODO: Also need to check king isn't in check, doesn't pass through check,
                // and doesn't land in check
                result = result.Set(Square(6, rank))
            }
        }
 
        // Queenside castling
        if board.castling.CanCastleQueenside(color) {
            // Check squares between king and rook are empty
            if !board.blockers.Check(Square(3, rank)) && 
               !board.blockers.Check(Square(2, rank)) && 
               !board.blockers.Check(Square(1, rank)) {
                // TODO: Also need to check king isn't in check, doesn't pass through check,
                // and doesn't land in check
                result = result.Set(Square(2, rank))
            }
        }
    }

		return result
}
func PawnMoves(square int, board *Board) Bitboard {
    var result Bitboard = 0
    file, rank := Coords(square)

    // Determine pawn color from board.pieces
    piece := board.pieces[square]
    color := piece.Color()

    var forward Direction
    var startRank int
    var enemyPieces Bitboard

    if color == White {
        forward = North
        startRank = 1
        enemyPieces = board.black
    } else {
        forward = South
        startRank = 6
        enemyPieces = board.white
    }

    // Single push forward
    pushFile := file + forward[0]
    pushRank := rank + forward[1]
    if pushFile >= 0 && pushFile < 8 && pushRank >= 0 && pushRank < 8 {
        pushSquare := Square(pushFile, pushRank)
        if !board.blockers.Check(pushSquare) {
            result = result.Set(pushSquare)

            // Double push from starting rank
            if rank == startRank {
                doublePushFile := file + forward[0]*2
                doublePushRank := rank + forward[1]*2
                doublePushSquare := Square(doublePushFile, doublePushRank)
                if !board.blockers.Check(doublePushSquare) {
                    result = result.Set(doublePushSquare)
                }
            }
        }
    }

    // Diagonal captures
    for _, offset := range []int{-1, 1} {
        captureFile := file + offset
        captureRank := rank + forward[1]

        if captureFile >= 0 && captureFile < 8 && captureRank >= 0 && captureRank < 8 {
            captureSquare := Square(captureFile, captureRank)

            // Normal capture
            if enemyPieces.Check(captureSquare) {
                result = result.Set(captureSquare)
            }

            // En passant
            if board.enPassant >= 0 && captureSquare == board.enPassant {
                result = result.Set(captureSquare)
            }
        }
    }

    return result
}

// LegalMoves returns an iterator over all legal moves for a color
func LegalMoves(b *Board) iter.Seq[Move] {
	color := b.activeColor
    return func(yield func(Move) bool) {
        b.ForEachColorPiece(color, func(piece Piece, from int) bool {
            // Generate pseudo-legal moves
            var moves Bitboard
            switch piece.Type() {
            case Pawn:
                moves = PawnMoves(from, b)
            case Knight:
                moves = KnightMoves(from, b)
            case Bishop:
                moves = BishopMoves(from, b)
            case Rook:
                moves = RookMoves(from, b)
            case Queen:
                moves = QueenMoves(from, b)
            case King:
                moves = KingMoves(from, b)
            }

						// Try each move
            for moves != 0 {
                var to int
                moves, to = moves.PopLSB()

								// Make the move
                captured := b.Read(to)
                b.Clear(from)
                b.Clear(to)
                b.Place(piece, to)

                // Check if legal
                legal := !b.IsInCheck(color)

                // Undo the move
                b.Clear(to)
                b.Place(piece, from)
                if captured.Type() != None {
                    b.Place(captured, to)
                }

								// Yield if legal
                if legal {
                    if !yield(Move{from, to}) {
                        return false
                    }
                }
            }

            return true
        })
    }
}
