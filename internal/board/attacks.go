package board

// Precalculated llookup table array for all 64 squares
var knightAttacksTable [64]uint64


// Clear file masks to catch edge wrapping vioations
const (
	notAFile uint64 = 0xFEFEFEFEFEFEFEFE // Everywhere EXCEPT File A
	notBFile uint64 = 0xFDFDFDFDFDFDFDFD // Everywhere EXCEPT File B
	notGFile uint64 = 0xBFBFBFBFBFBFBFBF // Everywhere EXCEPT File G
	notHFile uint64 = 0x7F7F7F7F7F7F7F7F // Everywhere EXCEPT File H
)

func init() {
	precalculateKnightAttacks()
	precalculateKingAttacks()
}

func GetKnightAttacks(sq Square) uint64 {
		return knightAttacksTable[sq]
}

func precalculateKnightAttacks() {
	for sq := 0; sq < 64; sq++ {
		b := uint64(1) << sq
		var attacks uint64

		// 2 Ranks Up, 1 File Left
		if (b & notAFile) != 0 { attacks |= (b << 16) >> 1 }
		// 2 Ranks Up, 1 File Right
		if (b & notHFile) != 0 { attacks |= (b << 16) << 1 }
		// 1 Rank Up, 2 Files Left
		if (b & notAFile & notBFile) != 0 { attacks |= (b << 8) >> 2 }
		// 1 Rank Up, 2 Files Right
		if (b & notGFile & notHFile) != 0 { attacks |= (b << 8) << 2 }

		// 2 Ranks Down, 1 File Left
		if (b & notAFile) != 0 { attacks |= (b >> 16) >> 1 }
		// 2 Ranks Down, 1 File Right
		if (b & notHFile) != 0 { attacks |= (b >> 16) << 1 }
		// 1 Rank Down, 2 Files Left
		if (b & notAFile & notBFile) != 0 { attacks |= (b >> 8) >> 2 }
		// 1 Rank Down, 2 Files Righ
		if (b & notGFile & notHFile) != 0 { attacks |= (b >> 8) << 2 }

		knightAttacksTable[sq] = attacks
	}
}

var kingAttacksTable [64]uint64

// GetKingAttacks returns the pre-computed attack bitmask for a King on a given square.
func GetKingAttacks(sq Square) uint64 {
	return kingAttacksTable[sq]
}

func precalculateKingAttacks() {
	for sq := 0; sq < 64; sq++ {
		b := uint64(1) << sq
		var attacks uint64

		// Vertical steps (No file checks required)
		attacks |= b << 8  // North
		attacks |= b >> 8  // South

		// Lateral & Diagonal steps moving Right (Safe from H-File wrap)
		if (b & notHFile) != 0 {
			attacks |= b << 1  // East
			attacks |= b << 9  // Northeast
			attacks |= b >> 7  // Southeast
		}

		// Lateral & Diagonal steps moving Left (Safe from A-File wrap)
		if (b & notAFile) != 0 {
			attacks |= b >> 1  // West
			attacks |= b << 7  // Northwest
			attacks |= b >> 9  // Southwest
		}

		kingAttacksTable[sq] = attacks
	}
}

var whitePawnAttacksTable [64]uint64
var blackPawnAttacksTable [64]uint64

// GetPawnAttacks returns the pre-computed diagonal attack bitmask for a pawn.
func GetPawnAttacks(sq Square, isWhite bool) uint64 {
	if isWhite {
		return whitePawnAttacksTable[sq]
	}
	return blackPawnAttacksTable[sq]
}

func init() {
	precalculateKnightAttacks()
	precalculateKingAttacks()
	precalculatePawnAttacks() // Add this
}

func precalculatePawnAttacks() {
	for sq := 0; sq < 64; sq++ {
		b := uint64(1) << sq

		// --- WHITE PAWN ATTACKS (Diagonally Up) ---
		var wAttacks uint64
		if (b & notAFile) != 0 { wAttacks |= b << 7 } // Northwest diagonal
		if (b & notHFile) != 0 { wAttacks |= b << 9 } // Northeast diagonal
		whitePawnAttacksTable[sq] = wAttacks

		// --- BLACK PAWN ATTACKS (Diagonally Down) ---
		var bAttacks uint64
		if (b & notAFile) != 0 { bAttacks |= b >> 9 } // Southwest diagonal
		if (b & notHFile) != 0 { bAttacks |= b >> 7 } // Southeast diagonal
		blackPawnAttacksTable[sq] = bAttacks
	}
}

