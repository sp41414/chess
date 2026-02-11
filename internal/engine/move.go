package engine

// From: 0-5 bits
// To: 6-11 bits
// Flags: 12-15 bits
type Move uint16

const (
	// Basic move flags
	QuietMove  = 0
	DoublePush = 1
	KCastle    = 2
	QCastle    = 3
	Capture    = 4
	EPCapture  = 5

	// Promotion flags
	NPromotion = 8
	BPromotion = 9
	RPromotion = 10
	QPromotion = 11

	// Promotion + Capture flags
	NPromotionCapture = 12
	BPromotionCapture = 13
	RPromotionCapture = 14
	QPromotionCapture = 15
)

// NewMove composes a move uint16 from from, to, and flags.
func NewMove(from, to, flags int) Move {
	return Move((from & 0x3F) | ((to & 0x3F) << 6) | ((flags & 0xF) << 12))
}

// From returns the from square, 0-5th bits of the move.
func (m Move) From() int {
	return int(m & 0x3F)
}

// To returns the to square, 6-11th bits of the move.
func (m Move) To() int {
	return int((m >> 6) & 0x3F)
}

// Flags returns the flags, 12-15th bits of the move.
func (m Move) Flags() int {
	return int((m >> 12) & 0xF)
}

// IsCapture returns true if Capture flag bit is set.
func (m Move) IsCapture() bool {
	return (m.Flags() & Capture) != 0
}

// IsPromotion returns true if any promotion flag bit is set.
func (m Move) IsPromotion() bool {
	return m.Flags() >= 8
}

// IsSpecial returns true if the move is not a quiet move
// from flags.
func (m Move) IsSpecial() bool {
	return m.Flags() != QuietMove
}
