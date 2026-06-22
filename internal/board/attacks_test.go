package board

import "testing"

func TestPrecalculated_KnightAttacks(t *testing.T) {
	// Initialize the precalculated knight attack tables manually if needed or rely on package init()
	d4, _ := SquareFromAlgebraic("d4")

	// Fetch the pre-calced uint64 bitmask for a Knight on d4
	attacks := GetKnightAttacks(d4)

	// Define our expected target squares
	expectedSquares := []string{"b3", "b5", "c2", "c6", "e2", "e6", "f3", "f5"}

	// Verify each expected square is flipped 'on' (1) in the bitmask
	for _, sqStr := range expectedSquares {
		sq, _ := SquareFromAlgebraic(sqStr)
		mask := uint64(1) << sq
		if attacks&mask == 0 {
			t.Errorf("Expected square %s to be in knight attacks for d4, but it was 0 in mask.", sqStr)
		}
	}
	
	// Verify an irrelevant square is 'off' (0)
	a1, _ := SquareFromAlgebraic("a1")
	if attacks&(1<<a1) != 0 {
		t.Errorf("Expected square a1 to NOT be in knight attacks for d4, but it was 1.")
	}	
}

func TestPrecalculated_KingAttacks(t *testing.T) {
	e4, _ := SquareFromAlgebraic("e4")
	attacks := GetKingAttacks(e4)

	// A King on e4 atacks in all surrounding squares: d3, d4, d5, e3, e5, f3, f4, f5
	expectedSquares := []string{"d3", "d4", "d5", "e3", "e5", "f3", "f4", "f5"}
		
	for _, sqStr := range expectedSquares {
		sq, _ := SquareFromAlgebraic(sqStr)
		if attacks & (1 << sq) == 0 {
			t.Errorf("Expected King on e4 to attack %s, but bit was 0", sqStr)
		}
	}

	// Boundary check: King on a1 shouldn't wrap around to the h file when moving left
	a1, _ := SquareFromAlgebraic("a1")
	a1Attacks := GetKingAttacks(a1)
	h1, _ := SquareFromAlgebraic("h1")
	if a1Attacks & (1 << h1) != 0 {
		t.Errorf("Guard failure: King on a1 wrapped around to attack h1")
	}
}

func TestPrecalculated_PawnAttacks(t *testing.T) {
	e4, _:= SquareFromAlgebraic("e4")

	//White pawns attack diagonally up (d5 and f5)
	wAttacks := GetPawnAttacks(e4, true) // True for white
	d5, _ := SquareFromAlgebraic("d5")
	f5, _ := SquareFromAlgebraic("f5")

	if (wAttacks & (1 << d5)) == 0 || (wAttacks & (1 << f5)) == 0 {
		t.Errorf("Expected white pawn on e4 to attack d5 and f5, but one or both were missing.")
	}

	//Black pawns attack diagonally down (d3 and f3)
	bAttacks := GetPawnAttacks(e4, false) // False for black
	d3, _ := SquareFromAlgebraic("d3")
	f3, _ := SquareFromAlgebraic("f3")

	if (bAttacks & (1 << d3)) == 0 || (bAttacks & (1 << f3)) == 0 {
		t.Errorf("Expected black pawn on e4 to attack d3 and f3, but one or both were missing.")
	}

	// Boundary check: Pawn on a4 should not wrap left to h file when attacking
	a4, _ := SquareFromAlgebraic("a4")
	a4Attacks := GetPawnAttacks(a4, true) // White pawn
	b5, _ := SquareFromAlgebraic("b5")
	h5, _ := SquareFromAlgebraic("h5")

	if (a4Attacks & (1 << b5)) == 0 {
		t.Errorf("Expected white pawn on a4 to attack b5, but it was missing.")
	}
	if a4Attacks & (1 << h5) != 0 {
		t.Errorf("Guard failure: White pawn on a4 wrapped around to attack h5")
	}
}
