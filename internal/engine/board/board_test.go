package board

import (
	"testing"
)

func TestKnightMoves(t *testing.T) {
	tests := []struct {
		square   int
		expected []int
	}{
		{
			square:   28, // e4
			expected: []int{11, 13, 18, 22, 34, 38, 43, 45},
		},
		{
			square:   0, // a1 (corner)
			expected: []int{10, 17},
		},
		{
			square:   7, // h1 (corner)
			expected: []int{13, 22},
		},
		{
			square:   63, // h8 (corner)
			expected: []int{46, 53},
		},
	}

	for _, tt := range tests {
		moves := CalculateKnightMoves(tt.square)

		count := 0
		for i := range 64 {
			if moves&(1<<i) != 0 {
				count++
			}
		}

		if count != len(tt.expected) {
			t.Errorf("Knight on square %d: got %d moves, want %d", tt.square, count, len(tt.expected))
		}

		for _, exp := range tt.expected {
			if moves&(1<<exp) == 0 {
				t.Errorf("Knight on square %d: missing move to square %d", tt.square, exp)
			}
		}
	}
}

func TestRookMoves(t *testing.T) {
	tests := []struct {
		name      string
		square    int
		occupancy Bitboard
		expected  []int
	}{
		{
			name:      "d4 empty board",
			square:    27,
			occupancy: 0,
			expected:  []int{3, 11, 19, 35, 43, 51, 59, 24, 25, 26, 28, 29, 30, 31},
		},
		{
			name:      "d4 blocked north at d6",
			square:    27,
			occupancy: 1 << 43,
			expected:  []int{3, 11, 19, 35, 43, 24, 25, 26, 28, 29, 30, 31},
		},
		{
			name:      "a1 corner empty",
			square:    0,
			occupancy: 0,
			expected:  []int{1, 2, 3, 4, 5, 6, 7, 8, 16, 24, 32, 40, 48, 56},
		},
		{
			name:      "h8 corner empty",
			square:    63,
			occupancy: 0,
			expected:  []int{56, 57, 58, 59, 60, 61, 62, 7, 15, 23, 31, 39, 47, 55},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moves := CalculateRookMoves(tt.square, tt.occupancy)

			count := 0
			for i := range 64 {
				if moves&(1<<i) != 0 {
					count++
				}
			}

			if count != len(tt.expected) {
				t.Errorf("got %d moves, want %d", count, len(tt.expected))
			}

			for _, exp := range tt.expected {
				if moves&(1<<exp) == 0 {
					t.Errorf("missing move to square %d", exp)
				}
			}
		})
	}
}

func TestMove(t *testing.T) {
	tests := []struct {
		from  int
		to    int
		flags int
	}{
		{from: 12, to: 28, flags: 0},
		{from: 0, to: 63, flags: 15},
		{from: 27, to: 35, flags: 5},
	}

	for _, tt := range tests {
		move := NewMove(tt.from, tt.to, tt.flags)

		if move.From() != tt.from {
			t.Errorf("From: got %d, want %d", move.From(), tt.from)
		}
		if move.To() != tt.to {
			t.Errorf("To: got %d, want %d", move.To(), tt.to)
		}
		if move.Flags() != tt.flags {
			t.Errorf("Flags: got %d, want %d", move.Flags(), tt.flags)
		}
	}
}

func TestBishopMoves(t *testing.T) {
	tests := []struct {
		name      string
		square    int
		occupancy Bitboard
		expected  []int
	}{
		{
			name:      "d4 empty board",
			square:    27,
			occupancy: 0,
			expected:  []int{0, 9, 18, 36, 45, 54, 63, 6, 13, 20, 34, 41, 48},
		},
		{
			name:      "d4 blocked NE at f6",
			square:    27,
			occupancy: 1 << 45,
			expected:  []int{0, 9, 18, 36, 45, 6, 13, 20, 34, 41, 48},
		},
		{
			name:      "a1 corner empty",
			square:    0,
			occupancy: 0,
			expected:  []int{9, 18, 27, 36, 45, 54, 63},
		},
		{
			name:      "h8 corner empty",
			square:    63,
			occupancy: 0,
			expected:  []int{0, 9, 18, 27, 36, 45, 54},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moves := CalculateBishopMoves(tt.square, tt.occupancy)

			count := 0
			for i := range 64 {
				if moves&(1<<i) != 0 {
					count++
				}
			}

			if count != len(tt.expected) {
				t.Errorf("got %d moves, want %d", count, len(tt.expected))
			}

			for _, exp := range tt.expected {
				if moves&(1<<exp) == 0 {
					t.Errorf("missing move to square %d", exp)
				}
			}
		})
	}
}

