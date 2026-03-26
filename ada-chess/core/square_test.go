package core
import "testing"

// define bounds to test
var tests = []struct {
	name string
	rank int
	file int
	expect Square
	expectRank int
	expectFile int
	expectString string
}{
	{"subzero-file", -1, 0, InvalidSquare, -1, -1, InvalidSquareName},
	{"subzero-rank", 0, -1, InvalidSquare, -1, -1, InvalidSquareName},
	{"zero-square", 0, 0, 0, 0, 0, "a1"},
	{"mid-square", 3, 3, 27, 3, 3, "d4"},
	{"end-square", 7, 7, 63, 7, 7, "h8"},
	{"high-file", 8, 0, InvalidSquare, -1, -1, InvalidSquareName},
	{"high-rank", 0, 8, InvalidSquare, -1, -1, InvalidSquareName},
}

// test new square bounds
func TestNewSquareBounds(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func (t *testing.T) {
			got := NewSquare(tt.rank, tt.file)
			if got != tt.expect {
				t.Errorf("Square(%d, %d) = %d, expected %d", tt.rank, tt.file, got, tt.expect)
			}
		})
	}
}

// test the rank function's bounds
func TestRankFileBounds(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func (t *testing.T) {
			got := NewSquare(tt.rank, tt.file)
			rank := got.Rank()
			file := got.File()
			if rank != tt.expectRank || file != tt.expectFile {
				t.Errorf("Rank(%d) = %d, File(%d) = %d, expected %d and %d", got, rank, got, file, tt.expectRank, tt.expectFile)
			}
		})
	}
}

// test string bounds
func TestStringBounds(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func (t *testing.T) {
			got := NewSquare(tt.rank, tt.file).String()
			if got != tt.expectString {
				t.Errorf("NewSquare(%d, %d).String() = %s, expected %s", tt.rank, tt.file, got, tt.expectString)
			}
		})
	}
}
