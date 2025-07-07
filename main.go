package main

import (
	"bot/board"
	"fmt"
)

func main() {
	b := &board.Board{}
	b.FromFen("8/8/8/0N7/8/8/8/8")
	b.SetTurn(true)

	moves := b.GenMoves()
	fmt.Println(moves)
	fmt.Println(len(moves))
}
