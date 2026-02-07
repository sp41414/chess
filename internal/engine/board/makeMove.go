package board

type Undo struct {
	Captured     int
	EnPassant    int
	CastleRights int
	HalfMove     int
	Hash         uint64
}

// checkPieceMove returns the piece index from the board which is moving.
func checkPieceMove(m Move, b *Board) int {
	pIdx := -1

	for i, p := range b.Pieces {
		if p&Bitboard(1<<m.From()) != 0 {
			pIdx = i
			break
		}
	}

	return pIdx
}

// checkPieceCapture returns the piece index from the board which is getting captured.
func checkPieceCapture(m Move, b *Board) int {
	captureIdx := -1
	start := -1
	end := -1

	if b.SideToMove == White {
		start = BlackPawn
		end = BlackKing
	} else {
		start = WhitePawn
		end = WhiteKing
	}

	for i := start; i <= end; i++ {
		if b.Pieces[i]&Bitboard(1<<m.To()) != 0 {
			captureIdx = i
			break
		}
	}

	return captureIdx
}

// movePiece updates position, occupancy, and castle rights
// depending on if the piece moving is a Rook or King.
func movePiece(pIdx int, m Move, b *Board) {
	// Update piece position and occupancy
	b.Pieces[pIdx] &^= (1 << m.From())
	b.Pieces[pIdx] |= (1 << m.To())
	b.Occupancy[b.SideToMove] &^= (1 << m.From())
	b.Occupancy[b.SideToMove] |= (1 << m.To())
	b.Occupancy[All] &^= (1 << m.From())
	b.Occupancy[All] |= (1 << m.To())

	// Remove castle rights if king moves or rook moves from starting position
	if pIdx == WhiteKing {
		b.CastleRights &^= (WhiteQueenSide | WhiteKingSide)
	}

	if pIdx == BlackKing {
		b.CastleRights &^= (BlackQueenSide | BlackKingSide)
	}

	switch m.From() {
	case 0:
		b.CastleRights &^= WhiteQueenSide
	case 7:
		b.CastleRights &^= WhiteKingSide
	case 56:
		b.CastleRights &^= BlackQueenSide
	case 63:
		b.CastleRights &^= BlackKingSide
	}
}

// removeCaptured removes a piece from the position, updates occupancy, hash,
// and castle rights if the removed piece is a Rook on starting position.
func removeCaptured(captureIdx int, m Move, b *Board) {
	// Remove piece from opposite side and update occupancy
	b.Pieces[captureIdx] &^= (1 << m.To())
	b.Occupancy[b.SideToMove^1] &^= (1 << m.To())
	b.Occupancy[All] &^= (1 << m.To())

	// Update zobrist hash for capture
	b.Hash ^= zobristPieces[captureIdx][m.To()]

	// Remove castle rights if rook is captured from starting position
	switch m.To() {
	case 0:
		b.CastleRights &^= WhiteQueenSide
	case 7:
		b.CastleRights &^= WhiteKingSide
	case 56:
		b.CastleRights &^= BlackQueenSide
	case 63:
		b.CastleRights &^= BlackKingSide
	}
}

// updateEnPassant checks for double pushes and updates the board EnPassant
// and zobrist hash EnPassant.
func updateEnPassant(pIdx int, m Move, b *Board) {
	// Store old enpassant value and update to new enpassant
	enPassant := b.EnPassant
	b.EnPassant = -1

	// Check for double push from starting position
	if (pIdx == WhitePawn && (m.To()-m.From()) == 16) || (pIdx == BlackPawn && (m.To()-m.From()) == -16) {
		// The square inbetween
		b.EnPassant = (m.From() + m.To()) / 2
	}

	// Update zobrist hash
	if enPassant != -1 {
		b.Hash ^= zobristEnPassant[enPassant]
	}
	if b.EnPassant != -1 {
		b.Hash ^= zobristEnPassant[b.EnPassant]
	}
}

// updateState updates the half move rule counter, and the full move counter on Board.
func updateState(pIdx int, captureIdx int, b *Board) {
	// Half move counter
	if captureIdx != -1 || pIdx == WhitePawn || pIdx == BlackPawn {
		b.HalfMove = 0
	} else {
		b.HalfMove++
	}

	// Full move counter
	if b.SideToMove == Black {
		b.FullMove++
	}
}

const (
	FlagEnPassant = 1
	FlagCastle    = 2

	FlagPromoteQueen  = 8
	FlagPromoteRook   = 9
	FlagPromoteBishop = 10
	FlagPromoteKnight = 11
)

