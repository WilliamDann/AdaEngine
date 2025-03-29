package game

type Side = string
type Direction = Coord

// direction defs
var (
	// caridinal directions
	North     Direction = Coord{0, 1}
	Northeast Direction = Coord{1, 1}
	East      Direction = Coord{1, 0}
	Southeast Direction = Coord{1, -1}
	South     Direction = Coord{0, -1}
	Southwest Direction = Coord{-1, -1}
	West      Direction = Coord{-1, 0}
	Northwest Direction = Coord{-1, 1}

	// knight move directions
	Knight_NNE Direction = Coord{1, -2}
	Knight_ENE Direction = Coord{2, -1}
	Knight_ESE Direction = Coord{2, 1}
	Knight_SSE Direction = Coord{1, 2}
	Knight_SSW Direction = Coord{-1, 2}
	Knight_WSW Direction = Coord{-2, 1}
	Knight_WNW Direction = Coord{-2, -1}
	Knight_NNW Direction = Coord{-1, -2}

	// board sides
	Kingside  Side = "kingside"
	Queenside Side = "queenside"
)

var PromoteOpts = []PieceType{Knight, Bishop, Rook, Queen}

// defines a single rule for a piece's movment
type MoveRule = func(Position, Coord) []Move

// set of functions that define how the pieces move
type RuleSet = map[PieceType][]MoveRule

func ApplyRuleSet(position *Position, ruleSet RuleSet) []Move {
	var moves []Move

	// apply the rule set to the pieces on the board
	for piece, coords := range position.board.pieces {
		if piece.Color == position.fen.ActiveColor {
			for _, start := range coords {
				for _, rule := range ruleSet[piece.Type] {
					moves = append(moves, rule(*position, start)...)
				}
			}
		}
	}

	return moves
}

// standard chess rule set
var StandardRules RuleSet = RuleSet{
	Pawn: {
		pawn,
	},

	Knight: {
		func(position Position, start Coord) []Move { return step(position, start, Knight_NNE) },
		func(position Position, start Coord) []Move { return step(position, start, Knight_ENE) },
		func(position Position, start Coord) []Move { return step(position, start, Knight_ESE) },
		func(position Position, start Coord) []Move { return step(position, start, Knight_SSE) },
		func(position Position, start Coord) []Move { return step(position, start, Knight_SSW) },
		func(position Position, start Coord) []Move { return step(position, start, Knight_WSW) },
		func(position Position, start Coord) []Move { return step(position, start, Knight_WNW) },
		func(position Position, start Coord) []Move { return step(position, start, Knight_NNW) },
	},

	Rook: {
		func(position Position, start Coord) []Move { return slider(position, start, North) },
		func(position Position, start Coord) []Move { return slider(position, start, East) },
		func(position Position, start Coord) []Move { return slider(position, start, South) },
		func(position Position, start Coord) []Move { return slider(position, start, West) },
	},

	Bishop: {
		func(position Position, start Coord) []Move { return slider(position, start, Northeast) },
		func(position Position, start Coord) []Move { return slider(position, start, Southeast) },
		func(position Position, start Coord) []Move { return slider(position, start, Southwest) },
		func(position Position, start Coord) []Move { return slider(position, start, Northwest) },
	},

	Queen: {
		func(position Position, start Coord) []Move { return slider(position, start, North) },
		func(position Position, start Coord) []Move { return slider(position, start, East) },
		func(position Position, start Coord) []Move { return slider(position, start, South) },
		func(position Position, start Coord) []Move { return slider(position, start, West) },
		func(position Position, start Coord) []Move { return slider(position, start, Northeast) },
		func(position Position, start Coord) []Move { return slider(position, start, Southeast) },
		func(position Position, start Coord) []Move { return slider(position, start, Southwest) },
		func(position Position, start Coord) []Move { return slider(position, start, Northwest) },
	},

	King: {
		func(position Position, start Coord) []Move { return step(position, start, North) },
		func(position Position, start Coord) []Move { return step(position, start, East) },
		func(position Position, start Coord) []Move { return step(position, start, South) },
		func(position Position, start Coord) []Move { return step(position, start, West) },
		func(position Position, start Coord) []Move { return step(position, start, Northeast) },
		func(position Position, start Coord) []Move { return step(position, start, Southeast) },
		func(position Position, start Coord) []Move { return step(position, start, Southwest) },
		func(position Position, start Coord) []Move { return step(position, start, Northwest) },
		castling,
	},
}

var CaptureOnlyRules RuleSet = RuleSet{
	Pawn: {
		func(position Position, start Coord) []Move { return Captures(pawn(position, start)) },
	},

	Knight: {
		func(position Position, start Coord) []Move { return Captures(step(position, start, Knight_NNE)) },
		func(position Position, start Coord) []Move { return Captures(step(position, start, Knight_ENE)) },
		func(position Position, start Coord) []Move { return Captures(step(position, start, Knight_ESE)) },
		func(position Position, start Coord) []Move { return Captures(step(position, start, Knight_SSE)) },
		func(position Position, start Coord) []Move { return Captures(step(position, start, Knight_SSW)) },
		func(position Position, start Coord) []Move { return Captures(step(position, start, Knight_WSW)) },
		func(position Position, start Coord) []Move { return Captures(step(position, start, Knight_WNW)) },
		func(position Position, start Coord) []Move { return Captures(step(position, start, Knight_NNW)) },
	},

	Rook: {
		func(position Position, start Coord) []Move { return Captures(slider(position, start, North)) },
		func(position Position, start Coord) []Move { return Captures(slider(position, start, East)) },
		func(position Position, start Coord) []Move { return Captures(slider(position, start, South)) },
		func(position Position, start Coord) []Move { return Captures(slider(position, start, West)) },
	},

	Bishop: {
		func(position Position, start Coord) []Move { return Captures(slider(position, start, Northeast)) },
		func(position Position, start Coord) []Move { return Captures(slider(position, start, Southeast)) },
		func(position Position, start Coord) []Move { return Captures(slider(position, start, Southwest)) },
		func(position Position, start Coord) []Move { return Captures(slider(position, start, Northwest)) },
	},

	Queen: {
		func(position Position, start Coord) []Move { return Captures(slider(position, start, North)) },
		func(position Position, start Coord) []Move { return Captures(slider(position, start, East)) },
		func(position Position, start Coord) []Move { return Captures(slider(position, start, South)) },
		func(position Position, start Coord) []Move { return Captures(slider(position, start, West)) },
		func(position Position, start Coord) []Move { return Captures(slider(position, start, Northeast)) },
		func(position Position, start Coord) []Move { return Captures(slider(position, start, Southeast)) },
		func(position Position, start Coord) []Move { return Captures(slider(position, start, Southwest)) },
		func(position Position, start Coord) []Move { return Captures(slider(position, start, Northwest)) },
	},

	King: {
		func(position Position, start Coord) []Move { return Captures(step(position, start, North)) },
		func(position Position, start Coord) []Move { return Captures(step(position, start, East)) },
		func(position Position, start Coord) []Move { return Captures(step(position, start, South)) },
		func(position Position, start Coord) []Move { return Captures(step(position, start, West)) },
		func(position Position, start Coord) []Move { return Captures(step(position, start, Northeast)) },
		func(position Position, start Coord) []Move { return Captures(step(position, start, Southeast)) },
		func(position Position, start Coord) []Move { return Captures(step(position, start, Southwest)) },
		func(position Position, start Coord) []Move { return Captures(step(position, start, Northwest)) },
	},
}
