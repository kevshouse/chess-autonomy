# Chess Autonomy

A multi-agent chess system where pieces decide their own moves. Each piece has its own personality, limited world knowledge, and the ability to communicate with teammates. Pieces volunteer for or decline orders based on courage, aggression, loyalty, risk tolerance, and alignment with strategic goals.

Human input is restricted to adjusting weights and biases before each game—once play begins, the pieces are on their own.

## Design

### Piece Autonomy

Every piece on the board is an independent agent with five core traits:

| Trait          | Range  | Meaning                                          |
|----------------|--------|--------------------------------------------------|
| Courage        | [0, 1] | Willingness to accept risky orders               |
| Aggression     | [0, 1] | Preference for attacking over defending          |
| Loyalty        | [0, 1] | Adherence to received orders                     |
| Risk Tolerance | [0, 1] | Accepts moves with probability of capture        |
| Goal Alignment | [0, 1] | Alignment with the side's strategic objectives   |

### World Knowledge

Each piece perceives the board through a limited lens. Possible knowledge levels:

- **Self** — Only aware of its own position and type
- **Adjacent** — Sees the 8 neighboring squares
- **Line of Sight** — Sees along ranks, files, and diagonals (rook/bishop/queen lines)
- **Radius** — Sees all squares within a Chebyshev distance
- **Full** — Complete board awareness

A piece's legal move generation operates on its *perceived* board, not the true board. It may miss opportunities or walk into threats it cannot see.

### Communication

Pieces on the same side share information through a configurable network topology:

| Topology | Description                                      |
|----------|--------------------------------------------------|
| Full     | Every piece talks to every other piece           |
| Chain    | Pieces only communicate with immediate neighbors |
| Monarch  | All pieces report to and receive from the King   |

Messages carry intent (attack, defend, scout, block, retreat, support), a target square, and a confidence level. Trust between agents is weighted—a piece may discount information from a teammate it distrusts.

### Move Selection

Each turn follows a five-phase cycle:

1. **Perception** — Each agent observes the board through its knowledge mask
2. **Communication** — Agents share observations through the network
3. **Solicitation** — Strategic orders are presented; agents volunteer for moves they are willing to execute, each with a confidence score
4. **Arbitration** — A move is selected from volunteers using one of four strategies:
   - Highest confidence
   - Highest order priority × confidence
   - Weighted random
   - Consensus (most volunteered move)
5. **Execution** — The chosen move is applied to the board

If no agent volunteers, the most loyal piece is forced to move.

### Configuration

All parameters are set before the game via a JSON configuration file. No human intervention occurs during play.

```json
{
  "white_personalities": {
    "pawn":   { "courage": 0.4, "aggression": 0.3, "loyalty": 0.8, "risk_tolerance": 0.2, "goal_alignment": 0.6 },
    "knight": { "courage": 0.8, "aggression": 0.7, "loyalty": 0.5, "risk_tolerance": 0.7, "goal_alignment": 0.5 },
    "bishop": { "courage": 0.5, "aggression": 0.5, "loyalty": 0.6, "risk_tolerance": 0.4, "goal_alignment": 0.6 },
    "rook":   { "courage": 0.6, "aggression": 0.5, "loyalty": 0.7, "risk_tolerance": 0.5, "goal_alignment": 0.5 },
    "queen":  { "courage": 0.9, "aggression": 0.8, "loyalty": 0.4, "risk_tolerance": 0.6, "goal_alignment": 0.7 },
    "king":   { "courage": 0.1, "aggression": 0.1, "loyalty": 1.0, "risk_tolerance": 0.0, "goal_alignment": 1.0 }
  },
  "white_knowledge": {
    "pawn": "adjacent",
    "knight": "adjacent",
    "bishop": "line_of_sight",
    "rook": "line_of_sight",
    "queen": "full",
    "king": "radius_2"
  },
  "white_network_topology": "monarch",
  "white_orders": [
    { "type": "attack",  "target_square": "e4", "priority": 0.9 },
    { "type": "defend",  "target_square": "e1", "priority": 0.8 }
  ],
  "black_personalities": {},
  "black_knowledge": {},
  "black_network_topology": "full",
  "black_orders": [],
  "arbiter_strategy": "highest_confidence"
}
```
## Project Sructure

```chess-autonomy/
├── cmd/
│   └── chess-autonomy/
│       └── main.go
├── internal/
│   ├── board/        # Board representation, squares, moves, legal move generation
│   ├── piece/        # Agent definition, personality, knowledge masks
│   ├── comms/        # Communication networks and message passing
│   ├── orders/       # Strategic orders and volunteer solicitation
│   ├── arbiter/      # Move selection from volunteered candidates
│   ├── game/         # Game orchestration and turn management
│   └── config/       # Configuration loading, validation, and defaults
├── go.mod
└── README.md

```
## Usage

```bash
# Run with default configuration
go run ./cmd/chess-autonomy/

# Run with a custom configuration file
go run ./cmd/chess-autonomy/ -config ./my_config.json

# Run tests
go test ./internal/...
```

## Development Approach

This project follows test-driven development. Each package has a corresponding _test.go file with specifications that guide implementation.

Implementation order:

board/ — Square parsing, board representation, legal move generation

piece/ — Agent struct, personality, knowledge masks, perceived boards

comms/ — Network topologies, message passing, trust-weighted fusion

orders/ — Order types, volunteer solicitation

arbiter/ — Move selection strategies

config/ — JSON configuration loading and validation

game/ — Full game loop orchestration

cmd/ — CLI entry point