// checkFlags checks for board Flags and modifies the moves for castling,
// EnPassant, and promotion. Updates occupancy, Hash, and positions.
func checkFlags(pIdx int, m Move, b *Board) {
	switch m.Flags() {
	case FlagEnPassant:
		if b.SideToMove == White {
			b.Pieces[BlackPawn] &^= (1 << (m.To() - 8))
			b.Occupancy[Black] &^= (1 << (m.To() - 8))
			b.Occupancy[All] &^= (1 << (m.To() - 8))
			b.Hash ^= zobristPieces[BlackPawn][m.To()-8]
		} else {
			b.Pieces[WhitePawn] &^= (1 << (m.To() + 8))
			b.Occupancy[White] &^= (1 << (m.To() + 8))
			b.Occupancy[All] &^= (1 << (m.To() + 8))
			b.Hash ^= zobristPieces[WhitePawn][m.To()+8]
		}
	case FlagCastle:
		var rookFrom, rookTo, rookIdx int
		switch m.To() {
		case 6:
			rookFrom, rookTo, rookIdx = 7, 5, WhiteRook
		case 2:
			rookFrom, rookTo, rookIdx = 0, 3, WhiteRook
		case 62:
			rookFrom, rookTo, rookIdx = 63, 61, BlackRook
		case 58:
			rookFrom, rookTo, rookIdx = 56, 59, BlackRook
		}

		b.Pieces[rookIdx] &^= (1 << rookFrom)
		b.Pieces[rookIdx] |= (1 << rookTo)
		b.Hash ^= zobristPieces[rookIdx][rookFrom]
		b.Hash ^= zobristPieces[rookIdx][rookTo]
		b.Occupancy[b.SideToMove] &^= (1 << rookFrom)
		b.Occupancy[b.SideToMove] |= (1 << rookTo)
		b.Occupancy[All] &^= (1 << rookFrom)
		b.Occupancy[All] |= (1 << rookTo)
	case FlagPromoteQueen, FlagPromoteBishop, FlagPromoteKnight, FlagPromoteRook:
		b.Pieces[pIdx] &^= (1 << m.To())
		b.Hash ^= zobristPieces[pIdx][m.To()]
		promotionType := m.Flags() - 8

		if b.SideToMove == Black {
			b.Pieces[BlackQueen+promotionType] |= (1 << m.To())
			b.Hash ^= zobristPieces[BlackQueen+promotionType][m.To()]
		} else {
			b.Pieces[WhiteQueen+promotionType] |= (1 << m.To())
			b.Hash ^= zobristPieces[WhiteQueen+promotionType][m.To()]
		}
	}
}

func (b *Board) MakeMove(m Move) Undo {
	undo := Undo{
		Captured:     -1,
		EnPassant:    b.EnPassant,
		CastleRights: b.CastleRights,
		HalfMove:     b.HalfMove,
		Hash:         b.Hash,
	}

	pIdx := checkPieceMove(m, b)
	captureIdx := checkPieceCapture(m, b)

	movePiece(pIdx, m, b)

	if captureIdx != -1 {
		undo.Captured = captureIdx
		removeCaptured(captureIdx, m, b)
	}

	updateEnPassant(pIdx, m, b)
	updateState(pIdx, captureIdx, b)
	checkFlags(pIdx, m, b)

	b.Hash ^= zobristCastle[b.CastleRights]
	b.Hash ^= zobristPieces[pIdx][m.From()]
	b.Hash ^= zobristPieces[pIdx][m.To()]
	b.Hash ^= zobristSideToMove
	b.SideToMove ^= 1

	return undo
}

func (b *Board) UnmakeMove(m Move, undo Undo) {
	// Restore state
	if b.SideToMove == White {
		b.FullMove--
	}
	b.SideToMove ^= 1
	b.Hash = undo.Hash
	b.HalfMove = undo.HalfMove

	// Reverse special flags
	b.CastleRights = undo.CastleRights
	b.EnPassant = undo.EnPassant
	pIdx := checkPieceMove(m, b)
	switch m.Flags() {
	case FlagEnPassant:
		if b.SideToMove == White {
			b.Pieces[BlackPawn] |= (1 << (m.To() - 8))
			b.Occupancy[Black] |= (1 << (m.To() - 8))
			b.Occupancy[All] |= (1 << (m.To() - 8))
		} else {
			b.Pieces[WhitePawn] |= (1 << (m.To() + 8))
			b.Occupancy[White] |= (1 << (m.To() + 8))
			b.Occupancy[All] |= (1 << (m.To() + 8))
		}
	case FlagCastle:
		var rookFrom, rookTo, rookIdx int
		switch m.To() {
		case 6:
			rookFrom, rookTo, rookIdx = 7, 5, WhiteRook
		case 2:
			rookFrom, rookTo, rookIdx = 0, 3, WhiteRook
		case 62:
			rookFrom, rookTo, rookIdx = 63, 61, BlackRook
		case 58:
			rookFrom, rookTo, rookIdx = 56, 59, BlackRook
		}

		b.Pieces[rookIdx] |= (1 << rookFrom)
		b.Pieces[rookIdx] &^= (1 << rookTo)
		b.Occupancy[b.SideToMove] |= (1 << rookFrom)
		b.Occupancy[b.SideToMove] &^= (1 << rookTo)
		b.Occupancy[All] |= (1 << rookFrom)
		b.Occupancy[All] &^= (1 << rookTo)
	case FlagPromoteQueen, FlagPromoteBishop, FlagPromoteKnight, FlagPromoteRook:
		promotionType := m.Flags() - 8

		if b.SideToMove == Black {
			b.Pieces[BlackQueen+promotionType] &^= (1 << m.To())
			b.Pieces[BlackPawn] |= (1 << m.To())
			pIdx = BlackPawn
		} else {
			b.Pieces[WhiteQueen+promotionType] &^= (1 << m.To())
			b.Pieces[WhitePawn] |= (1 << m.To())
			pIdx = WhitePawn
		}
	}

	// Restore captured piece
	if undo.Captured != -1 {
		b.Pieces[undo.Captured] |= (1 << m.To())
		b.Occupancy[b.SideToMove^1] |= (1 << m.To())
		b.Occupancy[All] |= (1 << m.To())

	}

	// Move piece back
	b.Pieces[pIdx] |= (1 << m.From())
	b.Pieces[pIdx] &^= (1 << m.To())
	b.Occupancy[b.SideToMove] |= (1 << m.From())
	b.Occupancy[b.SideToMove] &^= (1 << m.To())
	b.Occupancy[All] |= (1 << m.From())
	b.Occupancy[All] &^= (1 << m.To())
}
