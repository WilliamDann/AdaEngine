package movegen

import (
	"errors"
	"math/rand"
	"math/bits"

	"github.com/WilliamDann/AdaEngine/ada-chess/core"
)

// magic bitboards for sliding pieces
type MagicEntry struct {
	mask       core.Bitboard
	magic      uint64
	indexBits  uint8
}

// type of function to get the squares a piece can see
type AttackFn func(core.Square, core.Bitboard) core.Bitboard

var rookMagics   [64]MagicEntry
var bishopMagics [64]MagicEntry

var rookMoves    [64][]core.Bitboard
var bishopMoves  [64][]core.Bitboard

// find magic numbers
func findMagic(
	square    core.Square,
	mask      core.Bitboard,
	attackFn  AttackFn,
	indexBits uint8,
) (MagicEntry, []core.Bitboard) {
	for {
		entry      := MagicEntry{
			mask: mask,
			magic: rand.Uint64() & rand.Uint64() & rand.Uint64(),
			indexBits: indexBits,
		}
		table, err := tryMakeTable(square, entry, attackFn)
		if err == nil {
			return entry, table
		}
	}
}

func tryMakeTable(
	square core.Square,
	entry MagicEntry,
	attackFn AttackFn,
) ([]core.Bitboard, error) {
	table := make([]core.Bitboard, 1 << entry.indexBits)
	for blockers := core.Bitboard(0); ; {
		moves := attackFn(square, blockers)
		index := magicIndex(entry, blockers)

		if table[index] == 0 {
			table[index] = moves
		} else if table[index] != moves {
			return nil, errors.New("magic table collision")
		}

		blockers = core.Bitboard((uint64(blockers) - uint64(entry.mask)) & uint64(entry.mask))
		if blockers == 0 {
			break
		}
	}

	return table, nil
}

// find bitboard for legal moves using a magic number
func magicIndex(entry MagicEntry, blockers core.Bitboard) int {
	blockers = blockers.Intersection(entry.mask)
	hash    := uint64(blockers) * entry.magic
	return int(hash >> (64 - entry.indexBits))
}

func RookMoves(square core.Square, blockers core.Bitboard) core.Bitboard {
	entry := rookMagics[square]
	return rookMoves[square][magicIndex(entry, blockers)]
}
func BishopMoves(square core.Square, blockers core.Bitboard) core.Bitboard {
	entry := bishopMagics[square]
	return bishopMoves[square][magicIndex(entry, blockers)]
}
func QueenMoves(square core.Square, blockers core.Bitboard) core.Bitboard {
	return RookMoves(square, blockers).Union(BishopMoves(square, blockers))
}

// GenerateMagics computes magic bitboard tables for all 64 squares.
// Called by init(); exported so it can be benchmarked.
func GenerateMagics() ([64]MagicEntry, [64][]core.Bitboard, [64]MagicEntry, [64][]core.Bitboard) {
	var rMagics, bMagics [64]MagicEntry
	var rMoves, bMoves   [64][]core.Bitboard

	for sq := core.Square(0); sq < 64; sq++ {
		rMask := rookMask(sq)
		rBits := uint8(bits.OnesCount64(uint64(rMask)))
		rMagics[sq], rMoves[sq] = findMagic(sq, rMask, rookAttacks, rBits)

		bMask := bishopMask(sq)
		bBits := uint8(bits.OnesCount64(uint64(bMask)))
		bMagics[sq], bMoves[sq] = findMagic(sq, bMask, bishopAttacks, bBits)
	}

	return rMagics, rMoves, bMagics, bMoves
}

// find magic bitboards on startup
func init() {
	rookMagics, rookMoves, bishopMagics, bishopMoves = GenerateMagics()
}
