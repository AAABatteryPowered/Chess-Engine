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
		return uint64(len(moves))
	}

	var nodes uint64
	for _, move := range moves {
		newBoard := b.Copy()

		/*if depth == 2 && i < 5 { // Just print first 5 moves
			fmt.Printf("Move %d: From=%d To=%d, Turn before=%v\n",
				i, move.From, move.To, newBoard.Turn)
		}*/

		newBoard.PlayMove(move)

		/*if depth == 2 && i < 5 {
			fmt.Printf("  After move: Turn=%v, Next moves=%d\n",
				newBoard.Turn, len(newBoard.Moves()))
		}*/

		nodes += Perft(newBoard, depth-1)
	}

	return nodes
}
func main() {
	b := &board.Board{}
	b.FromFen("rnbqkbnr/pppppppp/8/8/4P3/8/PPPPPPPP/RNBQKBNR b KQkq e3 0 1")
	b.SetTurn(true)

	fmt.Println("Running Perft test:")
	fmt.Println(Perft(b, 3))

}

//put a super-piece on the square the king is on and if it can attack a piece like a rook or bishop then the king is in check aura farm
