package game

import (
	"math"
)

type MoveGenerator struct {
	ruleSet  RuleSet
	position Position

	moves []Move

	activePieces []Coord
	activeRules  []MoveRule

	finished bool
}

// get all the pieces for the active side
func (mg MoveGenerator) getActivePicees() []Coord {
	var active []Coord

	for piece, coords := range mg.position.board.pieces {
		if piece.Color == mg.position.fen.ActiveColor {
			active = append(active, coords...)
		}
	}

	return active
}

func NewMoveGenerator(position Position, ruleSet RuleSet) *MoveGenerator {
	var mg MoveGenerator

	mg.moves = []Move{}
	mg.position = position

	mg.ruleSet = ruleSet

	// find all the pieces for the side to move
	mg.activePieces = mg.getActivePicees()
	if len(mg.activePieces) == 0 {
		mg.finished = true
		return &mg
	}

	// queue first set of move rules
	mg.activeRules = mg.ruleSet[mg.getCurrentPiece().Type]

	mg.finished = false

	return &mg
}

// get the Coord for the piece being worked on
func (mg MoveGenerator) getCurrentCoord() Coord {
	return mg.activePieces[len(mg.activePieces)-1]
}

// get the Piece value for the piece the generator is working on
func (mg MoveGenerator) getCurrentPiece() Piece {
	return mg.position.board.Get(mg.getCurrentCoord())
}

// get the piece the generator is currently working on
func (mg *MoveGenerator) nextPiece() {
	n := len(mg.activePieces)
	mg.activePieces = mg.activePieces[:n-1]

	if n-1 == 0 {
		mg.finished = true
		return // TODO finished
	}

	mg.activeRules = mg.ruleSet[mg.getCurrentPiece().Type]
}

// get rule the generator is currently working on
func (mg MoveGenerator) getCurrentRule() MoveRule {
	return mg.activeRules[len(mg.activeRules)-1]
}

// pop rule from active rules
func (mg *MoveGenerator) nextRule() {
	n := len(mg.activeRules)
	mg.activeRules = mg.activeRules[:n-1]

	if n-1 == 0 {
		mg.nextPiece()
		return
	}
}

// all promotion opts
func (mg *MoveGenerator) promoteMove(start Coord, cursor Coord, capture bool) {
	piece := mg.getCurrentPiece()
	for _, opt := range PromoteOpts {
		move := MoveBuilder().
			Piece(piece).
			From(start).
			To(cursor).
			Capture(capture).
			Promote(true, &opt).
			Castle(false, nil).
			Move()
		mg.moves = append(mg.moves, move)
	}
}

func (mg *MoveGenerator) applyRule(next MoveRule) {
	piece := mg.getCurrentPiece()
	start := mg.getCurrentCoord()
	cursor := next(piece, mg.position, start)

	for {
		// if the rule is complete, finish
		if cursor == nil {
			mg.nextRule()
			return
		}

		castle := piece.Type == King && math.Abs(float64(start.X)-float64(cursor.X)) > 1

		capture := !mg.position.board.IsEmpty(*cursor)
		promote := false

		if piece.Type == Pawn {
			if cursor.Y == 7 {
				promote = true
			}
		}

		if promote {
			mg.promoteMove(start, *cursor, capture)
		} else if castle {
			side := Queenside
			if cursor.X == 6 {
				side = Kingside
			}

			move := MoveBuilder().
				Piece(piece).
				From(start).
				To(*cursor).
				Capture(false).
				Promote(false, nil).
				Castle(true, &side).
				Move()
			mg.moves = append(mg.moves, move)
		} else {
			move := MoveBuilder().
				Piece(piece).
				From(start).
				To(*cursor).
				Capture(capture).
				Promote(false, nil).
				Castle(false, nil).
				Move()
			mg.moves = append(mg.moves, move)
		}

		// get next element
		cursor = next(piece, mg.position, *cursor)
	}
}

// gets a lost of psuedolegal moves from a ruleset
func (mg *MoveGenerator) applyRuleSet() []Move {
	for {
		if mg.finished {
			break
		}

		mg.applyRule(mg.getCurrentRule())
	}

	return mg.moves
}

func (mg *MoveGenerator) Generate() []Move {
	return mg.applyRuleSet()
}
