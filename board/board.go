package board

import (
	"fmt"
	"math/bits"
	"slices"
)

const (
	FileA  Bitboard = 0x0101010101010101
	FileH  Bitboard = 0x8080808080808080
	FileAB Bitboard = 0x303030303030303
	FileGH Bitboard = 0xC0C0C0C0C0C0C0C0

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
	Knight [64]Bitboard
	Rook   [64]Bitboard
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
	if !b.Turn {
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
	for rank := 0; rank < 8; rank++ {
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

/*func (b *Board) FollowRay(turn bool, piecetype int, direction int, recurringmoves *[]Move) *[]Move {
	var raymoves []Move = *recurringmoves
	if recurringmoves == nil {
		raymoves = make([]Move, 1)
	}

	if turn {
		if piecetype == 1 {
			if direction == 1 {
				for (b.WRooks << 8 &^ b.FilledSquares) > 0 {
					b.WRooks = b.WRooks << 8 &^ b.FilledSquares
					append(raymoves)
					b.FollowRay(turn, piecetype, direction, &raymoves)
				}
			}
		}
	}

	return &raymoves
}*/

func RookDepth(startsquare int, depth int) []Bitboard {
	var b Bitboard = 0
	b.Set(startsquare)

	var bbs []Bitboard = make([]Bitboard, 0)

	north := b << (8 * depth) // &^ b.w
	south := b >> (8 * depth)
	east := (b << depth) &^ FileA
	west := (b >> depth) &^ FileH

	bbs = append(bbs, north)
	bbs = append(bbs, south)
	bbs = append(bbs, east)
	bbs = append(bbs, west)

	return bbs
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

		l1 := (b >> 1) & Bitboard(0x7f7f7f7f7f7f7f7f)
		l2 := (b >> 2) & Bitboard(0x3f3f3f3f3f3f3f3f)
		r1 := (b << 1) & Bitboard(0xfefefefefefefefe)
		r2 := (b << 2) & Bitboard(0xfcfcfcfcfcfcfcfc)
		h1 := l1 | r1
		h2 := l2 | r2

		precomp.Knight[i] = (h1 << 16) | (h1 >> 16) | (h2 << 8) | (h2 >> 8)
		precomp.King[i] = north | south | east | west | northEast | northWest | southEast | southWest
		precomp.Rook[i] = north | south | east | west
	}

	return precomp
}

var precomped *preCompTables = GeneratePrecomputedTables()

func TrailingZerosLoop(b Bitboard) []int {
	var squareslist []int
	bb := uint64(b)
	for bb != 0 {
		square := bits.TrailingZeros64(bb)
		squareslist = append(squareslist, square)

		bb &= bb - 1
	}
	return squareslist
}

func (b *Board) IsInCheck() bool {

	return false
}

func RecurringRookDepth(bb *Board, turn bool, moves *[]Move) *[]Move {
	var recurringrookmoves *[]Move = moves

	var squares []int
	if !turn {
		squares = TrailingZerosLoop(bb.BRooks)
	} else {
		squares = TrailingZerosLoop(bb.WRooks)
	}

	for _, startsquare := range squares {
		if startsquare < 64 {
			blocking := false
			blockeddir := 1
			for n := range 8 {
				moveboards := RookDepth(startsquare, n)
				if len(moveboards) > 0 {
					for i := range 4 {
						targetsquare := bits.TrailingZeros64(uint64(moveboards[i]))
						//fmt.Println(blocking && ((targetsquare-blockeddir)%8 == 0))
						fmt.Println(blocking, targetsquare, blockeddir)
						if blocking && ((targetsquare-blockeddir)%8 == 0) {
							continue
						}
						if targetsquare < 64 {
							pieceat := bb.PieceAt(targetsquare)
							if pieceat > 5 || pieceat < 0 {
								//fmt.Println(fmt.Sprintf("%d wants to go to %d", startsquare, targetsquare))
								Move := Move{startsquare, targetsquare}
								*recurringrookmoves = append(*recurringrookmoves, Move)
								//blocking = false
							} else {
								blockeddir = targetsquare
								blocking = true
							}
						}
					}
				}
			}
		}
	}
	return recurringrookmoves
}

func (b *Board) GenMoves() []Move {
	allMoves := []Move{}
	ourpieces, opponentpieces := b.Pieces()

	if b.Turn {
		if b.IsInCheck() {
			return nil
		}
		for i, bb := range []Bitboard{b.WKings, b.WQueens, b.WRooks, b.WBishops, b.WKnights, b.WPawns} {
			switch i {
			case 0:
				square := bits.TrailingZeros64(uint64(bb))
				if square < 64 {
					adjacentking := precomped.King[square] &^ ourpieces
					after := adjacentking.ToSquares()
					for _, v := range after {
						Move := Move{square, v}
						allMoves = append(allMoves, Move)
					}
				}

			case 2:
				var placeholder []Move = make([]Move, 0)
				n := RecurringRookDepth(b, b.Turn, &placeholder)
				allMoves = slices.Concat(allMoves, *n)
			case 4:
				squares := TrailingZerosLoop(bb)
				for _, square := range squares {
					if square < 64 {
						adjacentknight := precomped.Knight[square] &^ ourpieces
						after := adjacentknight.ToSquares()
						for _, v := range after {
							Move := Move{square, v}
							allMoves = append(allMoves, Move)
						}
					}
				}

			case 5:
				push1pawns := (bb >> 8) &^ b.FilledSquares
				after := push1pawns.ToSquares()
				for _, v := range after {
					Move := Move{v + 8, v}
					allMoves = append(allMoves, Move)
				}
				push2pawns := ((bb & WPawnStartRank) >> 16) &^ b.FilledSquares
				after = push2pawns.ToSquares()
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
			}
		}
	}

	return allMoves
}
