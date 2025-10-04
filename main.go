package main

import (
	"bot/board"
	"fmt"
)

func movegenTest(b *board.Board, moves []board.Move, depth int) int {
	if depth == 0 {
		return 0
	}
	total := 0
	for _, mov := range moves {
		clone := &board.Board{}
		*clone = *b
		clone.PlayMove(mov)
		clonemoves := clone.Moves()
		total += movegenTest(clone, clonemoves, depth-1)
		total += len(clonemoves)
	}

	return total
}

func main() {
	b := &board.Board{}
	b.FromFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	b.SetTurn(true)
	fmt.Println("test")

	fmt.Println(b.Moves())
	fmt.Println(len(b.Moves()))
}

//put a super-piece on the square the king is on and if it can attack a piece like a rook or bishop then the king is in check aura farm
