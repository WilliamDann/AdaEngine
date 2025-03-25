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

// function that gives the next square given the current one
type MoveRule = func(Piece, Position, Coord) *Coord

// set of functions that define how the pieces move
type RuleSet = map[PieceType][]MoveRule

// standard chess ruleset
var StandardRuleSet RuleSet = RuleSet{
	Pawn: {
		pawnStep(),
		pawnTwoStep(),
		pawnCaptureStep(Northeast),
		pawnCaptureStep(Northwest),
	},

	Rook: {
		sliderRule(North),
		sliderRule(East),
		sliderRule(South),
		sliderRule(West),
	},

	Bishop: {
		sliderRule(Northeast),
		sliderRule(Southeast),
		sliderRule(Southwest),
		sliderRule(Northwest),
	},

	Queen: {
		sliderRule(North),
		sliderRule(East),
		sliderRule(South),
		sliderRule(West),
		sliderRule(Northeast),
		sliderRule(Southeast),
		sliderRule(Southwest),
		sliderRule(Northwest),
	},

	Knight: {
		stepRule(Knight_NNE),
		stepRule(Knight_ENE),
		stepRule(Knight_ESE),
		stepRule(Knight_SSE),
		stepRule(Knight_SSW),
		stepRule(Knight_WSW),
		stepRule(Knight_WNW),
		stepRule(Knight_NNW),
	},

	King: {
		stepRule(North),
		stepRule(East),
		stepRule(South),
		stepRule(West),
		stepRule(Northeast),
		stepRule(Southeast),
		stepRule(Southwest),
		stepRule(Northwest),

		castle(Kingside),
		castle(Queenside),
	},
}
