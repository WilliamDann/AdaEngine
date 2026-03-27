package search

import (
	"github.com/WilliamDann/AdaEngine/ada-chess/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/position"
)

// Piece values in centipawns.
var pieceValue = [7]int{
	0,    // None
	100,  // Pawn
	320,  // Knight
	330,  // Bishop
	500,  // Rook
	900,  // Queen
	0,    // King (not counted)
}

// Advancement bonus per rank in centipawns (for non-pawn pieces).
var advanceBonus = [7]int{
	0,  // None
	0,  // Pawn — uses pawnAdvance table instead
	3,  // Knight
	2,  // Bishop
	1,  // Rook
	1,  // Queen
	0,  // King
}

// Pawn advancement bonus by rank. Flat for early pushes, steep near promotion.
//            rank: 0   1   2   3   4    5    6    7
var pawnAdvance = [8]int{0, 0, 5, 5, 15, 30, 60, 100}

// Center bonus per piece type, scaled by closeness to center (0-3).
var centerBonus = [7]int{
	0,  // None
	2,  // Pawn
	5,  // Knight — strongest center preference
	3,  // Bishop
	1,  // Rook
	2,  // Queen
	0,  // King
}

// King zone coordination weight — multiplied by count^2 so a lone piece
// near the king gets almost nothing but a group gets a large bonus.
const kingZoneWeight = 3

// Bonus for a rook/queen on the same rank or file as the enemy king,
// or a bishop/queen on the same diagonal. Regardless of blockers —
// the piece is aimed at the king.
const lineAlignWeight = 10

// Precomputed tables, filled in init().
var (
	centerDist [64]int           // Chebyshev distance from center (0-3)
	kingZone   [64]core.Bitboard // 3x3 area around each square

	rankMask  [8]core.Bitboard    // all squares on rank r
	fileMask  [8]core.Bitboard    // all squares on file f
	diagMask  [15]core.Bitboard   // all squares on diagonal (r-f+7)
	adiagMask [15]core.Bitboard   // all squares on anti-diagonal (r+f)
	between   [64][64]core.Bitboard // squares strictly between two squares on the same line
)

func init() {
	for sq := 0; sq < 64; sq++ {
		r, f := sq/8, sq%8

		// Center distance
		rd := max(3-r, r-4)
		fd := max(3-f, f-4)
		centerDist[sq] = max(rd, fd)

		// King zone: 3x3 (Chebyshev distance 1)
		var mask core.Bitboard
		for zr := max(0, r-1); zr <= min(7, r+1); zr++ {
			for zf := max(0, f-1); zf <= min(7, f+1); zf++ {
				mask = mask.Set(core.Square(zr*8 + zf))
			}
		}
		kingZone[sq] = mask

		// Line masks
		bb := core.Bitboard(0).Set(core.Square(sq))
		rankMask[r] |= bb
		fileMask[f] |= bb
		diagMask[r-f+7] |= bb
		adiagMask[r+f] |= bb
	}
}

func sign(x int) int {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}

