package game

import (
	"strings"
)

type Board struct {
	squares map[Coord]Piece
	pieces  map[Piece][]Coord
}

func NewBoard() *Board {
	var b Board

	b.pieces = map[Piece][]Coord{}
	b.squares = map[Coord]Piece{}

	return &b
}

func (b Board) IsEmpty(coord Coord) bool {
	return b.Get(coord).Is(*NoPiece())
}

func (b Board) Get(coord Coord) Piece {
	piece, found := b.squares[coord]
	if !found {
		return *NoPiece()
	}
	return piece
}

// add a piece to the board
func (b *Board) Set(piece Piece, coord Coord) bool {
	if !b.IsEmpty(coord) {
		return false
	}

	b.pieces[piece] = append(b.pieces[piece], coord)
	b.squares[coord] = piece

	return true
}

// remove the piece at a coord
func (b *Board) Clear(coord Coord) {
	// get the type of piece at the given square
	if b.IsEmpty(coord) {
		return
	}
	piece := b.Get(coord)

	// find the piece and remove it
	for i, found := range b.pieces[piece] {
		if coord.Equ(found) {
			// remove from peice list
			b.pieces[piece] = arrRemove(b.pieces[piece], i)

			// remove from square map
			delete(b.squares, coord)

			return
		}
	}

	// the piece did not exist
}

func (b Board) FenPieceData() string {
	var sb strings.Builder

	empty := 0
	for y := 7; y >= 0; y-- {
		for x := 7; x >= 0; x-- {
			coord := Coord{x, y}

			// count empty square to write at next non empty sqyare
			if b.IsEmpty(coord) {
				empty++
			} else {
				if empty != 0 {
					// write empty number
					sb.WriteRune(rune('0' + empty))
					empty = 0

				}
				// write piece char
				sb.WriteString(b.Get(coord).String())
			}
		}

		if empty != 0 {
			sb.WriteRune(rune('0' + empty))
			empty = 0
		}

		sb.WriteRune('/')
	}

	str := sb.String()
	return str[:len(str)-1]
}

func (b Board) String() string {
	var sb strings.Builder

	sb.WriteRune(' ')
	sb.WriteRune(' ')

	for x := 0; x <= 7; x++ {
		sb.WriteRune(rune('a' + x))
		sb.WriteRune(' ')
	}
	sb.WriteRune('\n')

	for y := 0; y <= 7; y++ {
		sb.WriteRune(rune('8' - y))
		sb.WriteRune(' ')
		for x := 0; x <= 7; x++ {
			sb.WriteString(b.Get(*NewCoord(x, y)).String())
			sb.WriteRune(' ')
		}
		sb.WriteRune('\n')
	}

	return sb.String()
}

func arrRemove(s []Coord, i int) []Coord {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
