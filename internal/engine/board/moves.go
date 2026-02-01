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

// CalculateKnightMoves goes through every L direction
// and checks if it goes off the Bitboard
// around FileH, FileA, FileGH, FileAB
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

// CalculateRookMoves goes through up, down, left, and right directions
// and checks if it hits a piece or is off the board
func CalculateRookMoves(square int, occupancy Bitboard) Bitboard {
	attacks := Bitboard(0)
	current := square

	// Loop north
	for {
		current += 8
		if current > 63 {
			break
		}
		attacks |= Bitboard(1 << current)
		if occupancy&(1<<current) != 0 {
			break
		}
	}

	// Loop south
	current = square
	for {
		current -= 8
		if current < 0 {
			break
		}
		attacks |= Bitboard(1 << current)
		if occupancy&(1<<current) != 0 {
			break
		}
	}

	// Loop east
	current = square
	for {
		current += 1
		if current > 63 {
			break
		}
		attacks |= Bitboard(1 << current)
		if Bitboard(1<<current)&Bitboard(FileH) != 0 {
			break
		}
		if occupancy&(1<<current) != 0 {
			break
		}
	}

	// Loop west
	current = square
	for {
		current -= 1
		if current < 0 {
			break
		}
		attacks |= Bitboard(1 << current)
		if Bitboard(1<<current)&Bitboard(FileA) != 0 {
			break
		}
		if occupancy&(1<<current) != 0 {
			break
		}
	}

	return attacks
}

// CalculateBishopMoves goes through every diagonal direction and checks
// if it goes off board or hits a piece
func CalculateBishopMoves(square int, occupancy Bitboard) Bitboard {
	attacks := Bitboard(0)
	current := square

	// Loop north-east
	for {
		if Bitboard(1<<current)&Bitboard(FileH) != 0 {
			break
		}
		current += 9
		if current > 63 {
			break
		}
		attacks |= Bitboard(1 << current)
		if occupancy&(1<<current) != 0 {
			break
		}
	}

	// Loop north-west
	current = square
	for {
		if Bitboard(1<<current)&Bitboard(FileA) != 0 {
			break
		}
		current += 7
		if current > 63 {
			break
		}
		attacks |= Bitboard(1 << current)
		if occupancy&(1<<current) != 0 {
			break
		}
	}

	// Loop south-east
	current = square
	for {
		if Bitboard(1<<current)&Bitboard(FileH) != 0 {
			break
		}
		current -= 7
		if current < 0 {
			break
		}
		attacks |= Bitboard(1 << current)
		if occupancy&(1<<current) != 0 {
			break
		}
	}

	// Loop south-west
	current = square
	for {
		if Bitboard(1<<current)&Bitboard(FileA) != 0 {
			break
		}
		current -= 9
		if current < 0 {
			break
		}
		attacks |= Bitboard(1 << current)
		if occupancy&(1<<current) != 0 {
			break
		}
	}

	return attacks
}

// CalculateQueenMoves combines the rook moves and bishop moves, which basically means all directions.
func CalculateQueenMoves(square int, occupancy Bitboard) Bitboard {
	return CalculateRookMoves(square, occupancy) | CalculateBishopMoves(square, occupancy)
}
