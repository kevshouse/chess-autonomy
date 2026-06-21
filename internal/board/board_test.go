package board

import (
	"fmt"
	"strings"
	"testing"

	"chess-autonomy/internal/piece" // Connects the board domain to the piece domain
)

func TestSquare_FileAndRank(t *testing.T) {
	tests := []struct {
		name     string
		input    Square
		wantFile int
		wantRank int
	}{
		{"a1 (index 0)", 0, 0, 0},
		{"b1 (index 1)", 1, 1, 0},
		{"h1 (index 7)", 7, 7, 0},
		{"a2 (index 8)", 8, 0, 1},
		{"h8 (index 63)", 63, 7, 7},

		// 2. Mid-board strategic check: (e4 coordinate)
		// e4 is File 4 Rank 3 -> (3 << 3) | 4 = 24 + 4 = 28
		{"e4 (index 28)", 28, 4, 3},

		// 3. Additional specific coordinates to verify alignment
		{"d5 (index 35)", 35, 3, 4}, // (4 << 3) | 3 = 32 + 3 = 35
		{"c3 (index 18)", 18, 2, 2}, // (2 << 3) | 2 = 16 + 2 = 18
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFile := tt.input.File()
			if gotFile != tt.wantFile {
				t.Errorf("%s: File() = %d, want %d", tt.name, gotFile, tt.wantFile)
			}

			gotRank := tt.input.Rank()
			if gotRank != tt.wantRank {
				t.Errorf("%s: Rank() = %d, want %d", tt.name, gotRank, tt.wantRank)
			}
		})
	}
}

func TestBoard_PutAndGetPiece(t *testing.T) {
	// 1. Initialize a completely empty board struct
	b := NewBoard()

	e4, _ := SquareFromAlgebraic("e4")

	// 2. ASSERT: Verify that the square is initially empty
	pt, colour, occupied := b.GetPieceAt(e4)
	if occupied {
		t.Errorf("Expected e4 to be empty, but got occupied with piece type %d and colour %d", pt, colour)
	}

	// 3. Place a white Pawn on e4.
	b.PutPieceAt(e4, piece.Pawn, piece.White)

	// FIX: Re-query the board state to fetch the updated values into your variables!
	pt, colour, occupied = b.GetPieceAt(e4)

	if !occupied {
		t.Errorf("Expected e4 to be occupied after placing a piece, but it is still empty")
	}
	if pt != piece.Pawn {
		t.Errorf("GetPieceAt(e4) type = %v, want %v", pt, piece.Pawn)
	}
	if colour != piece.White {
		t.Errorf("GetPieceAt(e4) colour = %v, want %v", colour, piece.White)
	}
}

func TestBoard_LoadFEN(t *testing.T) {
	b := NewBoard()
	startFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	err := b.LoadFEN(startFEN)
	if err != nil {
		t.Fatalf("LoadFEN failed: %v", err)
	}

	// Verify White King in on e1 (index 4, rank 0, file 4)
	e1, _ := SquareFromAlgebraic("e1")
	pt, colour, occupied := b.GetPieceAt(e1)
	if !occupied || pt != piece.King || colour != piece.White {
		t.Errorf("Expected White King on e1, got type %v, colour %v, occupied %v", pt, colour, occupied)
	}

	// Verify Balck King in on e8 (index 60, rank 7, file 4)
	e8, _ := SquareFromAlgebraic("e8")
	pt, colour, occupied = b.GetPieceAt(e8)
	if !occupied || pt != piece.King || colour != piece.Black {
		t.Errorf("Expected Black King on e8, got type %v, colour %v, occupied %v", pt, colour, occupied)
	}

	// Verify an empty mid-board square like e4
	e4, _ := SquareFromAlgebraic("e4")
	_, _, occupied = b.GetPieceAt(e4)
	if occupied {
		t.Errorf("Expected square e4 to be empty in the starting position")
	}
}

func TestBoard_StringVisualizer(t *testing.T) {
	b := NewBoard()
	startFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	err := b.LoadFEN(startFEN)
	if err != nil {
		t.Fatalf("LoadFEN failed: %v", err)
	}

	got := b.String()
	// Define our expected visual 8x8 layout output string.
	want := strings.TrimSpace(`
	r n b q k b n r 
p p p p p p p p 
. . . . . . . . 
. . . . . . . . 
. . . . . . . . 
. . . . . . . . 
P P P P P P P P 
R N B Q K B N R 
`)
	// Clean up carriage returns or leading/trailing whitespace for an exact comparison.
	gotClean := strings.TrimSpace(got)
	if gotClean != want {
		t.Errorf("Board.String() output did not match expected grid.\nGot:\n%s\n\nWant:\n%s", gotClean, want)
	}
}


