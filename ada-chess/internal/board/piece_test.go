package board
import "testing"

func TestPieceBounds(t *testing.T) {
	tests := []struct {
		name  string
		ptype PieceType
		color Color
	}{
		{"white pawn", Pawn, White},
		{"black pawn", Pawn, Black},
		{"white knight", Knight, White},
		{"black knight", Knight, Black},
		{"white bishop", Bishop, White},
		{"black bishop", Bishop, Black},
		{"white rook", Rook, White},
		{"black rook", Rook, Black},
		{"white queen", Queen, White},
		{"black queen", Queen, Black},
		{"white king", King, White},
		{"black king", King, Black},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPiece(tt.ptype, tt.color)
			if p.Type() != tt.ptype {
				t.Errorf("type: got %d want %d", p.Type(), tt.ptype)
			}
			if p.Color() != tt.color {
				t.Errorf("color: got %d want %d", p.Color(), tt.color)
			}
		})
	}
}
