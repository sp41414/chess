package engine

import "fmt"

func Perft(b *Board, depth int) uint64 {
	if depth == 0 {
		return 1
	}

	moves := b.GenerateMoves()
	nodes := uint64(0)

	for _, move := range moves {
		undo := b.MakeMove(move)
		nodes += Perft(b, depth-1)
		b.UnmakeMove(move, undo)
	}

	return nodes
}

func PerftTest(b *Board, maxDepth int) {
	for depth := 1; depth <= maxDepth; depth++ {
		nodes := Perft(b, depth)
		fmt.Printf("Depth %d: %d\n", depth, nodes)
	}
}
