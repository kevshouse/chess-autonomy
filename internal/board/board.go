package board

import (
	"fmt"
	"strconv"
	"strings"

	"chess-autonomy/internal/piece"
)

// Castling rights bitmask definitions
const (
	CastlingWK uint8 = 1 << 0 // White King-side (1)
	CastlingWQ uint8 = 1 << 1 // White Queen-side (2)
	CastlingBK uint8 = 1 << 2 // Black King-side (4)
	CastlingBQ uint8 = 1 << 3 // Black Queen-side (8)
)


// castlingSpoliationMAsks tracks which castling privileges are lost when a square is altered.
// By default, every square is 0x0F (preserves all rights)
var castlingSpoliationMasks [64]uint8

func init() {
	// Inits all squares to preserve rights
	for i := 0; i < 64; i++ {
		castlingSpoliationMasks[i] = 0x0F
	}
	
	// White King movement or capture strips White rights entirely (~3 = 12)
	castlingSpoliationMasks[4] = ^(CastlingWK | CastlingWQ)
	// White Rook movements or captures strip respective rights
	castlingSpoliationMasks[0] = ^CastlingWQ // a1 rook
	castlingSpoliationMasks[7] = ^CastlingWK // h1 rook

	// Black King movement or capture strips Black rights entirely (~12 = 3)
	castlingSpoliationMasks[7] = ^CastlingBK | CastlingBQ
	// Black Rook movements or captures strip respective rights
	castlingSpoliationMasks[56] = ^CastlingBQ // a8 rook
	castlingSpoliationMasks[63] = ^CastlingBK // h8 rook
}

// Board manages a 64-bit chess position layout using one-hot bitboards.
type Board struct {
	pieces  [7]uint64 // 0: None, 1: Pawn, 2: Knight, 3: Bishop, 4: Rook, 5: Queen, 6: King
	colours [2]uint64 // 0: White, 1: Black
	// epSquare records the active en passant target square (0-63), or NoSquare (64)
	epSquare Square
	// castlingRights tracks the available castling rights using a 4-bit mask
	castling uint8
	halfmoveClock int // Halfmove clock for the fifty-move rule
	fullmoveCounter int // Fullmove counter for the game
}

// NewBoard initializes a new empty chess board state.
func NewBoard() *Board {
	return &Board{
		epSquare: NoSquare,
		castling:        0,
		halfmoveClock:   0,
		fullmoveCounter: 1,
	}
}

func (b *Board) HalfmoveClock() int { return b.halfmoveClock }

func (b *Board) FullmoveCounter() int { return b.fullmoveCounter }

// CastlingRights returns the raw 4-bit mask representing the current castling rights for both sides.
func (b *Board) CastlingRights() uint8 {
	return b.castling
}

