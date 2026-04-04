package search

import (
	"github.com/WilliamDann/AdaEngine/ada-chess/core"
)

// represents bound type
type SearchFlag uint8
const (
	Exact SearchFlag = iota // score is an exact eval
	LowerBound              // score is >= beta
	UpperBound              // score is <= alpha
)

// an entry in the transposition table
type TTEntry struct {
	Key   uint64
	Move  core.Move
	Depth int8
	Score int16
	Flag  SearchFlag
}

// the transposition table
type TT struct {
	entries []TTEntry
	mask    uint64       // size - 1
}

func NewTT(size int) *TT {
	for size&(size-1) != 0 {
		size &= size - 1
	}

	return &TT{
		entries: make([]TTEntry, size),
		mask:    uint64(size - 1),
	}
}

func packData(e TTEntry) uint64 {
	return uint64(e.Move) | uint64(e.Depth)<<16 | uint64(e.Score) << 24 | uint64(e.Flag)<<40
}

// check the transposition table
func (tt *TT) Probe(key uint64) (TTEntry, bool) {
	if tt == nil {
		return TTEntry{}, false
	}
	entry := tt.entries[key&tt.mask]
	entry.Key ^= packData(entry)
	if entry.Key == key {
		return entry, true
	}
	return TTEntry{}, false
}

// add to the transposition table
func (tt *TT) Store(entry TTEntry) {
	if tt == nil {
		return
	}
	idx        := entry.Key & tt.mask
	packed     := entry
	packed.Key ^= packData(entry)
	tt.entries[idx] = packed
}
