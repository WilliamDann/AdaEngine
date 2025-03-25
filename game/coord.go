package game

import (
	"strconv"
)

type Coord struct {
	X int
	Y int
}

func (s Coord) Valid() bool {
	return s.X >= 0 && s.X <= 7 && s.Y >= 0 && s.Y <= 7
}

func (s Coord) Add(other Coord) Coord {
	s.X += other.X
	s.Y += other.Y
	return s
}

func (s Coord) Sub(other Coord) Coord {
	s.X -= other.X
	s.Y -= other.Y
	return s
}

func (s Coord) Equ(other Coord) bool {
	return s.X == other.X && s.Y == other.Y
}

func CoordNumToLetter(row int) rune {
	return rune('a' + row)
}

func CoordLetterToNum(letter rune) int {
	return int(letter - 'a')
}

func (s Coord) String() string {
	return string(CoordNumToLetter(s.X)) + strconv.FormatInt(int64(s.Y+1), 10)
}

func NewCoordSan(san string) *Coord {
	x := CoordLetterToNum(rune(san[0]))
	y := int(san[1] - '1')

	return NewCoord(x, y)
}

func NewCoord(x int, y int) *Coord {
	coord := Coord{x, 7 - y}
	return &coord
}
