package piece

import "testing"

func TestPiece_Symbol(t *testing.T) {
	tests := []struct {
		name      string
		pieceType Type
		colour    Colour
		want      string
	}{
		{"White Pawn", Pawn, White, "P"},
		{"White Knight", Knight, White, "N"},
		{"Black Queen", Queen, Black, "q"},
		{"Black King", King, Black, "k"},
		{"No Piece", None, White, "."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Symbol(tt.pieceType, tt.colour)
			if got != tt.want {
				t.Errorf("Symbol(%d, %d) = %q, want %q", tt.pieceType, tt.colour, got, tt.want)
			}
		})
	}
}