func TestQueenMoves(t *testing.T) {
	tests := []struct {
		name      string
		square    int
		occupancy Bitboard
		expected  int
	}{
		{
			name:      "d4 empty board",
			square:    27,
			occupancy: 0,
			expected:  27,
		},
		{
			name:      "d4 with blocking",
			square:    27,
			occupancy: (1 << 35) | (1 << 19),
			expected:  22,
		},
		{
			name:      "a1 corner empty",
			square:    0,
			occupancy: 0,
			expected:  21,
		},
		{
			name:      "h8 corner empty",
			square:    63,
			occupancy: 0,
			expected:  21,
		},
		{
			name:      "e4 center",
			square:    28,
			occupancy: 0,
			expected:  27,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moves := CalculateQueenMoves(tt.square, tt.occupancy)

			count := 0
			for i := range 64 {
				if moves&(1<<i) != 0 {
					count++
				}
			}

			if count != tt.expected {
				t.Errorf("got %d moves, want %d", count, tt.expected)
			}
		})
	}
}

func TestKingMoves(t *testing.T) {
	tests := []struct {
		name     string
		square   int
		expected []int
	}{
		{
			name:     "e4 center",
			square:   28,
			expected: []int{19, 20, 21, 27, 29, 35, 36, 37},
		},
		{
			name:     "a1 corner",
			square:   0,
			expected: []int{1, 8, 9},
		},
		{
			name:     "h1 corner",
			square:   7,
			expected: []int{6, 14, 15},
		},
		{
			name:     "a8 corner",
			square:   56,
			expected: []int{48, 49, 57},
		},
		{
			name:     "h8 corner",
			square:   63,
			expected: []int{54, 55, 62},
		},
		{
			name:     "d4",
			square:   27,
			expected: []int{18, 19, 20, 26, 28, 34, 35, 36},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moves := CalculateKingMoves(tt.square)

			count := 0
			for i := range 64 {
				if moves&(1<<i) != 0 {
					count++
				}
			}

			if count != len(tt.expected) {
				t.Errorf("got %d moves, want %d", count, len(tt.expected))
			}

			for _, exp := range tt.expected {
				if moves&(1<<exp) == 0 {
					t.Errorf("missing move to square %d", exp)
				}
			}
		})
	}
}

func TestPawnMoves(t *testing.T) {
	tests := []struct {
		name            string
		square          int
		color           int
		occupancy       Bitboard
		enPassantSquare int
		expected        []int
	}{
		{
			name:            "white e2 empty board",
			square:          12,
			color:           White,
			occupancy:       0,
			enPassantSquare: -1,
			expected:        []int{19, 20, 21, 28},
		},
		{
			name:            "white e4 empty",
			square:          28,
			color:           White,
			occupancy:       0,
			enPassantSquare: -1,
			expected:        []int{35, 36, 37},
		},
		{
			name:            "white e2 e3 blocked",
			square:          12,
			color:           White,
			occupancy:       1 << 20,
			enPassantSquare: -1,
			expected:        []int{19, 21},
		},
		{
			name:            "black e7 empty",
			square:          52,
			color:           Black,
			occupancy:       0,
			enPassantSquare: -1,
			expected:        []int{43, 44, 45, 36},
		},
		{
			name:            "white d5 en passant e6",
			square:          35,
			color:           White,
			occupancy:       0,
			enPassantSquare: 44,
			expected:        []int{42, 43, 44},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moves := CalculatePawnMoves(tt.square, tt.color, tt.occupancy, tt.enPassantSquare)

			count := 0

			for i := range 64 {
				if moves&(1<<i) != 0 {
					count++
				}
			}

			if count != len(tt.expected) {
				t.Errorf("got %d moves, want %d", count, len(tt.expected))
			}

			for _, exp := range tt.expected {
				if moves&(1<<exp) == 0 {
					t.Errorf("missing move to square %d", exp)
				}
			}
		})
	}
}

