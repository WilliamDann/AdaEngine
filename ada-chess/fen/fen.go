package fen

import (
	"errors"
	"strconv"
	"strings"

	"github.com/WilliamDann/AdaEngine/ada-chess/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/position"
)

// map of rune -> piece
var pieces = map[rune]core.Piece{
	'P': core.NewPiece(core.Pawn, core.White),
	'N': core.NewPiece(core.Knight, core.White),
	'B': core.NewPiece(core.Bishop, core.White),
	'R': core.NewPiece(core.Rook, core.White),
	'Q': core.NewPiece(core.Queen, core.White),
	'K': core.NewPiece(core.King, core.White),
	'p': core.NewPiece(core.Pawn, core.Black),
	'n': core.NewPiece(core.Knight, core.Black),
	'b': core.NewPiece(core.Bishop, core.Black),
	'r': core.NewPiece(core.Rook, core.Black),
	'q': core.NewPiece(core.Queen, core.Black),
	'k': core.NewPiece(core.King, core.Black),
}

// get piece data from a FEN string
func parsePieceData(data string) (*core.Chessboard, error) {
	segments := strings.Split(data, "/")
	b        := core.NewChessboard()

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
				b.Set(core.Square(ptr), piece)
				ptr++
			} else {
				return nil, errors.New("invalid row data: unknown char")
			}
		}
	}

	return b, nil
}

// get active color from FEN string
func parseActiveColor(data string) (core.Color, error) {
	if data == "w" {
		return core.White, nil
	} else if data == "b" {
		return core.Black, nil
	}
	return 0, errors.New("invalid color segment")
}

// get casting rights from FEN string
func parseCastling(data string) (position.CastlingRights, error) {
	if data == "-" {
		return 0, nil
	}

	var rights position.CastlingRights
	for _, ch := range data {
		switch ch {
		case 'K':
			rights |= position.WhiteKingside
		case 'Q':
			rights |= position.WhiteQueenside
		case 'k':
			rights |= position.BlackKingside
		case 'q':
			rights |= position.BlackQueenside
		default:
			return 0, errors.New("invalid castling segment")
		}
	}

	return rights, nil
}

// get en passant square from FEN string
func parseEnPassant(data string) (core.Square, error) {
	if data == "-" {
		return core.InvalidSquare, nil
	}

	if len(data) != 2 {
		return core.InvalidSquare, errors.New("invalid en passant square")
	}

	file := int(data[0] - 'a')
	rank := int(data[1] - '1')
	sq := core.NewSquare(rank, file)
	if !sq.Valid() {
		return core.InvalidSquare, errors.New("invalid en passant square")
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
func Parse(fen string) (*position.Position, error) {
	segments := strings.Split(fen, " ")
	if len(segments) != 6 {
		return nil, errors.New("invalid fen format: incorrect segment number")
	}

	pos := position.NewPosition()

	// pieces
	board, err := parsePieceData(segments[0])
	if err != nil {
		return nil, err
	}
	pos.Board = board

	// active color
	color, err := parseActiveColor(segments[1])
	if err != nil {
		return nil, err
	}
	pos.ActiveColor = color

	// castling
	castling, err := parseCastling(segments[2])
	if err != nil {
		return nil, err
	}
	pos.Castling = castling

	// en passant
	epSquare, err := parseEnPassant(segments[3])
	if err != nil {
		return nil, err
	}
	pos.EnPassant = epSquare

	// halfmove
	halfmove, err := parseHalfmoveClock(segments[4])
	if err != nil {
		return nil, err
	}
	pos.Halfmoves = halfmove

	// fullmove
	fullmove,err := parseFullmoveClock(segments[5])
	if err != nil {
		return nil, err
	}
	pos.Fullmoves = fullmove

	// ok
	return pos, nil
}

// output a Position into a fen string
func Format(pos *position.Position) string {
	return "NOT IMPLEMENTED"
}
