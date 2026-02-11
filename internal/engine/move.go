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

func NewMove(from, to, flags int) Move {
	return Move((from & 0x3F) | ((to & 0x3F) << 6) | ((flags & 0xF) << 12))
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

func (m Move) IsCapture() bool {
	return (m.Flags() & Capture) != 0
}

func (m Move) IsPromotion() bool {
	return m.Flags() >= 8
}

func (m Move) IsSpecial() bool {
	return m.Flags() != QuietMove
}
