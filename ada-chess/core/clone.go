package core

// Clone returns a deep copy of the chessboard.
func (b *Chessboard) Clone() *Chessboard {
	c := *b
	return &c
}
