package board

import "testing"

func TestMove_PackingAndUnpacking(t *testing.T) {
	from, _ := SquareFromAlgebraic("e2") // Index 12
	to, _ := SquareFromAlgebraic("e4")   // Index 28

	// Create a new Move instruction
	m := NewMove(from, to)

	// Verify the unpacking extraction layer
	if m.From() != from {
		t.Errorf("Move From() = %v, want %v", m.From(), from)
	}

	if m.To() != to {
		t.Errorf("Move To() = %v, want %v", m.To(), to)
	}
}

func TestMove_Flags(t *testing.T) {
	from, _ := SquareFromAlgebraic("f7")
	to, _ := SquareFromAlgebraic("f8")

	// Create a new Move instruction representing a pawn promotion to Queen
	m := NewMoveWithFlag(from, to, FlagPromoteQueen)

	if m.From() != from {
		t.Errorf("Flag Move From() = %v, want %v", m.From(), from)
	}

	if m.To() != to {
		t.Errorf("Flag Move To() = %v, want %v", m.To(), to)
	}

	// Verify the flag extraction matches
	if m.Flag() != FlagPromoteQueen {
		t.Errorf("Move Flag() = %v, want %v", m.Flag(), FlagPromoteQueen)
	}

	if !m.IsPromotion() {
		t.Errorf("Expected IsPromotion() to be true for Queen promotion")
	}
}
