package core

// represents an index in the chess board
type Square int

// sentinel value for invalid squares
const InvalidSquare Square = 255
const InvalidSquareName string = "??"

// Square is an index on the chess board mapped from rank and file:
//
//   file  0  1  2  3  4  5  6  7
//         a  b  c  d  e  f  g  h
//   rank +--------------------------
//   7  8 | 56 57 58 59 60 61 62 63
//   6  7 | 48 49 50 51 52 53 54 55
//   5  6 | 40 41 42 43 44 45 46 47
//   4  5 | 32 33 34 35 36 37 38 39
//   3  4 | 24 25 26 27 28 29 30 31
//   2  3 | 16 17 18 19 20 21 22 23
//   1  2 |  8  9 10 11 12 13 14 15
//   0  1 |  0  1  2  3  4  5  6  7
//
//   index = rank * 8 + file
func NewSquare(rank, file int) Square {
	if !validAxis(rank) || !validAxis(file) {
		return InvalidSquare
	}

	return Square(rank * 8 + file)
}

func (sq Square) Valid() bool {
	return sq >= 0 && sq < 64
}

// checks if a rank or file is valid
func validAxis(value int) bool {
	return value >= 0 && value < 8
}

func (sq Square) Rank() int {
	if !sq.Valid() {
		return -1
	}
	return int(sq) / 8
}
func (sq Square) File() int {
	if !sq.Valid() {
		return -1
	}
	return int(sq) % 8
}

func (sq Square) String() string {
	if !sq.Valid() {
		return InvalidSquareName
	}
	return string(rune('a' + sq.File())) + string(rune('1' + sq.Rank()))
}
