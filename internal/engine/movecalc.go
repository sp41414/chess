package engine

var (
	KnightMoves    [64]Bitboard
	BlackPawnMoves [64]Bitboard
	WhitePawnMoves [64]Bitboard
	KingMoves      [64]Bitboard
)

func InitLookupTables() {
	for sq := range 64 {
		KnightMoves[sq] = KnightAttacks(sq)
		BlackPawnMoves[sq] = PawnAttacks(sq, Black)
		WhitePawnMoves[sq] = PawnAttacks(sq, White)
		KingMoves[sq] = KingAttacks(sq)
	}
}

func KnightAttacks(sq int) Bitboard {
	var bitboard Bitboard
	r := sq / 8
	f := sq % 8

	offsets := []struct {
		dr, df int
	}{
		{2, 1}, {2, -1}, {-2, 1}, {-2, -1},
		{1, 2}, {1, -2}, {-1, 2}, {-1, -2},
	}

	for _, o := range offsets {
		dr, df := r+o.dr, f+o.df

		if dr >= 0 && dr < 8 && df >= 0 && df < 8 {
			bitboard.Set(dr*8 + df)
		}
	}
	return bitboard
}

func PawnAttacks(sq, color int) Bitboard {
	var bitboard Bitboard
	r := sq / 8
	f := sq % 8

	if color == White {
		if r < 7 {
			if f > 0 {
				bitboard.Set((r+1)*8 + f - 1)
			}
			if f < 7 {
				bitboard.Set((r+1)*8 + (f + 1))
			}
		}
	} else {
		if r > 0 {
			if f > 0 {
				bitboard.Set((r-1)*8 + (f - 1))
			}
			if f < 7 {
				bitboard.Set((r-1)*8 + (f + 1))
			}
		}
	}

	return bitboard
}

func KingAttacks(sq int) Bitboard {
	var bitboard Bitboard
	r := sq / 8
	f := sq % 8

	offsets := []struct{ dr, df int }{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1},
		{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
	}

	for _, o := range offsets {
		dr, df := r+o.dr, f+o.df

		if dr >= 0 && dr <= 7 && df >= 0 && df <= 7 {
			bitboard.Set(dr*8 + df)
		}
	}

	return bitboard
}

func RookAttacks(sq int, occupancy Bitboard) Bitboard {
	occ := occupancy & RookMagics[sq].Mask
	idx := (uint64(occ) * RookMagics[sq].MagicNum) >> RookMagics[sq].Shift
	return RookMoves[sq][idx]
}

func BishopAttacks(sq int, occupancy Bitboard) Bitboard {
	occ := occupancy & BishopMagics[sq].Mask
	idx := (uint64(occ) * BishopMagics[sq].MagicNum) >> BishopMagics[sq].Shift
	return BishopMoves[sq][idx]
}

func QueenAttacks(sq int, occupancy Bitboard) Bitboard {
	return RookAttacks(sq, occupancy) | BishopAttacks(sq, occupancy)
}