func TestMakeMove(t *testing.T) {
	InitZobrist()

	b := InitBoard()

	b.Pieces[WhitePawn] |= (1 << 12)
	b.Occupancy[White] |= (1 << 12)
	b.Occupancy[All] |= (1 << 12)
	b.SideToMove = White
	b.EnPassant = -1
	b.CastleRights = 0
	b.Hash = b.CalculateHash()

	move := NewMove(12, 28, 0)
	b.MakeMove(move)

	if b.Pieces[WhitePawn]&(1<<12) != 0 {
		t.Error("Pawn still on e2")
	}

	if b.Pieces[WhitePawn]&(1<<28) == 0 {
		t.Error("Pawn not on e4")
	}

	if b.EnPassant != 20 {
		t.Errorf("En passant: got %d, want 20", b.EnPassant)
	}

	if b.SideToMove != Black {
		t.Error("Side to move should be Black")
	}
}

func TestMakeMoveWithCapture(t *testing.T) {
	InitZobrist()

	b := InitBoard()

	b.Pieces[WhitePawn] |= (1 << 28)
	b.Occupancy[White] |= (1 << 28)
	b.Occupancy[All] |= (1 << 28)

	b.Pieces[BlackPawn] |= (1 << 35)
	b.Occupancy[Black] |= (1 << 35)
	b.Occupancy[All] |= (1 << 35)

	b.SideToMove = White
	b.EnPassant = -1
	b.Hash = b.CalculateHash()

	move := NewMove(28, 35, 0)
	b.MakeMove(move)

	if b.Pieces[WhitePawn]&(1<<28) != 0 {
		t.Error("White pawn still on e4")
	}
	if b.Pieces[WhitePawn]&(1<<35) == 0 {
		t.Error("White pawn not on d5")
	}

	if b.Pieces[BlackPawn]&(1<<35) != 0 {
		t.Error("Black pawn still on d5")
	}

	if b.Occupancy[Black]&(1<<35) != 0 {
		t.Error("Black occupancy still on d5")
	}

	if b.HalfMove != 0 {
		t.Errorf("HalfMove should be 0 after capture, got %d", b.HalfMove)
	}
}

func TestCastleRightsRemoval(t *testing.T) {
	InitZobrist()

	b := InitBoard()

	b.Pieces[WhiteKing] |= (1 << 4)
	b.Occupancy[White] |= (1 << 4)
	b.Occupancy[All] |= (1 << 4)

	b.SideToMove = White
	b.CastleRights = WhiteKingSide | WhiteQueenSide | BlackKingSide | BlackQueenSide
	b.Hash = b.CalculateHash()

	move := NewMove(4, 12, 0)
	b.MakeMove(move)

	if b.CastleRights&WhiteKingSide != 0 {
		t.Error("White should lose kingside castle rights")
	}
	if b.CastleRights&WhiteQueenSide != 0 {
		t.Error("White should lose queenside castle rights")
	}

	if b.CastleRights&BlackKingSide == 0 {
		t.Error("Black should keep kingside castle rights")
	}
	if b.CastleRights&BlackQueenSide == 0 {
		t.Error("Black should keep queenside castle rights")
	}
}

