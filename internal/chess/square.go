package chess

// get the square number of given coords
func Square(file, rank int) int {
	return rank * 8 + file
}
// get the coords of a given squrae number
func Coords(square int) (file, rank int) {
	return square % 8, square / 8
}

// get the Standard Algebraic Notation (a1, c3, d8, etc) of a given square num
func SAN(square int) string {
	file, rank := Coords(square)
	return string('a' + file) + string('1' + rank)
}
