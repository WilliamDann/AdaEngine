package core

import "testing"

func TestNewMove(t *testing.T) {
	m := NewMove(NewSquare(1, 4), NewSquare(3, 4)) // e2e4
	if m.From() != NewSquare(1, 4) {
		t.Errorf("From: got %v, want e2", m.From())
	}
	if m.To() != NewSquare(3, 4) {
		t.Errorf("To: got %v, want e4", m.To())
	}
	if m.MoveType() != MoveNormal {
		t.Errorf("Type: got %d, want Normal", m.MoveType())
	}
	if m.String() != "e2e4" {
		t.Errorf("String: got %q, want %q", m.String(), "e2e4")
	}
}

func TestPromotion(t *testing.T) {
	cases := []struct {
		piece PieceType
		want  string
	}{
		{Knight, "e7e8n"},
		{Bishop, "e7e8b"},
		{Rook, "e7e8r"},
		{Queen, "e7e8q"},
	}
	for _, tc := range cases {
		m := NewPromotion(NewSquare(6, 4), NewSquare(7, 4), tc.piece)
		if m.MoveType() != MovePromotion {
			t.Errorf("%s: type got %d, want Promotion", tc.want, m.MoveType())
		}
		if m.PromoPiece() != tc.piece {
			t.Errorf("%s: promo piece got %d, want %d", tc.want, m.PromoPiece(), tc.piece)
		}
		if m.String() != tc.want {
			t.Errorf("String: got %q, want %q", m.String(), tc.want)
		}
	}
}

func TestEnPassant(t *testing.T) {
	m := NewEnPassant(NewSquare(4, 4), NewSquare(5, 3)) // e5d6
	if m.From() != NewSquare(4, 4) {
		t.Errorf("From: got %v, want e5", m.From())
	}
	if m.To() != NewSquare(5, 3) {
		t.Errorf("To: got %v, want d6", m.To())
	}
	if m.MoveType() != MoveEnPassant {
		t.Errorf("Type: got %d, want EnPassant", m.MoveType())
	}
	if m.String() != "e5d6" {
		t.Errorf("String: got %q, want %q", m.String(), "e5d6")
	}
}

func TestCastling(t *testing.T) {
	m := NewCastling(NewSquare(0, 4), NewSquare(0, 6)) // e1g1
	if m.MoveType() != MoveCastling {
		t.Errorf("Type: got %d, want Castling", m.MoveType())
	}
	if m.String() != "e1g1" {
		t.Errorf("String: got %q, want %q", m.String(), "e1g1")
	}
}

func TestMoveFromTo_AllSquares(t *testing.T) {
	for from := Square(0); from < 64; from++ {
		for to := Square(0); to < 64; to++ {
			m := NewMove(from, to)
			if m.From() != from {
				t.Fatalf("from=%d to=%d: From() = %d", from, to, m.From())
			}
			if m.To() != to {
				t.Fatalf("from=%d to=%d: To() = %d", from, to, m.To())
			}
		}
	}
}

func TestMoveList(t *testing.T) {
	var ml MoveList
	if ml.Count() != 0 {
		t.Fatalf("empty list count: got %d", ml.Count())
	}

	m1 := NewMove(NewSquare(1, 4), NewSquare(3, 4))
	m2 := NewMove(NewSquare(0, 1), NewSquare(2, 2))
	ml.Add(m1)
	ml.Add(m2)

	if ml.Count() != 2 {
		t.Fatalf("count: got %d, want 2", ml.Count())
	}
	if ml.Get(0) != m1 {
		t.Errorf("Get(0): got %v, want %v", ml.Get(0), m1)
	}
	if ml.Get(1) != m2 {
		t.Errorf("Get(1): got %v, want %v", ml.Get(1), m2)
	}

	ml.Clear()
	if ml.Count() != 0 {
		t.Errorf("after Clear: count got %d", ml.Count())
	}
}

func TestMoveEquality(t *testing.T) {
	a := NewMove(NewSquare(1, 4), NewSquare(3, 4))
	b := NewMove(NewSquare(1, 4), NewSquare(3, 4))
	c := NewMove(NewSquare(1, 4), NewSquare(2, 4))

	if a != b {
		t.Error("identical moves should be equal")
	}
	if a == c {
		t.Error("different moves should not be equal")
	}
}
