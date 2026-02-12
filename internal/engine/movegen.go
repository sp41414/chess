package engine

// IsSqAttacked returns true if the square given is attacked by a piece
// from the given color
func (b *Board) IsSqAttacked(sq, color int) bool {
	oc := b.Occupancy[All]

	// 0 for White, 6 for Black
	offset := color * 6

	// For bishop and rook we also check for Queen
	// since it moves similarly to both the bishop and the rook.
	bishop := BishopAttacks(sq, oc)
	if bishop&(b.Pieces[offset+2]|b.Pieces[offset+4]) != 0 {
		return true
	}

	rook := RookAttacks(sq, oc)
	if rook&(b.Pieces[offset+3]|b.Pieces[offset+4]) != 0 {
		return true
	}

	if KnightMoves[sq]&b.Pieces[offset+1] != 0 {
		return true
	}

	if KingMoves[sq]&b.Pieces[offset+5] != 0 {
		return true
	}

	pawn := PawnAttacks(sq, 1-color)
	if pawn&b.Pieces[offset] != 0 {
		return true
	}

	return false
}

func (b *Board) genKnightMoves(moves *[]Move) {
	offset := b.SideToMove * 6
	knights := b.Pieces[offset+1]

	for knights != 0 {
		sq := knights.PopLSB()
		attacks := KnightMoves[sq] & ^b.Occupancy[b.SideToMove]

		for attacks != 0 {
			to := attacks.PopLSB()
			flag := QuietMove
			if b.Occupancy[b.SideToMove^1].Occupied(to) {
				flag = Capture
			}
			*moves = append(*moves, NewMove(sq, to, flag))
		}
	}
}

func (b *Board) genBishopMoves(moves *[]Move) {
	offset := b.SideToMove * 6
	bishops := b.Pieces[offset+2]
	for bishops != 0 {
		sq := bishops.PopLSB()
		attacks := BishopAttacks(sq, b.Occupancy[All]) & ^b.Occupancy[b.SideToMove]
		for attacks != 0 {
			to := attacks.PopLSB()
			flag := QuietMove
			if b.Occupancy[b.SideToMove^1].Occupied(to) {
				flag = Capture
			}
			*moves = append(*moves, NewMove(sq, to, flag))
		}
	}
}

func (b *Board) genRookMoves(moves *[]Move) {
	offset := b.SideToMove * 6
	rooks := b.Pieces[offset+3]
	for rooks != 0 {
		sq := rooks.PopLSB()
		attacks := RookAttacks(sq, b.Occupancy[All]) & ^b.Occupancy[b.SideToMove]
		for attacks != 0 {
			to := attacks.PopLSB()
			flag := QuietMove
			if b.Occupancy[b.SideToMove^1].Occupied(to) {
				flag = Capture
			}
			*moves = append(*moves, NewMove(sq, to, flag))
		}
	}
}

func (b *Board) genQueenMoves(moves *[]Move) {
	offset := b.SideToMove * 6
	queens := b.Pieces[offset+4]
	for queens != 0 {
		sq := queens.PopLSB()
		attacks := QueenAttacks(sq, b.Occupancy[All]) & ^b.Occupancy[b.SideToMove]
		for attacks != 0 {
			to := attacks.PopLSB()
			flag := QuietMove
			if b.Occupancy[b.SideToMove^1].Occupied(to) {
				flag = Capture
			}
			*moves = append(*moves, NewMove(sq, to, flag))
		}
	}
}

func (b *Board) genKingMoves(moves *[]Move) {
	offset := b.SideToMove * 6
	kings := b.Pieces[offset+5]
	for kings != 0 {
		sq := kings.PopLSB()
		attacks := KingAttacks(sq) & ^b.Occupancy[b.SideToMove]
		for attacks != 0 {
			to := attacks.PopLSB()
			flag := QuietMove
			if b.Occupancy[b.SideToMove^1].Occupied(to) {
				flag = Capture
			}
			*moves = append(*moves, NewMove(sq, to, flag))
		}
	}
}

