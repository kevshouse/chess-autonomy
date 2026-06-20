package board

import (
	"testing"
)

func TestBitboard_SetAndHas(t *testing.T) {
	// We create a blank 64- bit board state.
	var bb uint64 = 0

	// We parse a square using our existing code
	e4, _ := SquareFromAlgebraic("e4") // index 28

	// Assert that e4 is currently empty
	if (bb & (1 << e4)) != 0 {
		t.Errorf("Expected e4 to be empty, but it was set")
	}

	// Set the e4 square
	bb |= (1 << e4)

	// Assert that e4 is now set
	if (bb & (1 << e4)) == 0 {
		t.Errorf("Expected e4 to be set (occupied state), but it was empty")
	}
}
