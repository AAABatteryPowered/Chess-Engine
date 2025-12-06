package main

import (
	"bot/board"
	"bot/evaluation"
	"bot/moves"
	"fmt"
	"os"
	"runtime/pprof"
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
		b.PlayMove(move)
		nodes += Perft(b, depth-1)
		b.UndoMove(move)
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
		b.PlayMove(move)

		var nodes uint64
		if depth == 1 {
			nodes = 1
		} else {
			nodes = Perft(b, depth-1)
		}

		b.UndoMove(move)

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
	ranks := "87654321"

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

	f, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	board.InitMagicBitboards()
	b := board.Board{
		UndoCount: 0,
	}
	b.FromFen("5rk1/2q3p1/p1b1p2p/1p1p3P/5r1Q/2N3R1/PPP2PP1/5RK1 w - - 4 23")

	//fmt.Println(Perft(&b, 6))
	move := evaluation.FindBestMove(&b, 6)
	fmt.Println(MoveToString(move, &b))

	//fmt.Println(MoveToString(moves.NewMove(51, 59, moves.FlagPromotionKnight), &b))
	//b.PlayMove(move)
	//b.DebugPrint()
	//evaluation.BestMove(&b)
}
