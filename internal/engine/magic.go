package engine

// Magic struct holds the magic number used to generate the moves
// for a piece, the mask of blockers, and the shift used to generate
// the index for the move table
type Magic struct {
	MagicNum uint64
	Mask     Bitboard
	Shift    uint8
}

var (
	RookMagics   [64]Magic
	BishopMagics [64]Magic
	RookMoves    [64][4096]Bitboard
	BishopMoves  [64][512]Bitboard
)

// https://github.com/mcostalba/Stockfish/blob/master/src/bitboard.cpp#L147
var Seeds [8]uint64 = [8]uint64{255, 728, 54343, 32803, 12281, 10316, 15100, 16645}

// PRNG is a pseudo random number generator for generating magic numbers
type PRNG struct {
	state uint64
}

// NewPRNG returns a new PRNG with the given seed from Seeds array
// used for generating magic numbers
func NewPRNG(seed uint64) *PRNG {
	return &PRNG{state: seed}
}

// Rand64 returns a pseudo random 64-bit number using the PRNG state
// after applying bit operations
func (p *PRNG) Rand64() uint64 {
	p.state ^= p.state >> 12
	p.state ^= p.state << 25
	p.state ^= p.state >> 27
	return p.state * 2685821657736338717
}

// SparseRand64 ANDs the PRNG state with itself and itself again
// to generate a pseudo random 64-bit number
func (p *PRNG) SparseRand64() uint64 {
	return p.Rand64() & p.Rand64() & p.Rand64()
}

// InitMagic initializes the Rook and Bishop magic numbers
// and the move tables RookMoves and BishopMoves.
func InitMagic() {
	// THE ROOK
	for sq := range 64 {
		rng := NewPRNG(Seeds[sq%8])
		mask := MaskRookAttacks(sq)
		bits := mask.Count()
		shift := uint8(64 - bits)

		numBlockers := 1 << bits
		blockers := make([]Bitboard, numBlockers)
		attacks := make([]Bitboard, numBlockers)

		b := Bitboard(0)
		for i := range numBlockers {
			blockers[i] = b
			attacks[i] = GenRookAttacks(sq, b)
			b = (b - mask) & mask
		}

		for {
			magic := rng.SparseRand64()
			// A good magic number has atleast 6 bits set in the LSB
			if ((mask * Bitboard(magic)) & 0xFF00000000000000).Count() < 6 {
				continue
			}

			used := make([]Bitboard, 4096)
			success := true

			for i := range blockers {
				index := (uint64(blockers[i]) * magic) >> shift
				if used[index] == 0 {
					used[index] = attacks[i]
					RookMoves[sq][index] = attacks[i]
				} else if used[index] != attacks[i] {
					success = false
					break
				}
			}

			if success {
				RookMagics[sq] = Magic{MagicNum: magic, Mask: mask, Shift: shift}
				break
			}
		}
	}

	// Bishops
	for sq := range 64 {
		rng := NewPRNG(Seeds[sq%8])
		mask := MaskBishopAttacks(sq)
		bits := mask.Count()
		shift := uint8(64 - bits)

		numBlockers := 1 << bits
		blockers := make([]Bitboard, numBlockers)
		attacks := make([]Bitboard, numBlockers)

		b := Bitboard(0)
		for i := range numBlockers {
			blockers[i] = b
			attacks[i] = GenBishopAttacks(sq, b)
			b = (b - mask) & mask
		}

		for {
			magic := rng.SparseRand64()
			if ((mask * Bitboard(magic)) & 0xFF00000000000000).Count() < 6 {
				continue
			}

			used := make([]Bitboard, 512)
			success := true

			for i := range blockers {
				index := (uint64(blockers[i]) * magic) >> shift
				if used[index] == 0 {
					used[index] = attacks[i]
					BishopMoves[sq][index] = attacks[i]
				} else if used[index] != attacks[i] {
					success = false
					break
				}
			}

			if success {
				BishopMagics[sq] = Magic{MagicNum: magic, Mask: mask, Shift: shift}
				break
			}
		}
	}
}

