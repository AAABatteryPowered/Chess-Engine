package evaluation

import (
	"bot/board"
	"bot/moves"
	"math"
	"math/bits"
)

var pieceValues [5]int = [5]int{900, 500, 320, 301, 100}
var PieceSquareTables = [6][64]int{
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

//var TranspositionTables

var edges board.Bitboard = 0xff818181818181ff
var center board.Bitboard = 0x1818000000

//var kingHeatMap [64]int =

func Evaluate(b *board.Board) int {
	var score int = 0

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

	return score
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

func Search(b *board.Board, depth int, alpha int, beta int) int {
	if depth == 0 {
		return SearchAllCaptures(b, alpha, beta)
	}

	moves := b.Moves(false)
	for i := 0; i < moves.Count; i++ {
		move := moves.Moves[i]

		b.PlayMove(move)
		value := -Search(b, depth-1, -beta, -alpha)
		b.UndoMove(move)

		if value >= beta {
			return value
		}

		alpha = max(alpha, value)
	}

	return alpha
}

func SearchAllCaptures(b *board.Board, alpha int, beta int) int {
	eval := Evaluate(b)
	if eval >= beta {
		return beta
	}
	alpha = max(alpha, eval)

	capturemoves := b.Moves(true)
	for i := 0; i < capturemoves.Count; i++ {
		move := capturemoves.Moves[i]

		b.PlayMove(move)
		value := -SearchAllCaptures(b, -beta, -alpha)
		b.UndoMove(move)

		if value >= beta {
			return value
		}

		alpha = max(alpha, value)
	}

	return alpha
}

func FindBestMove(b *board.Board, depth int) moves.Move {

	var bestMove moves.Move
	bestValue := math.MinInt

	moves := b.Moves(false)
	//OrderMoves(b, &moves)

	for i := 0; i < moves.Count; i++ {
		move := moves.Moves[i]

		b.PlayMove(move)
		moveValue := -Search(b, depth-1, math.MinInt, math.MaxInt)
		b.UndoMove(move)

		if moveValue > bestValue {
			bestValue = moveValue
			bestMove = move
		}
	}

	return bestMove
}
