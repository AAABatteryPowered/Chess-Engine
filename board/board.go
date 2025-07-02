package board

import (
	"fmt"
	"math/bits"
)

const (
	FileA Bitboard = 0x0101010101010101
	FileH Bitboard = 0x8080808080808080

	WPawnStartRank Bitboard = 0x00FF000000000000
)

type BoardMethods interface {
	FromFen(string)
	GenMoves() []Move
	SetTurn(bool)
}

type Board struct {
	WKings   Bitboard
	WQueens  Bitboard
	WRooks   Bitboard
	WBishops Bitboard
	WKnights Bitboard
	WPawns   Bitboard

	BKings   Bitboard
	BQueens  Bitboard
	BRooks   Bitboard
	BBishops Bitboard
	BKnights Bitboard
	BPawns   Bitboard

	FilledSquares Bitboard
	Turn          bool
}

type preCompTables struct {
	King   [64]Bitboard
	Rook   [64]Bitboard
	Knight [64]Bitboard
	Bishop [64]Bitboard
}

type Move struct {
	From int
	To   int
}

type Bitboard uint64

func (b *Bitboard) Set(pos int) {
	*b |= 1 << pos
}

func (b *Bitboard) Clear(pos int) {
	*b &^= 1 << pos
}

func (b *Bitboard) Toggle(pos int) {
	*b ^= 1 << pos
}

func (b *Board) Pieces() (Bitboard, Bitboard) {
	if b.Turn {
		return b.BKings | b.BQueens | b.BRooks | b.BBishops | b.BKnights | b.BPawns, b.WKings | b.WQueens | b.WRooks | b.WBishops | b.WKnights | b.WPawns
	} else {
		return b.WKings | b.WQueens | b.WRooks | b.WBishops | b.WKnights | b.WPawns, b.BKings | b.BQueens | b.BRooks | b.BBishops | b.BKnights | b.BPawns
	}
}

func (b Bitboard) IsSet(pos int) bool {
	return (b>>pos)&1 == 1
}

func (b *Board) SetTurn(t bool) {
	b.Turn = t
}

func (b *Board) FromFen(s string) {
	posPointer := 0
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			num := int(ch - '0')
			posPointer += num
		}
		switch ch {
		case 'K':
			b.WKings.Set(posPointer)
			b.FilledSquares.Set(posPointer)
			posPointer += 1
		case 'Q':
			b.WQueens.Set(posPointer)
			b.FilledSquares.Set(posPointer)
			posPointer += 1
		case 'R':
			b.WRooks.Set(posPointer)
			b.FilledSquares.Set(posPointer)
			posPointer += 1
		case 'B':
			b.WBishops.Set(posPointer)
			b.FilledSquares.Set(posPointer)
			posPointer += 1
		case 'N':
			b.WKnights.Set(posPointer)
			b.FilledSquares.Set(posPointer)
			posPointer += 1
		case 'P':
			b.WPawns.Set(posPointer)
			b.FilledSquares.Set(posPointer)
			posPointer += 1

		case 'k':
			b.BKings.Set(posPointer)
			b.FilledSquares.Set(posPointer)
			posPointer += 1
		case 'q':
			b.BQueens.Set(posPointer)
			b.FilledSquares.Set(posPointer)
			posPointer += 1
		case 'r':
			b.BRooks.Set(posPointer)
			b.FilledSquares.Set(posPointer)
			posPointer += 1
		case 'b':
			b.BBishops.Set(posPointer)
			b.FilledSquares.Set(posPointer)
			posPointer += 1
		case 'n':
			b.BKnights.Set(posPointer)
			b.FilledSquares.Set(posPointer)
			posPointer += 1
		case 'p':
			b.BPawns.Set(posPointer)
			b.FilledSquares.Set(posPointer)
			posPointer += 1
		}
	}
}

func (b *Board) PieceAt(square int) int {
	allbb := b.AllBitboards()
	for index, board := range allbb {
		if board.IsSet(square) {
			return index
		}
	}
	return -1
}

func (b Bitboard) ToSquares() []int {
	var squares []int
	bb := uint64(b)
	for bb != 0 {
		square := bits.TrailingZeros64(bb)
		squares = append(squares, square)
		bb &= bb - 1
	}
	return squares
}

