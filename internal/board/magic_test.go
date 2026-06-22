package board

import "testing"

func TestMagic_RookEmptyBoard(t *testing.T) {
	e4, _ := SquareFromAlgebraic("e4")
	
	// With 0 occupancy (empty board), fetch Rook attacks from our table
	attacks := GetRookAttacks(e4, 0)

	// Expected squares along rank 4 and file e
	expectedSquares := []string{
		"e1", "e2", "e3", "e5", "e6", "e7", "e8",
		"a4", "b4", "c4", "d4", "f4", "g4", "h4",
	}

	for _, sqStr := range expectedSquares {
		sq, _ := SquareFromAlgebraic(sqStr)
		if (attacks & (1 << sq)) == 0 {
			t.Errorf("Expected Rook on e4 to attack %s on an empty board, but bit was 0", sqStr)
		}
	}
}

func TestMagic_BishopWithBlockers(t *testing.T) {
	d4, _ := SquareFromAlgebraic("d4")
	f6, _ := SquareFromAlgebraic("f6")
	c3, _ := SquareFromAlgebraic("c3")

	// Set up an occupancy mask containing pieces on f6 and c3
	var occupancy uint64 = (1 << f6) | (1 << c3)

	attacks := GetBishopAttacks(d4, occupancy)

	// Bishop should see c3 and f6 (the blockers), but NOT pass them to b2 or g7
	if (attacks & (1 << c3)) == 0 { t.Errorf("Expected Bishop to attack blocker on c3") }
	if (attacks & (1 << f6)) == 0 { t.Errorf("Expected Bishop to attack blocker on f6") }
	
	b2, _ := SquareFromAlgebraic("b2")
	g7, _ := SquareFromAlgebraic("g7")
	if (attacks & (1 << b2)) != 0 { t.Errorf("Error: Bishop went through blocker to b2") }
	if (attacks & (1 << g7)) != 0 { t.Errorf("Error: Bishop went through blocker to g7") }
}