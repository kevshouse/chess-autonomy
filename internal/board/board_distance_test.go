package board

import (
	"testing"
)

func TestSquare_Distance(t *testing.T) {
	tests := []struct {
		name          string
		from          Square
		to            Square
		wantManhattan int
		wantChebyshev int
	}{
		// 1. Identical squares
		{"e4 to e4 (same square)", 28, 28, 0, 0},

		// 2. Straight line (e4 to e7) -> 3 ranks up, same file
		// e4 = 28, e7 = 52 -> Δfile = 0, Δrank = 3
		{"e4 to e7 (straight line)", 28, 52, 3, 3},

		// 3. Diagonal step (e4 to f5) -> 1 file right, 1 rank up
		// e4 = 28, f5 = 37 -> Δfile = 1, Δrank = 1
		{"e4 to f5 (diagonal move)", 28, 37, 2, 1},

		// 4. Large distance (a1 to h8) -> 7 files right, 7 ranks up
		// a1 = 0, h8 = 63 -> Δfile = 7, Δrank = 7
		{"a1 to h8 (max diagonal)", 0, 63, 14, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotManhattan := tt.from.ManhattanDistance(tt.to)
			if gotManhattan != tt.wantManhattan {
				t.Errorf("%s: ManhattanDistance() = %d, want %d", tt.name, gotManhattan, tt.wantManhattan)
			}

			gotChebyshev := tt.from.ChebyshevDistance(tt.to)
			if gotChebyshev != tt.wantChebyshev {
				t.Errorf("%s: ChebyshevDistance() = %d, want %d", tt.name, gotChebyshev, tt.wantChebyshev)
			}
		})
	}
}
