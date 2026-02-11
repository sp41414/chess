package engine

var (
	KnightMoves    []uint64
	BlackPawnMoves []uint64
	WhitePawnMoves []uint64
	KingMoves      []uint64
)

func InitLookupTables() {
	KnightMoves = make([]uint64, 64)
	BlackPawnMoves = make([]uint64, 64)
	WhitePawnMoves = make([]uint64, 64)
	KingMoves = make([]uint64, 64)

	for sq := range 64 {
		// TODO: calculate knight moves first
	}
}
