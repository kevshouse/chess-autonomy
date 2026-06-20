package board

import (
	"fmt"
	"strings"

	"chess-autonomy/internal/piece"
)

// Board manages a 64-bit chess position layout using one-hot bitboards.
type Board struct {
	pieces  [7]uint64 // 0: None, 1: Pawn, 2: Knight, 3: Bishop, 4: Rook, 5: Queen, 6: King
	colours [2]uint64 // 0: White, 1: Black
}

// NewBoard initializes a new empty chess board state.
func NewBoard() *Board {
	return &Board{}
}

// PutPieceAt places a piece of a given type and colour at the specified square.
func (b *Board) PutPieceAt(sq Square, pt piece.Type, c piece.Colour) {
	// Clear any existing piece at the square
	var mask uint64 = 1 << sq

	b.pieces[pt] |= mask
	b.colours[c] |= mask
}

// GetPieceAt retrieves the piece type and colour at the specified square.
func (b *Board) GetPieceAt(sq Square) (piece.Type, piece.Colour, bool) {
	var mask uint64 = 1 << sq

	var occupiedByColour bool
	var c piece.Colour

	if (b.colours[piece.White] & mask) != 0 {
		occupiedByColour = true
		c = piece.White
	} else if (b.colours[piece.Black] & mask) != 0 {
		occupiedByColour = true
		c = piece.Black
	}

	if !occupiedByColour {
		return piece.None, piece.White, false // No piece present
	}

	for pt := piece.Pawn; pt <= piece.King; pt++ {
		if (b.pieces[pt] & mask) != 0 {
			return pt, c, true
		}
	}

	return piece.None, piece.White, false // Should not reach here if occupiedByColour is true
}

// LoadFEN loads a chess position from a FEN string into the board.
func (b *Board) LoadFEN(fen string) error {
	// Clear any existing bitboard data before reloading
	b.pieces = [7]uint64{}
	b.colours = [2]uint64{}

	parts := strings.Fields(fen)
	if len(parts) < 1 {
		return fmt.Errorf("invalid FEN: missing board layout")
	}

	boardPart := parts[0]
	ranks := strings.Split(boardPart, "/")
	if len(ranks) != 8 {
		return fmt.Errorf("invalid FEN: expected 8 ranks, found %d", len(ranks))
	}

	// FEN starts from rank 8 (index 7) down to rank 1 (index 0)
	for r := 7; r >= 0; r-- {
		rankStr := ranks[7-r] // Map FEN row to our internal rank index
		file := 0

		for i := 0; i < len(rankStr); i++ {
			ch := rankStr[i]
			if ch >= '1' && ch <= '8' {
				// Empty squares
				emptySquares := int(ch - '0')
				file += emptySquares
				continue
			}

			// Map the character to piece type and colour
			var pt piece.Type
			var c piece.Colour

			// Uppercase is White, lowercase is Black
			if ch >= 'A' && ch <= 'Z' {
				c = piece.White
			} else {
				c = piece.Black
			}

			// Determine piece type character
			normCh := ch
			if c == piece.Black {
				normCh = ch - 32 // Convert to uppercase for mapping
			}

			switch normCh {
			case 'P':
				pt = piece.Pawn
			case 'N':
				pt = piece.Knight
			case 'B':
				pt = piece.Bishop
			case 'R':
				pt = piece.Rook
			case 'Q':
				pt = piece.Queen
			case 'K':
				pt = piece.King
			default:
				return fmt.Errorf("invalid FEN piece character: %q", ch)
			}

			// Calculate the square index based on rank and file
			sq := Square((r << 3) | file)
			b.PutPieceAt(sq, pt, c)
			file++
		}
	}

	return nil
}

// String returns a formatted 8x8 ASCII representation grid of the current board state.
// Satisfies the fmt.Stringer interface for clean visual console logging.
func (b *Board) String() string {
	var sb strings.Builder

	// Chess boards print from the 8th rank down to the 1st rank.
	for r := 7; r >= 0; r-- {
		for f := 0; f < 8; f++ {
			// Construct the square index cleanly using our branchless bitfield geometry
			sq := Square((r << 3) | f)
			pt, colour, occupied := b.GetPieceAt(sq)
			if !occupied {
				sb.WriteString(". ")
			} else {
				sb.WriteString(piece.Symbol(pt, colour) + " ")
			}
		}
		// Append a newline character at the end of every rank line
		sb.WriteString("\n")
	}
	return sb.String()
}

// MakeMove updates the internal bitboard layout by applying a packed Move instruction.
// Currently handles basic quiet moves.
func (b *Board) MakeMove(m Move) {
	fromSq := m.From()
	toSq := m.To()
	
	pt, colour, occupied := b.GetPieceAt(fromSq)
	if !occupied {
		return // No piece to move; in a real implementation, you might want to return an error
	}

	// Generate our bitmasks for the source and target locations.
	var fromMask uint64 = 1 << fromSq
	var toMask uint64 = 1 << toSq

	// Mutate the Piece Type bitboard
	b.pieces[pt] ^= fromMask // Clear the piece from the source square
	b.pieces[pt] |= toMask   // Set the piece at the target square

	// Mutate the Colour bitboard
	b.colours[colour] ^= fromMask // Clear the colour from the source square
	b.colours[colour] ^= toMask   // Set the colour at the target square

	// TODO: Handle Pawn Promotion flags (intercept targetPieceType calculation)

	// 4. STEPPING LAYER: Execute standard XOR relocation for the moving piece
	// Clear source square bit
	// Set destination square bit
	// Clear source square bit
	// Set destination square bit
}

