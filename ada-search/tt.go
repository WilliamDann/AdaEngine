package search

import (
	"sync/atomic"

	"github.com/WilliamDann/AdaEngine/ada-chess/core"
)

// TTFlag describes what kind of bound the stored score represents.
type TTFlag uint8

const (
	FlagExact TTFlag = iota // Score is exact (PV node)
	FlagAlpha               // Score is an upper bound (all moves failed low)
	FlagBeta                // Score is a lower bound (beta cutoff)
)

// TTEntry is returned by Probe.
type TTEntry struct {
	Hash  uint64
	Move  core.Move
	Score int16
	Depth int8
	Flag  TTFlag
}

// ttSlot is a single lock-free slot in the transposition table.
// The key is stored as hash XOR data so that a torn read (one half
// written by a different goroutine) is detected as a miss rather
// than returning corrupt data.
type ttSlot struct {
	key  atomic.Uint64
	data atomic.Uint64
}

// TT is a fixed-size, lock-free transposition table safe for
// concurrent use by multiple goroutines.
type TT struct {
	slots []ttSlot
	mask  uint64
}

// NewTT creates a table with the given number of entries (rounded up to power of 2).
func NewTT(size int) *TT {
	n := uint64(1)
	for n < uint64(size) {
		n <<= 1
	}
	return &TT{
		slots: make([]ttSlot, n),
		mask:  n - 1,
	}
}

func packData(move core.Move, score int16, depth int8, flag TTFlag) uint64 {
	return uint64(move) |
		uint64(uint16(score))<<16 |
		uint64(uint8(depth))<<32 |
		uint64(flag)<<40
}

// Probe looks up a position by hash. Returns the entry and whether it was a hit.
func (tt *TT) Probe(hash uint64) (TTEntry, bool) {
	slot := &tt.slots[hash&tt.mask]
	data := slot.data.Load()
	key := slot.key.Load()
	if key^data != hash {
		return TTEntry{}, false
	}
	return TTEntry{
		Hash:  hash,
		Move:  core.Move(data & 0xFFFF),
		Score: int16(uint16(data >> 16)),
		Depth: int8(uint8(data >> 32)),
		Flag:  TTFlag(data >> 40),
	}, true
}

// Store writes an entry. For the same position, newer entries always replace
// (iterative deepening ensures they are at least as deep). For different
// positions mapping to the same slot, the deeper entry is kept.
func (tt *TT) Store(hash uint64, move core.Move, score int, depth int, flag TTFlag) {
	data := packData(move, int16(score), int8(depth), flag)
	slot := &tt.slots[hash&tt.mask]

	oldData := slot.data.Load()
	oldKey := slot.key.Load()
	oldHash := oldKey ^ oldData
	if oldHash != hash {
		// Different position — only replace if new depth >= old depth
		oldDepth := int8(uint8(oldData >> 32))
		if int8(depth) < oldDepth {
			return
		}
	}

	slot.data.Store(data)
	slot.key.Store(hash ^ data)
}
