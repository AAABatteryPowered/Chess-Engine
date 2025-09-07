package main

import (
	"bot/board"
	"fmt"
)

func main() {
	b := &board.Board{}
	b.FromFen("8/8/8/8/8/8/8/R3K3 b KQkq")
	b.SetTurn(true)

	moves := b.Moves()
	fmt.Println(moves)
	fmt.Println(len(moves))
}

//put a super-piece on the square the king is on and if it can attack a piece like a rook or bishop then the king is in check aura farm
