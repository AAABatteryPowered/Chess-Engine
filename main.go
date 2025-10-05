package main

import (
	"bot/board"
	"fmt"
)

func Perft(b *board.Board, depth int) uint64 {
	if depth == 0 {
		return 1
	}

	moves := b.Moves()

	if depth == 1 {
		fmt.Printf("Depth 1: %d moves\n", len(moves))
		return uint64(len(moves))
	}

	var nodes uint64
	for i, move := range moves {
		newBoard := &board.Board{}
		*newBoard = *b

		movesBeforePlay := len(newBoard.Moves())
		newBoard.PlayMove(move)
		movesAfterPlay := len(newBoard.Moves())

		if depth == 2 {
			fmt.Printf("Move %d: %d moves before, %d moves after\n",
				i, movesBeforePlay, movesAfterPlay)
		}

		nodes += Perft(newBoard, depth-1)
	}

	return nodes
}
func main() {
	b := &board.Board{}
	b.FromFen("8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1 ")
	b.SetTurn(true)

	fmt.Println("Running Perft test:")
	fmt.Println(len(b.Moves()))

}

//put a super-piece on the square the king is on and if it can attack a piece like a rook or bishop then the king is in check aura farm
