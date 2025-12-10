package evaluation

import (
	"bot/board"
	"bot/moves"
	"math/bits"
)

var pieceValues [5]int = [5]int{900, 500, 320, 301, 100}
var PieceSquareTables = [6][64]int{
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
	//Queen
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
}

type EntryFlag int

const (
	Exact EntryFlag = iota
	Alpha
	Beta
)

type TTEntry struct {
	Depth int
	Score int
	Flag  EntryFlag // exact, alpha, beta
}

var TranspositionTable map[board.Bitboard]TTEntry = make(map[board.Bitboard]TTEntry)

func StoreTT(hash board.Bitboard, entry TTEntry) {
	TranspositionTable[hash] = entry
}

func LookupTT(hash board.Bitboard, depth int) (TTEntry, bool) {
	entry, ok := TranspositionTable[hash]
	if ok && entry.Depth >= depth {
		return entry, true
	}
	return TTEntry{}, false
}

var edges board.Bitboard = 0xff818181818181ff
var center board.Bitboard = 0x1818000000

//var kingHeatMap [64]int =

func Evaluate(b *board.Board) int {
	var score int = 0

	for i := 0; i < 6; i++ {
		for bb := *b.AllBitboards[i]; bb != 0; bb &= bb - 1 {
			square := bits.TrailingZeros64(uint64(bb))
			if i == 5 {
				//fmt.Println("Piece", i, "square", square, "PST", PieceSquareTables[i][square])
			}

			score += PieceSquareTables[i][square] //pieceValues[i] +
		}
	}

	for i := 6; i < 12; i++ {
		for bb := *b.AllBitboards[i]; bb != 0; bb &= bb - 1 {
			square := bits.TrailingZeros64(uint64(bb))
			score -= PieceSquareTables[i-6][square] //pieceValues[i-6] +
		}
	}

	return score * -1
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
	ogalpha := alpha
	if entry, ok := LookupTT(b.Hash, depth); ok {
		switch entry.Flag {
		case Exact:
			return entry.Score
		case Alpha:
			if entry.Score <= alpha {
				return alpha
			}
		case Beta:
			if entry.Score >= beta {
				return beta
			}
		}
	}

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
			StoreTT(b.Hash, TTEntry{
				Depth: depth,
				Score: value,
				Flag:  Beta,
			})
			return value
		}

		alpha = max(alpha, value)
	}

	flag := Exact
	if alpha <= ogalpha {
		flag = Alpha
	} else if alpha >= beta {
		flag = Beta
	}

	StoreTT(b.Hash, TTEntry{
		Depth: depth,
		Score: alpha,
		Flag:  flag,
	})

	return alpha
}

func SearchAllCaptures(b *board.Board, alpha int, beta int) int {
	if entry, ok := LookupTT(b.Hash, 1); ok {
		switch entry.Flag {
		case Exact:
			return entry.Score
		case Alpha:
			if entry.Score <= alpha {
				return alpha
			}
		case Beta:
			if entry.Score >= beta {
				return beta
			}
		}
	}

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
	alpha := -9999
	beta := 9999

	moves := b.Moves(false)
	//OrderMoves(b, &moves)

	for i := 0; i < moves.Count; i++ {
		move := moves.Moves[i]

		b.PlayMove(move)
		moveValue := -Search(b, depth-1, -beta, -alpha)
		//fmt.Println("Move:", move.MoveToString(), "Value:", moveValue)
		b.UndoMove(move)

		if moveValue > alpha {
			alpha = moveValue
			bestMove = move
		}
	}

	return bestMove
}