// MaskRookAttacks returns a bitboard with all the squares attacked by a rook
// and where blockers could exist
func MaskRookAttacks(sq int) Bitboard {
	var bitboard Bitboard
	f := sq % 8
	r := sq / 8

	// Up
	for tr := r + 1; tr < 7; tr++ {
		bitboard.Set(tr*8 + f)
	}
	// Down
	for tr := r - 1; tr > 0; tr-- {
		bitboard.Set(tr*8 + f)
	}
	// Right
	for tf := f + 1; tf < 7; tf++ {
		bitboard.Set(r*8 + tf)
	}
	// Left
	for tf := f - 1; tf > 0; tf-- {
		bitboard.Set(r*8 + tf)
	}
	return bitboard
}

// MaskBishopAttacks returns a bitboard with all the squares attacked by a bishop
// and where blockers could exist
func MaskBishopAttacks(sq int) Bitboard {
	var bitboard Bitboard
	f := sq % 8
	r := sq / 8
	tr := r + 1
	tf := f + 1

	// Up-Right
	for tf < 7 && tr < 7 {
		bitboard.Set(tr*8 + tf)
		tr++
		tf++
	}
	tr, tf = r-1, f+1
	// Down-Right
	for tr > 0 && tf < 7 {
		bitboard.Set(tr*8 + tf)
		tr--
		tf++
	}
	tr, tf = r+1, f-1
	// Up-Left
	for tr < 7 && tf > 0 {
		bitboard.Set(tr*8 + tf)
		tr++
		tf--
	}
	tr, tf = r-1, f-1
	// Down-Left
	for tr > 0 && tf > 0 {
		bitboard.Set(tr*8 + tf)
		tr--
		tf--
	}
	return bitboard
}

// GenRookAttacks calculates all the squares a rook can move to or capture on,
// given a square and a bitboard of blocker pieces. Unlike the mask function,
// this function will include the board edges and stop immediately after encountering a blocker.
func GenRookAttacks(sq int, blockers Bitboard) Bitboard {
	var attacks Bitboard
	f, r := sq%8, sq/8

	for tr := r + 1; tr <= 7; tr++ {
		target := tr*8 + f
		attacks.Set(target)
		if blockers.Occupied(target) {
			break
		}
	}
	for tr := r - 1; tr >= 0; tr-- {
		target := tr*8 + f
		attacks.Set(target)
		if blockers.Occupied(target) {
			break
		}
	}
	for tf := f + 1; tf <= 7; tf++ {
		target := r*8 + tf
		attacks.Set(target)
		if blockers.Occupied(target) {
			break
		}
	}
	for tf := f - 1; tf >= 0; tf-- {
		target := r*8 + tf
		attacks.Set(target)
		if blockers.Occupied(target) {
			break
		}
	}
	return attacks
}

// GenBishopAttacks calculates all the squares a bishop can move to or capture on,
// given a square and a bitboard of blocker pieces. It walks the four diagonals until
// it hits a blocker piece or the edge of the board.
func GenBishopAttacks(sq int, blockers Bitboard) Bitboard {
	var attacks Bitboard
	f, r := sq%8, sq/8
	tr, tf := r+1, f+1

	for tr < 8 && tf < 8 {
		target := tr*8 + tf
		attacks.Set(target)
		tr++
		tf++
		if blockers.Occupied(target) {
			break
		}
	}
	tr, tf = r-1, f+1
	for tr >= 0 && tf < 8 {
		target := tr*8 + tf
		attacks.Set(target)
		tr--
		tf++
		if blockers.Occupied(target) {
			break
		}
	}
	tr, tf = r+1, f-1
	for tr < 8 && tf >= 0 {
		target := tr*8 + tf
		attacks.Set(target)
		tr++
		tf--
		if blockers.Occupied(target) {
			break
		}
	}
	tr, tf = r-1, f-1
	for tr >= 0 && tf >= 0 {
		target := tr*8 + tf
		attacks.Set(target)
		tr--
		tf--
		if blockers.Occupied(target) {
			break
		}
	}

	return attacks
}
