package main

import (
	"bot/board"
	"fmt"
)

func main() {
	b := &board.Board{}
	b.FromFen("8/8/8/3K4/8/8/8/8")
	b.SetTurn(true)

	moves := b.GenMoves()
	fmt.Println(moves)
	fmt.Println(len(moves))
}
