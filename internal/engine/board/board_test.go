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
