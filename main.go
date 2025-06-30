package main

import (
	"bot/board"
	"fmt"
)

func main() {
	b := &board.Board{}
	b.FromFen("rnbqkbnr/1p1p1p1p/8/p1p1p1p1/1P1P1P1P/8/P1P1P1P1/RNBQKBNR")
	b.SetTurn(true)

	fmt.Println(len(b.GenMoves()))
	fmt.Println(b.GenMoves())
	b.DebugPrint()
}