func TestBoard_MakeMove(t *testing.T) {
	b := NewBoard()
	startFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	err := b.LoadFEN(startFEN)
	if err != nil {
		t.Fatalf("LoadFEN failed: %v", err)
	}
	
	e2, _ := SquareFromAlgebraic("e2")
	e4, _ := SquareFromAlgebraic("e4")

	// 1. Create our packed quiet move instruction
	m := NewMove(e2, e4)

	// 2. Execut state mutation
	b.MakeMove(m)

	// Assert: Verify old square is now empty
	_, _, occupied := b.GetPieceAt(e2)
	if occupied {
		t.Errorf("Expected e2 to be empty after move, but it is still occupied")
	}

	// Assert: Verify new square has the white pawn
	pt, colour, occupiedE4 := b.GetPieceAt(e4)
	if !occupiedE4 {
		t.Errorf("Expected square e4 to be occupied after move")
	}
	if pt != piece.Pawn || colour != piece.White {
		t.Errorf("Expected White Pawn on e4, got type %v, colour %v", pt, colour)
	}
}

func TestBoard_MakeMove_DefensiveCapture(t *testing.T) {
	b := NewBoard()
	// Custom tactical setup string: White pawn on e4, Black pawn on d5
	// Rest of the board is empty for clean isolation
	testFEN := "8/8/8/3p4/4P3/8/8/8 w - - 0 1"
	if err := b.LoadFEN(testFEN); err != nil {
		t.Fatalf("Failed to load FEN: %v", err)
	}

	e4, _ := SquareFromAlgebraic("e4")
	d5, _ := SquareFromAlgebraic("d5")

	// Create a standard move instruction targeting the occupied square
	m := NewMove(e4, d5)

	// Execute state mutation
	b.MakeMove(m)

	// 1. ASSERT: Verify old square e4 is vacant
	_, _, occupiedE4 := b.GetPieceAt(e4)
	if occupiedE4 {
		t.Errorf("Expected source square e4 to be empty after capture")
	}

	// 2. ASSERT: Verify target square d5 no longer contains a Black Pawn
	pt, colour, occupiedD5 := b.GetPieceAt(d5)
	if !occupiedD5 {
		t.Fatalf("Expected destination square d5 to be occupied by the capturing piece")
	}
	if pt != piece.Pawn || colour != piece.White {
		t.Errorf("Expected White Pawn on d5, but got piece type %v and colour %v", pt, colour)
	}
	// This line prints the layout
	//t.Logf("\n%s", b.String())
}

func TestBoard_VisualDiagnostic(t *testing.T) {
	b := NewBoard()
	// Let's load the full starting position
	b.LoadFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	// Print a clean header and pass the board object directly to fmt.Println
	fmt.Println("\n--- CURRENT BITBOARD VISUALIZER STATE ---")
	fmt.Println(b)
	fmt.Println("-----------------------------------------")
}

func TestBoard_MakeMove_Promotion(t * testing.T) {
	b := NewBoard()
	// White pawn on a7 ready to step up, rest empty for clarity
	testFEN := "8/P7/8/8/8/8/8/8 w - - 0 1"
	if err := b.LoadFEN(testFEN); err != nil {
		t.Fatalf("Failed to load FEN: %v", err)
	}

	a7, _ := SquareFromAlgebraic("a7")
	a8, _ := SquareFromAlgebraic("a8")

	// Create a promotion move instruction using our flag constant (Pawn to Queen).
	m := NewMoveWithFlag(a7, a8, FlagPromoteQueen)

	// Execute the promotion move
	b.MakeMove(m)
	
	// 1. ASSERT: Verify old square a7 is now empty
	_, _, occupiedA7 := b.GetPieceAt(a7)
	if occupiedA7 {
		t.Errorf("Expected square a7 to be empty after promotion")
	}

	// 2. ASSERT: Verify new square a8 has the promoted Queen
	pt, colour, occupiedA8 := b.GetPieceAt(a8)
	if !occupiedA8 {
		t.Errorf("Expected square a8 to be occupied after promotion")
	}
	if pt != piece.Queen || colour != piece.White {
		t.Errorf("Expected White Queen on a8, got type %v, colour %v", pt, colour)
	}
	if pt != piece.Queen {
		t.Errorf("Expected piece type to be Queen after promotion, got %v", pt)
	}
	if colour != piece.White {
		t.Errorf("Expected piece colour to be White after promotion, got %v", colour)
	}
}

