package main

import "bot/board"

func main() {
	b := &board.Board{}
	b.FromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR")
	b.DebugPrint()
}
