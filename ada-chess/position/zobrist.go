package position

import (
	"math/rand"

	"github.com/WilliamDann/AdaEngine/ada-chess/core"
)

// Zobrist hash keys — initialized once with a fixed seed for determinism.
var (
	// Random key for each (piece, square) pair.
	// Piece values 0..14 (4-bit encoding), squares 0..63.
	pieceKeys [15][64]uint64

	// One key per castling-rights combination (4 bits → 16 values).
	castlingKeys [16]uint64

	// One key per en-passant file (0-7), plus index 8 for "no EP".
	epKeys [9]uint64

	// XOR'd in when it's black to move.
	sideKey uint64
)

func init() {
	rng := rand.New(rand.NewSource(0xADA))

	for p := range pieceKeys {
		for sq := range pieceKeys[p] {
			pieceKeys[p][sq] = rng.Uint64()
		}
	}
	for i := range castlingKeys {
		castlingKeys[i] = rng.Uint64()
	}
	for i := range epKeys {
		epKeys[i] = rng.Uint64()
	}
	sideKey = rng.Uint64()
}

// ComputeHash builds a Zobrist hash from scratch for the given position.
func ComputeHash(pos *Position) uint64 {
	var h uint64
	for sq := core.Square(0); sq < 64; sq++ {
		piece := pos.Board.Check(sq)
		if piece != core.None {
			h ^= pieceKeys[piece][sq]
		}
	}
	h ^= castlingKeys[pos.Castling]
	if pos.EnPassant.Valid() {
		h ^= epKeys[pos.EnPassant.File()]
	} else {
		h ^= epKeys[8]
	}
	if pos.ActiveColor == core.Black {
		h ^= sideKey
	}
	return h
}
