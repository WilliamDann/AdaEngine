package perft

// holds Perft info
type PerftResults struct {
	Nodes            uint64
	Captures         uint64
	Enpassant        uint64
	Castles          uint64
	Promos           uint64
	Checks           uint64
	DiscoveredChecks uint64
	DoubleChecks     uint64
	Checkmates       uint64
}

func (r *PerftResults) Equ(other *PerftResults) bool {
	result := true

	result = result && (r.Nodes == other.Nodes)
	result = result && (r.Captures == other.Captures)
	result = result && (r.Enpassant == other.Enpassant)
	result = result && (r.Castles == other.Castles)
	result = result && (r.Promos == other.Promos)
	result = result && (r.Checks == other.Checks)
	result = result && (r.DiscoveredChecks == other.DiscoveredChecks)
	result = result && (r.DoubleChecks == other.DoubleChecks)
	result = result && (r.Checkmates == other.Checkmates)

	return result
}

func (r *PerftResults) Add(other *PerftResults) {
	// this could probably be done with reflect
	r.Nodes += other.Nodes
	r.Captures += other.Captures
	r.Enpassant += other.Enpassant
	r.Castles += other.Castles
	r.Promos += other.Promos
	r.Checks += other.Checks
	r.DiscoveredChecks += other.DiscoveredChecks
	r.DoubleChecks += other.DoubleChecks
	r.Checkmates += other.Checkmates
}

func NewPerftResultsNodes(nodes uint64) *PerftResults {
	var perft PerftResults
	perft.Nodes = nodes
	return &perft
}

func NewPerftResults() *PerftResults {
	return &PerftResults{}
}
