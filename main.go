package main

import (
	"bot/board"
	"fmt"
)

func main() {
	b := &board.Board{}
	b.FromFen("rnbqkbnr/pppppppp/8/8/4R3/8/PPPPPPPP/RNBQKBNR")
	b.SetTurn(true)
	b.DebugPrint()
	//fmt.Println(len(b.GenMoves()))
	fmt.Println(b.BPawns)

	fmt.Println(b.BPawns << 8)
}
