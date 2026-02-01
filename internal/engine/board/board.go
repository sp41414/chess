package board

import (
	"math/bits"
	"math/rand"
)

type Bitboard uint64

type Board struct {
	Pieces [12]Bitboard
	// [White, Black, All]
	Occupancy [3]Bitboard

	// Game State
	// 0: White, 1: Black
	SideToMove int
	// Index 0-63, -1 if no target
	EnPassant int
	// 4-bit mask: 0001 (WK) 0010 (WQ) 0100 (BK), 1000 (BQ)
	CastleRights int
	// 50-move rule
	HalfMove int
	// Increment after every black move
	FullMove int
	// Zobrist hashing for transposition tables
	Hash uint64
}

// Constants to fulfill the Board datastructure
const (
	WhitePawn = iota
	WhiteKnight
	WhiteBishop
	WhiteRook
	WhiteQueen
	WhiteKing
	BlackPawn
	BlackKnight
	BlackBishop
	BlackRook
	BlackQueen
	BlackKing
)

const (
	White = 0
	Black = 1
	All   = 2
)

const (
	WhiteKingSide  = 1 << 0 // 0001
	WhiteQueenSide = 1 << 1 // 0010
	BlackKingSide  = 1 << 2 // 0100
	BlackQueenSide = 1 << 3 // 1000
)

var (
	zobristPieces     [12][64]uint64
	zobristCastle     [16]uint64
	zobristEnPassant  [64]uint64
	zobristSideToMove uint64
)

// Computes a random Uint64 and fills the zobrist constants with it.
// The seed is constant for reproducability.
func InitZobrist() {
	rng := rand.New(rand.NewSource(66666))

	for piece := range 12 {
		for square := range 64 {
			zobristPieces[piece][square] = rng.Uint64()
		}
	}

	for i := range 16 {
		zobristCastle[i] = rng.Uint64()
	}

	for square := range 64 {
		zobristEnPassant[square] = rng.Uint64()
	}

	zobristSideToMove = rng.Uint64()
}

// Calculates the hash for Zobrist using every piece and game state
func (b *Board) CalculateHash() uint64 {
	hash := uint64(0)

	for piece := range 12 {
		bb := b.Pieces[piece]
		for bb != 0 {
			square := bits.TrailingZeros64(uint64(bb))
			bb &= bb - 1
			hash ^= zobristPieces[piece][square]
		}
	}

	hash ^= zobristCastle[b.CastleRights]
	if b.EnPassant != -1 {
		hash ^= zobristEnPassant[b.EnPassant]
	}

	if b.SideToMove == Black {
		hash ^= zobristSideToMove
	}

	return hash
}

// Initializes new empty board and calculates Zobrist hash
func InitBoard() *Board {
	b := &Board{}
	b.Hash = b.CalculateHash()
	return b
}