// EnPassantSquare returns the currently active en passant target square, or NoSquare if none is set.
func (b *Board) EnPassantSquare() Square {
	return b.epSquare
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

// LoadFEN parses a standard FEN string and populates the internal bitboards.
func (b *Board) LoadFEN(fen string) error {
	b.pieces = [7]uint64{}
	b.colours = [2]uint64{}
	b.epSquare = NoSquare
	b.castling = 0
	b.halfmoveClock = 0
	b.fullmoveCounter = 1

	parts := strings.Fields(fen)
	if len(parts) == 0 {
		return fmt.Errorf("empty FEN string")
	}

	boardPart := parts[0]
	ranks := strings.Split(boardPart, "/")
	if len(ranks) != 8 {
		return fmt.Errorf("invalid FEN: expected 8 ranks, found %d", len(ranks))
	}

	// [original 8-rank board parsing loop]
	for r := 7; r >= 0; r-- {
		rankStr := ranks[7-r] 
		file := 0
		for i := 0; i < len(rankStr); i++ {
			ch := rankStr[i]
			if ch >= '1' && ch <= '8' {
				file += int(ch - '0')
				continue
			}
			var pt piece.Type
			var c piece.Colour
			if ch >= 'A' && ch <= 'Z' { c = piece.White } else { c = piece.Black }
			normCh := ch
			if c == piece.Black { normCh = ch - 32 }
			switch normCh {
			case 'P': pt = piece.Pawn
			case 'N': pt = piece.Knight
			case 'B': pt = piece.Bishop
			case 'R': pt = piece.Rook
			case 'Q': pt = piece.Queen
			case 'K': pt = piece.King
			default: return fmt.Errorf("invalid FEN piece character: %q", ch)
			}
			sq := Square((r << 3) | file)
			b.PutPieceAt(sq, pt, c)
			file++
		}
	}

	// Parse Castling Rights (typically the 3rd field in FEN, e.g., "KQkq")
	if len(parts) >= 3 {
		castlingStr := parts[2]
		if castlingStr != "-" {
			for i := 0; i < len(castlingStr); i++ {
				switch castlingStr[i] {
				case 'K': b.castling |= CastlingWK
				case 'Q': b.castling |= CastlingWQ
				case 'k': b.castling |= CastlingBK
				case 'q': b.castling |= CastlingBQ
				}
			}
		}
	}

	// Parse En Passant Target
	if len(parts) >= 4 {
		epStr := parts[3]
		if epStr != "-" {
			sq, err := SquareFromAlgebraic(epStr)
			if err == nil { b.epSquare = sq }
		}
	}

	// Parse Halfmove Clock (Field 5)
	if len(parts) >= 5 {
		if val, err := strconv.Atoi(parts[4]); err == nil {
			b.halfmoveClock = val
		}		
	}

	// Parse Fullmove Counter (Field 6)
	if len(parts) >= 6 {
		if val, err := strconv.Atoi(parts[5]); err == nil {
			b.fullmoveCounter = val
		}
	}
	
	return nil
}


// MakeMove updates the internal bitboard layout by applying a packed Move instruction.
func (b *Board) MakeMove(m Move) {
	from := m.From()
	to := m.To()
	flag := m.Flag()

	// 1. Identify what piece is moving
	pt, colour, occupied := b.GetPieceAt(from)
	if !occupied {
		return 
	}
	// 1a. Check if this move constitutes a natural capture before resetting
	_, _, isCapture := b.GetPieceAt(to)

	// 2. Compute raw location bitmasks
	var fromMask uint64 = 1 << from
	var toMask uint64 = 1 << to
	
	// 3. DEFENSIVE LAYER: Unconditionally clear the destination bit
	inverseToMask := ^toMask

	opponentColour := piece.White
	if colour == piece.White {
		opponentColour = piece.Black
	}

	b.colours[opponentColour] &= inverseToMask
	for pType := piece.Pawn; pType <= piece.King; pType++ {
		b.pieces[pType] &= inverseToMask
	}

	// 4. STEPPING LAYER: Check for promotion transformations
	targetPieceType := pt
	if m.IsPromotion() {
		targetPieceType = m.PromotionType()
	}

	// Remove moving piece from source square
	b.pieces[pt] ^= fromMask
	b.colours[colour] ^= fromMask

	// Place final piece on the target square
	b.pieces[targetPieceType] |= toMask
	b.colours[colour] |= toMask

	// GAME CLOCK MANAGEMENT ///////////////
	// Halfmove Clock resets to 0 if a Pawn moves or a Capture occurs; otherwise, it increments.
	if pt == piece.Pawn || isCapture || flag == FlagEnPassant {
		b.halfmoveClock = 0
	} else {
		b.halfmoveClock++
	}

	// Fullmove counter increments immediately following any move completed by Black
	if colour == piece.Black {
		b.fullmoveCounter++
	}
	// GAME CLOCK MANAGEMENT ///////////////

	// 5. EN PASSANT TRACKING LAYER
	b.epSquare = NoSquare

	if flag == FlagDoublePawnPush {
		if colour == piece.White {
			b.epSquare = from + 8 
		} else {
			b.epSquare = from - 8 
		}
	}

// CASTLING PRIVILEGE REFRESH (Defensive Layer Bitwise Spoliation)
	b.castling &= castlingSpoliationMasks[from]
	b.castling &= castlingSpoliationMasks[to]
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
