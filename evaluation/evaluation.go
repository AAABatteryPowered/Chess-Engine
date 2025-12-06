package evaluation

import (
	"bot/board"
	"bot/moves"
	"fmt"
	"math"
	"math/bits"
	"math/rand"
)

var pieceValues [5]float32 = [5]float32{900, 500, 320, 301, 100}
var PieceSquareTables = [6][64]float32{
	// Pawn
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		5, 10, 10, -20, -20, 10, 10, 5,
		5, -5, -10, 0, 0, -10, -5, 5,
		0, 0, 0, 20, 20, 0, 0, 0,
		5, 5, 10, 25, 25, 10, 5, 5,
		10, 10, 20, 30, 30, 20, 10, 10,
		50, 50, 50, 50, 50, 50, 50, 50,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	// Knight
	{
		-50, -40, -30, -30, -30, -30, -40, -50,
		-40, -20, 0, 5, 5, 0, -20, -40,
		-30, 5, 10, 15, 15, 10, 5, -30,
		-30, 0, 15, 20, 20, 15, 0, -30,
		-30, 5, 15, 20, 20, 15, 5, -30,
		-30, 0, 10, 15, 15, 10, 0, -30,
		-40, -20, 0, 0, 0, 0, -20, -40,
		-50, -40, -30, -30, -30, -30, -40, -50,
	},
	// Bishop
	{
		-20, -10, -10, -10, -10, -10, -10, -20,
		-10, 5, 0, 0, 0, 0, 5, -10,
		-10, 10, 10, 10, 10, 10, 10, -10,
		-10, 0, 10, 10, 10, 10, 0, -10,
		-10, 5, 5, 10, 10, 5, 5, -10,
		-10, 0, 5, 10, 10, 5, 0, -10,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-20, -10, -10, -10, -10, -10, -10, -20,
	},
	// Rook
	{
		0, 0, 0, 5, 5, 0, 0, 0,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		5, 10, 10, 10, 10, 10, 10, 5,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	// Queen
	{
		-20, -10, -10, -5, -5, -10, -10, -20,
		-10, 0, 5, 0, 0, 0, 0, -10,
		-10, 5, 5, 5, 5, 5, 0, -10,
		0, 0, 5, 5, 5, 5, 0, -5,
		-5, 0, 5, 5, 5, 5, 0, -5,
		-10, 0, 5, 5, 5, 5, 0, -10,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-20, -10, -10, -5, -5, -10, -10, -20,
	},
	// King
	{
		20, 30, 10, 0, 0, 10, 30, 20,
		20, 20, 0, 0, 0, 0, 20, 20,
		-10, -20, -20, -20, -20, -20, -20, -10,
		-20, -30, -30, -40, -40, -30, -30, -20,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
	},
}

var edges board.Bitboard = 0xff818181818181ff
var center board.Bitboard = 0x1818000000

//var kingHeatMap [64]int =

func Evaluate(b *board.Board) float64 {
	var score float32 = 0

	for i := 0; i < 5; i++ {
		for bb := *b.AllBitboards[i]; bb != 0; bb &= bb - 1 {
			square := bits.TrailingZeros64(uint64(bb))
			score += pieceValues[i] + PieceSquareTables[i][square]
		}
	}

	for i := 6; i < 11; i++ {
		for bb := *b.AllBitboards[i]; bb != 0; bb &= bb - 1 {
			square := bits.TrailingZeros64(uint64(bb))
			// Flip square for black pieces
			flippedSquare := square ^ 56
			score -= pieceValues[i-6] + PieceSquareTables[i-6][flippedSquare]
		}
	}

	return float64(score)
}

func evaluatePiecePositions(b *board.Board) int {
	score := 0

	score += bits.OnesCount64(uint64(b.WKnights&center)) * 10
	score -= bits.OnesCount64(uint64(b.BKnights&center)) * 10

	score += bits.OnesCount64(uint64(b.WPawns&center)) * 5
	score -= bits.OnesCount64(uint64(b.BPawns&center)) * 5

	score -= bits.OnesCount64(uint64(b.WKnights&edges)) * 20
	score += bits.OnesCount64(uint64(b.BKnights&edges)) * 20

	return score
}

func BestMove(b *board.Board) {
	Moves := b.Moves()

	fmt.Println(Moves.Moves[rand.Intn(Moves.Count)])
}

func MiniMax(b *board.Board, depth int, alpha float64, beta float64, maximising bool) float64 {
	if depth == 0 {
		return Evaluate(b)
	}

	moves := b.Moves()

	if maximising {
		value := math.Inf(-1)
		for i := 0; i < moves.Count; i++ {
			move := moves.Moves[i]
			b.PlayMove(move)
			value = math.Max(value, float64(MiniMax(b, depth-1, alpha, beta, false)))
			b.UndoMove(move)
			//b.UndoCount--
			alpha = math.Max(alpha, value)
			if alpha >= beta {
				break
			}

		}
		return value
	} else {
		value := math.Inf(1)
		for i := 0; i < moves.Count; i++ {
			move := moves.Moves[i]
			b.PlayMove(move)
			value = min(value, MiniMax(b, depth-1, alpha, beta, true))
			b.UndoMove(move)
			//b.UndoCount--
			beta = math.Min(beta, value)
			if alpha >= beta {
				break
			}

		}
		return value
	}
}

func FindBestMove(b *board.Board, depth int) moves.Move {
	var bestMove moves.Move
	bestValue := math.Inf(-1)
	alpha := math.Inf(-1)
	beta := math.Inf(1)

	moves := b.Moves()

	for i := 0; i < moves.Count; i++ {
		move := moves.Moves[i]

		b.PlayMove(move)
		moveValue := MiniMax(b, depth-1, alpha, beta, false) // false because opponent's turn
		b.UndoMove(move)

		if moveValue > bestValue {
			bestValue = moveValue
			bestMove = move
		}

		alpha = math.Max(alpha, bestValue)
	}

	return bestMove
}