func (b *Board) genPawnMoves(moves *[]Move) {
	offset := b.SideToMove * 6
	pawns := b.Pieces[offset]

	for pawns != 0 {
		sq := pawns.PopLSB()
		rank := sq / 8

		// Single push
		if b.SideToMove == White {
			to := sq + 8
			if to <= 63 && !b.Occupancy[All].Occupied(to) {
				if rank == 6 {
					*moves = append(*moves, NewMove(sq, to, NPromotion))
					*moves = append(*moves, NewMove(sq, to, BPromotion))
					*moves = append(*moves, NewMove(sq, to, RPromotion))
					*moves = append(*moves, NewMove(sq, to, QPromotion))
				} else {
					*moves = append(*moves, NewMove(sq, to, QuietMove))
				}
			}
		} else {
			to := sq - 8
			if to >= 0 && !b.Occupancy[All].Occupied(to) {
				if rank == 1 {
					*moves = append(*moves, NewMove(sq, to, NPromotion))
					*moves = append(*moves, NewMove(sq, to, BPromotion))
					*moves = append(*moves, NewMove(sq, to, RPromotion))
					*moves = append(*moves, NewMove(sq, to, QPromotion))
				} else {
					*moves = append(*moves, NewMove(sq, to, QuietMove))
				}
			}
		}

		// Double push
		if b.SideToMove == White && rank == 1 {
			if !b.Occupancy[All].Occupied(sq+8) && !b.Occupancy[All].Occupied(sq+16) {
				*moves = append(*moves, NewMove(sq, sq+16, DoublePush))
			}
		} else if b.SideToMove == Black && rank == 6 {
			if !b.Occupancy[All].Occupied(sq-8) && !b.Occupancy[All].Occupied(sq-16) {
				*moves = append(*moves, NewMove(sq, sq-16, DoublePush))
			}
		}

		// Diagonal captures
		attacks := PawnAttacks(sq, b.SideToMove) & b.Occupancy[b.SideToMove^1]
		for attacks != 0 {
			to := attacks.PopLSB()
			if b.SideToMove == White && rank == 6 {
				*moves = append(*moves, NewMove(sq, to, NPromotionCapture))
				*moves = append(*moves, NewMove(sq, to, BPromotionCapture))
				*moves = append(*moves, NewMove(sq, to, RPromotionCapture))
				*moves = append(*moves, NewMove(sq, to, QPromotionCapture))
			} else if b.SideToMove == Black && rank == 1 {
				*moves = append(*moves, NewMove(sq, to, NPromotionCapture))
				*moves = append(*moves, NewMove(sq, to, BPromotionCapture))
				*moves = append(*moves, NewMove(sq, to, RPromotionCapture))
				*moves = append(*moves, NewMove(sq, to, QPromotionCapture))
			} else {
				*moves = append(*moves, NewMove(sq, to, Capture))
			}
		}

		// En passant
		if b.EnPassant != -1 {
			epAttacks := PawnAttacks(sq, b.SideToMove)
			if epAttacks.Occupied(b.EnPassant) {
				*moves = append(*moves, NewMove(sq, b.EnPassant, EPCapture))
			}
		}
	}
}

func (b *Board) GenerateMoves() []Move {
	pseudoLegalMoves := make([]Move, 0, 64)

	b.genPawnMoves(&pseudoLegalMoves)
	b.genKnightMoves(&pseudoLegalMoves)
	b.genBishopMoves(&pseudoLegalMoves)
	b.genRookMoves(&pseudoLegalMoves)
	b.genQueenMoves(&pseudoLegalMoves)
	b.genKingMoves(&pseudoLegalMoves)
	// TODO: castling moves: KS, QS, check if enemy controlled

	moves := make([]Move, 0, len(pseudoLegalMoves))
	offset := b.SideToMove * 6

	for _, move := range pseudoLegalMoves {
		undo := b.MakeMove(move)
		king := b.Pieces[offset+5].LSB()

		if !b.IsSqAttacked(king, b.SideToMove) {
			moves = append(moves, move)
		}

		b.UnmakeMove(move, undo)
	}

	return moves
}
