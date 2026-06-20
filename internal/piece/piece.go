package piece

// Colour represents the player's side
type Colour uint8

const (
	White Colour = iota
	Black
)

// Type represents the chess piece type.
type Type uint8

const (
	None Type = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

// Symbol returns the standard single-character string representation (FEN notation)
// White is uppercase, Black is lowercase, empty/None is "".
func Symbol(pieceType Type, colour Colour) string {
	if pieceType == None {
		return "."
	}

	var symbol string
	switch pieceType {
	case Pawn:
		symbol = "P"
	case Knight:
		symbol = "N"
	case Bishop:
		symbol = "B"
	case Rook:
		symbol = "R"
	case Queen:
		symbol = "Q"
	case King:
		symbol = "K"
	default:
		return "."
	}

	if colour == Black {
		// Fast ASCII shift to lowercase: 'p' - 'P' = 32
		return string(symbol[0] + 32)
	}

	return symbol
}
