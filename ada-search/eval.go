package search

import (
	"github.com/WilliamDann/AdaEngine/ada-chess/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/movegen"
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

// Pawn advancement bonus by rank. Negligible until very advanced.
//            rank: 0  1  2  3  4   5   6   7
var pawnAdvance = [8]int{0, 0, 0, 2, 5, 10, 25, 50}

// Center bonus per piece type, scaled by closeness to center (0-3).
var centerBonus = [7]int{
	0,  // None
	5,  // Pawn — control the center
	5,  // Knight — strongest center preference
	3,  // Bishop
	1,  // Rook
	2,  // Queen
	0,  // King
}

// King zone coordination weight — multiplied by count^2 so a lone piece
// near the king gets almost nothing but a group gets a large bonus.
const kingZoneWeight = 3


// Board region masks: queenside (files a-c), center (files d-e), kingside (files f-h).
var regionMask [3]core.Bitboard

const (
	regionQueenside = 0
	regionCenter    = 1
	regionKingside  = 2
)

// Weight for pieces in the broad region around the enemy king.
// Lighter than the tight king zone — rewards directing forces to the right side.
const kingRegionWeight = 2

// Precomputed tables, filled in init().
var (
	centerDist [64]int           // Chebyshev distance from center (0-3)
	kingZone   [64]core.Bitboard // 3x3 area around each square
)

func init() {
	// Board regions
	for sq := 0; sq < 64; sq++ {
		f := sq % 8
		bb := core.Bitboard(0).Set(core.Square(sq))
		switch {
		case f <= 2:
			regionMask[regionQueenside] |= bb
		case f <= 4:
			regionMask[regionCenter] |= bb
		default:
			regionMask[regionKingside] |= bb
		}
	}

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
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sqRegion(sq int) int {
	f := sq % 8
	if f <= 2 {
		return regionQueenside
	}
	if f <= 4 {
		return regionCenter
	}
	return regionKingside
}

func findKingSq(pos *position.Position, color core.Color) int {
	bb := pos.Board.Pieces(core.NewPiece(core.King, color))
	for sq := range bb.Squares() {
		return int(sq)
	}
	return 0
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
			r := int(sq) / 8
			if pt == core.Pawn {
				ws += val + pawnAdvance[r] + cb*(3-centerDist[sq])
			} else {
				ws += val + adv*r + cb*(3-centerDist[sq])
			}
		}

		bs := 0
		for sq := range black.Squares() {
			r := int(sq) / 8
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

	// Broad region pressure: count non-king pieces in the same board third
	// (queenside/center/kingside) as the enemy king.
	bkRegion := sqRegion(bkSq)
	wkRegion := sqRegion(wkSq)
	wRegionCount := regionMask[bkRegion].Intersection(wNonKing).Count()
	bRegionCount := regionMask[wkRegion].Intersection(bNonKing).Count()

	if pos.ActiveColor == core.White {
		score += wZoneCount*wZoneCount*kingZoneWeight - bZoneCount*bZoneCount*kingZoneWeight
		score += wRegionCount*kingRegionWeight - bRegionCount*kingRegionWeight
	} else {
		score += bZoneCount*bZoneCount*kingZoneWeight - wZoneCount*wZoneCount*kingZoneWeight
		score += bRegionCount*kingRegionWeight - wRegionCount*kingRegionWeight
	}

	// Space control: pawn attacks are strong permanent control,
	// piece attacks only count where enemy pawns don't cover.
	// Enemy outposts (squares they attack with pawns that we don't) are penalized.
	occupied := pos.Board.Occupied()

	var wPawnAtk, bPawnAtk core.Bitboard
	for sq := range pos.Board.Pieces(core.NewPiece(core.Pawn, core.White)).Squares() {
		wPawnAtk = wPawnAtk.Union(movegen.PawnAttacks(sq, core.White))
	}
	for sq := range pos.Board.Pieces(core.NewPiece(core.Pawn, core.Black)).Squares() {
		bPawnAtk = bPawnAtk.Union(movegen.PawnAttacks(sq, core.Black))
	}

	var wPieceAtk, bPieceAtk core.Bitboard
	for _, pt := range []core.PieceType{core.Knight, core.Bishop, core.Rook, core.Queen, core.King} {
		for sq := range pos.Board.Pieces(core.NewPiece(pt, core.White)).Squares() {
			switch pt {
			case core.Knight:
				wPieceAtk = wPieceAtk.Union(movegen.KnightMoves(sq))
			case core.Bishop:
				wPieceAtk = wPieceAtk.Union(movegen.BishopMoves(sq, occupied))
			case core.Rook:
				wPieceAtk = wPieceAtk.Union(movegen.RookMoves(sq, occupied))
			case core.Queen:
				wPieceAtk = wPieceAtk.Union(movegen.BishopMoves(sq, occupied).Union(movegen.RookMoves(sq, occupied)))
			case core.King:
				wPieceAtk = wPieceAtk.Union(movegen.KingMoves(sq))
			}
		}
		for sq := range pos.Board.Pieces(core.NewPiece(pt, core.Black)).Squares() {
			switch pt {
			case core.Knight:
				bPieceAtk = bPieceAtk.Union(movegen.KnightMoves(sq))
			case core.Bishop:
				bPieceAtk = bPieceAtk.Union(movegen.BishopMoves(sq, occupied))
			case core.Rook:
				bPieceAtk = bPieceAtk.Union(movegen.RookMoves(sq, occupied))
			case core.Queen:
				bPieceAtk = bPieceAtk.Union(movegen.BishopMoves(sq, occupied).Union(movegen.RookMoves(sq, occupied)))
			case core.King:
				bPieceAtk = bPieceAtk.Union(movegen.KingMoves(sq))
			}
		}
	}

	// Safe piece control: squares attacked by pieces but not defended by enemy pawns
	wSafe := wPieceAtk.Subtract(bPawnAtk).Count()
	bSafe := bPieceAtk.Subtract(wPawnAtk).Count()

	// Enemy outposts: squares enemy pawns attack but ours don't
	wOutposts := bPawnAtk.Subtract(wPawnAtk).Count() // squares black controls with pawns, we can't contest
	bOutposts := wPawnAtk.Subtract(bPawnAtk).Count() // squares white controls with pawns, black can't contest

	const (
		pawnControlWeight = 3
		safeControlWeight = 1
		outpostPenalty    = 2
	)

	if pos.ActiveColor == core.White {
		score += wPawnAtk.Count()*pawnControlWeight - bPawnAtk.Count()*pawnControlWeight
		score += wSafe*safeControlWeight - bSafe*safeControlWeight
		score -= wOutposts*outpostPenalty - bOutposts*outpostPenalty
	} else {
		score += bPawnAtk.Count()*pawnControlWeight - wPawnAtk.Count()*pawnControlWeight
		score += bSafe*safeControlWeight - wSafe*safeControlWeight
		score -= bOutposts*outpostPenalty - wOutposts*outpostPenalty
	}

	return score
}
