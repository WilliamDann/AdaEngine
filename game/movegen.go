package game

// direction defs
var (
	// caridinal directions
	North     = Coord{0, 1}
	Northeast = Coord{1, 1}
	East      = Coord{1, 0}
	Southeast = Coord{1, -1}
	South     = Coord{0, -1}
	Southwest = Coord{-1, -1}
	West      = Coord{-1, 0}
	Northwest = Coord{-1, 1}

	// knight move directions
	Knight_NNE = Coord{1, -2}
	Knight_ENE = Coord{2, -1}
	Knight_ESE = Coord{2, 1}
	Knight_SSE = Coord{1, 2}
	Knight_SSW = Coord{-1, 2}
	Knight_WSW = Coord{-2, 1}
	Knight_WNW = Coord{-2, -1}
	Knight_NNW = Coord{-1, -2}
)

var PromoteOpts = []PieceType{Knight, Bishop, Rook, Queen}

// function that gives the next square given the current one
type MoveRule = func(Coord) *Coord

type MoveGenerator struct {
	position Position
	moves    []Move

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

func NewMoveGenerator(position Position) *MoveGenerator {
	var mg MoveGenerator

	mg.moves = []Move{}

	// for pawn moves to be simpler
	if !position.fen.ActiveColor {
		position.board = position.board.Flip()
	}
	mg.position = position

	// find all the pieces for the side to move
	mg.activePieces = mg.getActivePicees()
	if len(mg.activePieces) == 0 {
		mg.finished = true
		return &mg
	}

	// queue first set of move rules
	mg.activeRules = mg.rules()[mg.getCurrentPiece().Type]

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

	mg.activeRules = mg.rules()[mg.getCurrentPiece().Type]
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
			Move()
		mg.moves = append(mg.moves, move)
	}
}

func (mg *MoveGenerator) applyRule(next MoveRule) {
	piece := mg.getCurrentPiece()
	start := mg.getCurrentCoord()
	cursor := next(start)

	for {
		// if the rule is complete, finish
		if cursor == nil {
			mg.nextRule()
			return
		}

		capture := !mg.position.board.IsEmpty(*cursor)
		promote := false

		if piece.Type == Pawn {
			if cursor.Y == 7 {
				promote = true
			}
		}

		if promote {
			mg.promoteMove(start, *cursor, capture)
		} else {
			move := MoveBuilder().
				Piece(piece).
				From(start).
				To(*cursor).
				Capture(capture).
				Promote(false, nil).
				Move()
			mg.moves = append(mg.moves, move)
		}

		// get next element
		cursor = next(*cursor)
	}
}

func (mg *MoveGenerator) Generate() []Move {
	for {
		if mg.finished {
			break
		}

		mg.applyRule(mg.getCurrentRule())
	}

	// if board was flipped, adjust coords
	if mg.position.fen.ActiveColor == Black {
		for i, move := range mg.moves {
			mg.moves[i].From = Coord{7, 7}.Sub(move.From)
			mg.moves[i].To = Coord{7, 7}.Sub(move.To)

		}
	}

	return mg.moves
}

// moverule for a slider in a given direction
func (mg MoveGenerator) sliderRule(direction Coord) MoveRule {
	blocked := false
	return func(cursor Coord) *Coord {
		// if we've been blocked, this diection is complete
		if blocked {
			return nil
		}

		// take a step in the direction
		cursor = cursor.Add(direction)

		// if we're off board, direction is complete
		if !cursor.Valid() {
			blocked = true
			return nil
		}

		// if there is a piece in our way
		if !mg.position.board.IsEmpty(cursor) {
			blocker := mg.position.board.Get(cursor)
			blocked = true

			// if it's our piece we've reached the end of this direction
			if blocker.Color == mg.getCurrentPiece().Color {
				return nil
			}
		}

		return &cursor
	}
}

