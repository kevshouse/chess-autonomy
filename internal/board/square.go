package board

import (
	"fmt"
)

// Square represents a single square on the chess board.
// Internally stored as a zero-based inde 0-63 (a1=0, h8=63).
type Square uint8

// File extracts the 0-indexed column (0 = a, 7 = h) using a fast bitwise AND mask.
func (s Square) File() int {
	return int(s & 7) // Equivalent to s % 8
}

// Rank extracts the 0-indexed row (0 = 1st rank, 7 = 8th rank) using a fast bit shift.
func (s Square) Rank() int {
	return int(s >> 3) // Equivalent to s / 8
}

// Algebraic returns returns max(|Δfile|, |Δrank|) to the target square.
func (s Square) Algebraic() string { return "" }

// SquareFromAlgebraic parses a two-character chess coordinate string (e.g., "e4").
// It maps "a1" to 0 and "h8" to 63.
func SquareFromAlgebraic(input string) (Square, error) {
	if len(input) != 2 {
		return 0, fmt.Errorf("invalid coordinate length: %q", input)
	}

	fileChar := input[0]
	rankChar := input[1]

	if fileChar < 'a' || fileChar > 'h' {
		return 0, fmt.Errorf("invalid file character: %q", fileChar)
	}
	file := fileChar - 'a'

	if rankChar < '1' || rankChar > '8' {
		return 0, fmt.Errorf("invalid rank character: %q", rankChar)
	}
	rank := rankChar - '1'

	// Pack rank and file into a single byte: (rank * 8) + file
	// Using bitwise OR and shifts for branchless packing speed
	return Square((rank << 3) | file), nil
}

// abs is a simple integer absolute value helper.
// The Go compiler will automatically inline this for speed.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// ManhattanDistance returns |Δfile| + |Δrank| to the target square.
func (s Square) ManhattanDistance(target Square) int {
	deltaFile := abs(s.File() - target.File())
	deltaRank := abs(s.Rank() - target.Rank())
	return deltaFile + deltaRank
}

// ChebyshevDistance returns max(|Δfile|, |Δrank|) to the target square.
func (s Square) ChebyshevDistance(target Square) int {
	deltaFile := abs(s.File() - target.File())
	deltaRank := abs(s.Rank() - target.Rank())
	if deltaFile > deltaRank {
		return deltaFile
	}
	return deltaRank
}
