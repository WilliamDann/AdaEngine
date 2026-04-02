package pgn

import (
	"github.com/WilliamDann/AdaEngine/ada-chess/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/movegen"
	"github.com/WilliamDann/AdaEngine/ada-chess/position"
)

// pieceChar returns the SAN character for a piece type.
// Pawns have no prefix in SAN.
var pieceChar = [7]byte{0, 0, 'N', 'B', 'R', 'Q', 'K'}

// SAN converts a move to Standard Algebraic Notation given the position
// before the move is made. The position must have the move as a legal move.
func SAN(pos *position.Position, m core.Move) string {
	// Castling
	if m.MoveType() == core.MoveCastling {
		if m.To().File() > m.From().File() {
			return "O-O"
		}
		return "O-O-O"
	}

	from := m.From()
	to := m.To()
	piece := pos.Board.Check(from)
	pt := piece.Type()
	isCapture := pos.Board.HasPiece(to) || m.MoveType() == core.MoveEnPassant

	var buf [10]byte
	n := 0

	if pt != core.Pawn {
		buf[n] = pieceChar[pt]
		n++

		// Disambiguation: check if another piece of the same type can reach the same square
		sameFile, sameRank, ambiguous := disambiguation(pos, m, pt)
		if ambiguous {
			if !sameFile {
				buf[n] = byte('a' + from.File())
				n++
			} else if !sameRank {
				buf[n] = byte('1' + from.Rank())
				n++
			} else {
				buf[n] = byte('a' + from.File())
				n++
				buf[n] = byte('1' + from.Rank())
				n++
			}
		}
	} else if isCapture {
		// Pawn captures include the departure file
		buf[n] = byte('a' + from.File())
		n++
	}

	if isCapture {
		buf[n] = 'x'
		n++
	}

	// Destination square
	buf[n] = byte('a' + to.File())
	n++
	buf[n] = byte('1' + to.Rank())
	n++

	// Promotion
	if m.MoveType() == core.MovePromotion {
		buf[n] = '='
		n++
		buf[n] = pieceChar[m.PromoPiece()]
		n++
	}

	// Check / checkmate suffix
	next := position.MakeMove(pos, m)
	legal := movegen.LegalMoves(next)
	if inCheck(next) {
		if legal.Count() == 0 {
			buf[n] = '#'
		} else {
			buf[n] = '+'
		}
		n++
	}

	return string(buf[:n])
}

// disambiguation checks whether another piece of the same type can move to
// the same target square. Returns whether any ambiguous piece shares the same
// file, same rank, and whether any ambiguity exists at all.
func disambiguation(pos *position.Position, m core.Move, pt core.PieceType) (sameFile, sameRank, ambiguous bool) {
	from := m.From()
	to := m.To()
	legal := movegen.LegalMoves(pos)

	for i := 0; i < legal.Count(); i++ {
		other := legal.Get(i)
		if other == m {
			continue
		}
		if other.To() != to {
			continue
		}
		otherPiece := pos.Board.Check(other.From())
		if otherPiece.Type() != pt {
			continue
		}

		ambiguous = true
		if other.From().File() == from.File() {
			sameFile = true
		}
		if other.From().Rank() == from.Rank() {
			sameRank = true
		}
	}
	return
}

// inCheck returns true if the side to move is in check.
func inCheck(pos *position.Position) bool {
	color := pos.ActiveColor
	enemy := color.Flip()
	kingBB := pos.Board.Pieces(core.NewPiece(core.King, color))
	for sq := range kingBB.Squares() {
		return !movegen.Attackers(pos, sq, enemy).Empty()
	}
	return false
}
