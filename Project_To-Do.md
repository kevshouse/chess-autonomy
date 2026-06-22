# Project To-Do List
[x] Establish 64-bit One-Hot Bitboard Structures

[x] Implement standard FEN String Position Parsing

[x] Build the 8x8 ASCII Console Visual Debugger

[x] Implement packed 16-bit Move instructions

[x] Implement bitwise State Mutation (MakeMove)

[x] Quiet Moves

[x] Defensive Capture Masking

[x] Pawn Promotions

[x] Implement En Passant Target Square Tracking

[x] Track Active Castling Rights States

[x] Track Game Counters (Halfmove Clock and Fullmove Counter)

[x] Primitive Move Generation

	[x] Pre-calculate Knight Jump Bitboards

	[x] Pre-calculate King Step Bitboards

	[x] Map out Pawn attack/push arrays

[x] Design Sliding Attack Vectors (Rooks, Bishops, Queens)

Current Step:
- Implement the Move Generation Coordinator Loop (compiling the actual pseudo-legal move lists for an active position using our pre-calculated lookup tables).