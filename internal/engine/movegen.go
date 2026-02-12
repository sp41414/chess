package engine

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
