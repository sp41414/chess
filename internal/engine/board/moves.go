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