func (b *Board) AllBitboards() []Bitboard {
	return []Bitboard{
		b.WKings,
		b.WQueens,
		b.WRooks,
		b.WBishops,
		b.WKnights,
		b.WPawns,

		b.BKings,
		b.BQueens,
		b.BRooks,
		b.BBishops,
		b.BKnights,
		b.BPawns,
	}
}

func (b *Board) DebugPrint() {
	var finalstr string
	newlinesadded := 0
	allbbs := b.AllBitboards()
	for pos := 0; pos < 64; pos++ {
		piecefound := false
		for i, bb := range allbbs {
			if bb.IsSet(pos) {
				switch i {
				case 0:
					finalstr += "K "
				case 1:
					finalstr += "Q "
				case 2:
					finalstr += "R "
				case 3:
					finalstr += "B "
				case 4:
					finalstr += "N "
				case 5:
					finalstr += "P "
				case 6:
					finalstr += "k "
				case 7:
					finalstr += "q "
				case 8:
					finalstr += "r "
				case 9:
					finalstr += "b "
				case 10:
					finalstr += "n "
				case 11:
					finalstr += "p "
				}
				piecefound = true
				break
			}
		}

		if !piecefound {
			finalstr += ". "
		}

		if (len(finalstr)-newlinesadded)%16 == 0 {
			newlinesadded += 1
			finalstr += "\n"
		}
	}
	fmt.Println(finalstr)
}

func (b Bitboard) DebugPrint() {
	for rank := 7; rank >= 0; rank-- {
		for file := 0; file < 8; file++ {
			square := rank*8 + file
			if (b>>square)&1 == 1 {
				fmt.Print("1 ")
			} else {
				fmt.Print("0 ")
			}
		}
		fmt.Println()
	}
}

func GeneratePrecomputedTables() *preCompTables {
	precomp := &preCompTables{}

	for i := 0; i < 64; i++ {
		b := Bitboard(1 << i)

		north := b << 8
		south := b >> 8
		east := (b << 1) &^ FileA
		west := (b >> 1) &^ FileH
		northEast := (b << 9) &^ FileA
		northWest := (b << 7) &^ FileH
		southEast := (b >> 7) &^ FileA
		southWest := (b >> 9) &^ FileH

		precomp.King[i] = north | south | east | west | northEast | northWest | southEast | southWest
	}

	return precomp
}

var precomped *preCompTables = GeneratePrecomputedTables()

func (b *Board) GenMoves() []Move {
	allMoves := []Move{}
	opponentpieces, ourpieces := b.Pieces()

	if b.Turn {
		for i, bb := range []Bitboard{b.WKings, b.WQueens, b.WRooks, b.WBishops, b.WKnights, b.WPawns} {
			switch i {
			case 1:
				//adjacentking := precomped.King[i] &^ ourpieces
			case 5:
				push1pawns := (bb >> 8) &^ b.FilledSquares
				after := push1pawns.ToSquares()
				for _, v := range after {
					Move := Move{v + 8, v}
					allMoves = append(allMoves, Move)
				}
				push2pawns := ((bb & WPawnStartRank) >> 16) &^ b.FilledSquares
				after = push2pawns.ToSquares()
				//fmt.Println(after)
				for _, v := range after {
					Move := Move{v + 16, v}
					allMoves = append(allMoves, Move)
				}

				leftcapture := ((bb &^ FileH) >> 9) & opponentpieces
				after = leftcapture.ToSquares()
				for _, v := range after {
					Move := Move{v + 9, v}
					allMoves = append(allMoves, Move)
				}

				rightcapture := ((bb &^ FileA) >> 7) & opponentpieces
				after = rightcapture.ToSquares()
				for _, v := range after {
					Move := Move{v + 7, v}
					allMoves = append(allMoves, Move)
				}

				/*allMoves = append(allMoves, bb<<8)
				if !(b.FilledSquares.IsSet(pos - 8)) {
					allMoves = append(allMoves, Move{pos, pos - 8})
					if !(b.FilledSquares.IsSet(pos - 16)) {
						allMoves = append(allMoves, Move{pos, pos - 16})
					}
				}*/
			}
		}
	}

	return allMoves
}