func init() {
	// Precompute between masks: squares strictly between a and b
	// on the same rank, file, or diagonal.
	for a := 0; a < 64; a++ {
		ar, af := a/8, a%8
		for b := 0; b < 64; b++ {
			br, bf := b/8, b%8
			dr, df := br-ar, bf-af
			// Must be on same rank, file, or diagonal
			if dr == 0 && df == 0 {
				continue
			}
			if dr != 0 && df != 0 && abs(dr) != abs(df) {
				continue
			}
			sr, sf := sign(dr), sign(df)
			var mask core.Bitboard
			r, f := ar+sr, af+sf
			for r != br || f != bf {
				mask = mask.Set(core.Square(r*8 + f))
				r += sr
				f += sf
			}
			between[a][b] = mask
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func findKingSq(pos *position.Position, color core.Color) int {
	bb := pos.Board.Pieces(core.NewPiece(core.King, color))
	for sq := range bb.Squares() {
		return sq
	}
	return 0
}

// lineAttacks counts how many rooks/queens share a rank or file with the
// king, and how many bishops/queens share a diagonal, for one color
// attacking the given king square. Only counted if there is at most one
// piece between the attacker and the king.
func lineAttacks(pos *position.Position, color core.Color, kingSq int) int {
	kr, kf := kingSq/8, kingSq%8
	occupied := pos.Board.Occupied()
	count := 0

	// Rooks and queens on same rank or file
	rooks := pos.Board.Pieces(core.NewPiece(core.Rook, color)).
		Union(pos.Board.Pieces(core.NewPiece(core.Queen, color)))
	onLine := rankMask[kr].Union(fileMask[kf]).Intersection(rooks)
	for sq := range onLine.Squares() {
		if between[sq][kingSq].Intersection(occupied).Count() <= 1 {
			count++
		}
	}

	// Bishops and queens on same diagonal
	bishops := pos.Board.Pieces(core.NewPiece(core.Bishop, color)).
		Union(pos.Board.Pieces(core.NewPiece(core.Queen, color)))
	onDiag := diagMask[kr-kf+7].Union(adiagMask[kr+kf]).Intersection(bishops)
	for sq := range onDiag.Squares() {
		if between[sq][kingSq].Intersection(occupied).Count() <= 1 {
			count++
		}
	}

	return count
}

// Evaluate returns a score in centipawns from the active color's perspective.
// Positive means the active color is better.
func Evaluate(pos *position.Position) int {
	wkSq := findKingSq(pos, core.White)
	bkSq := findKingSq(pos, core.Black)

	score := 0
	for pt := core.PieceType(1); pt <= 5; pt++ {
		val := pieceValue[pt]
		adv := advanceBonus[pt]
		cb := centerBonus[pt]

		white := pos.Board.Pieces(core.NewPiece(pt, core.White))
		black := pos.Board.Pieces(core.NewPiece(pt, core.Black))

		ws := 0
		for sq := range white.Squares() {
			r := sq / 8
			if pt == core.Pawn {
				ws += val + pawnAdvance[r] + cb*(3-centerDist[sq])
			} else {
				ws += val + adv*r + cb*(3-centerDist[sq])
			}
		}

		bs := 0
		for sq := range black.Squares() {
			r := sq / 8
			if pt == core.Pawn {
				bs += val + pawnAdvance[7-r] + cb*(3-centerDist[sq])
			} else {
				bs += val + adv*(7-r) + cb*(3-centerDist[sq])
			}
		}

		if pos.ActiveColor == core.White {
			score += ws - bs
		} else {
			score += bs - ws
		}
	}

	// King zone coordination: count non-king pieces in 3x3 around enemy king.
	// Bonus scales with count^2 so only group attacks are rewarded.
	wNonKing := pos.Board.ColorPieces(core.White).Subtract(
		pos.Board.Pieces(core.NewPiece(core.King, core.White)))
	bNonKing := pos.Board.ColorPieces(core.Black).Subtract(
		pos.Board.Pieces(core.NewPiece(core.King, core.Black)))

	wZoneCount := kingZone[bkSq].Intersection(wNonKing).Count()
	bZoneCount := kingZone[wkSq].Intersection(bNonKing).Count()

	// Line alignment: rooks/queens on same rank/file, bishops/queens on same diagonal.
	wLineCount := lineAttacks(pos, core.White, bkSq)
	bLineCount := lineAttacks(pos, core.Black, wkSq)

	if pos.ActiveColor == core.White {
		score += wZoneCount*wZoneCount*kingZoneWeight - bZoneCount*bZoneCount*kingZoneWeight
		score += wLineCount*lineAlignWeight - bLineCount*lineAlignWeight
	} else {
		score += bZoneCount*bZoneCount*kingZoneWeight - wZoneCount*wZoneCount*kingZoneWeight
		score += bLineCount*lineAlignWeight - wLineCount*lineAlignWeight
	}

	return score
}
