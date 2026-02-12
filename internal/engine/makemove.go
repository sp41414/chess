package engine

type Undo struct {
	Captured     int
	CastleRights int
	EnPassant    int
	HalfMove     int
}

// MakeMove makes a move on the board and returns an undo struct
// for UnmakeMove
func (b *Board) MakeMove(m Move) Undo {
	from, to, flags := m.From(), m.To(), m.Flags()

	undo := Undo{
		Captured:     -1,
		CastleRights: b.CastleRights,
		EnPassant:    b.EnPassant,
		HalfMove:     b.HalfMove,
	}

	pIdx := -1
	for i := range 12 {
		if b.Pieces[i].Occupied(from) {
			pIdx = i
			break
		}
	}

	cIdx := -1
	if m.IsCapture() && flags != EPCapture {
		for i := range 12 {
			if b.Pieces[i].Occupied(to) {
				cIdx = i
				break
			}
		}
	}
	undo.Captured = cIdx

	// Basic quiet move
	b.Pieces[pIdx].Clear(from)
	b.Pieces[pIdx].Set(to)

	// Capture, remove the captured piece and update occupancy
	if flags == Capture {
		b.Pieces[cIdx].Clear(to)
		if b.SideToMove == White {
			b.Occupancy[Black].Clear(to)
		} else {
			b.Occupancy[White].Clear(to)
		}
	}

	// EnPassant check on double push
	if (pIdx == WhitePawn || pIdx == BlackPawn) && flags == DoublePush {
		// The middle square/inbetween
		b.EnPassant = (from + to) / 2
	} else {
		b.EnPassant = -1
	}

	// EnPassant check on actual EnPassant flag
	if flags == EPCapture {
		if b.SideToMove == White {
			b.Pieces[BlackPawn].Clear(b.EnPassant - 8)
			b.Occupancy[Black].Clear(b.EnPassant - 8)
		} else {
			b.Pieces[WhitePawn].Clear(b.EnPassant + 8)
			b.Occupancy[White].Clear(b.EnPassant + 8)
		}
	}

	// CastleRights check, remove on rook capture
	switch to {
	case 0:
		b.CastleRights &^= WhiteQueenSide
	case 7:
		b.CastleRights &^= WhiteKingSide
	case 56:
		b.CastleRights &^= BlackQueenSide
	case 63:
		b.CastleRights &^= BlackKingSide
	}

	// Remove on rook move
	switch from {
	case 0:
		b.CastleRights &^= WhiteQueenSide
	case 7:
		b.CastleRights &^= WhiteKingSide
	case 56:
		b.CastleRights &^= BlackQueenSide
	case 63:
		b.CastleRights &^= BlackKingSide
	}

	// Remove on king move
	switch pIdx {
	case WhiteKing:
		b.CastleRights &^= (WhiteKingSide | WhiteQueenSide)
	case BlackKing:
		b.CastleRights &^= (BlackKingSide | BlackQueenSide)
	}

	// Actually castle, move the rook
	if flags == KCastle || flags == QCastle {
		var rookFrom, rookTo int
		if flags == KCastle {
			if b.SideToMove == White {
				rookFrom, rookTo = 7, 5
			} else {
				rookFrom, rookTo = 63, 61
			}
		} else {
			// QCastle
			if b.SideToMove == White {
				rookFrom, rookTo = 0, 3
			} else {
				rookFrom, rookTo = 56, 59
			}
		}

		rookIdx := WhiteRook
		if b.SideToMove == Black {
			rookIdx = BlackRook
		}

		b.Pieces[rookIdx].Clear(rookFrom)
		b.Pieces[rookIdx].Set(rookTo)
		b.Occupancy[b.SideToMove].Clear(rookFrom)
		b.Occupancy[b.SideToMove].Set(rookTo)
	}

	// Promotion
	if m.IsPromotion() {
		b.Pieces[pIdx].Clear(to)
		pType := flags - 8

		if b.SideToMove == White {
			b.Pieces[WhiteQueen+pType].Set(to)
		} else {
			b.Pieces[BlackQueen+pType].Set(to)
		}
	}

	// Game state: HalfMove, FullMove
	if cIdx != -1 || pIdx == WhitePawn || pIdx == BlackPawn {
		b.HalfMove = 0
	} else {
		b.HalfMove++
	}

	if b.SideToMove == Black {
		b.FullMove++
	}

	// Update occupancy for move
	if b.SideToMove == White {
		b.SideToMove = Black
		b.Occupancy[White].Clear(from)
		b.Occupancy[White].Set(to)
	} else {
		b.SideToMove = White
		b.Occupancy[Black].Clear(from)
		b.Occupancy[Black].Set(to)
	}

	b.Occupancy[All] = b.Occupancy[White] | b.Occupancy[Black]

	return undo
}

// UnmakeMove undoes a move on the board, using the undo struct
// returned by MakeMove and the move to be unmade.
func (b *Board) UnmakeMove(m Move, undo Undo) {
	if b.SideToMove == White {
		b.FullMove--
		b.SideToMove = Black
	} else {
		b.SideToMove = White
	}

	from, to, flags := m.From(), m.To(), m.Flags()

	pIdx := -1
	for i := range 12 {
		if b.Pieces[i].Occupied(to) {
			pIdx = i
			break
		}
	}

	b.Pieces[pIdx].Clear(to)
	b.Occupancy[b.SideToMove].Clear(to)

	if m.IsPromotion() {
		pawnIdx := WhitePawn
		if b.SideToMove == Black {
			pawnIdx = BlackPawn
		}
		b.Pieces[pawnIdx].Set(from)
		b.Occupancy[b.SideToMove].Set(from)
	} else {
		b.Pieces[pIdx].Set(from)
		b.Occupancy[b.SideToMove].Set(from)
	}

	if flags == EPCapture {
		capSq := to - 8
		capIdx := BlackPawn
		if b.SideToMove == Black {
			capSq = to + 8
			capIdx = WhitePawn
			b.Occupancy[White].Set(capSq)
		} else {
			b.Occupancy[Black].Set(capSq)
		}
		b.Pieces[capIdx].Set(capSq)
	} else if undo.Captured != -1 {
		b.Pieces[undo.Captured].Set(to)
		if b.SideToMove == White {
			b.Occupancy[Black].Set(to)
		} else {
			b.Occupancy[White].Set(to)
		}
	}

	if flags == KCastle || flags == QCastle {
		var rFrom, rTo int
		if flags == KCastle {
			if b.SideToMove == White {
				rFrom, rTo = 7, 5
			} else {
				rFrom, rTo = 63, 61
			}
		} else {
			if b.SideToMove == White {
				rFrom, rTo = 0, 3
			} else {
				rFrom, rTo = 56, 59
			}
		}

		rIdx := WhiteRook
		if b.SideToMove == Black {
			rIdx = BlackRook
		}
		b.Pieces[rIdx].Clear(rTo)
		b.Pieces[rIdx].Set(rFrom)
		b.Occupancy[b.SideToMove].Clear(rTo)
		b.Occupancy[b.SideToMove].Set(rFrom)
	}

	b.CastleRights = undo.CastleRights
	b.EnPassant = undo.EnPassant
	b.HalfMove = undo.HalfMove

	b.Occupancy[All] = b.Occupancy[White] | b.Occupancy[Black]
}
