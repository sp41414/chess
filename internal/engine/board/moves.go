package board

type Move uint16

func NewMove(from, to, flags int) Move {
	// Use & to ensure each value fits in bit range
	from &= 0x3F // 6 bits
	to &= 0x3F   // 6 bits
	flags &= 0xF // 4 bits

	fromBits := from        // Bits 0 - 5
	toBits := to << 6       // Bits 6 - 11
	flagBits := flags << 12 // Bits 12 - 15

	// Use | to combine it all into a move
	return Move(fromBits) | Move(toBits) | Move(flagBits)
}

// Extracts From from move within 6 bit range
func (m Move) From() int {
	return int(m & 0x3F)
}

// Extracts To from move within 6 - 11 bit range
func (m Move) To() int {
	return int((m >> 6) & 0x3F)
}

// Extracts Flags from move within 12 - 15 bit range
func (m Move) Flags() int {
	return int((m >> 12) & 0xF)
}

const (
	FileA uint64 = 0x0101010101010101
	FileH uint64 = 0x8080808080808080
)

func CalculateKnightMoves(square int) Bitboard {
	attacks := Bitboard(0)
	knight := Bitboard(1 << square)

	FileAB := Bitboard(FileA | (FileA << 1))
	FileGH := Bitboard(FileH | (FileH >> 1))

	attacks |= (knight &^ Bitboard(FileH)) << 17 // up 2, right 1
	attacks |= (knight &^ Bitboard(FileA)) << 15 // up 2, left 1
	attacks |= (knight &^ FileGH) << 10          // up 1, right 2
	attacks |= (knight &^ FileAB) << 6           // up 1, left 2
	attacks |= (knight &^ FileGH) >> 6           // down 1, right 2
	attacks |= (knight &^ FileAB) >> 10          // down 1, left 2
	attacks |= (knight &^ Bitboard(FileH)) >> 15 // down 2, right 1
	attacks |= (knight &^ Bitboard(FileA)) >> 17 // down 2, left 1

	return attacks
}
