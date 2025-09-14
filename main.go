package main

import (
	"bot/board"
	"fmt"
)

func main() {
	fmt.Println("chess engine is a go")
	b := &board.Board{}
	b.FromFen("8/8/8/3pP3/8/8/8/8 w - d6 0 1")
	b.SetTurn(true)
	b.DebugPrint()

	moves := b.Moves()
	fmt.Println(moves)
	fmt.Println(len(moves))
}

//put a super-piece on the square the king is on and if it can attack a piece like a rook or bishop then the king is in check aura farm
