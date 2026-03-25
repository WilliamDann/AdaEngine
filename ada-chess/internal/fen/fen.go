package fen

import (
	"errors"
	"strconv"
	"strings"

	"github.com/WilliamDann/AdaEngine/ada-chess/internal/board"
	"github.com/WilliamDann/AdaEngine/ada-chess/internal/game"
)

// map of rune -> piece
var pieces = map[rune]board.Piece{
	'P': board.NewPiece(board.Pawn, board.White),
	'N': board.NewPiece(board.Knight, board.White),
	'B': board.NewPiece(board.Bishop, board.White),
	'R': board.NewPiece(board.Rook, board.White),
	'Q': board.NewPiece(board.Queen, board.White),
	'K': board.NewPiece(board.King, board.White),
	'p': board.NewPiece(board.Pawn, board.Black),
	'n': board.NewPiece(board.Knight, board.Black),
	'b': board.NewPiece(board.Bishop, board.Black),
	'r': board.NewPiece(board.Rook, board.Black),
	'q': board.NewPiece(board.Queen, board.Black),
	'k': board.NewPiece(board.King, board.Black),
}

// get piece data from a FEN string
func parsePieceData(data string) (*board.Chessboard, error) {
	segments := strings.Split(data, "/")
	b        := board.NewChessboard()

	if len(segments) != 8 {
		return nil, errors.New("invalid number of peice data segments")
	}

	for rank, segment := range segments {
		ptr := (7-rank) * 8
		
		for _, chr := range segment {
			if chr >= '0' && chr <= '8' {
				ptr += int(chr - '0')
			} else if strings.ContainsRune("pnbrqkPNBRQK", chr) {
				piece, _ := pieces[chr]
				b.Set(board.Square(ptr), piece)
				ptr++
			} else {
				return nil, errors.New("invalid row data: unknown char")
			}
		}
	}

	return b, nil
}

// get active color from FEN string
func parseActiveColor(data string) (board.Color, error) {
	if data == "w" {
		return board.White, nil
	} else if data == "b" {
		return board.Black, nil
	}
	return 0, errors.New("invalid color segment")
}

// get casting rights from FEN string
func parseCastling(data string) (game.CastlingRights, error) {
	if data == "-" {
		return 0, nil
	}

	var rights game.CastlingRights
	for _, ch := range data {
		switch ch {
		case 'K':
			rights |= game.WhiteKingside
		case 'Q':
			rights |= game.WhiteQueenside
		case 'k':
			rights |= game.BlackKingside
		case 'q':
			rights |= game.BlackQueenside
		default:
			return 0, errors.New("invalid castling segment")
		}
	}

	return rights, nil
}

// get en passant square from FEN string
func parseEnPassant(data string) (board.Square, error) {
	if data == "-" {
		return board.InvalidSquare, nil
	}

	if len(data) != 2 {
		return board.InvalidSquare, errors.New("invalid en passant square")
	}

	file := int(data[0] - 'a')
	rank := int(data[1] - '1')
	sq := board.NewSquare(rank, file)
	if !sq.Valid() {
		return board.InvalidSquare, errors.New("invalid en passant square")
	}

	return sq, nil
}

// get halfmvoes from FEN string
func parseHalfmoveClock(data string) (int, error) {
	n, err := strconv.Atoi(data)
	if err != nil {
		return 0, errors.New("invalid halfmove clock")
	}
	if n < 0 {
		return 0, errors.New("halfmove clock cannot be negative")
	}
	return n, nil
}

// get fullmoves in game from FEN string
func parseFullmoveClock(data string) (int, error) {
	n, err := strconv.Atoi(data)
	if err != nil {
		return 0, errors.New("invalid fullmove counter")
	}
	if n < 1 {
		return 0, errors.New("fullmove counter must be at least 1")
	}
	return n, nil
}


// parse a FEN string into a chessboard 
func Parse(fen string) (*game.Position, error) {
	segments := strings.Split(fen, " ")
	if len(segments) != 6 {
		return nil, errors.New("invalid fen format: incorrect segment number")
	}

	position := game.NewPosition()

	// pieces
	board, err := parsePieceData(segments[0])
	if err != nil {
		return nil, err
	}
	position.Board = board

	// active color
	color, err := parseActiveColor(segments[1])
	if err != nil {
		return nil, err
	}
	position.ActiveColor = color

	// castling
	castling, err := parseCastling(segments[2])
	if err != nil {
		return nil, err
	}
	position.Castling = castling

	// en passant
	epSquare, err := parseEnPassant(segments[3])
	if err != nil {
		return nil, err
	}
	position.EnPassant = epSquare

	// halfmove
	halfmove, err := parseHalfmoveClock(segments[4])
	if err != nil {
		return nil, err
	}
	position.Halfmoves = halfmove

	// fullmove
	fullmove,err := parseFullmoveClock(segments[5])
	if err != nil {
		return nil, err
	}
	position.Fullmoves = fullmove

	// ok
	return position, nil
}

// output a Position into a fen string
func Format(pos *game.Position) string {
	return "NOT IMPLEMENTED"
}
