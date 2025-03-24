package game

import "testing"

func TestCoordString(t *testing.T) {
	values := map[Coord]string{
		{0, 0}: "a1",
		{2, 2}: "c3",
		{7, 0}: "h1",
		{0, 7}: "a8",
		{7, 7}: "h8",
	}

	for coord, expected := range values {
		got := coord.String()
		if got != expected {
			t.Errorf("Coord string incorrect got '%s' expected '%s'", got, expected)
		}
	}
}
