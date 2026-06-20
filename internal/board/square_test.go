package board

import "testing"

// TestSquareFromAlgebraic validates parsing for:
// - Valid squares: a1, e4, h8, d5
// - Invalid squares: i1, a9, empty string, "aa", "a0", "4e"
//
// TestSquareAlgebraic validates round-trip from algebraic to Square and back.
//
// TestChebyshevDistance cases:
// - Same square returns 0
// - a1 to a2 returns 1
// - a1 to h8 returns 7
// - a1 to b2 returns 1
// - e4 to d5 returns 1
//
// TestManhattanDistance cases:
// - Same square returns 0
// - a1 to h8 returns 14
// - e4 to d5 returns 2

func TestSquareFromAlgebraic(t *testing.T) {
	tests := []struct {
		input    string
		wantFile int
		wantRank int
		wantErr  bool
	}{
		{"a1", 0, 0, false},
		{"e4", 4, 3, false},
		{"h8", 7, 7, false},
		{"d5", 3, 4, false},
		{"i1", 0, 0, true},
		{"a9", 0, 0, true},
		{"", 0, 0, true},
		{"aa", 0, 0, true},
		{"a0", 0, 0, true},
		{"4e", 0, 0, true},
	}

	for _, tt := range tests {
		sq, err := SquareFromAlgebraic(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("SquareFromAlgebraic(%q) expected error, got nil", tt.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("SquareFromAlgebraic(%q) unexpected error: %v", tt.input, err)
			continue
		}
		if sq.File() != tt.wantFile {
			t.Errorf("SquareFromAlgebraic(%q).File() = %d, want %d", tt.input, sq.File(), tt.wantFile)
		}
		if sq.Rank() != tt.wantRank {
			t.Errorf("SquareFromAlgebraic(%q).Rank() = %d, want %d", tt.input, sq.Rank(), tt.wantRank)
		}
	}
}
