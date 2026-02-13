package engine

import "math/bits"

// Init initializes the board with the given FEN string
// in format "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
// for starting position.
func Init(fen string) *Board {
	InitLookupTables()
	InitMagic()
	return InitBoard(fen)
}

// IsInCheck returns true if the current side to move's king is attacked
// by opposite side's piece.
func (b *Board) IsInCheck() bool {
	var kingSq int
	if b.SideToMove == White {
		kingSq = bits.TrailingZeros64(uint64(b.Pieces[WhiteKing]))
	} else {
		kingSq = bits.TrailingZeros64(uint64(b.Pieces[BlackKing]))
	}
	return b.IsSqAttacked(kingSq, b.SideToMove^1)
}

// IsCheckmate returns true if the current side to move's king is not attacked
// with no legal moves.
func (b *Board) IsCheckmate() bool {
	return len(b.GenerateMoves()) == 0 && b.IsInCheck()
}

// IsStalemate returns true if the current side to move's king is attacked
// with no legal moves.
func (b *Board) IsStalemate() bool {
	return len(b.GenerateMoves()) == 0 && !b.IsInCheck()
}

// IsFiftyMoveRule returns true if the current half move state is 50
func (b *Board) IsFiftyMoveRule() bool {
	return b.HalfMove >= 100
}

func (b *Board) IsInsufficientMaterial() bool {
	white := b.Occupancy[White].Count()
	black := b.Occupancy[Black].Count()
	// K vs K
	if white == 1 && black == 1 {
		return true
	}

	// K vs K+N or K vs K+B
	if white == 1 && black == 2 {
		if b.Pieces[BlackKnight].Count() == 1 || b.Pieces[BlackBishop].Count() == 1 {
			return true
		}
	}
	if black == 1 && white == 2 {
		if b.Pieces[WhiteKnight].Count() == 1 || b.Pieces[WhiteBishop].Count() == 1 {
			return true
		}
	}

	// K+B vs K+B on same color squares
	if white == 2 && black == 2 && b.Pieces[WhiteBishop].Count() == 1 && b.Pieces[BlackBishop].Count() == 1 {
		whiteBishopSq := b.Pieces[WhiteBishop].LSB()
		blackBishopSq := b.Pieces[BlackBishop].LSB()

		whiteColor := (whiteBishopSq/8 + whiteBishopSq%8) % 2
		blackColor := (blackBishopSq/8 + blackBishopSq%8) % 2

		if whiteColor == blackColor {
			return true
		}
	}

	return false
}

// IsThreefoldRepetition returns true if the current position has been
// repeated 3 times.
func (b *Board) IsThreefoldRepetition() bool {
	return b.PositionCount[b.GetPositionKey()] >= 3
}

// IsDraw returns true if the current position is either a draw by
// stalemate, insufficient material, or threefold repetition.
func (b *Board) IsDraw() bool {
	return b.IsStalemate() || b.IsInsufficientMaterial() || b.IsThreefoldRepetition()
}

// GetFEN returns the current board state in FEN format
func (b *Board) GetFEN() string {
	return b.ExportFEN()
}

// GetMoves returns a list of legal moves for the current board state.
func (b *Board) GetMoves() []Move {
	return b.GenerateMoves()
}

// GetPieces returns a map of squares to piece types found in the frontend as SVGs
func (b *Board) GetPieces() map[int]string {
	pieces := make(map[int]string)
	chars := [12]string{"wP", "wN", "wB", "wR", "wQ", "wK", "bP", "bN", "bB", "bR", "bQ", "bK"}

	for sq := range 64 {
		for i := range 12 {
			if b.Pieces[i].Occupied(sq) {
				pieces[sq] = chars[i]
				break
			}
		}
	}

	return pieces
}

// PlayMove plays the given move on the board and returns the resulting
// board state to undo the move.
func (b *Board) PlayMove(m Move) Undo {
	return b.MakeMove(m)
}

func (b *Board) UndoMove(m Move, u Undo) {
	b.UnmakeMove(m, u)
}