func TestBoard_EnPassantTracking(t *testing.T) {
	b := NewBoard()
	// Custom setup: White pawn on e5, Black pawn on d7 ready to move two squares
	testFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	if err := b.LoadFEN(testFEN); err != nil {
		t.Fatalf("Failed to load FEN: %v", err)
	}

	e2, _ := SquareFromAlgebraic("e2")
	e4, _ := SquareFromAlgebraic("e4")
	e3, _ := SquareFromAlgebraic("e3")

	// 1. Assert: Initially, there should be no active en passant target square
	if b.EnPassantSquare() != NoSquare {
		t.Errorf("Expected no active en passant square at start, got %v", b.EnPassantSquare())
	}
	
	// 2. Execute a two-square pawn advance from e2 to e4
	move := NewMoveWithFlag(e2, e4, FlagDoublePawnPush)
	b.MakeMove(move)

	// 3. Assert: The square skipped over (e3) must now be flagged as the target
	if b.EnPassantSquare() != e3 {
		t.Errorf("Expected en passant target square to be e3 after double pawn push, got %v", b.EnPassantSquare())
	}

	// 4. ACTION: Play a subsequent quiet move (e.g., Black playing a7 to a6)
	a7, _ := SquareFromAlgebraic("a7")
	a6, _ := SquareFromAlgebraic("a6")
	m2 := NewMove(a7, a6)
	b.MakeMove(m2)

	// 5. ASSERT: The en passant privilege expires and resets to NoSquare
	if b.EnPassantSquare() != NoSquare {
		t.Errorf("Expected en passant target square to reset to NoSquare after a non-pawn move, got %v", b.EnPassantSquare())
	}
}

func TestBoard_CastlingRightsTracking(t *testing.T) {
	b := NewBoard()
	startFen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	if err := b.LoadFEN(startFen); err != nil {
		t.Fatalf("Failed to load FEN: %v", err)
	}

	// 1. Assert: Initially, all castling rights should be available (4 rights)
	if b.CastlingRights() != (CastlingWK | CastlingWQ | CastlingBK | CastlingBQ) {
		t.Errorf("Expected full castling rights, got mask: %b", b.CastlingRights())
	}

	// 2. ACTION: Move the White King from e1 to e2 (breaks both White castling rights)
	e1, _ := SquareFromAlgebraic("e1")
	e2, _ := SquareFromAlgebraic("e2")
	move := NewMove(e1, e2)
	b.MakeMove(move)

	// 3. ASSERT: White's rights must be wiped out entirely, but Black's must remain untouched
	expectedMask := CastlingBK | CastlingBQ
	if b.CastlingRights() != expectedMask {
		t.Errorf("Expected only Black castling rights %b, got %b", expectedMask, b.CastlingRights())
	}
}

// Counters Test
func TestBoard_GameCountersTracking(t *testing.T) {
	b := NewBoard()
	// Fen specifies 4 halfmoves (sincs last pawn/capture) and 12 full moves total.
	testFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 4 12"
	if err := b.LoadFEN(testFEN); err != nil {
		t.Fatalf("Failed to load FEN: %v", err)
	}

	// 1. ASSERT: Verify parsing initialization matches FEN spec
	if b.HalfmoveClock() != 4 {
		t.Errorf("Expected halfmove clock to be 4, got %d", b.HalfmoveClock())
	}
	if b.FullmoveCounter() != 12 {
		t.Errorf("Expected fullmove counter to be 12, got %d", b.FullmoveCounter())
	}

	// 2. ACTION: Move a Knight quietly (increments halfmove, fullmove remains 12 because White is moving)
	g1, _ := SquareFromAlgebraic("g1")
	f3, _ := SquareFromAlgebraic("f3")
	mQuiet := NewMove(g1, f3)
	b.MakeMove(mQuiet)

	if b.HalfmoveClock() != 5 {
		t.Errorf("Expected halfmove clock to increment to 5 on quiet piece move, got %d", b.HalfmoveClock())
	}

	// 3. ACTION: Reset check via Pawn Push (Should wipe halfmove clock to 0)
	e2, _ := SquareFromAlgebraic("e2")
	e4, _ := SquareFromAlgebraic("e4")
	mPawn := NewMoveWithFlag(e2, e4, FlagDoublePawnPush)
	b.MakeMove(mPawn)

	if b.HalfmoveClock() != 0 {
		t.Errorf("Expected halfmove clock to reset to 0 after pawn move, got %d", b.HalfmoveClock())
	}
}
