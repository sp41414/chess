package engine

import (
	"fmt"
	"math/bits"
	"strconv"
	"strings"
)

type Bitboard uint64

type Board struct {
	// Every piece type for both colors
	Pieces [12]Bitboard
	// Occupancy[0]: White, Occupancy[1]: Black, Occupancy[2]: All
	Occupancy [3]Bitboard
	// Undo stack
	StateHistory []Undo
	// Moves made in the current game
	MoveHistory []Move
	// Index 0-63, -1 if no target
	EnPassant int
	// 0: White, 1: Black
	SideToMove int
	// 50-move rule
	HalfMove int
	// Increment after every black move
	FullMove int
	// 4-bit Mask: 0001 (WK) 0010 (WQ) 0100 (BK), 1000 (BQ)
	CastleRights int
}

// Indexes of Board.Pieces Bitboard array
const (
	WhitePawn int = iota
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

// 4-bit masks for Board.CastleRights
const (
	WhiteKingSide  int = 1 << 0 // 0001
	WhiteQueenSide     = 1 << 1 // 0010
	BlackKingSide      = 1 << 2 // 0100
	BlackQueenSide     = 1 << 3 // 1000
)

// White and Black And All, 0 and 1 and 2, for Board.SideToMove and Board.Occupancy
const (
	White int = iota
	Black
	All
)

// Helper functions for Bitboard

// Occupied returns true if the square is occupied
func (b Bitboard) Occupied(sq int) bool {
	return b&(1<<sq) != 0
}

// Set sets a bit at square to 1
func (b *Bitboard) Set(sq int) {
	*b |= 1 << sq
}

// Clear clears a bit at square to 0
func (b *Bitboard) Clear(sq int) {
	*b &^= (1 << sq)
}

// Toggle toggles a bit at square to its opposite
func (b *Bitboard) Toggle(sq int) {
	*b ^= (1 << sq)
}

// Count returns the number of bits set in the Bitboard
func (b Bitboard) Count() int {
	return bits.OnesCount64(uint64(b))
}

// Get least significant bit, or first piece, of bitboard
func (b Bitboard) LSB() int {
	return bits.TrailingZeros64(uint64(b))
}

func (b *Bitboard) PopLSB() int {
	sq := bits.TrailingZeros64(uint64(*b))
	*b &= *b - 1
	return sq
}

func (b Bitboard) String() string {
	s := strconv.FormatUint(uint64(b), 2)
	return strings.Repeat("0", 64-len(s)) + s
}

// ParseFEN takes a FEN string with format
// "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
// and fills the Board struct with the data parsed
func (b *Board) ParseFEN(fen string) error {
	parts := strings.Split(fen, " ")
	if len(parts) != 6 {
		return fmt.Errorf("invalid FEN: expected 6 parts got %d", len(parts))
	}

	switch parts[1] {
	case "w":
		b.SideToMove = White
	case "b":
		b.SideToMove = Black
	default:
		return fmt.Errorf("invalid FEN: expected side to move to be w or b, got %s", parts[1])
	}

	ranks := strings.Split(parts[0], "/")
	if len(ranks) != 8 {
		return fmt.Errorf("invalid FEN: expected number of ranks to be 8, got %d", len(ranks))
	}

	square := 56
	for _, rank := range ranks {
		file := 0
		for _, char := range rank {
			if char >= '1' && char <= '8' {
				file += int(char - '0')
			} else {
				sq := square + file
				var pieceIdx int
				switch char {
				case 'P':
					pieceIdx = WhitePawn
				case 'N':
					pieceIdx = WhiteKnight
				case 'B':
					pieceIdx = WhiteBishop
				case 'R':
					pieceIdx = WhiteRook
				case 'Q':
					pieceIdx = WhiteQueen
				case 'K':
					pieceIdx = WhiteKing
				case 'p':
					pieceIdx = BlackPawn
				case 'n':
					pieceIdx = BlackKnight
				case 'b':
					pieceIdx = BlackBishop
				case 'r':
					pieceIdx = BlackRook
				case 'q':
					pieceIdx = BlackQueen
				case 'k':
					pieceIdx = BlackKing
				default:
					return fmt.Errorf("invalid FEN: unknown piece %c", char)
				}
				b.Pieces[pieceIdx] |= (1 << Bitboard(sq))
				file++
			}
		}
		square -= 8
	}

outer:
	for _, char := range parts[2] {
		switch char {
		case 'K':
			b.CastleRights |= WhiteKingSide
		case 'Q':
			b.CastleRights |= WhiteQueenSide
		case 'k':
			b.CastleRights |= BlackKingSide
		case 'q':
			b.CastleRights |= BlackQueenSide
		case '-':
			break outer
		default:
			return fmt.Errorf("invalid FEN: unknown castling right %c", char)
		}
	}

	if parts[3] == "-" {
		b.EnPassant = -1
	} else {
		file := int(parts[3][0] - 'a')
		rank := int(parts[3][1] - '1')
		b.EnPassant = rank*8 + file
	}

	_, err := fmt.Sscanf(parts[4], "%d", &b.HalfMove)
	if err != nil {
		return fmt.Errorf("invalid FEN: halfmove clock: %v", err)
	}

	_, err = fmt.Sscanf(parts[5], "%d", &b.FullMove)
	if err != nil {
		return fmt.Errorf("invalid FEN: fullmove number: %v", err)
	}

	for i := WhitePawn; i <= WhiteKing; i++ {
		b.Occupancy[White] |= b.Pieces[i]
	}

	for i := BlackPawn; i <= BlackKing; i++ {
		b.Occupancy[Black] |= b.Pieces[i]
	}

	b.Occupancy[All] = b.Occupancy[White] | b.Occupancy[Black]

	return nil
}

// ExportFEN extracts board data and returns a string in format
// "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
func (b *Board) ExportFEN() string {
	fen := strings.Builder{}

	for rank := 7; rank >= 0; rank-- {
		emptyCount := 0

		for file := range 8 {
			square := rank*8 + file
			pieceIdx := -1
			for i := range 12 {
				if b.Pieces[i]&(1<<square) != 0 {
					pieceIdx = i
					break
				}
			}

			if pieceIdx == -1 {
				emptyCount++
			} else {
				if emptyCount > 0 {
					fen.WriteString(strconv.Itoa(emptyCount))
					emptyCount = 0
				}

				switch pieceIdx {
				case WhitePawn:
					fen.WriteRune('P')
				case WhiteKnight:
					fen.WriteRune('N')
				case WhiteBishop:
					fen.WriteRune('B')
				case WhiteRook:
					fen.WriteRune('R')
				case WhiteQueen:
					fen.WriteRune('Q')
				case WhiteKing:
					fen.WriteRune('K')
				case BlackPawn:
					fen.WriteRune('p')
				case BlackKnight:
					fen.WriteRune('n')
				case BlackBishop:
					fen.WriteRune('b')
				case BlackRook:
					fen.WriteRune('r')
				case BlackQueen:
					fen.WriteRune('q')
				case BlackKing:
					fen.WriteRune('k')
				}
			}
		}
		if emptyCount > 0 {
			fen.WriteString(strconv.Itoa(emptyCount))
		}

		if rank > 0 {
			fen.WriteRune('/')
		}
	}

	return fen.String()
}

// InitBoard returns an empty baord with optional
// fen argument to fill the board with that data
func InitBoard(fen string) *Board {
	b := &Board{}
	if fen != "" {
		b.ParseFEN(fen)
	}

	return b
}