func TestRookMoveCastleRights(t *testing.T) {
	InitZobrist()

	b := InitBoard()

	b.Pieces[WhiteRook] |= (1 << 0)
	b.Occupancy[White] |= (1 << 0)
	b.Occupancy[All] |= (1 << 0)

	b.SideToMove = White
	b.CastleRights = WhiteKingSide | WhiteQueenSide
	b.Hash = b.CalculateHash()

	move := NewMove(0, 8, 0)
	b.MakeMove(move)

	if b.CastleRights&WhiteQueenSide != 0 {
		t.Error("White should lose queenside castle rights")
	}
	if b.CastleRights&WhiteKingSide == 0 {
		t.Error("White should keep kingside castle rights")
	}
}

func TestRookCaptureRemovesCastleRights(t *testing.T) {
	InitZobrist()

	b := InitBoard()

	b.Pieces[WhiteQueen] |= (1 << 35)
	b.Occupancy[White] |= (1 << 35)
	b.Occupancy[All] |= (1 << 35)

	b.Pieces[BlackRook] |= (1 << 56)
	b.Occupancy[Black] |= (1 << 56)
	b.Occupancy[All] |= (1 << 56)

	b.SideToMove = White
	b.CastleRights = BlackKingSide | BlackQueenSide
	b.Hash = b.CalculateHash()

	move := NewMove(35, 56, 0)
	b.MakeMove(move)

	if b.CastleRights&BlackQueenSide != 0 {
		t.Error("Black should lose queenside castle rights after rook captured on a8")
	}
	if b.CastleRights&BlackKingSide == 0 {
		t.Error("Black should keep kingside castle rights")
	}
}

func TestEnPassantUpdatesOccupancyAll(t *testing.T) {
	InitZobrist()

	b := InitBoard()

	b.Pieces[WhitePawn] |= (1 << 35)
	b.Occupancy[White] |= (1 << 35)
	b.Occupancy[All] |= (1 << 35)

	b.Pieces[BlackPawn] |= (1 << 36)
	b.Occupancy[Black] |= (1 << 36)
	b.Occupancy[All] |= (1 << 36)

	b.SideToMove = White
	b.EnPassant = 44 // e6
	b.Hash = b.CalculateHash()

	move := NewMove(35, 44, FlagEnPassant)
	b.MakeMove(move)

	if b.Occupancy[All]&(1<<36) != 0 {
		t.Error("Occupancy[All] should not have piece on e5 after en passant")
	}
	if b.Occupancy[All]&(1<<44) == 0 {
		t.Error("Occupancy[All] should have piece on e6")
	}
}

func TestCastlingUpdatesOccupancyAll(t *testing.T) {
	InitZobrist()

	b := InitBoard()

	// White king on e1, rook on h1
	b.Pieces[WhiteKing] |= (1 << 4)
	b.Occupancy[White] |= (1 << 4)
	b.Occupancy[All] |= (1 << 4)

	b.Pieces[WhiteRook] |= (1 << 7)
	b.Occupancy[White] |= (1 << 7)
	b.Occupancy[All] |= (1 << 7)

	b.SideToMove = White
	b.CastleRights = WhiteKingSide
	b.Hash = b.CalculateHash()

	move := NewMove(4, 6, FlagCastle)
	b.MakeMove(move)

	if b.Occupancy[All]&(1<<4) != 0 {
		t.Error("Occupancy[All] should not have piece on e1")
	}
	if b.Occupancy[All]&(1<<7) != 0 {
		t.Error("Occupancy[All] should not have piece on h1")
	}
	if b.Occupancy[All]&(1<<6) == 0 {
		t.Error("Occupancy[All] should have piece on g1 (king)")
	}
	if b.Occupancy[All]&(1<<5) == 0 {
		t.Error("Occupancy[All] should have piece on f1 (rook)")
	}
}

func TestZobristHash(t *testing.T) {
	InitZobrist()

	b1 := &Board{}
	b2 := &Board{}

	hash1 := b1.CalculateHash()
	hash2 := b2.CalculateHash()

	if hash1 != hash2 {
		t.Error("Empty boards should have same hash")
	}

	b1.Pieces[WhitePawn] |= 1 << 12
	hash1 = b1.CalculateHash()

	if hash1 == hash2 {
		t.Error("Different positions should have different hashes")
	}
}
