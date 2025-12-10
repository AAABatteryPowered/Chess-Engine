package main

import (
	"bot/board"
	"bot/evaluation"
	"bot/moves"
	"bufio"
	"fmt"
	"os"
	"runtime/pprof"
	"sort"
	"strings"
	"time"
)

func Perft(b *board.Board, depth int) uint64 {
	if depth == 0 {
		return 1
	}

	moves := b.Moves(false)

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
	Moves := b.Moves(false)
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

		moveStr := move.MoveToString()
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

func StringToMove(moveStr string) moves.Move {
	files := "abcdefgh"
	ranks := "87654321"

	fromFile := strings.IndexByte(files, moveStr[0])
	fromRank := strings.IndexByte(ranks, moveStr[1])
	from := int8(fromRank*8 + fromFile)

	toFile := strings.IndexByte(files, moveStr[2])
	toRank := strings.IndexByte(ranks, moveStr[3])
	to := int8(toRank*8 + toFile)

	if len(moveStr) == 5 {
		switch moveStr[4] {
		case 'q':
			return moves.NewMove(from, to, moves.FlagPromotionKnight)
		case 'r':
			return moves.NewMove(from, to, moves.FlagPromotionBishop)
		case 'b':
			return moves.NewMove(from, to, moves.FlagPromotionRook)
		case 'n':
			return moves.NewMove(from, to, moves.FlagPromotionQueen)
		}
	}

	return moves.NewMove(from, to, moves.FlagNone)
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
	reader := bufio.NewScanner(os.Stdin)

	b := board.Board{
		UndoCount: 0,
	}
	b.FromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	board.InitZobrist()
	fmt.Println(board.CalculateHash(&b))

	for reader.Scan() {
		line := strings.TrimSpace(reader.Text())

		switch {
		case line == "start":
			move := evaluation.FindBestMove(&b, 7)
			fmt.Println(move.MoveToString())
			b.PlayMove(move)
		case strings.HasPrefix(line, "go"):
			b.PlayMove(StringToMove(strings.Split(line, " ")[1]))
			move := evaluation.FindBestMove(&b, 7)
			moveStr := move.MoveToString()
			b.PlayMove(move)

			fmt.Println(moveStr)
		case line == "playurself":
			for {
				move := evaluation.FindBestMove(&b, 7)
				fmt.Println(move.MoveToString())
				b.PlayMove(move)
			}
		case line == "quit":
			return
		case line == "test":
			fmt.Println(evaluation.Evaluate(&b))
			undo := moves.NewMove(1, 16, moves.FlagNone)
			b.PlayMove(undo)
			fmt.Println(evaluation.Evaluate(&b))
			b.UndoMove(undo)
			undo = moves.NewMove(1, 18, moves.FlagNone)
			b.PlayMove(undo)
			fmt.Println(evaluation.Evaluate(&b))
			b.UndoMove(undo)
		}
	}

	fmt.Fprintln(os.Stderr, "Loop exited!")
	if err := reader.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	} else {
		fmt.Fprintln(os.Stderr, "EOF received")
	}

	//fmt.Println(b.Mailbox)
	//fmt.Println(Perft(&b, 6))

	//fmt.Println(MoveToString(moves.NewMove(51, 59, moves.FlagPromotionKnight), &b))
	//b.PlayMove(move)
	//b.DebugPrint()
	//evaluation.BestMove(&b)
}
