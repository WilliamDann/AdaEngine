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

	// pieces
	for sq := range pos.Board.Occupied().Squares() {
		piece := pos.Board.Check(sq)
		if piece != core.None {
			hsh ^= pieceSquareKeys[piece][sq]
		}
	}

	// side to move
	if pos.ActiveColor == core.Black {
		hsh ^= sideToMoveKey
	}

	// castling
	for i := 0; i < 4; i++ {
		if pos.Castling.Has(CastlingRights(1 << i)) {
			hsh ^= castlingKeys[i]
		}
	}

	// en peasant
	if pos.EnPassant.Valid() {
		hsh ^= enPassantKeys[pos.EnPassant.File()]
	}

	return hsh
}