// moverule for a slider for a single step
func (mg MoveGenerator) stepRule(direction Coord) MoveRule {
	fired := false
	return func(cursor Coord) *Coord {
		// if this diection has already been examined, we're done
		if fired {
			return nil
		}
		fired = true

		// check for a capture
		cursor = cursor.Add(direction)

		// if the square is invalid, do not return the mvoe
		if !cursor.Valid() {
			return nil
		}

		if !mg.position.board.IsEmpty(cursor) {
			blocker := mg.position.board.Get(cursor)

			// if it's our piece we've reached the end of this direction
			if blocker.Color == mg.getCurrentPiece().Color {
				return nil
			}
		}

		// return the move
		return &cursor
	}
}

// moverule for single steps that MUST be a capture
func (mg MoveGenerator) pawnCaptureStep(direction Coord) MoveRule {
	fired := false
	return func(cursor Coord) *Coord {
		// if this diection has already been examined, we're done
		if fired {
			return nil
		}
		fired = true

		// check for a capture
		cursor = cursor.Add(direction)
		if !mg.position.board.IsEmpty(cursor) {
			blocker := mg.position.board.Get(cursor)

			// the only case were this is a valid move is capturing an opponent piece
			if blocker.Color != mg.getCurrentPiece().Color {
				return &cursor
			}
		}

		// not a capture
		return nil
	}
}

// generator for pawn moving a single step
//  this is different from the base single step as pawns cannot capture forward
func (mg MoveGenerator) pawnStep() MoveRule {
	fired := false
	return func(cursor Coord) *Coord {
		// if this diection has already been examined, we're done
		// or if the pawn is not in the 2nd rand, the move is invalid
		if fired {
			return nil
		}
		fired = true

		// check for a capture
		cursor = cursor.Add(North)
		if !mg.position.board.IsEmpty(cursor) {
			return nil
		}

		// not a capture
		return &cursor
	}
}

// generator for pawn moving up two steps
func (mg MoveGenerator) pawnTwoStep() MoveRule {
	fired := false
	return func(cursor Coord) *Coord {
		// if this diection has already been examined, we're done
		// or if the pawn is not in the 2nd rand, the move is invalid
		if fired || cursor.Y != 1 {
			return nil
		}
		fired = true

		// check for a capture
		cursor = cursor.Add(North).Add(North)
		if !mg.position.board.IsEmpty(cursor) {
			return nil
		}

		// not a capture
		return &cursor
	}
}

// generate move rules
func (mg MoveGenerator) rules() map[PieceType][]MoveRule {
	return map[PieceType][]MoveRule{
		Pawn: {
			mg.pawnStep(),
			mg.pawnTwoStep(),
			mg.pawnCaptureStep(Northeast),
			mg.pawnCaptureStep(Northwest),
		},

		Rook: {
			mg.sliderRule(North),
			mg.sliderRule(East),
			mg.sliderRule(South),
			mg.sliderRule(West),
		},

		Bishop: {
			mg.sliderRule(Northeast),
			mg.sliderRule(Southeast),
			mg.sliderRule(Southwest),
			mg.sliderRule(Northwest),
		},

		Queen: {
			mg.sliderRule(North),
			mg.sliderRule(East),
			mg.sliderRule(South),
			mg.sliderRule(West),
			mg.sliderRule(Northeast),
			mg.sliderRule(Southeast),
			mg.sliderRule(Southwest),
			mg.sliderRule(Northwest),
		},

		Knight: {
			mg.stepRule(Knight_NNE),
			mg.stepRule(Knight_ENE),
			mg.stepRule(Knight_ESE),
			mg.stepRule(Knight_SSE),
			mg.stepRule(Knight_SSW),
			mg.stepRule(Knight_WSW),
			mg.stepRule(Knight_WNW),
			mg.stepRule(Knight_NNW),
		},

		King: {
			mg.stepRule(North),
			mg.stepRule(East),
			mg.stepRule(South),
			mg.stepRule(West),
			mg.stepRule(Northeast),
			mg.stepRule(Southeast),
			mg.stepRule(Southwest),
			mg.stepRule(Northwest),
		}}
}
