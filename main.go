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
		total += movegenTest(b, clonemoves, depth-1)
		total += len(clonemoves)
	}

	return total
}

func main() {
	b := &board.Board{}
	b.FromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	b.SetTurn(true)
	fmt.Println(movegenTest(b, b.GenMoves(), 4))
}

//put a super-piece on the square the king is on and if it can attack a piece like a rook or bishop then the king is in check aura farm
