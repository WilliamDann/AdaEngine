package position

import (
	"github.com/WilliamDann/AdaEngine/ada-chess/core"
	"math/rand"
)

// unique keys for each hash value
var rng = rand.New(rand.NewSource(0x1234567890ABCDEF))
var (
	pieceSquareKeys [15][64]uint64  // one key per (piece, square)
	sideToMoveKey   uint64          // key for black to move
	castlingKeys    [4]uint64       // one per casting rights bit
	enPassantKeys   [8]uint64       // one per file
)

// helper for hasing castling keys
func castlingKey(p *Position) uint64 {
	var hsh uint64

	for i := 0; i < 4; i++ {
		if p.Castling.Has(CastlingRights(1 << i)) {
			hsh ^= castlingKeys[i]
		}
	}

	return hsh
}

// helper for active color
func activeColorKey(p *Position) uint64 {
	if p.ActiveColor == core.Black {
		return sideToMoveKey
	}
	return 0
}

// helper for en passant
func enPassantKey(p *Position) uint64 {
	var hsh uint64

	// en peasant
	if p.EnPassant.Valid() {
		hsh ^= enPassantKeys[p.EnPassant.File()]
	}

	return hsh
}

// helper for board state
func piecesKey(p *Position) uint64 {
	var hsh uint64

	for sq := range p.Board.Occupied().Squares() {
		piece := p.Board.Check(sq)
		if piece != core.None {
			hsh ^= pieceSquareKeys[piece][sq]
		}
	}

	return hsh
}

// determine keys at program start
func init() {
	// pieces
	for piece := range pieceSquareKeys {
		for sq := range pieceSquareKeys[piece] {
			pieceSquareKeys[piece][sq] = rng.Uint64()
		}
	}

	// black to move
	sideToMoveKey = rng.Uint64()

	// castling bits
	for i := range castlingKeys {
		castlingKeys[i] = rng.Uint64()
	}

	// en passant files
	for i := range enPassantKeys {
		enPassantKeys[i] = rng.Uint64()
	}
}

// compute a zobrist hash for a position
func (pos *Position) ComputeZobrist() uint64 {
	var hsh uint64

	hsh ^= piecesKey(pos)
	hsh ^= activeColorKey(pos)
	hsh ^= castlingKey(pos)
	hsh ^= enPassantKey(pos)

	return hsh
}

