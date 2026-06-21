package board

import "chess-autonomy/internal/piece"

// Move represents a chess move packed tightly into a single 16-bit word.
// Bit layout:
// Bits 0-5:   Source Square (0-63)
// Bits 6-11:  Destination Square (0-63)
// Bits 12-15: Special flags (Promotions, Castling, En Passant)
type Move uint16

// Move flag definitions (values 0-15 fit perfectly in the 4 bits reserved for flags)
const (
	FlagQuiet            Move = 0  // Normal move, no special action
	FlagDoublePawnPush   Move = 1  // Pawn moves two squares forward
	FlagKingCastle       Move = 2  // King-side castling
	FlagQueenCastle      Move = 3  // Queen-side castling
	FlagCapture          Move = 4  // Capture move
	FlagEnPassant        Move = 5  // En passant capture
	FlagPromoteKnight    Move = 6  // Pawn promotion to Knight
	FlagPromoteBishop    Move = 7  // Pawn promotion to Bishop
	FlagPromoteRook      Move = 8  // Pawn promotion to Rook
	FlagPromoteQueen     Move = 9  // Pawn promotion to Queen
	FlagPromoteKnightCap Move = 12 // Promotion to Knight with capture
	FlagPromoteBishopCap Move = 13 // Promotion to Bishop with capture
	FlagPromoteRookCap   Move = 14 // Promotion to Rook with capture
	FlagPromoteQueenCap  Move = 15 // Promotion to Queen with capture
)

// NewMove creates a new Move instruction from the given source and destination squares.
func NewMove(from, to Square) Move {
	return Move(uint16(from) | (uint16(to) << 6))
}

// NewMoveWithFlag packs source, destination, and a 4-bit operational flag together.
func NewMoveWithFlag(from, to Square, flag Move) Move {
	return Move(uint16(from) | (uint16(to) << 6) | (uint16(flag) << 12))
}

// From extracts the source square using a 6-bit mask (0x3F = 00111111 in binary).
func (m Move) From() Square {
	return Square(uint16(m) & 0x3F)
}

// To extracts the destination square by shifting right 6 bits and applying the 6-bit mask.
func (m Move) To() Square {
	return Square((uint16(m) >> 6) & 0x3F)
}

// Flag extracts the special flags (top 4 bits) by shifting right 12 bits.
func (m Move) Flag() Move {
	return Move(uint16(m) >> 12)
}

// IsPromotion checks if the move is a promotion by examining the flag bits (Bit 3 of the flag is set).
func (m Move) IsPromotion() bool {
	return (m.Flag() & 0x8) != 0
}

// IsCapture checks if the flag matches standard or promotion captures.
func (m Move) IsCapture() bool {
	flag := m.Flag()
	return flag == FlagCapture || flag == FlagEnPassant || (flag >= FlagPromoteKnightCap)
}

// Promotion maps the 4-bit move flag to the corresponding piece type for promotions.
func (m Move) PromotionType() piece.Type {
	flag := m.Flag()
	switch flag {
	case FlagPromoteKnight, FlagPromoteKnightCap:
		return piece.Knight
	case FlagPromoteBishop, FlagPromoteBishopCap:
		return piece.Bishop
	case FlagPromoteRook, FlagPromoteRookCap:
		return piece.Rook
	case FlagPromoteQueen, FlagPromoteQueenCap:
		return piece.Queen
	default:
		return piece.None // Not a promotion move
	}
}