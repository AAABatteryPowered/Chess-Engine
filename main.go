package main

import (
	"bot/board"
	"bot/moves"
	"fmt"
	"sort"
	"time"
)

func Perft(b *board.Board, depth int) uint64 {
	if depth == 0 {
		return 1
	}

	moves := b.Moves()

	if depth == 1 {
		return uint64(moves.Count)
	}

	var nodes uint64
	for i := 0; i < moves.Count; i++ {
		move := moves.Moves[i]
		undo := b.PlayMove(move)
		nodes += Perft(b, depth-1)
		b.UndoMove(move, undo)
	}

	return nodes
}

func PerftDivide(b *board.Board, depth int) uint64 {
	Moves := b.Moves()
	var totalNodes uint64

	type MoveResult struct {
		move  moves.Move
		nodes uint64
		str   string
	}
	results := make([]MoveResult, 0, Moves.Count)

	fmt.Printf("\nPerft Divide (depth %d):\n", depth)
	fmt.Println("------------------------")

	for i := 0; i < Moves.Count; i++ {
		move := Moves.Moves[i]
		undo := b.PlayMove(move)

		var nodes uint64
		if depth == 1 {
			nodes = 1
		} else {
			nodes = Perft(b, depth-1)
		}

		b.UndoMove(move, undo)

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

func MoveToString(m moves.Move, b *board.Board) string {
	files := "abcdefgh"
	ranks := "12345678"

	from := fmt.Sprintf("%c%c", files[m.From()%8], ranks[m.From()/8])
	to := fmt.Sprintf("%c%c", files[m.To()%8], ranks[m.To()/8])

	promotion := ""
	if m.IsPromotion() {
		switch m.PromotionPiece() {
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

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func main() {
	defer timer("main")()
	board.InitMagicBitboards()
	b := board.Board{}
	b.FromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	fmt.Println(Perft(&b, 5))
}
