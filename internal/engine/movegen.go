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

func (b *Board) GenerateMoves() []Move {
	moves := make([]Move, 0, 64)

	// TODO: pawn moves: en passant, promotion
	b.genKnightMoves(&moves)
	b.genBishopMoves(&moves)
	b.genRookMoves(&moves)
	b.genQueenMoves(&moves)
	b.genKingMoves(&moves)
	// TODO: castling moves: KS, QS, check if enemy controlled
	// TODO: check if king is in check, make sure no illegal moves
	// unless a piece defends the king or the king moves away
	// later on, we can check if the king is not in check with no legal moves
	// for stalemate, and then we can check if the king is in check with no legal moves
	// for checkmate.

	return moves
}
