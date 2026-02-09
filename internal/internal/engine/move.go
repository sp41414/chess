package engine

// From: Bits 0-5
// To: Bits 6-11
// Flags: Bits 12-15
type Move uint16

func NewMove(from, to, flags int) Move {
	from &= 0x3F
	to &= 0x3F
	flags &= 0xF

	return Move(from) | Move(to<<6) | Move(flags<<12)
}

func (m Move) From() int {
	return int(m & 0x3F)
}

func (m Move) To() int {
	return int((m >> 6) & 0x3F)
}

func (m Move) Flags() int {
	return int((m >> 12) & 0xF)
}
