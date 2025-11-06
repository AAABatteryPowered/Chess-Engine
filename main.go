package main

import (
	"bot/board"
	"fmt"
	"sort"
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
		newBoard.PlayMove(move)
		nodes += Perft(newBoard, depth-1)
	}

	return nodes
}

func PerftDivide(b *board.Board, depth int) uint64 {
	moves := b.Moves()
	var totalNodes uint64

	type MoveResult struct {
		move  board.Move
		nodes uint64
		str   string
	}
	results := make([]MoveResult, 0, len(moves))

	fmt.Printf("\nPerft Divide (depth %d):\n", depth)
	fmt.Println("------------------------")

	for _, move := range moves {
		newBoard := b.Copy()
		newBoard.PlayMove(move)

		var nodes uint64
		if depth == 1 {
			nodes = 1
		} else {
			nodes = Perft(newBoard, depth-1)
		}

		moveStr := MoveToString(move, b)
		results = append(results, MoveResult{move, nodes, moveStr})
		totalNodes += nodes
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].str < results[j].str
	})

	for _, result := range results {
		fmt.Printf("%s: %d\n", result.str, result.nodes)
	}

	fmt.Println("------------------------")
	fmt.Printf("Total: %d\n", totalNodes)

	return totalNodes
}

func MoveToString(m board.Move, b *board.Board) string {
	files := "abcdefgh"
	ranks := "12345678"

	from := fmt.Sprintf("%c%c", files[m.From%8], ranks[m.From/8])
	to := fmt.Sprintf("%c%c", files[m.To%8], ranks[m.To/8])

	promotion := ""
	if m.Promotion != 0 {
		switch m.Promotion {
		case 1:
			promotion = "q"
		case 2:
			promotion = "r"
		case 3:
			promotion = "b"
		case 4:
			promotion = "n"
		}
	}

	return from + to + promotion
}

func main() {
	b := &board.Board{}
	b.FromFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")

	ifsfd := b.WPawns
	ifsfd.DebugPrint()

	PerftDivide(b, 2)
}

/*
CASTLING WORKS EVEN WITHOUT A KING AND ROOK
CASTLING WORKS EVEN WITHOUT A KING AND ROOKCASTLING WORKS EVEN WITHOUT A KING AND ROOKCASTLING WORKS EVEN WITHOUT A KING AND ROOKCASTLING WORKS EVEN WITHOUT A KING AND ROOKCASTLING WORKS EVEN WITHOUT A KING AND ROOKCASTLING WORKS EVEN WITHOUT A KING AND ROOKCASTLING WORKS EVEN WITHOUT A KING AND ROOKCASTLING WORKS EVEN WITHOUT A KING AND ROOKCASTLING WORKS EVEN WITHOUT A KING AND ROOKCASTLING WORKS EVEN WITHOUT A KING AND ROOKCASTLING WORKS EVEN WITHOUT A KING AND ROOKCASTLING WORKS EVEN WITHOUT A KING AND ROOKCASTLING WORKS EVEN WITHOUT A KING AND ROOK
*/
